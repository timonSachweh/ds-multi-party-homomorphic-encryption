import torch
from torchvision import datasets, transforms

def get_datasets(transforms: transforms.Compose):
    return datasets.MNIST('../data', train=True, download=True, transform=transforms), \
        datasets.MNIST('../data', train=False, transform=transforms)


def get_dataloaders(transforms: transforms.Compose, train_kwargs: dict = {}, test_kwargs: dict = {}):
    d1, d2 = get_datasets(transforms)

    return torch.utils.data.DataLoader(d1, **train_kwargs), \
        torch.utils.data.DataLoader(d2, **test_kwargs)