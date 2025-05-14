import torch
from torchvision import datasets, transforms

def get_datasets(transforms: transforms.Compose):
    return datasets.MNIST('../data', train=True, download=True, transform=transforms), \
        datasets.MNIST('../data', train=False, transform=transforms)


def get_dataloaders(transforms: transforms.Compose, party_idx: int = None, num_parties: int = None, train_kwargs: dict = {}, test_kwargs: dict = {}):
    d1, d2 = get_datasets(transforms)
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