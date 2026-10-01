import argparse
import os
from pathlib import Path
import random
import shutil
from typing import Optional, Union

VALID_EXTENSIONS = ('.jpg', '.jpeg', '.png', '.bmp', '.webp')


def tomato_dataset_splitter(
    train_dir: Union[str, Path],
    test_dir: Union[str, Path],
    output_dir: Optional[Union[str, Path]] = None,
    val_ratio: float = 0.2,
    stratified: bool = True,
    seed: int = 42,
    val_folder_name: str = "validation",
) -> Path:
    train_path = Path(train_dir).resolve()
    test_path = Path(test_dir).resolve()

    if output_dir is None:
        out_path = train_path.parent / "dataset_split"
    else:
        out_path = Path(output_dir).resolve()

    if not train_path.exists():
        raise FileNotFoundError(f"Train directory not found: {train_path}")

    rng = random.Random(seed)

    classes = sorted([d.name for d in train_path.iterdir() if d.is_dir() and not d.name.startswith('.')])
    if not classes:
        raise ValueError(f"No class folders found in: {train_path}")

    for cls in classes:
        (out_path / "train" / cls).mkdir(parents=True, exist_ok=True)
        (out_path / val_folder_name / cls).mkdir(parents=True, exist_ok=True)

    split_mode = "stratified" if stratified else "random (non-stratified)"
    print(f"Splitting train data ({len(classes)} classes) using {split_mode} mode (val_ratio={val_ratio})...")

    if stratified:
        for cls in classes:
            cls_dir = train_path / cls
            images = sorted([f for f in cls_dir.iterdir() if f.is_file() and f.suffix.lower() in VALID_EXTENSIONS])
            rng.shuffle(images)

            n_val = max(1, int(len(images) * val_ratio)) if len(images) > 1 else 0
            val_imgs = images[:n_val]
            train_imgs = images[n_val:]

            for img in train_imgs:
                shutil.copy2(img, out_path / "train" / cls / img.name)
            for img in val_imgs:
                shutil.copy2(img, out_path / val_folder_name / cls / img.name)
    else:
        all_images = []
        for cls in classes:
            for f in (train_path / cls).iterdir():
                if f.is_file() and f.suffix.lower() in VALID_EXTENSIONS:
                    all_images.append((cls, f))

        rng.shuffle(all_images)
        n_val = max(1, int(len(all_images) * val_ratio)) if len(all_images) > 1 else 0
        val_imgs = all_images[:n_val]
        train_imgs = all_images[n_val:]

        for cls, img in train_imgs:
            shutil.copy2(img, out_path / "train" / cls / img.name)
        for cls, img in val_imgs:
            shutil.copy2(img, out_path / val_folder_name / cls / img.name)

    if test_path.exists() and test_path.is_dir():
        print("Copying test data...")
        test_classes = sorted([d.name for d in test_path.iterdir() if d.is_dir() and not d.name.startswith('.')])
        for cls in test_classes:
            dest_test_cls = out_path / "test" / cls
            dest_test_cls.mkdir(parents=True, exist_ok=True)
            for img in (test_path / cls).iterdir():
                if img.is_file() and img.suffix.lower() in VALID_EXTENSIONS:
                    shutil.copy2(img, dest_test_cls / img.name)

    print(f"Dataset successfully prepared at: {out_path}")
    return out_path