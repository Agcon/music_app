FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/music_app

FROM frolvlad/alpine-glibc:alpine-3.18_glibc-2.35

WORKDIR /app

COPY --from=builder /app/server /app/
COPY --from=builder /app/web /app/web
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/.env /app/.env

EXPOSE 8086

CMD ["./server"]
