import cv2
import numpy as np

def bytes_to_ndarray(image_bytes: bytes):
    # Преобразование байтов в numpy массив (массив байт)
    nparr = np.frombuffer(image_bytes, np.uint8)

    # Декодирование изображения из numpy массива в формате OpenCV
    image = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

    return image

