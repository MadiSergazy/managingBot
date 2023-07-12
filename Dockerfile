FROM golang:1.20-alpine3.18 AS builder 

COPY . /madi_telegram_bot/

WORKDIR /madi_telegram_bot/

RUN go mod download

RUN GOOS=linux go build -o ./.bin/bot ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /madi_telegram_bot/bin/bot .

EXPOSE 80

CMD ["./bot"]

# COPY --from=0 it means we want to copy from previous step of building
# EXPOSE 80 - port for getting application in inside