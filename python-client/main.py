import json
import config
from flask import Flask
from flask import request
import ml
import numpy as np

app = Flask(__name__)
c = config.Config()

model_service = ml.ModelService(c.model)

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

if __name__ == '__main__':
    app.run(debug=c.server.debug, port=c.server.port)