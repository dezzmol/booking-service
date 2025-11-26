FROM golang:1.23.1 as builder

WORKDIR /app

COPY go.mod go.sum ./internal/config ./
COPY internal/generated ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o server ./cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/internal/config /root/internal/config
RUN mv -f /root/internal/config/docker-config.yaml /root/internal/config/config.yaml
RUN mkdir -p /root/internal/generated
COPY --from=builder /app/internal/generated /root/internal/generated

EXPOSE 8081
EXPOSE 50050

CMD ["./server"]