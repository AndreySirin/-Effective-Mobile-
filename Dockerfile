FROM golang:1.24.1 AS builder

WORKDIR /app

COPY . ./

RUN go mod tidy

WORKDIR /app/cmd

RUN go build -o /app/binary


FROM ubuntu:latest

WORKDIR /root

COPY --from=builder /app/binary .

COPY config.yaml /root/

COPY internal/storage/migrate /root/migrate

EXPOSE 8080

CMD ["/root/binary"]