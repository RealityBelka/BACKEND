import json
import os.path

from src.convertors import bytes_to_ndarray, binary_to_bin_numbers_pair, aac_to_wav
from src.messaging import NatsClient
from src.analyze import analyze_photo, analyze_voice

nats_client = NatsClient()

async def process_photo(msg):
    data = msg.data

    # Десериализуем сообщение
    photo, numbers = binary_to_bin_numbers_pair(data)

    nd_array = bytes_to_ndarray(photo)

    result = await analyze_photo(nd_array, numbers)

    pr = json.dumps(result).encode()

    await nats_client.get_nc().publish(msg.reply, pr)
    await nats_client.get_nc().flush()


async def process_audio(msg):
    data = msg.data

    audio_aac, numbers = binary_to_bin_numbers_pair(data)

    wav_audio_path = aac_to_wav(audio_aac)
    if wav_audio_path:
        try:
            result = await analyze_voice(wav_audio_path, numbers)  # place numbers here

            pr = json.dumps(result).encode()

            await nats_client.get_nc().publish(msg.reply, pr)
            await nats_client.get_nc().flush()
        finally:
            if os.path.exists(wav_audio_path):
                os.remove(wav_audio_path)


