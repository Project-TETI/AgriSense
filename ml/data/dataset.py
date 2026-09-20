import os
from dataclasses import dataclass
from typing import Callable, Optional
from PIL import Image
import torch
from torch.utils.data import Dataset, DataLoader
import lightning as pl
import torchvision.transforms as T

class LeafImageDataset(Dataset):
    def __init__(self, root_dir: str, transform: Optional[Callable] = None):
        self.root_dir = root_dir
        self.transform = transform
        self.image_paths = []
        self.labels = []
        self.classes = sorted([
          d for d in os.listdir(root_dir)
          if os.path.isdir(os.path.join(root_dir, d)) and not d.startswith('.')
        ])
        self.class_to_idx = {cls_name: i for i, cls_name in enumerate(self.classes)}
        valid_extensions = ('.jpg', '.jpeg', '.png', '.bmp', '.webp')

        for class_name in self.classes:
            class_dir = os.path.join(root_dir, class_name)
            label = self.class_to_idx[class_name]
            for img_name in sorted(os.listdir(class_dir)):
                if img_name.lower().endswith(valid_extensions):
                    self.image_paths.append(os.path.join(class_dir, img_name))
                    self.labels.append(label)
    
    def __len__(self) -> int:
        return len(self.image_paths)
    
    def __getitem__(self, idx : int):
        img_path = self.image_paths[idx]
        label = self.labels[idx]
        image = Image.open(img_path).convert('RGB')

        if self.transform:
            image = self.transform(image)
        return image, label

