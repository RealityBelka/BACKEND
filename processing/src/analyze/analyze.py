from typing import List, Union
from src.analyze.analyzers import FaceParams
from src.convertors import bytes_to_ndarray

face_analyzer = FaceParams()

async def analyze_photo(msg: bytes) -> dict[str, Union[bool, List[str]]]:
    image = bytes_to_ndarray(msg)

    count = face_analyzer.get_faces_count(image)
    if count != 1:
        return {
            "ok": False,
            "message": f"there is {count} faces on photo"
        }

    return {
        "ok": True,
        "message": "photo is correct"
    }


async def analyze_voice(msg: bytes) -> dict[str, Union[bool, List[str]]]:
    return {
        "ok": True,
        "result": "audio processed successfully"
    }
