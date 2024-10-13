from nats.aio.client import Client as Nats

from src.config import Config

config = Config()

class NatsClient:
    def __init__(self):
        if not hasattr(self, 'nc'):
            self.nc = Nats()

    async def connect(
            self,
            url=f"nats://{config.nats_host}:{config.nats_port}",
            max_pending_size=1024*1024,
    ):
        print(f"connecting to nats://{config.nats_host}:{config.nats_port}...")
        await self.nc.connect(url)

        self.nc._max_pending_size = max_pending_size

    def get_nc(self):
        return self.nc

    async def close_connection(self):
        await self.nc.drain()
