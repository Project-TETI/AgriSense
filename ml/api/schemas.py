from typing import List
from pydantic import BaseModel, Field


class PredictResponse(BaseModel):
    plant: str = Field(..., description="Target plant species evaluated")
    class_name: str = Field(..., description="Detected disease class name or 'unknown'")
    confidence_score: float = Field(..., description="Probability confidence score (0.0 - 1.0)")
    status: str = Field(..., description="Inference status: 'SUCCESS' or 'UNCERTAIN'")


class SupportedPlantsResponse(BaseModel):
    plants: List[str] = Field(..., description="List of supported plant species with available models")


class HealthResponse(BaseModel):
    status: str = Field(..., description="Service health status")
    available_plants: List[str] = Field(..., description="List of supported plant species")
    device: str = Field(..., description="Compute device used for inference ('cpu' or 'cuda')")