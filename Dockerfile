FROM golang:1.20-alpine

WORKDIR /app

RUN apk update && apk upgrade && apk add --no-cache bash

COPY go.* ./
RUN go mod download

ADD . .
RUN go build -o build/parserd cmd/parserd/*.go

EXPOSE 3000
CMD ["./build/parserd"]