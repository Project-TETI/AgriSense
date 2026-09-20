import os
from dataclasses import dataclass
from typing import Callable, Optional
from PIL import Image
import torch
from torch.utils.data import Dataset, DataLoader
import lightning as pl
import torchvision.transforms as T

@dataclass
class DataConfig:
    batch_size: int = 32
    num_workers: int = 4
    pin_memory: bool = True
    shuffle: bool = True
    train_transform: Optional[Callable] = None
    eval_transform: Optional[Callable] = None

class LeafDataModule(pl.LightningDataModule):
    def __init__(self, data_dir: str, config: DataConfig):
        super().__init__()
        self.save_hyperparameters(ignore=['config'])
        self.data_dir = data_dir
        self.config = config

        self.train_dataset: Optional[Dataset] = None
        self.val_dataset: Optional[Dataset] = None
        self.test_dataset: Optional[Dataset] = None
    
    def setup(self, stage: Optional[str] = None):
        if stage == 'fit' or stage is None:
            self.train_dataset = LeafImageDataset(os.path.join(self.data_dir, 'train'), transform=self.config.train_transform)
            self.val_dataset = LeafImageDataset(os.path.join(self.data_dir, 'validation'), transform=self.config.eval_transform)
        
        if stage == 'test' or stage is None:
            self.test_dataset = LeafImageDataset(os.path.join(self.data_dir, 'test'), transform=self.config.eval_transform)
    
    def train_dataloader(self) -> DataLoader:
        return DataLoader(
            self.train_dataset,
            batch_size=self.config.batch_size,
            shuffle=self.config.shuffle,
            num_workers=self.config.num_workers,
            pin_memory=self.config.pin_memory,
            persistent_workers=self.config.num_workers > 0
        )
    
    def val_dataloader(self) -> DataLoader:
        return DataLoader(
            self.val_dataset,
            batch_size=self.config.batch_size,
            shuffle=False,
            num_workers=self.config.num_workers,
            pin_memory=self.config.pin_memory,
            persistent_workers=self.config.num_workers > 0
        )

    def test_dataloader(self) -> DataLoader:
        return DataLoader(
            self.test_dataset,
            batch_size=self.config.batch_size,
            num_workers=self.config.num_workers,
            shuffle=False,
            pin_memory=self.config.pin_memory,
            persistent_workers=self.config.num_workers > 0
        )
    