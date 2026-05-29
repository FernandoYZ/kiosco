#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Variables
BINARY_PATH="bin/kiosco"
CMD_PATH="./cmd/kiosco"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
GO_VERSION=$(go version | awk '{print $3}')

# Show project info
show_info() {
    echo -e "${BLUE}Información del Proyecto${NC}"
    echo -e "${BLUE}=======================${NC}"
    echo "Nombre: Kiosco"
    echo "Versión: ${VERSION}"
    echo "Compilado: ${BUILD_TIME}"
    echo "Go: ${GO_VERSION}"
    echo "Binario: ${BINARY_PATH}"
    echo ""
    echo "Archivos relevantes:"
    echo "  - Código: ${CMD_PATH}/main.go"
    echo "  - Config: internal/config/"
    echo "  - Modelos: internal/models/"
    echo "  - Repositorios: internal/repositories/"
    echo "  - Templates: templates/"
    echo "  - Assets: assets/"
}

# Show system status
show_status() {
    echo -e "${BLUE}Estado del Sistema${NC}"
    echo -e "${BLUE}=================${NC}"
    echo -n "Go: "
    go version
    echo -n "Templ: "
    templ version 2>/dev/null || echo "no instalado"
    echo -n "TailwindCSS: "
    ./tailwindcss --version 2>/dev/null || echo "no instalado (run: make setup)"
    echo -n "SQLite: "
    sqlite3 --version 2>/dev/null || echo "no instalado"
    echo -n "Binario compilado: "
    if [ -f "${BINARY_PATH}" ]; then
        echo -e "${GREEN}sí${NC}"
    else
        echo -e "${RED}no${NC}"
    fi
}

main() {
    case "${1:-info}" in
        info)
            show_info
            ;;
        status)
            show_status
            ;;
        *)
            echo "Usage: $0 {info|status}"
            exit 1
            ;;
    esac
}

main "$@"
