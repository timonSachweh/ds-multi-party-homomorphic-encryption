from functools import partial, wraps
import logging
from flask import request
import torch
from torchvision import transforms
import config

def get_kwargs(config: config.ModelConfig, device: torch.device):
    train_kwargs = {'batch_size': config.batch_size}
    test_kwargs = {'batch_size': config.test_batch_size}
    if device.type == 'cuda':
        cuda_kwargs = {'num_workers': 1,
                       'pin_memory': True,
                       'shuffle': True}
        train_kwargs.update(cuda_kwargs)
        test_kwargs.update(cuda_kwargs)
    return train_kwargs, test_kwargs


def get_transforms():
    return transforms.Compose([
            transforms.ToTensor(),
            transforms.Normalize((0.1307,), (0.3081,))
            ])


def get_device(no_cuda, no_mps):
    use_cuda = not no_cuda and torch.cuda.is_available()
    use_mps = not no_mps and torch.backends.mps.is_available()

    if use_cuda:
        device = torch.device("cuda")
    elif use_mps:
        device = torch.device("mps")
    else:
        device = torch.device("cpu")
    return device


class HealthCheckFilter(logging.Filter):
    """Filter for logging output"""

    def __init__(self, path, name=''):
        """Class constructor.
        We pass 'path' argument to instance  which is
        used by to filter logging for Flask routes.
        """
        self.path = path
        super().__init__(name)

    def filter(self, record):
        """Main filter function.
        We add a space after path here to ensure subpaths
        are not unintentionally excluded from logging"""
        return f"{self.path} " not in record.getMessage()


def disable_logging(func=None, *args, **kwargs):
    """Disable log messages for werkzeug log handler
    for a specific Flask routes.

    :param (function) func: wrapped function
    :param (list) args: decorator arguments
    :param (dict) kwargs: decorator keyword arguments
    :return (function) wrapped function
    """
    _logger = 'werkzeug'
    if not func:
        return partial(disable_logging, *args, **kwargs)

    @wraps(func)
    def wrapper(*args, **kwargs):
        path = request.environ['PATH_INFO']
        log = logging.getLogger(_logger)
        log.addFilter(HealthCheckFilter(path))
        return func(*args, **kwargs)
    return wrapper