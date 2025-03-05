
import os


class Config:
    def __init__(self):
        self.server = ServerConfig()
        self.model = ModelConfig()

class ServerConfig:
    def __init__(self):
        self.api = os.getenv("PYTHON_API_PATH", "/api")
        self.debug = bool(os.getenv("DEBUG", True))
    
class ModelConfig:
    def __init__(self):
        self.name = os.getenv("ML_MODEL_NAME", "mnist")
        self.version = os.getenv("ML_MODEL_VERSION", "1.0")
        self.model_path = os.getenv("ML_MODEL_PATH", "./model.pt")
        self.batch_size = int(os.getenv("ML_BATCH_SIZE", 64))
        self.test_batch_size = int(os.getenv("ML_TEST_BATCH_SIZE", 1000))
        self.epochs = int(os.getenv("ML_EPOCHS", 1))
        self.lr = float(os.getenv("ML_LR", 1.0))
        self.gamma = float(os.getenv("ML_GAMMA", 0.7))
        self.no_cuda = bool(os.getenv("ML_NO_CUDA", False))
        self.no_mps = bool(os.getenv("ML_NO_MPS", False))
        self.seed = int(os.getenv("ML_SEED", 1))
        self.log_interval = int(os.getenv("ML_LOG_INTERVAL", 10))