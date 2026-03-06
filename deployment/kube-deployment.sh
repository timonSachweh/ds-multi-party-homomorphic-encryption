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
    export NODE_PORT=$($skube get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services ds-aggregation-server)

    sleep 20
    for (( i=1; i<=$1; i++ ))
    do
        $shelm install $prefix-c${i} ${BASE_PATH}/deployment/he-client --set clientId=$((i-1)) --set totalClients=$1 --set service.aggregationService=http://$prefix-aggregation-server:8080/v1
        sleep 60
    done

    sleep 2

    echo "Running sessions:"
    $shelm ls
}

# Function to stop services
stop_services() {
    echo "Stopping services..."
    $shelm delete $($shelm ls -q | grep ds)

    echo "All services have been stopped."
}

request_training() {
    echo "Requesting training"
    curl -X GET http://localhost:$NODE_PORT/v1/clients/train
}

save_logs() {
    local outdir="$BASE_PATH/logs/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$outdir"

    $skube get pods --all-namespaces --no-headers \
        | awk '$2 ~ /^ds/ {print $1 " " $2}' \
        | while read -r ns name; do
            echo "Processing $name"

            # containers
            containers=$($skube get pod -n "$ns" "$name" -o jsonpath='{.spec.containers[*].name}' 2>/dev/null)
            if [ -n "$containers" ]; then
                for c in $containers; do
                    $skube logs -n "$ns" "$name" -c "$c" > "$outdir/${name}.${c}.log" 2>&1 || true
                done
            fi
        done

    echo "Saved logs to $outdir"
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
elif [ "$1" = "train" ] || [ "$1" = "training" ] || [ "$1" = "request_training" ]; then
    request_training
elif [ "$1" = "save_logs" ] || [ "$1" = "save" ]; then
    save_logs
else
    echo "Usage: $0 [run|start|up] [num_clients]"
    echo "       $0 [stop|down]"
    exit 1
fi
