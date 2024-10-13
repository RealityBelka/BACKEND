import io
import time

import cv2

from typing import List, Union
from src.analyze.analyzers import FaceParams, AudioParams

face_analyzer = FaceParams()
voice_analyzer = AudioParams()

async def analyze_photo(photo, numbers: List[int]) -> dict[str, Union[bool, List[str]]]:
    flag = False  # True, если проверка пройдена успешно

    img_rgb = cv2.cvtColor(photo, cv2.COLOR_BGR2RGB)

    '''1'''
    faces_count = await face_analyzer.get_faces_count(img_rgb)
    if faces_count > 1:
        return {"ok": flag, "message": "Больше одного лица на изображении"}
    if faces_count < 1:
        return {"ok": flag, "message": "Лицо не обнаружено"}

    '''2'''
    # is_in_frame = face_analyzer.is_face_in_frame(photo, numbers)
    # if not is_in_frame:
        # return {"ok": flag, "message": "Лицо должно полностью помещаться в рамку"}


    '''3'''
    eyes_distance = await face_analyzer.get_eye_distance(img_rgb)
    if eyes_distance < 0.1:
        return {"ok": flag, "message": "Приблизьте телефон к лицу"}

    '''3.5'''
    mono_background = await face_analyzer.calculate_background_uniformity(photo)
    if mono_background > 18:
        return {"ok": flag, "message": "Слишком пёстрый задний фон"}

    '''4'''
    head_pose = await face_analyzer.get_head_pose(img_rgb)
    if not ((140 <= head_pose["yaw"] <= 220) and (80 <= head_pose["pitch"] <= 180) and (-40 <= head_pose["roll"] <= 40)):
        return {"ok": flag, "message": "Держите голову прямо"}

    '''5'''
    is_obstructed = face_analyzer.check_face_obstruction(img_rgb)
    if is_obstructed:
        return {"ok": flag, "message": "Лицо должно быть полностью открыто"}

    '''6'''
    brightness, CV = await face_analyzer.calculate_face_illumination(img_rgb)
    if brightness < 50 or 150 < brightness:
        return {"ok": flag, "message": "Обеспечьте равномерное освещение лица (brightness)"}
    if CV > 15:
        return {"ok": flag, "message": "Обеспечьте равномерное освещение лица (CV)"}

    '''7'''
    _, is_blurred = await face_analyzer.calculate_blurriness(photo)
    if is_blurred:
        return {"ok": flag, "message": "Отодвиньте телефон от лица для фокусировки"}

    '''8'''
    is_neutral = await face_analyzer.check_neutral_status(img_rgb)
    if not is_neutral:
        return {"ok": flag, "message": "Выражение лица должно быть нейтральным"}


    '''10'''
    is_real = await face_analyzer.check_spoofing(img_rgb)
    if not is_real:
        return {"ok": flag, "message": "Кажется, в кадре не реальный человек"}

    flag = True

    return {"ok": flag, "message": None}



async def analyze_voice(audio_path: str, numbers: List[int]) -> dict[str, Union[bool, List[str]]]:
    if not audio_path:
        print("empty input value. cannot proceed")
        return {
            "ok": False,
            "message": "empty input"
        }

    print("done file collect")
    time.sleep(2)

    flag = False  # True, если проверка пройдена успешно

    noise_level = await voice_analyzer.analyze_audio_quality(audio_path)
    if noise_level == "тихо":
        return {"ok": flag, "message": "Говорите громче или переместитесь в более тихое место"}
    if noise_level == "громко":
        return {"ok": flag, "message": "Говорите тише или отодвиньте телефон от лица"}

    text = await voice_analyzer.text_recognition(audio_path)

    language = await voice_analyzer.check_language(text)
    if language != "ru":
        return {"ok": flag, "message": "Произносите указанные цифры на русском языке"}

    only_numerals = voice_analyzer.only_numerals(text)
    # print("4) only_numerals: ", only_numerals)
    if only_numerals:
        recognized_text = voice_analyzer.replace_word_numerals(text)
    else:
        return {
            "ok": flag,
            "message": "Произносите только указанные на экране цифры"
        }

    check_single_speaker = await voice_analyzer.check_single_speaker(audio_path)
    if not check_single_speaker:
        return {
            "ok": flag,
            "message": "Посторонние шумы. Переместитесь в более тихое место",
        }

    recognized_numerals = [int(x) for x in recognized_text]
    if numbers != recognized_numerals:
        return {
            "ok": flag,
            "message": "Произносите только указанные на экране цифры"
        }

    flag = True

    return {"ok": flag, "message": None}
