.PHONY: help dev dev-static build build-quick build-prod build-linux run run-prod clean clean-full templ assets assets-css assets-js lint fmt test test-coverage db-wal db-verify info status setup verify
.DEFAULT_GOAL := help

# Colors for output
BLUE=\033[0;34m
GREEN=\033[0;32m
NC=\033[0m

help: ## Mostrar este mensaje de ayuda
	@echo "$(BLUE)Kiosco - Sistema de Control de Consumo Escolar$(NC)"
	@echo "$(BLUE)===========================================$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""

# ==================== SETUP ====================

setup: ## Descargar binario tailwindcss (idempotente)
	@bash scripts/setup.sh download

verify: ## Verificar que todas las dependencias están instaladas
	@bash scripts/setup.sh verify

# ==================== DEVELOPMENT ====================

dev: templ assets ## Modo desarrollo con hot reload
	@echo "$(BLUE)🚀 Iniciando modo desarrollo...$(NC)"
	@go run ./cmd/kiosco

dev-static: assets ## Solo generar assets (sin ejecutar servidor)
	@echo "$(GREEN)✓ Assets generados$(NC)"

# ==================== BUILD & COMPILATION ====================

build: clean templ assets lint ## Compilar para producción (con validaciones)
	@bash scripts/build.sh standard

build-quick: templ ## Compilación rápida (sin validaciones ni assets)
	@bash scripts/build.sh quick

build-prod: clean templ assets lint test ## Build de producción con tests
	@bash scripts/build.sh prod

build-linux: clean templ assets ## Compilar binario Linux/amd64 para VPS
	@bash scripts/build.sh linux

# ==================== RUNTIME ====================

run: ## Ejecutar el binario compilado
	@if [ ! -f bin/kiosco ]; then \
		echo "$(RED)❌ Binario no encontrado. Ejecuta 'make build' primero$(NC)"; \
		exit 1; \
	fi
	@PORT=3200 bin/kiosco

run-prod: ## Ejecutar binario de producción (sin logs verbosos)
	@PORT=3200 bin/kiosco 2>&1 | tee kiosco.log

# ==================== CODE GENERATION ====================

templ: ## Regenerar templates (templ generate)
	@templ generate

assets: ## Regenerar assets (CSS/JS)
	@bash scripts/assets.sh all

assets-css: ## Solo construir CSS (Tailwind)
	@bash scripts/assets.sh css

assets-js: ## Solo descargar/compilar JS
	@bash scripts/assets.sh js

# ==================== CODE QUALITY ====================

lint: ## Ejecutar linter (go vet)
	@bash scripts/test.sh lint

fmt: ## Formatear código (gofmt)
	@bash scripts/test.sh fmt

test: ## Ejecutar tests
	@bash scripts/test.sh test

test-coverage: ## Tests con coverage
	@bash scripts/test.sh coverage

# ==================== CLEANUP ====================

clean: ## Limpiar artifacts (binario, coverage, etc)
	@bash scripts/cleanup.sh clean

clean-full: ## Limpieza profunda (+ tailwindcss binary, node_modules, dist)
	@bash scripts/cleanup.sh full

# ==================== DATABASE ====================

db-verify: ## Verificar integridad de la DB y estado de WAL
	@bash scripts/db.sh verify

# ==================== INFO ====================

info: ## Mostrar información del proyecto
	@bash scripts/info.sh info

status: ## Verificar estado del sistema
	@bash scripts/info.sh status

# ==================== ALIASES (atajos) ====================

d: dev ## Atajo: make d = make dev
b: build ## Atajo: make b = make build
r: run ## Atajo: make r = make run
c: clean ## Atajo: make c = make clean
f: fmt ## Atajo: make f = make fmt
