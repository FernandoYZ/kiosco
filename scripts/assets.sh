#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Build CSS using Tailwind
build_css() {
    if [ ! -f ./tailwindcss ]; then
        echo -e "${RED}❌ tailwindcss no encontrado — run: make setup${NC}"
        exit 1
    fi

    echo -e "${BLUE}🎨 Compilando CSS...${NC}"
    mkdir -p public/dist
    ./tailwindcss -i ./assets/main.css -o ./public/dist/styles.css --minify
    echo -e "${GREEN}✓ CSS compilado${NC}"
}

# Download/copy JavaScript libraries and app code
build_js() {
    echo -e "${BLUE}🎨 Descargando/copiando JavaScript...${NC}"
    mkdir -p public/dist

    # External libraries from CDN
    echo "  Descargando Alpine.js..."
    curl -fsSL https://cdn.jsdelivr.net/npm/alpinejs@3.15.4/dist/cdn.min.js \
        -o public/dist/alpine.min.js

    echo "  Descargando HTMX..."
    curl -fsSL https://cdnjs.cloudflare.com/ajax/libs/htmx/2.0.7/htmx.min.js \
        -o public/dist/htmx.min.js

    echo "  Descargando html-to-image..."
    curl -fsSL https://cdn.jsdelivr.net/npm/html-to-image@1.11.11/dist/html-to-image.min.js \
        -o public/dist/canvas.min.js

    echo "  Descargando Alpine Collapse..."
    curl -fsSL https://cdn.jsdelivr.net/npm/@alpinejs/collapse@3.x.x/dist/cdn.min.js \
        -o public/dist/collapse.min.js

    # App code
    cp assets/main.js public/dist/bundle.min.js

    echo -e "${GREEN}✓ JavaScript compilado${NC}"
}

# Build all assets
build_all() {
    echo -e "${BLUE}🎨 Construyendo assets (CSS/JS)...${NC}"
    build_css
    build_js
    echo -e "${GREEN}✓ Assets construidos${NC}"
}

main() {
    case "${1:-all}" in
        css)
            build_css
            ;;
        js)
            build_js
            ;;
        all)
            build_all
            ;;
        *)
            echo "Usage: $0 {css|js|all}"
            exit 1
            ;;
    esac
}

main "$@"
