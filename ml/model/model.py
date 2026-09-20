import torch
import torch.nn as nn
import lightning as pl
from torchvision import models
from torchmetrics.classification import MulticlassAccuracy, MulticlassF1Score

class LeafClassifier(pl.LightningModule):
    def __init__(self, model_name: str, num_classes: int, learning_rate: float = 1e-3, weight_decay: float = 1e-3, freeze_backbone: bool = True):
        super().__init__()
        self.save_hyperparameters()
        self.learning_rate = learning_rate
        self.weight_decay = weight_decay
        self.model_name = model_name

        if model_name == 'mobilenet_v3_small':
            base = models.mobilenet_v3_small(weights=models.MobileNet_V3_Small_Weights.DEFAULT)
            self.backbone = nn.Sequential(base.features, base.avgpool)
            in_features = base.classifier[3].in_features
            self.classifier = base.classifier
            self.classifier[3] = nn.Linear(in_features, num_classes)

        elif model_name == 'efficientnet_b0':
            base = models.efficientnet_b0(weights=models.EfficientNet_B0_Weights.DEFAULT)
            self.backbone = nn.Sequential(base.features, base.avgpool)
            in_features = base.classifier[1].in_features
            self.classifier = base.classifier
            self.classifier[1] = nn.Linear(in_features, num_classes)
        else:
            raise ValueError(f"Model {model_name} is not supported!")        

        self.criterion = nn.CrossEntropyLoss()
        self.train_acc = MulticlassAccuracy(num_classes=num_classes)
        self.val_acc = MulticlassAccuracy(num_classes=num_classes)
        self.test_acc = MulticlassAccuracy(num_classes=num_classes)
        self.train_f1 = MulticlassF1Score(num_classes=num_classes, average='macro')
        self.val_f1 = MulticlassF1Score(num_classes=num_classes, average='macro')
        self.test_f1 = MulticlassF1Score(num_classes=num_classes, average='macro')
    
    def forward(self, x: torch.Tensor) -> torch.Tensor:
        x = self.backbone(x)
        x = torch.flatten(x, 1)
        x = self.classifier(x)
        return x

    def training_step(self, batch, batch_idx):
        x, y = batch
        logits = self(x)
        loss = self.criterion(logits, y)
        self.train_acc(logits, y)
        self.train_f1(logits, y)
        self.log("train_loss", loss, on_epoch=True, prog_bar=True)
        self.log("train_acc", self.train_acc, on_epoch=True, prog_bar=True)
        self.log("train_f1", self.train_f1, on_epoch=True, prog_bar=True)
        return loss

    def validation_step(self, batch, batch_idx):
        x, y = batch
        logits = self(x)
        loss = self.criterion(logits, y)
        self.val_acc(logits, y)
        self.val_f1(logits, y)
        self.log("val_loss", loss, on_epoch=True, prog_bar=True)
        self.log("val_acc", self.val_acc, on_epoch=True, prog_bar=True)
        self.log("val_f1", self.val_f1, on_epoch=True, prog_bar=True)

    def test_step(self, batch, batch_idx):
        x, y = batch
        logits = self(x)
        loss = self.criterion(logits, y)
        self.test_acc(logits, y)
        self.test_f1(logits, y)
        self.log("test_loss", loss, on_epoch=True, prog_bar=True)
        self.log("test_acc", self.test_acc, on_epoch=True, prog_bar=True)
        self.log("test_f1", self.test_f1, on_epoch=True, prog_bar=True)

    def on_train_epoch_start(self):
        for m in self.backbone.modules():
            if isinstance(m, (nn.BatchNorm2d, nn.SyncBatchNorm)):
                m.eval()

    def configure_optimizers(self):
        backbone_decay, backbone_no_decay = [], []
        head_decay, head_no_decay = [], []

        for name, param in self.named_parameters():
            is_no_decay = any(nd in name for nd in ['bias', 'bn', 'norm'])

            if 'backbone' in name:
                if is_no_decay:
                    backbone_no_decay.append(param)
                else:
                    backbone_decay.append(param)
            else:
                if is_no_decay:
                    head_no_decay.append(param)
                else:
                    head_decay.append(param)
        
        optimizer_grouped_parameters = [
            {"params": backbone_decay, "weight_decay": self.weight_decay, "lr": self.learning_rate, "is_backbone": True},
            {"params": backbone_no_decay, "weight_decay": 0.0, "lr": self.learning_rate, "is_backbone": True},
            {"params": head_decay, "weight_decay": self.weight_decay, "lr": self.learning_rate, "is_backbone": False},
            {"params": head_no_decay, "weight_decay": 0.0, "lr": self.learning_rate, "is_backbone": False},
        ]

        optimizer_grouped_parameters = [g for g in optimizer_grouped_parameters if len(g["params"]) > 0]
        optimizer = torch.optim.AdamW(optimizer_grouped_parameters)
        scheduler = torch.optim.lr_scheduler.CosineAnnealingLR(
            optimizer,
            T_max=self.trainer.max_epochs,
            eta_min=1e-6
        )

        return {
            "optimizer": optimizer,
            "lr_scheduler": {
                "scheduler": scheduler,
                "interval": "epoch",
                "frequency": 1
            }
        }