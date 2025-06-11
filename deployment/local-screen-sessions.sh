#!/bin/zsh

BASE_PATH="$(dirname "$(readlink -f "$0")")/.."

# Function to start services
run_services() {
    echo "Starting services..."

    # Start a screen session for each service
    screen -dmS "he-aggregation-service" zsh -c "\
        source '${BASE_PATH}/env/server.env'; \
        go run ${BASE_PATH}/aggregation-server/cmd/main.go"

    sleep 2
    for (( i=1; i<=$1; i++ ))
    do
        GO_PORT=$((8080+i))
        PYTHON_PORT=$((9090+i))
        PARTY_IDX=$((i-1))
        echo "Starting client $i with ports: ${GO_PORT} and ${PYTHON_PORT}..."
        screen -dmS "he-client-${i}-go" zsh -c "\
            source '${BASE_PATH}/env/client1.env'; \
            export HTTP_PORT=${GO_PORT}; \
            export EXTERNAL_URL=http://localhost:${GO_PORT}; \
            export PYTHON_PORT=${PYTHON_PORT}; \
            export DATA_SPLIT_PARTY=${PARTY_IDX}; \
            export DATA_SPLIT_NUM_PARTIES=$1; \
            go run ${BASE_PATH}/he-client/cmd/main.go"
        screen -dmS "he-client-${i}-python" zsh -c "\
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

    screen -ls | grep -E "he-client-[0-9]+-(go|python)" | awk '{print $1}' | while read -r session; do
        echo "Stopping session: $session"
        screen -S "$session" -p 0 -X stuff $'\003'
    done
    
    echo "All services have been stopped."
}

# Check the first parameter
if [ "$1" = "start" ] || [ "$1" = "run" ] || [ "$1" = "up" ]; then
    NUM_CLIENTS=2
    if [ -n "$2" ]; then
        NUM_CLIENTS=$2
    fi
    run_services $NUM_CLIENTS
elif [ "$1" = "stop" ] || [ "$1" = "down" ]; then
    stop_services
else
    echo "Usage: $0 [run|start|up] [num_clients]"
    echo "       $0 [stop|down]"
    exit 1
fi