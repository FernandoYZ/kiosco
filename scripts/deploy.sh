#!/bin/bash

################################################################################
# Deploy Script para Kiosco
# Descripción: Maneja sincronización segura de DB y binarios a VPS
# Uso: ./scripts/deploy.sh [comando]
################################################################################

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Config
VPS_HOST="${VPS_HOST:-vps}"
VPS_USER="${VPS_USER:-root}"
VPS_PATH="/opt/kiosco"
DB_PATH="database/database.db"
BINARY_PATH="bin/kiosco"
AUDIT_LOG="deploy.audit.log"

# Functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $*"
}

success() {
    echo -e "${GREEN}✓${NC} $*"
}

error() {
    echo -e "${RED}✗${NC} $*" >&2
}

warning() {
    echo -e "${YELLOW}⚠${NC} $*"
}

audit() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" >> "$AUDIT_LOG"
}

confirm() {
    local prompt="$1"
    local response

    read -p "$(echo -e "${YELLOW}${prompt}${NC}") (s/n): " response
    [[ "$response" == "s" ]] && return 0 || return 1
}

check_dependencies() {
    log "Verificando dependencias..."

    for cmd in scp ssh date; do
        if ! command -v "$cmd" &> /dev/null; then
            error "$cmd no encontrado"
            exit 1
        fi
    done

    success "Dependencias OK"
}

check_vps_connectivity() {
    log "Verificando conectividad con VPS..."

    if ! ssh -o ConnectTimeout=5 "$VPS_HOST" "test -d $VPS_PATH" 2>/dev/null; then
        error "No se puede conectar a $VPS_HOST:$VPS_PATH"
        error "Verifica tu ~/.ssh/config tiene:"
        error "  Host $VPS_HOST"
        error "      HostName <ip-o-dominio>"
        error "      User <usuario>"
        error "      IdentityFile ~/.ssh/id_ed25519"
        exit 1
    fi

    success "VPS conectado: $VPS_HOST"
}

################################################################################
# COMANDOS DISPONIBLES
################################################################################

cmd_help() {
    cat << 'EOF'
╔════════════════════════════════════════════════════════════════════════════╗
║                     Kiosco Deploy Script - Ayuda                           ║
╚════════════════════════════════════════════════════════════════════════════╝

COMANDOS:

  db:backup         Hacer backup local de la DB antes de cualquier operación
  db:pull           Traer database.db del VPS a local (dev)
  db:push           Enviar database.db local al VPS (⚠️  DESTRUCTIVO)
  binary:push       Enviar binario compilado al VPS
  status            Ver estado de VPS
  help              Mostrar esta ayuda

FLUJO RECOMENDADO:

  1. make build-linux     # Compilar binario para Linux
  2. ./scripts/deploy.sh db:backup    # Backup de precaución
  3. ./scripts/deploy.sh binary:push  # Enviar binario
  4. ssh vps 'systemctl restart kiosco'  # Reiniciar servicio

CONFIGURACIÓN SSH:

  Agregar a ~/.ssh/config:

  Host vps
      HostName 1.2.3.4
      User root
      IdentityFile ~/.ssh/id_ed25519
      StrictHostKeyChecking accept-new

VARIABLES DE ENTORNO:

  VPS_HOST   Host SSH en ~/.ssh/config (default: vps)

EJEMPLOS:

  ./scripts/deploy.sh db:backup
  VPS_HOST=prod ./scripts/deploy.sh binary:push
  ./scripts/deploy.sh status

EOF
}

cmd_db_backup() {
    log "Creando backup local de DB..."

    if [ ! -f "$DB_PATH" ]; then
        error "DB no encontrada: $DB_PATH"
        exit 1
    fi

    local backup_file="database/database.db.backup.$(date +%Y%m%d_%H%M%S)"
    cp "$DB_PATH" "$backup_file"

    success "Backup creado: $backup_file"
    ls -lh "$backup_file"
    audit "DB backup: $backup_file"
}

cmd_db_pull() {
    log "Trayendo database.db del VPS..."
    warning "Esto SOBRESCRIBIRÁ tu DB local"

    if ! confirm "¿Estás seguro de traer la DB del VPS ($VPS_HOST)?"; then
        log "Cancelado"
        return 0
    fi

    check_vps_connectivity

    # Backup de precaución
    if [ -f "$DB_PATH" ]; then
        local backup="database/database.db.before-pull.$(date +%Y%m%d_%H%M%S)"
        cp "$DB_PATH" "$backup"
        success "Backup local creado: $backup"
    fi

    log "Descargando DB..."
    scp "$VPS_HOST:$VPS_PATH/$DB_PATH" "$DB_PATH"

    success "DB actualizada desde VPS"
    ls -lh "$DB_PATH"
    audit "DB pulled from $VPS_HOST"
}

cmd_db_push() {
    log "Preparando para enviar database.db al VPS..."

    if [ ! -f "$DB_PATH" ]; then
        error "DB no encontrada: $DB_PATH"
        exit 1
    fi

    warning "⚠️  ALERTA: ESTO SOBRESCRIBIRÁ LA DB DE PRODUCCIÓN"
    warning "Asegúrate de:"
    warning "  1. Haber hecho backup en VPS: ssh $VPS_HOST 'cp $VPS_PATH/$DB_PATH $VPS_PATH/$DB_PATH.backup'"
    warning "  2. Que esta DB sea correcta y completa"
    warning "  3. Que NO haya usuarios activos en el VPS"

    if ! confirm "¿REALMENTE deseas enviar la DB al VPS?"; then
        log "Cancelado"
        return 0
    fi

    if ! confirm "Última oportunidad para cancelar. ¿Continuar?"; then
        log "Cancelado"
        return 0
    fi

    check_vps_connectivity

    log "Haciendo backup en VPS como precaución..."
    ssh "$VPS_HOST" "cp $VPS_PATH/$DB_PATH $VPS_PATH/$DB_PATH.backup.$(date +%Y%m%d_%H%M%S)"

    log "Enviando DB..."
    scp "$DB_PATH" "$VPS_HOST:$VPS_PATH/$DB_PATH"

    success "DB enviada al VPS"
    audit "DB pushed to $VPS_HOST (backup creado automáticamente)"

    warning "Próximo paso: ssh $VPS_HOST 'systemctl restart kiosco'"
}

cmd_binary_push() {
    log "Preparando para enviar binario al VPS..."

    if [ ! -f "$BINARY_PATH" ]; then
        error "Binario no encontrado: $BINARY_PATH"
        error "Compila primero: make build-linux"
        exit 1
    fi

    local binary_size=$(du -h "$BINARY_PATH" | cut -f1)
    log "Binario: $BINARY_PATH ($binary_size)"

    if ! confirm "¿Enviar binario al VPS?"; then
        log "Cancelado"
        return 0
    fi

    check_vps_connectivity

    log "Haciendo backup del binario anterior en VPS..."
    ssh "$VPS_HOST" "test -f $VPS_PATH/kiosco && cp $VPS_PATH/kiosco $VPS_PATH/kiosco.backup.$(date +%Y%m%d_%H%M%S) || true"

    log "Enviando binario..."
    scp "$BINARY_PATH" "$VPS_HOST:$VPS_PATH/kiosco"

    log "Haciendo ejecutable..."
    ssh "$VPS_HOST" "chmod +x $VPS_PATH/kiosco"

    success "Binario enviado y configurado"
    audit "Binary pushed to $VPS_HOST (backup: kiosco.backup.*)"

    log ""
    log "$(echo -e "${YELLOW}Próximos pasos:${NC}")"
    echo "  1. ssh $VPS_HOST"
    echo "  2. systemctl stop kiosco"
    echo "  3. systemctl start kiosco"
    echo "  4. systemctl status kiosco"
    echo "  5. curl http://localhost:3200/login"
}

cmd_status() {
    log "Obteniendo estado de VPS..."
    check_vps_connectivity

    log "Información del VPS:"
    ssh "$VPS_HOST" << 'EOSSH'
        echo "  Versión del binario:"
        /opt/kiosco/kiosco --version 2>/dev/null || echo "    (no disponible)"
        echo ""
        echo "  Estado del servicio:"
        systemctl status kiosco --no-pager -n 5 2>/dev/null || echo "    (no configurado)"
        echo ""
        echo "  Última vez ejecutado:"
        ls -lh /opt/kiosco/kiosco
        echo ""
        echo "  DB:"
        ls -lh /opt/kiosco/database/database.db*
EOSSH
}

################################################################################
# MAIN
################################################################################

main() {
    local cmd="${1:-help}"

    # Validaciones básicas
    if [ ! -d "scripts" ]; then
        error "Ejecuta este script desde la raíz del proyecto"
        exit 1
    fi

    case "$cmd" in
        help)
            cmd_help
            ;;
        db:backup)
            check_dependencies
            cmd_db_backup
            ;;
        db:pull)
            check_dependencies
            cmd_db_pull
            ;;
        db:push)
            check_dependencies
            cmd_db_push
            ;;
        binary:push)
            check_dependencies
            cmd_binary_push
            ;;
        status)
            check_dependencies
            cmd_status
            ;;
        *)
            error "Comando desconocido: $cmd"
            echo ""
            cmd_help
            exit 1
            ;;
    esac
}

main "$@"
