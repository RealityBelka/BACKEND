import json

from src.messaging import NatsClient
from src.analyze import analyze_photo, analyze_voice

nats_client = NatsClient()

async def process_photo(msg):
    data = msg.data

    result = await analyze_photo(data)

    pr = json.dumps(result).encode()

    await nats_client.get_nc().publish(msg.reply, pr)
    await nats_client.get_nc().flush()


async def process_audio(msg):
    data = msg.data

    result = await analyze_voice(data)

    pr = json.dumps(result).encode()

    await nats_client.get_nc().publish(msg.reply, pr)
    await nats_client.get_nc().flush()
