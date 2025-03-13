import json
import os
import sys

import torch
import config
from flask import Flask
import logging
from flask import request
import ml
import numpy as np
from utils import disable_logging

app = Flask(__name__)
app.logger.setLevel(logging.ERROR)
c = config.Config()

model_service = ml.ModelService(c.model)

@app.route('/health')
@disable_logging
def health_route():
    return {
        'message': 'Service is ready'
    }

@app.route(c.server.api + '/train')
def train_route():
    model_service.train()
    return {
        'message': 'Training complete!'
    }

@app.route(c.server.api + '/model-params', methods=['GET', 'POST'])
def model_params_route():
    if request.method == 'GET':
        response = json.dumps({
            "model_name": c.model.name,
            "version": c.model.version,
            "weights": model_service.get_model_params().tolist()
        })
        return response
    elif request.method == 'POST':
        body = request.get_json()
        model_service.set_model_params(np.array(body['weights'], dtype=np.float32))
        return {
            'message': 'Model parameters updated!'
        }

@app.route(c.server.api + '/predict', methods=['POST'])
def predict_route():
    body = request.get_json()
    return {
        'prediction': model_service.predict(np.array(body['data'], dtype=np.float32))
    }

@app.route(c.server.api + '/about')
def show_about():
    """
    Get deployment information, for debugging
    """

    def bash(command):
        output = os.popen(command).read()
        return output

    return {
        "sys.version": sys.version,
        "torch.__version__": torch.__version__,
        "torch.mps.is_avaliable()": torch.backends.mps.is_available(),
        "torch.cuda.is_available()": torch.cuda.is_available(),
        "torch.version.cuda": torch.version.cuda,
        "torch.backends.cudnn.version()": torch.backends.cudnn.version(),
        "torch.backends.cudnn.enabled": torch.backends.cudnn.enabled,
        "nvidia-smi": bash('nvidia-smi')
    }

if __name__ == '__main__':
    app.run(debug=c.server.debug, host="0.0.0.0", port=c.server.port)