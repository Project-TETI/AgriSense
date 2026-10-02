from contextlib import asynccontextmanager
from fastapi import FastAPI, File, Form, HTTPException, UploadFile
from fastapi.middleware.cors import CORSMiddleware

from .schemas import HealthResponse, PredictResponse, SupportedPlantsResponse
from .validator import validate_image_upload
from .predictor import model_service


@asynccontextmanager
async def lifespan(app: FastAPI):
    model_service.load_all_available_models()
    yield


app = FastAPI(
    title="AgriSense ML Inference Service",
    description="Multi-plant leaf disease detection service",
    version="1.0.0",
    lifespan=lifespan
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/api/v1/ml/health", response_model=HealthResponse, tags=["Health"])
async def health_check():
    return HealthResponse(
        status="healthy",
        available_plants=model_service.get_available_plants(),
        device=str(model_service.device)
    )


@app.get("/api/v1/ml/plants", response_model=SupportedPlantsResponse, tags=["Plants"])
async def get_supported_plants():
    return SupportedPlantsResponse(plants=model_service.get_available_plants())


@app.post("/api/v1/ml/predict", response_model=PredictResponse, tags=["Inference"])
async def predict_leaf_disease(
    file: UploadFile = File(..., description="Leaf image (JPEG/PNG)"),
    plant: str = Form("tomato", description="Plant type: 'tomato', 'rice', 'corn'")
):
    image_bytes = await validate_image_upload(file)
    try:
        class_name, confidence, status = model_service.predict(image_bytes, plant=plant)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))

    return PredictResponse(
        plant=plant.lower(),
        class_name=class_name,
        confidence_score=confidence,
        status=status
    )