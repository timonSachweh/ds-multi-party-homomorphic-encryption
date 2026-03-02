import torch
import torch.nn as nn
import torch.nn.functional as F
import numpy as np


class MnistNet(nn.Module):
    def __init__(self):
        super(MnistNet, self).__init__()
        self.conv1 = nn.Conv2d(1, 32, 3, 1)
        self.conv2 = nn.Conv2d(32, 64, 3, 1)
        self.dropout1 = nn.Dropout(0.25)
        self.dropout2 = nn.Dropout(0.5)
        self.fc1 = nn.Linear(9216, 128)
        self.fc2 = nn.Linear(128, 10)

    def forward(self, x):
        x = self.conv1(x)
        x = F.relu(x)
        x = self.conv2(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        x = self.dropout1(x)
        x = torch.flatten(x, 1)
        x = self.fc1(x)
        x = F.relu(x)
        x = self.dropout2(x)
        x = self.fc2(x)
        output = F.log_softmax(x, dim=1)
        return output

class LeNet(nn.Module):
    def __init__(self):
        super(LeNet, self).__init__()
        self.cnn_layers = torch.nn.Sequential(
            torch.nn.Conv2d(in_channels=3, out_channels=128,
                      kernel_size=3, stride=2, padding=1),
            torch.nn.ReLU(True),
            torch.nn.BatchNorm2d(128),
            torch.nn.Conv2d(in_channels=128, out_channels=32,
                      kernel_size=3, stride=2, padding=1),
            torch.nn.ReLU(True),
            torch.nn.BatchNorm2d(32),
            torch.nn.Conv2d(in_channels=32, out_channels=16,
                      kernel_size=3, stride=2, padding=1),
            torch.nn.BatchNorm2d(16),
            torch.nn.Flatten()
        )
        
        self.mlp_layers = torch.nn.Sequential(
            torch.nn.Linear(256, 50),
            torch.nn.ReLU(True),
            torch.nn.Linear(50, 10),
            torch.nn.LogSoftmax(dim=1)
        )

    def forward(self, x):
        x = self.cnn_layers(x)
        y = self.mlp_layers(x)
        return y