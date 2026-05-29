#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Variables
BINARY_NAME="kiosco"
BINARY_PATH="bin/${BINARY_NAME}"
CMD_PATH="./cmd/kiosco"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
LDFLAGS="-X main.version=${VERSION} -X 'main.buildTime=${BUILD_TIME}'"

# Quick build without validations
build_quick() {
    echo -e "${BLUE}⚡ Compilación rápida...${NC}"
    go build -o "${BINARY_PATH}" "${CMD_PATH}"
    echo -e "${GREEN}✓ Binario compilado: ${BINARY_PATH}${NC}"
}

# Production build with validations and tests
build_prod() {
    echo -e "${BLUE}🏭 Build de PRODUCCIÓN con tests...${NC}"
    go test ./... -v
    go build -ldflags="${LDFLAGS}" -o "${BINARY_PATH}" "${CMD_PATH}"
    echo -e "${GREEN}✓ Build de producción completado${NC}"
}

# Standard production build (with lint, no tests)
build_standard() {
    echo -e "${BLUE}🔨 Compilando kiosco ${VERSION}...${NC}"
    go build -ldflags="${LDFLAGS}" -o "${BINARY_PATH}" "${CMD_PATH}"
    echo -e "${GREEN}✓ Binario compilado: ${BINARY_PATH}${NC}"
    ls -lh "${BINARY_PATH}"
}

# Linux/amd64 build for VPS
build_linux() {
    echo -e "${BLUE}🐧 Compilando para Linux (amd64)...${NC}"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${BINARY_PATH}" "${CMD_PATH}"
    echo -e "${GREEN}✓ Binario Linux compilado: ${BINARY_PATH}${NC}"
    ls -lh "${BINARY_PATH}"
    echo -e "${YELLOW}▶ Próximo paso: scripts/deploy.sh${NC}"
}

main() {
    case "${1:-standard}" in
        quick)
            build_quick
            ;;
        prod)
            build_prod
            ;;
        linux)
            build_linux
            ;;
        standard)
            build_standard
            ;;
        *)
            echo "Usage: $0 {quick|standard|prod|linux}"
            exit 1
            ;;
    esac
}

main "$@"
