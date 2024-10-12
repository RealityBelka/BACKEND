import os
from dotenv import load_dotenv

class Config:
    def __init__(self):
        load_dotenv(dotenv_path=os.path.join(os.getcwd(), '.env'))
        self.nats_host=os.getenv("NATS_HOST")
        self.nats_port=os.getenv("NATS_PORT")
