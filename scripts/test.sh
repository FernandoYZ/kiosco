#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Run tests
run_tests() {
    echo -e "${BLUE}🧪 Ejecutando tests...${NC}"
    go test ./... -v
    echo -e "${GREEN}✓ Tests completados${NC}"
}

# Run tests with coverage
run_coverage() {
    echo -e "${BLUE}🧪 Tests con coverage...${NC}"
    go test ./... -v -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}✓ Coverage generado en coverage.html${NC}"
}

# Run linter
run_lint() {
    echo -e "${BLUE}🔍 Linting...${NC}"
    go vet ./...
    echo -e "${GREEN}✓ Lint completado (sin errores)${NC}"
}

# Format code
run_fmt() {
    echo -e "${BLUE}📐 Formateando código...${NC}"
    go fmt ./...
    echo -e "${GREEN}✓ Código formateado${NC}"
}

main() {
    case "${1:-test}" in
        test)
            run_tests
            ;;
        coverage)
            run_coverage
            ;;
        lint)
            run_lint
            ;;
        fmt)
            run_fmt
            ;;
        *)
            echo "Usage: $0 {test|coverage|lint|fmt}"
            exit 1
            ;;
    esac
}

main "$@"
