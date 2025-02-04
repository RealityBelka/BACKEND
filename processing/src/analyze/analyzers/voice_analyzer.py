import io
import re
import librosa
import webrtcvad
import numpy as np
from langdetect import detect
import speech_recognition as sr


class AudioParams:
    def __init__(self):
        self.number_words = {
            "ноль": '0', "один": '1', "раз": '1', "два": '2', "три": '3', "четыре": '4',
            "пять": '5', "шесть": '6', "семь": '7', "восемь": '8', "девять": '9',
            "десять": '10'
        }

    async def analyze_audio_quality(self, audio_path: str, threshold_quite=0.1, threshold_overload=0.99):
        """
        Определение уровня шума и наличия перегрузки.

        :param audio_path: Путь к аудиофайлу
        :param threshold_quite: Порог тихого звучания (чем больше значение, тем тише порог)
        :param threshold_overload: Порог перегрузки (чем больше значение, тем громче порог)
        :return: Уровень шума (тихо, громко, норм)
        """
        y, _ = librosa.load(audio_path, sr=16000)

        # Расчет отношения сигнал/шум
        signal_energy = np.mean(y ** 2)
        noise_energy = np.var(y)

        snr = 1000 * np.log10(signal_energy / noise_energy)
        overload = np.max(np.abs(y)) >= threshold_overload

        if snr > threshold_quite:
            return "тихо"
        if overload:
            return "громко"

        return "норм"

    async def replace_non_alphanumeric(self, text):
        """
        Заменяет все символы в строке, которые не являются буквами (латинскими или русскими) или цифрами, на пробелы.

        :param text: исходный текст
        :return: текст с замененными символами
        """
        return re.sub(r"[^a-zA-Zа-яА-ЯёЁ0-9]", ' ', text)

    async def text_recognition(self, audio_path: str):
        """Производит распознавание текста из аудио."""
        recognizer = sr.Recognizer()
        with sr.AudioFile(audio_path) as source:
            audio = recognizer.record(source)
        try:
            text = recognizer.recognize_google(audio, language="ru_RU")
            # wh = WhisperHuggingface(audio)  # Вторая версия, на выбор (но надо отдебажить)
            # text = wh.process()
            clean_text = await self.replace_non_alphanumeric(text)
            return clean_text
        except:
            return None

    async def check_language(self, transcription):
        """Проверяет язык записи"""
        try:
            lang = detect(transcription)
            return lang
        except:
            lang = "ru" if transcription else None
        return lang

    async def only_numerals(self, input_string):
        """Проверяет, содержит ли строка только числа (словесные или символы)"""

        words = input_string.lower().split()

        for word in words:
            if word.isdigit() or word in self.number_words:
                continue
            else:
                return False
        return True

    async def get_vad_segments(self, audio_path: str, sample_rate=16000, frame_duration=30):
        """
        Возвращает список активных (с голосом) и неактивных (тихих) сегментов в аудио.

        :param audio_path: Путь к аудиофайлу
        :param sample_rate: Частота дискретизации (должна быть 8000, 16000, 32000 или 48000)
        :param frame_duration: Длительность кадра в миллисекундах (10, 20 или 30 мс)
        :return: Сегменты активных и неактивных участков
        """
        vad = webrtcvad.Vad(2)  # Агрессивность VAD (0-3)

        # Загружаем аудиофайл и приводим его к нужной частоте дискретизации
        y, _ = librosa.load(audio_path, sr=sample_rate)

        # Преобразуем аудио в 16-битный формат PCM (webrtcvad требует 16-битный формат)
        y = (y * 32768).astype(np.int16)

        # Рассчитываем длину кадра в сэмплах (например, 30 мс -> 30/1000 * sample_rate)
        frame_length = int(sample_rate * frame_duration / 1000)

        # Разбиваем аудио на кадры по длительности
        frames = [y[i:i + frame_length] for i in range(0, len(y), frame_length)]

        segments = {
            "active": [],
            "non_active": []
        }

        for frame in frames:
            if len(frame) != frame_length:
                # Если длина кадра не соответствует ожидаемой длине, пропускаем его
                continue

            # Проверяем, является ли кадр активным (есть ли в нем речь)
            is_speech = vad.is_speech(frame.tobytes(), sample_rate)

            if is_speech:
                segments["active"].append(frame)
            else:
                segments["non_active"].append(frame)

        return segments

    async def analyze_silence(self, segments):
        """
        Анализирует неактивные (тихие) участки аудиозаписи.

        :param segments: Сегменты аудио (результат функции get_vad_segments)
        :return: Дисперсия для неактивных сегментов
        """
        non_active_segments = segments["non_active"]

        # Вычисляем дисперсию для каждого неактивного сегмента
        variances = [np.var(segment) for segment in non_active_segments]

        # Возвращаем среднюю дисперсию по всем тихим сегментам
        avg_variance = np.mean(variances) if variances else 0

        return avg_variance

    async def check_single_speaker(self, audio_path: str, silence_threshold=5000):
        """
        Проверяет уровень шума (дисперсию) на тихих участках аудио.

        :param audio_path: Путь к аудиофайлу
        :param silence_threshold: Порог для определения шума в тихих сегментах
        :return: True, если запись проходит проверку на шум, иначе False
        """
        try:
            segments = await self.get_vad_segments(audio_path)

            avg_variance = await self.analyze_silence(segments)

            # print(f"Средняя дисперсия на тихих участках: {avg_variance}")

            if avg_variance > silence_threshold:
                return False
            else:
                return True
        except Exception:
            return False


    def replace_word_numerals(self, text):
        new_text = []
        for word in text.split():
            if word.isdigit():
                new_text.append(word)
            else:
                new_text.append(self.number_words[word])
        return new_text

