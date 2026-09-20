import torch
import torch.nn as nn
import lightning as pl
from lightning.pytorch.callbacks import BaseFinetuning

class FreezeBackbone(BaseFinetuning):
    def __init__(self, unfreeze_at_epoch: int = 5, fine_tune_lr: float = 1e-4):
        super().__init__()
        self.unfreeze_at_epoch = unfreeze_at_epoch
        self.fine_tune_lr = fine_tune_lr
    
    def freeze_before_training(self, pl_module: pl.LightningModule):
        self.freeze(pl_module.backbone)
    
    def finetune_function(self, pl_module: pl.LightningModule, current_epoch: int, optimizer: torch.optim.Optimizer):
        if current_epoch == self.unfreeze_at_epoch:
            self.unfreeze(pl_module.backbone)
            
            for m in pl_module.backbone.modules():
                if isinstance(m, (nn.BatchNorm2d, nn.SyncBatchNorm)):
                    m.eval()

            for param_group in optimizer.param_groups:
                if param_group.get('is_backbone', False):
                    param_group['lr'] = self.fine_tune_lr