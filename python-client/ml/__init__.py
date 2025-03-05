import torch
from torch.optim.lr_scheduler import StepLR
import torch.nn.functional as F
import config
from .model import Net
from .utils import get_device, get_transforms
from .dataloader import get_dataloaders


class ModelService():
    def __init__(self, model_config: config.ModelConfig):   
        self.config = model_config
        self.device = get_device(self.config.no_cuda, self.config.no_mps)
        self.model = Net().to(self.device)
        self.transform = get_transforms()
    
    def train(self):
        torch.manual_seed(self.config.seed)
        
        train_kwargs, test_kwargs = utils.get_kwargs(self.config, self.device)
        
        train_loader, test_loader = get_dataloaders(transforms=self.transform, train_kwargs=train_kwargs, test_kwargs=test_kwargs)

        self.model = self.model.to(self.device)
        optimizer = torch.optim.Adadelta(self.model.parameters(), lr=self.config.lr)
        scheduler = StepLR(optimizer, step_size=1, gamma=self.config.gamma)
        for epoch in range(1, self.config.epochs + 1):
            self.train_model(train_loader, optimizer, epoch)
            self.test_model(test_loader)
            scheduler.step()
        
        torch.save(self.model.state_dict(), self.config.model_path)
    
    def predict(self, data):
        with torch.no_grad():
            data = torch.tensor(data).to(self.device)
            data = data.reshape(1, 1, 28, 28)
            output = self.model(data)
            return output.argmax(dim=1).item()

    def get_model_params(self):
        return torch.nn.utils.parameters_to_vector(self.model.parameters()).detach().cpu().numpy()
    
    def set_model_params(self, params):
        torch.nn.utils.vector_to_parameters(torch.tensor(params).to(self.device), self.model.parameters())
    

    def train_model(self, train_loader, optimizer, epoch):
        self.model.train()
        for batch_idx, (data, target) in enumerate(train_loader):
            data, target = data.to(self.device), target.to(self.device)
            optimizer.zero_grad()
            output = self.model(data)
            loss = F.nll_loss(output, target)
            loss.backward()
            optimizer.step()
            if batch_idx % self.config.log_interval == 0:
                print('Train Epoch: {} [{}/{} ({:.0f}%)]\tLoss: {:.6f}'.format(
                    epoch, batch_idx * len(data), len(train_loader.dataset),
                    100. * batch_idx / len(train_loader), loss.item()))
    
    def test_model(self, test_loader):
        self.model.eval()
        test_loss = 0
        correct = 0
        with torch.no_grad():
            for data, target in test_loader:
                data, target = data.to(self.device), target.to(self.device)
                output = self.model(data)
                test_loss += F.nll_loss(output, target, reduction='sum').item()  # sum up batch loss
                pred = output.argmax(dim=1, keepdim=True)  # get the index of the max log-probability
                correct += pred.eq(target.view_as(pred)).sum().item()
        
        test_loss /= len(test_loader.dataset)

        print('\nTest set: Average loss: {:.4f}, Accuracy: {}/{} ({:.0f}%)\n'.format(
            test_loss, correct, len(test_loader.dataset),
            100. * correct / len(test_loader.dataset)))