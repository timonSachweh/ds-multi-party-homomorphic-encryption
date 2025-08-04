#!/bin/zsh

BASE_PATH="$(dirname "$(readlink -f "$0")")/.."



screen -L -Logfile "${BASE_PATH}/logs/global-python.log" -dmS "global-python" zsh -c "\
            source '${BASE_PATH}/env/client1.env'; \
            export PYTHON_PORT=5001; \
            export DATA_SPLIT_PARTY=0; \
            export DATA_SPLIT_NUM_PARTIES=1; \
            uv run ${BASE_PATH}/python-client/src/main.py"