FROM golang:1.23.1 as builder

WORKDIR /app

COPY go.mod go.sum config ./
COPY internal/generated ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/config /root/config
RUN mv -f /root/config/docker-config.yaml /root/config/config.yaml
RUN mkdir -p /root/internal/generated
COPY --from=builder /app/internal/generated /root/internal/generated

EXPOSE 8081
EXPOSE 50050

CMD ["./server"]