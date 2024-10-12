FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/

RUN go build -o main main.go

WORKDIR /app

ENTRYPOINT [ "./cmd/main" ]
