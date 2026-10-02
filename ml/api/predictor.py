import io
import json
import os
import sys
from pathlib import Path
from typing import Dict, List, Optional, Tuple

from PIL import Image
import torch
import torchvision.transforms as T

ML_ROOT = Path(__file__).resolve().parent.parent
if str(ML_ROOT) not in sys.path:
    sys.path.append(str(ML_ROOT))

from model.model import LeafClassifier

IMAGENET_MEAN = [0.485, 0.456, 0.406]
IMAGENET_STD = [0.229, 0.224, 0.225]
CONFIDENCE_THRESHOLD = 0.60


class MultiModelService:
    def __init__(self):
        self.device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        self.weights_dir = Path(__file__).resolve().parent.parent / "weights"
        
        self.models: Dict[str, LeafClassifier] = {}
        self.class_names: Dict[str, List[str]] = {}

        self.transform = T.Compose([
            T.Resize((224, 224)),
            T.ToTensor(),
            T.Normalize(mean=IMAGENET_MEAN, std=IMAGENET_STD),
        ])

    def get_available_plants(self) -> List[str]:
        if not self.weights_dir.exists():
            return []
        return sorted([
            d.name.lower() for d in self.weights_dir.iterdir()
            if d.is_dir() and not d.name.startswith('.')
        ])

    def load_plant_model(self, plant: str) -> bool:
        plant_key = plant.lower().strip()
        plant_dir = self.weights_dir / plant_key

        if not plant_dir.is_dir():
            return False

        classes_path = plant_dir / "class_names.json"
        if classes_path.exists():
            with open(classes_path, "r", encoding="utf-8") as f:
                classes = json.load(f)
        else:
            print(f"[WARN] class_names.json not found in {plant_dir}, skipping.")
            return False

        self.class_names[plant_key] = classes
        num_classes = len(classes)

        ckpt_files = list(plant_dir.glob("*.ckpt")) + list(plant_dir.glob("*.pt"))
        if ckpt_files:
            model_file = str(ckpt_files[0])
            if model_file.endswith(".ckpt"):
                model = LeafClassifier.load_from_checkpoint(model_file)
            else:
                model = LeafClassifier(model_name="mobilenet_v3_small", num_classes=num_classes)
                state = torch.load(model_file, map_location=self.device)
                model.load_state_dict(state)
            print(f"[INFO] Loaded model for '{plant_key}' from: {model_file}")
        else:
            print(f"[WARN] No weight file found in {plant_dir}. Using default weights.")
            model = LeafClassifier(model_name="mobilenet_v3_small", num_classes=num_classes)

        model.to(self.device).eval()
        self.models[plant_key] = model
        return True

    def load_all_available_models(self):
        plants = self.get_available_plants()
        print(f"[INFO] Found plant weight folders: {plants}")
        for p in plants:
            self.load_plant_model(p)

    def predict(self, image_bytes: bytes, plant: str) -> Tuple[str, float, str]:
        plant_key = plant.lower().strip()

        if plant_key not in self.models:
            success = self.load_plant_model(plant_key)
            if not success:
                raise ValueError(
                    f"Plant '{plant}' model not found. Available plants: {self.get_available_plants()}"
                )

        model = self.models[plant_key]
        classes = self.class_names[plant_key]

        image = Image.open(io.BytesIO(image_bytes)).convert("RGB")
        tensor = self.transform(image).unsqueeze(0).to(self.device)

        with torch.no_grad():
            logits = model(tensor)
            probs = torch.softmax(logits, dim=1)[0]
            conf, pred_idx = torch.max(probs, dim=0)

        confidence = round(float(conf.item()), 4)

        if confidence < CONFIDENCE_THRESHOLD:
            return "unknown", confidence, "UNCERTAIN"

        predicted_class = classes[pred_idx.item()]
        return predicted_class, confidence, "SUCCESS"


model_service = MultiModelService()