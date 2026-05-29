#!/bin/bash

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
NC='\033[0m'

# Verify database integrity and WAL status
# WAL is activated automatically by the Go server in internal/config/database.go
verify_db() {
    echo -e "${BLUE}🔍 Verificando integridad de DB...${NC}"

    # Check integrity
    if ! sqlite3 database.db "PRAGMA integrity_check;"; then
        echo "Database integrity check failed!"
        return 1
    fi

    echo -e "${GREEN}✓ DB intacta${NC}"
    echo ""

    # Show WAL status
    echo -e "${BLUE}Estado de WAL:${NC}"
    sqlite3 database.db "PRAGMA journal_mode;"
    echo ""

    # Show database files
    echo -e "${BLUE}Archivos de DB:${NC}"
    ls -lh database.db* 2>/dev/null || echo "No database files found"
}

main() {
    case "${1:-verify}" in
        verify)
            verify_db
            ;;
        *)
            echo "Usage: $0 verify"
            echo ""
            echo "Diagnostics for SQLite database."
            echo "WAL mode is activated automatically by the server on startup."
            exit 1
            ;;
    esac
}

main "$@"
