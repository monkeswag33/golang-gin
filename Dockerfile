FROM golang:1.18 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o bin/webserver .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /usr/src/app/bin/webserver ./
CMD ["./webserver"]