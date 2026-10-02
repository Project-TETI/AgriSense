import cv2
import numpy as np
from fastapi import HTTPException, UploadFile

MAX_FILE_SIZE = 5 * 1024 * 1024
ALLOWED_MIME_TYPES = {"image/jpeg", "image/png", "image/jpg"}
MIN_LAPLACIAN_VAR = 50.0

async def validate_image_upload(file: UploadFile) -> bytes:
    if file.content_type not in ALLOWED_MIME_TYPES:
        raise HTTPException(
            status_code=400,
            detail={
                "error_code": "INVALID_FORMAT",
                "message": "File format unsupported. Only JPEG and PNG images are allowed."
            }
        )

    image_bytes = await file.read()
    if len(image_bytes) > MAX_FILE_SIZE:
        raise HTTPException(
            status_code=400,
            detail={
                "error_code": "FILE_TOO_LARGE",
                "message": f"File size exceeds the maximum limit of {MAX_FILE_SIZE / (1024 * 1024)} MB."
            }
        )

    img_np = np.frombuffer(image_bytes, np.uint8)
    img_cv = cv2.imdecode(img_np, cv2.IMREAD_GRAYSCALE)
    if img_cv is not None:
        laplacian_var = float(cv2.Laplacian(img_cv, cv2.CV_64F).var())
        if laplacian_var < MIN_LAPLACIAN_VAR:
            raise HTTPException(
                status_code=422,
                detail={
                    "error_code": "BLURRY_IMAGE",
                    "message": "Image is too blurry. Please upload a clearer image."
                }
            )

    return image_bytes