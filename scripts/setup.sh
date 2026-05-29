#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

TAILWIND_VERSION="v4.1.10"

# Download tailwindcss binary if not present
download_tailwind() {
    if [ -f ./tailwindcss ]; then
        echo -e "${GREEN}✓ tailwindcss already present${NC}"
        return 0
    fi

    echo -e "${BLUE}⬇ Descargando tailwindcss ${TAILWIND_VERSION}...${NC}"
    if curl -fsSL "https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWIND_VERSION}/tailwindcss-linux-x64" \
        -o ./tailwindcss; then
        chmod +x ./tailwindcss
        echo -e "${GREEN}✓ tailwindcss instalado${NC}"
        return 0
    else
        rm -f ./tailwindcss
        echo -e "${RED}❌ Error descargando tailwindcss${NC}"
        exit 1
    fi
}

# Verify all dependencies
verify_deps() {
    local missing=0

    echo -e "${BLUE}Verificando dependencias...${NC}"

    for tool in go templ curl; do
        if command -v "$tool" >/dev/null 2>&1; then
            echo -e "  ${GREEN}[OK]${NC} $tool"
        else
            echo -e "  ${RED}[MISSING]${NC} $tool"
            missing=1
        fi
    done

    if [ -f ./tailwindcss ]; then
        echo -e "  ${GREEN}[OK]${NC} tailwindcss (local binary)"
    else
        echo -e "  ${RED}[MISSING]${NC} tailwindcss — run: make setup"
        missing=1
    fi

    return $missing
}

main() {
    case "${1:-download}" in
        download)
            download_tailwind
            ;;
        verify)
            verify_deps
            ;;
        all)
            download_tailwind
            verify_deps
            ;;
        *)
            echo "Usage: $0 {download|verify|all}"
            exit 1
            ;;
    esac
}

main "$@"
