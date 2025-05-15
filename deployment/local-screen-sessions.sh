#!/bin/bash

BASE_PATH="$(dirname "$(readlink -f "$0")")/.."

# Function to start services
run_services() {
    echo "Starting services..."

    # Start a screen session for each service
    screen -dmS "he-aggregation-service" bash -c "\
        source '${BASE_PATH}/env/server.env'; \
        go run ${BASE_PATH}/aggregation-server/cmd/main.go"

    sleep 2
    for (( i=1; i<=$1; i++ ))
    do
        GO_PORT=$((8080+i))
        PYTHON_PORT=$((9090+i))
        PARTY_IDX=$((i-1))
        echo "Starting client $i with ports: ${GO_PORT} and ${PYTHON_PORT}..."
        screen -dmS "he-client-${i}-go" bash -c "\
            source '${BASE_PATH}/env/client1.env'; \
            export HTTP_PORT=${GO_PORT}; \
            export EXTERNAL_URL=http://localhost:${GO_PORT}; \
            export PYTHON_PORT=${PYTHON_PORT}; \
            export DATA_SPLIT_PARTY=${PARTY_IDX}; \
            export DATA_SPLIT_NUM_PARTIES=$1; \
            go run ${BASE_PATH}/he-client/cmd/main.go"
        screen -dmS "he-client-${i}-python" bash -c "\
            source '${BASE_PATH}/env/client1.env'; \
            export HTTP_PORT=${GO_PORT}; \
            export EXTERNAL_URL=http://localhost:${GO_PORT}; \
            export PYTHON_PORT=${PYTHON_PORT}; \
            export DATA_SPLIT_PARTY=${PARTY_IDX}; \
            export DATA_SPLIT_NUM_PARTIES=$1; \
            pyenv activate ml; \
            python ${BASE_PATH}/python-client/src/main.py"
        sleep 8
    done
    
    sleep 2
    
    RUNNING_SESSIONS="$(screen -ls)"
    echo "Running sessions:"
    echo "${RUNNING_SESSIONS}"
}

# Function to stop services
stop_services() {
    echo "Stopping services..."
    # Kill the screen session
    screen -S "he-aggregation-service" -p 0 -X stuff $'\003'

    screen -S "he-client-1-go" -p 0 -X stuff $'\003'
    screen -S "he-client-1-python" -p 0 -X stuff $'\003'

    screen -S "he-client-2-go" -p 0 -X stuff $'\003'
    screen -S "he-client-2-python" -p 0 -X stuff $'\003'
    
    echo "All services have been stopped."
}

# Check the first parameter
if [ "$1" = "run" ]; then
    NUM_CLIENTS=2
    if [ -n "$2" ]; then
        NUM_CLIENTS=$2
    fi
    run_services $NUM_CLIENTS
elif [ "$1" = "down" ]; then
    stop_services
else
    echo "Usage: $0 [run|down]"
    exit 1
fi