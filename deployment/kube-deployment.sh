#!/bin/bash

BASE_PATH="$(dirname "$(readlink -f "$0")")/.."

skube="microk8s kubectl"
shelm="microk8s helm"

prefix="ds"

# Function to start services
run_services() {
    echo "Starting services..."

    # Start a screen session for each service
    $shelm install $prefix ${BASE_PATH}/deployment/aggregation-server --set numClients=$1

    sleep 2
    for (( i=1; i<=$1; i++ ))
    do
        $shelm install $prefix-c${i} ${BASE_PATH}/deployment/he-client --set clientIndex=$((i-1)) --set totalClients=$1 --set service.aggregationService=http://$prefix-aggregation-server:8080/v1
        sleep 8
    done

    sleep 2

    sessions=$shelm ls
    echo "Running sessions:"
    echo "${sessions}"
}

# Function to stop services
stop_services() {
    echo "Stopping services..."
    $shelm delete $($shelm ls -q | grep ds)

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
