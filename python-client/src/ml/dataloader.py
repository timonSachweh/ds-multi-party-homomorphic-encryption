import torch
from enum import Enum, auto
from torchvision import datasets, transforms


class DatasetType(Enum):
    """Supported dataset types for the client.

    Use this enum when calling :func:`get_datasets` or
    :func:`get_dataloaders` to select between MNIST and CIFAR-10.
    """

    MNIST = auto()
    CIFAR10 = auto()


def get_datasets(tf: transforms.Compose, type: DatasetType = DatasetType.MNIST):
    """Return (train_dataset, test_dataset) for the requested type.

    Args:
        transforms: torchvision transforms to apply.
        type: which DatasetType to load.
    """
    if type == DatasetType.MNIST:
        return datasets.MNIST('../data', train=True, download=True, transform=tf), \
            datasets.MNIST('../data', train=False, transform=tf)
    elif type == DatasetType.CIFAR10:
        return datasets.CIFAR10('../data', train=True, download=True, transform=tf), \
            datasets.CIFAR10('../data', train=False, transform=tf)
    else:
        raise ValueError(f"Unsupported dataset type: {type}")


def get_dataloaders(
    tf: transforms.Compose,
    dataset_type: DatasetType = DatasetType.MNIST,
    party_idx: int = None,
    num_parties: int = None,
    train_kwargs: dict = {},
    test_kwargs: dict = {},
):
    d1, d2 = get_datasets(tf, type=dataset_type)
    if party_idx is not None and num_parties is not None and num_parties is not 0 and party_idx <= num_parties:
        # Split the dataset into num_parties parts
        train_split_len = len(d1) // num_parties
        test_split_len = len(d2) // num_parties
        train_range = range(train_split_len*party_idx, min(train_split_len*(party_idx+1), len(d1)))
        test_range = range(test_split_len*party_idx, min(test_split_len*(party_idx+1), len(d2)))
        d1 = torch.utils.data.Subset(d1, train_range)
        d2 = torch.utils.data.Subset(d2, test_range)


    return torch.utils.data.DataLoader(d1, **train_kwargs), \
        torch.utils.data.DataLoader(d2, **test_kwargs)