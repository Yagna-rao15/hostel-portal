import cv2
import numpy as np
import pytesseract


def preprocess_image(image_path):
    # Read the image
    image = cv2.imread(image_path)

    # Convert to grayscale
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

    # Resize image to set a standard DPI (300 DPI equivalent size)
    height, width = gray.shape[:2]
    target_height = (
        3000  # Adjust based on actual image size, set to roughly match 300 DPI
    )
    scale = target_height / height
    resized = cv2.resize(
        gray, (int(width * scale), target_height), interpolation=cv2.INTER_AREA
    )

    # Apply Gaussian blur to reduce noise and improve edge detection
    blurred = cv2.GaussianBlur(resized, (5, 5), 0)

    # Use adaptive thresholding to create a binary image
    thresh = cv2.adaptiveThreshold(
        blurred, 255, cv2.ADAPTIVE_THRESH_GAUSSIAN_C, cv2.THRESH_BINARY_INV, 11, 2
    )

    # Find contours to detect regions of interest
    contours, _ = cv2.findContours(thresh, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

    processed_images = []

    for contour in contours:
        # Get bounding box coordinates of each contour
        x, y, w, h = cv2.boundingRect(contour)

        # Define conditions to ignore very small contours
        if w > 30 and h > 30:  # Adjust size thresholds as needed
            # Crop the detected region
            roi = resized[y : y + h, x : x + w]

            # Rotate the cropped image to ensure upright text
            rotated = rotate_if_needed(roi)

            processed_images.append(rotated)

    return processed_images


def rotate_if_needed(image):
    try:
        # Use pytesseract to detect orientation and skew angle
        data = pytesseract.image_to_osd(image, output_type=pytesseract.Output.DICT)
        angle = data["rotate"]

        if angle != 0:
            # Rotate the image to correct the skew
            (h, w) = image.shape[:2]
            center = (w // 2, h // 2)
            M = cv2.getRotationMatrix2D(center, angle, 1.0)
            rotated = cv2.warpAffine(
                image, M, (w, h), flags=cv2.INTER_CUBIC, borderMode=cv2.BORDER_REPLICATE
            )
            return rotated
    except pytesseract.TesseractError:
        print("Skipping rotation due to insufficient text or resolution issues.")
        # If an error occurs, return the image as-is
        return image

    return image


def run_ocr_on_chunks(processed_images):
    # Initialize list to collect all OCR results
    results = []

    for img in processed_images:
        # Run Tesseract OCR on each processed image chunk
        text = pytesseract.image_to_string(
            img, config="--psm 6"
        )  # psm 6 assumes a uniform block of text
        results.append(text)

    return results


# Example usage:
processed_images = preprocess_image("./3.jpg")
ocr_results = run_ocr_on_chunks(processed_images)

for idx, result in enumerate(ocr_results):
    print(f"Chunk {idx+1}: {result}")
