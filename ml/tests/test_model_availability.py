import asyncio
import tempfile
import unittest
from pathlib import Path
from unittest.mock import AsyncMock, MagicMock, patch

from fastapi import HTTPException

from ml.api.main import predict_leaf_disease
from ml.api.predictor import ModelUnavailableError, MultiModelService


class ModelAvailabilityTest(unittest.TestCase):
    def test_missing_weights_are_not_loaded_or_advertised(self):
        with tempfile.TemporaryDirectory() as directory:
            plant = Path(directory) / "tomato"
            plant.mkdir()
            (plant / "class_names.json").write_text('["healthy"]', encoding="utf-8")
            service = MultiModelService()
            service.weights_dir = Path(directory)

            self.assertEqual(service.get_available_plants(), [])
            with patch("ml.api.predictor.LeafClassifier") as classifier:
                with self.assertRaises(ModelUnavailableError):
                    service.predict(b"image", "tomato")
                classifier.assert_not_called()
            self.assertNotIn("tomato", service.models)
            self.assertNotIn("tomato", service.class_names)

    def test_checkpoint_is_loaded_before_model_becomes_available(self):
        with tempfile.TemporaryDirectory() as directory:
            plant = Path(directory) / "tomato"
            plant.mkdir()
            (plant / "class_names.json").write_text('["healthy"]', encoding="utf-8")
            (plant / "model.pt").touch()
            service = MultiModelService()
            service.weights_dir = Path(directory)
            model = MagicMock()
            model.to.return_value = model
            with patch("ml.api.predictor.LeafClassifier", return_value=model), patch("ml.api.predictor.torch.load", return_value={"weights": 1}):
                self.assertTrue(service.load_plant_model("tomato"))
            model.load_state_dict.assert_called_once_with({"weights": 1})
            self.assertEqual(service.get_available_plants(), ["tomato"])
            self.assertEqual(service.class_names["tomato"], ["healthy"])

    def test_missing_weights_return_service_unavailable(self):
        with patch("ml.api.main.validate_image_upload", new_callable=AsyncMock, return_value=b"image"), patch(
            "ml.api.main.model_service.predict", side_effect=ModelUnavailableError("Model weights unavailable")
        ):
            with self.assertRaises(HTTPException) as raised:
                asyncio.run(predict_leaf_disease(file=MagicMock(), plant="tomato"))
        self.assertEqual(raised.exception.status_code, 503)


if __name__ == "__main__":
    unittest.main()
