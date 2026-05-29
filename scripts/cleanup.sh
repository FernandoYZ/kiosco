#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BINARY_PATH="bin/kiosco"

# Clean artifacts (binario, coverage, etc)
clean() {
    echo -e "${YELLOW}🗑 Limpiando artifacts...${NC}"
    rm -f "${BINARY_PATH}"
    rm -f coverage.out coverage.html
    echo -e "${GREEN}✓ Limpieza completada${NC}"
}

# Deep cleanup (+ tailwindcss binary, node_modules, dist)
clean_full() {
    echo -e "${YELLOW}🗑 Limpieza profunda...${NC}"
    rm -f "${BINARY_PATH}"
    rm -f coverage.out coverage.html
    rm -rf node_modules
    rm -rf public/dist/*
    rm -f ./tailwindcss
    echo -e "${GREEN}✓ Limpieza profunda completada${NC}"
}

main() {
    case "${1:-clean}" in
        clean)
            clean
            ;;
        full)
            clean_full
            ;;
        *)
            echo "Usage: $0 {clean|full}"
            exit 1
            ;;
    esac
}

main "$@"
