import json
import random

import config
from flask import Flask
import logging
from flask import request
import numpy as np
from utils import disable_logging

app = Flask(__name__)
app.logger.setLevel(logging.ERROR)
c = config.Config()

weights = np.full(20, 0.326, dtype=np.float32)
for i in range(10, len(weights)):
    weights[i] = random.random() * 2

@app.route('/health')
@disable_logging
def health_route():
    return {
        'message': 'Service is ready'
    }

@app.route(c.server.api + '/train')
def train_route():
    return {
        'message': 'Training complete!'
    }

@app.route(c.server.api + '/model-params', methods=['GET', 'POST'])
def model_params_route():
    if request.method == 'GET':
        response = json.dumps({
            "model_name": c.model.name,
            "version": c.model.version,
            "weights": weights.tolist()
        })
        return response
    elif request.method == 'POST':
        body = request.get_json()
        print("original weights:")
        print(weights)
        print("new weights:")
        print(np.array(body['weights'], dtype=np.float32))
        return {
            'message': 'Model parameters updated!'
        }

if __name__ == '__main__':
    app.run(debug=c.server.debug, host="0.0.0.0", port=c.server.port)