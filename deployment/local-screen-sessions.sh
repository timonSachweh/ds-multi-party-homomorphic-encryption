#!/bin/bash

main_path="$(dirname "$0")/../"

# Function to start services
run_services() {
    echo "Starting services..."
    # Start a screen session for each service
    screen -dmS "he-aggregation-service" bash -c "source '${main_path}/env/server.env'; go run ${main_path}aggregation-server/cmd/main.go"

    screen -dmS "he-client-1-go" bash -c "
    source '${main_path}/env/client1.env'; go run ${main_path}he-client/cmd/main.go"
    screen -dmS "he-client-1-python" bash -c "source '${main_path}/env/client1.env'; pyenv activate ml; python ${main_path}python-client/src/main.py"

    screen -dmS "he-client-2-go" bash -c "
    source '${main_path}/env/client2.env'; go run ${main_path}he-client/cmd/main.go"
    screen -dmS "he-client-2-python" bash -c "source '${main_path}/env/client2.env'; pyenv activate ml; python ${main_path}python-client/src/main.py"

    
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
    run_services
elif [ "$1" = "down" ]; then
    stop_services
else
    echo "Usage: $0 [run|down]"
    exit 1
fi