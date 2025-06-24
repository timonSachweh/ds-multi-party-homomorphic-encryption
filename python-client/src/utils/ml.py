import config
import torch
from torchvision import transforms


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