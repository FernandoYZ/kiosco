.PHONY: help dev dev-static build build-quick build-prod build-linux run run-prod clean clean-full templ assets assets-css assets-js lint fmt test test-coverage db-wal db-verify info status setup verify
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=kiosco
BINARY_PATH=bin/$(BINARY_NAME)
CMD_PATH=./cmd/kiosco
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S')
GO_VERSION=$(shell go version | awk '{print $$3}')
LDFLAGS=-X main.version=$(VERSION) -X 'main.buildTime=$(BUILD_TIME)'
TAILWIND_VERSION=v4.1.10

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

help: ## Mostrar este mensaje de ayuda
	@echo "$(BLUE)Kiosco - Sistema de Control de Consumo Escolar$(NC)"
	@echo "$(BLUE)===========================================$(NC)"
	@echo ""
	@echo "$(YELLOW)Comandos disponibles:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Ejemplos:$(NC)"
	@echo "  $(GREEN)make dev$(NC)          # Iniciar en modo desarrollo con hot reload"
	@echo "  $(GREEN)make build$(NC)        # Compilar para producción"
	@echo "  $(GREEN)make run$(NC)          # Ejecutar el binario compilado"
	@echo "  $(GREEN)make clean$(NC)        # Limpiar artifacts"

# ==================== SETUP ====================

setup: ## Descargar binario tailwindcss (idempotente)
	@if [ -f ./tailwindcss ]; then \
		echo "$(GREEN)✓ tailwindcss already present$(NC)"; \
	else \
		echo "$(BLUE)⬇ Descargando tailwindcss $(TAILWIND_VERSION)...$(NC)"; \
		curl -fsSL https://github.com/tailwindlabs/tailwindcss/releases/download/$(TAILWIND_VERSION)/tailwindcss-linux-x64 \
			-o ./tailwindcss || { rm -f ./tailwindcss; exit 1; }; \
		chmod +x ./tailwindcss; \
		echo "$(GREEN)✓ tailwindcss instalado$(NC)"; \
	fi

verify: ## Verificar que todas las dependencias están instaladas
	@MISSING=0; \
	echo "$(BLUE)Verificando dependencias...$(NC)"; \
	for tool in go templ curl; do \
		if command -v $$tool >/dev/null 2>&1; then \
			echo "  $(GREEN)[OK]$(NC) $$tool"; \
		else \
			echo "  $(RED)[MISSING]$(NC) $$tool"; \
			MISSING=1; \
		fi; \
	done; \
	if [ -f ./tailwindcss ]; then \
		echo "  $(GREEN)[OK]$(NC) tailwindcss (local binary)"; \
	else \
		echo "  $(RED)[MISSING]$(NC) tailwindcss — run: make setup"; \
		MISSING=1; \
	fi; \
	[ $$MISSING -eq 0 ]

# ==================== DEVELOPMENT ====================

dev: templ assets ## Modo desarrollo: genera templates, assets y ejecuta con hot reload
	@echo "$(BLUE)🚀 Iniciando modo desarrollo...$(NC)"
	@echo "$(YELLOW)Go version: $(GO_VERSION)$(NC)"
	@clear
	@go run $(CMD_PATH)

dev-static: assets ## Solo generar assets (sin ejecutar servidor)
	@echo "$(GREEN)✓ Assets generados$(NC)"

# ==================== BUILD & COMPILATION ====================

build: clean templ assets lint ## Compilar para producción (con validaciones)
	@echo "$(BLUE)🔨 Compilando kiosco $(VERSION)...$(NC)"
	@go build -ldflags="$(LDFLAGS)" -o $(BINARY_PATH) $(CMD_PATH)
	@echo "$(GREEN)✓ Binario compilado: $(BINARY_PATH)$(NC)"
	@ls -lh $(BINARY_PATH)

build-quick: templ ## Compilación rápida (sin validaciones ni assets)
	@echo "$(BLUE)⚡ Compilación rápida...$(NC)"
	@go build -o $(BINARY_PATH) $(CMD_PATH)
	@echo "$(GREEN)✓ Binario compilado: $(BINARY_PATH)$(NC)"

build-prod: clean templ assets lint test ## Build de producción con tests
	@echo "$(BLUE)🏭 Build de PRODUCCIÓN con tests...$(NC)"
	@go test ./... -v
	@go build -ldflags="$(LDFLAGS)" -o $(BINARY_PATH) $(CMD_PATH)
	@echo "$(GREEN)✓ Build de producción completado$(NC)"

build-linux: clean templ assets ## Compilar binario Linux/amd64 para VPS
	@echo "$(BLUE)🐧 Compilando para Linux (amd64)...$(NC)"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_PATH) $(CMD_PATH)
	@echo "$(GREEN)✓ Binario Linux compilado: $(BINARY_PATH)$(NC)"
	@ls -lh $(BINARY_PATH)
	@echo "$(YELLOW)▶ Próximo paso: scripts/deploy.sh$(NC)"

# ==================== RUNTIME ====================

run: ## Ejecutar el binario compilado
	@if [ ! -f $(BINARY_PATH) ]; then \
		echo "$(RED)❌ Binario no encontrado. Ejecuta 'make build' primero$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)▶ Ejecutando $(BINARY_PATH)...$(NC)"
	@PORT=3200 $(BINARY_PATH)

run-prod: ## Ejecutar binario de producción (sin logs verbosos)
	@echo "$(BLUE)▶ Ejecutando en modo producción...$(NC)"
	@PORT=3200 $(BINARY_PATH) 2>&1 | tee kiosco.log

# ==================== CODE GENERATION ====================

templ: ## Regenerar templates (templ generate)
	@echo "$(BLUE)📝 Generando templates...$(NC)"
	@templ generate
	@echo "$(GREEN)✓ Templates generados$(NC)"

assets: ## Regenerar assets (CSS/JS)
	@if [ ! -f ./tailwindcss ]; then \
		echo "$(RED)❌ tailwindcss no encontrado — run: make setup$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)🎨 Construyendo assets (CSS/JS)...$(NC)"
	@$(MAKE) assets-css
	@$(MAKE) assets-js
	@echo "$(GREEN)✓ Assets construidos$(NC)"

assets-css: ## Solo construir CSS (Tailwind)
	@echo "$(BLUE)🎨 Compilando CSS...$(NC)"
	@mkdir -p public/dist
	@./tailwindcss -i ./assets/main.css -o ./public/dist/styles.css --minify
	@echo "$(GREEN)✓ CSS compilado$(NC)"

assets-js: ## Solo descargar/compilar JS
	@echo "$(BLUE)🎨 Descargando/copiando JavaScript...$(NC)"
	@mkdir -p public/dist
	@curl -fsSL https://cdn.jsdelivr.net/npm/alpinejs@3.15.4/dist/cdn.min.js -o public/dist/alpine.min.js
	@curl -fsSL https://cdnjs.cloudflare.com/ajax/libs/htmx/2.0.7/htmx.min.js -o public/dist/htmx.min.js
	@curl -fsSL https://cdn.jsdelivr.net/npm/html-to-image@1.11.11/dist/html-to-image.min.js -o public/dist/canvas.min.js
	@curl -fsSL https://cdn.jsdelivr.net/npm/@alpinejs/collapse@3.x.x/dist/cdn.min.js -o public/dist/collapse.min.js
	@cp assets/main.js public/dist/bundle.min.js
	@echo "$(GREEN)✓ JavaScript compilado$(NC)"

# ==================== CODE QUALITY ====================

lint: ## Ejecutar linter (go vet)
	@echo "$(BLUE)🔍 Linting...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Lint completado (sin errores)$(NC)"

fmt: ## Formatear código (gofmt)
	@echo "$(BLUE)📐 Formateando código...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Código formateado$(NC)"

test: ## Ejecutar tests
	@echo "$(BLUE)🧪 Ejecutando tests...$(NC)"
	@go test ./... -v
	@echo "$(GREEN)✓ Tests completados$(NC)"

test-coverage: ## Tests con coverage
	@echo "$(BLUE)🧪 Tests con coverage...$(NC)"
	@go test ./... -v -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage generado en coverage.html$(NC)"

# ==================== CLEANUP ====================

clean: ## Limpiar artifacts (binario, coverage, etc)
	@echo "$(YELLOW)🗑 Limpiando artifacts...$(NC)"
	@rm -f $(BINARY_PATH)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Limpieza completada$(NC)"

clean-full: clean ## Limpieza profunda (+ tailwindcss binary, node_modules, dist)
	@echo "$(YELLOW)🗑 Limpieza profunda...$(NC)"
	@rm -rf node_modules
	@rm -rf public/dist/*
	@rm -f ./tailwindcss
	@echo "$(GREEN)✓ Limpieza profunda completada$(NC)"

# ==================== DATABASE ====================

db-wal: ## Activar WAL mode en la DB de producción (LEER SETUP_WAL_PRODUCCION.md PRIMERO)
	@echo "$(RED)⚠️  ADVERTENCIA: Este comando afecta la DB de producción$(NC)"
	@echo "$(YELLOW)Asegúrate de haber leído SETUP_WAL_PRODUCCION.md$(NC)"
	@read -p "Continuar? (s/n): " confirm; \
	if [ "$$confirm" = "s" ]; then \
		sqlite3 database.db "PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL; PRAGMA busy_timeout=5000; PRAGMA foreign_keys=ON;"; \
		echo "$(GREEN)✓ WAL activado$(NC)"; \
		sqlite3 database.db "PRAGMA journal_mode;"; \
	else \
		echo "$(YELLOW)Cancelado$(NC)"; \
	fi

db-verify: ## Verificar integridad de la DB
	@echo "$(BLUE)🔍 Verificando integridad de DB...$(NC)"
	@sqlite3 database.db "PRAGMA integrity_check;" && echo "$(GREEN)✓ DB intacta$(NC)"

# ==================== INFO ====================

info: ## Mostrar información del proyecto
	@echo "$(BLUE)Información del Proyecto$(NC)"
	@echo "$(BLUE)======================$(NC)"
	@echo "Nombre: Kiosco"
	@echo "Versión: $(VERSION)"
	@echo "Compilado: $(BUILD_TIME)"
	@echo "Go: $(GO_VERSION)"
	@echo "Binario: $(BINARY_PATH)"
	@echo ""
	@echo "Archivos relevantes:"
	@echo "  - Código: $(CMD_PATH)/main.go"
	@echo "  - Config: internal/config/"
	@echo "  - Modelos: internal/models/"
	@echo "  - Repositorios: internal/repositories/"
	@echo "  - Templates: templates/"
	@echo "  - Assets: assets/"

status: ## Verificar estado del sistema
	@echo "$(BLUE)Estado del Sistema$(NC)"
	@echo "$(BLUE)=================$(NC)"
	@echo -n "Go: "; go version
	@echo -n "Templ: "; templ version 2>/dev/null || echo "no instalado"
	@echo -n "TailwindCSS: "; ./tailwindcss --version 2>/dev/null || echo "no instalado (run: make setup)"
	@echo -n "SQLite: "; sqlite3 --version 2>/dev/null || echo "no instalado"
	@echo -n "Binario compilado: "
	@if [ -f $(BINARY_PATH) ]; then echo "$(GREEN)sí$(NC)"; else echo "$(RED)no$(NC)"; fi

# ==================== ALIASES (atajos) ====================

d: dev ## Atajo: make d = make dev
b: build ## Atajo: make b = make build
r: run ## Atajo: make r = make run
c: clean ## Atajo: make c = make clean
f: fmt ## Atajo: make f = make fmt

.PHONY: docker-build docker-run
