FROM python:3.12-slim

WORKDIR /app

# Установка системных библиотек, необходимых для работы OpenCV
RUN apt-get update && apt-get install -y \
    gcc \
    ffmpeg \
    libgl1-mesa-glx \
    libglib2.0-0 \
    build-essential \
    python3-dev

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

# Установка переменной окружения для небуферизированного вывода в лог
ENV PYTHONUNBUFFERED=1

ENTRYPOINT [ "python", "./main.py" ]
