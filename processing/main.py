import asyncio
import signal

from src.handlers import nats_client
from src.handlers import process_audio, process_photo

async def main():
    try:
        await nats_client.connect()

        loop = asyncio.get_event_loop()
        stop_event = asyncio.Event()

        def signal_handler():
            stop_event.set()

        loop.add_signal_handler(signal.SIGINT, signal_handler)
        loop.add_signal_handler(signal.SIGTERM, signal_handler)

        await nats_client.get_nc().subscribe("photo_subject", cb=process_photo)
        await nats_client.get_nc().subscribe("audio_subject", cb=process_audio)

        await stop_event.wait()

    finally:
        print("закрытие соединения...")
        await nats_client.close_connection()
        print("соединение закрыто")


if __name__ == "__main__":
    asyncio.run(main())
