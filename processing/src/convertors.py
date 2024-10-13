import tempfile
import cv2
import numpy as np
import io
import struct

from msgpack.fallback import BytesIO
from pydub import AudioSegment
from typing import List, BinaryIO

int32_size = 4

def binary_to_bin_numbers_pair(data: bytes) -> tuple[bytes, List[int]]:
    try:
        numbers_len = struct.unpack('>I', data[:int32_size])[0]  # '>I' означает 4 байта в формате BigEndian
        offset = int32_size

        numbers = []
        for _ in range(numbers_len):
            num = struct.unpack('>I', data[offset:offset + int32_size])[0]
            numbers.append(num)
            offset += int32_size

        image_len = struct.unpack('>I', data[offset:offset + int32_size])[0]
        offset += int32_size

        image_data = data[offset:offset + image_len]

        return image_data, numbers

    except Exception as e:
        print(f"Error parsing binary data: {str(e)}")
        return b"", []



def bytes_to_ndarray(image_bytes: bytes):
    # Преобразование байтов в numpy массив (массив байт)
    nparr = np.frombuffer(image_bytes, np.uint8)

    # Декодирование изображения из numpy массива в формате OpenCV
    image = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

    return image


def aac_to_wav(aac_data: bytes) -> str:
    try:
        audio = AudioSegment.from_file(io.BytesIO(aac_data), format="aac")

        # converting to wav
        temp_wav_file = tempfile.NamedTemporaryFile(delete=False, suffix=".wav")
        audio.export(temp_wav_file.name, format="wav")
        temp_wav_file.close()

        return temp_wav_file.name

    except Exception as e:
        print(f"Error converting AAC to WAV: {str(e)}")
        return ""
