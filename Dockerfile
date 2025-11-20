ARG GO_VERSION=1.24.6
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o run-app ./cmd/kiosco


FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

RUN useradd -m -s /bin/bash appuser

WORKDIR /app

COPY --from=builder /usr/src/app/run-app /usr/local/bin/run-app
COPY --from=builder /usr/src/app/internal/views ./internal/views
COPY --from=builder /usr/src/app/static ./static

RUN chown -R appuser:appuser /app

ENV TZ=America/Lima

USER appuser

EXPOSE 3200

ENTRYPOINT ["run-app"]