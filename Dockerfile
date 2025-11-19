ARG GO_VERSION=1.24.6
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app ./cmd/kiosco


FROM debian:bookworm

WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /run-app /usr/local/bin/

# Copiar views (templates) y archivos estáticos
COPY --from=builder /usr/src/app/internal/views ./internal/views
COPY --from=builder /usr/src/app/static ./static

CMD ["run-app"]
