# Kiosco Deploy Guide

## 🚀 Flujo de Deploy Recomendado

```bash
# 1. Compilar para Linux
make build:linux

# 2. Backup local (precaución)
./scripts/deploy.sh db:backup

# 3. Enviar binario al VPS
./scripts/deploy.sh binary:push

# 4. SSH al VPS y reiniciar
ssh vps
systemctl restart kiosco
systemctl status kiosco
```

---

## 📋 Configuración SSH Previa

### Crear SSH key (si no existe)
```bash
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "kiosco-deploy"
ssh-copy-id -i ~/.ssh/id_ed25519.pub root@VPS_IP
```

### Agregar a ~/.ssh/config
```
Host vps
    HostName 1.2.3.4              # Reemplazar con IP real
    User root
    IdentityFile ~/.ssh/id_ed25519
    StrictHostKeyChecking accept-new
```

### Verificar conectividad
```bash
ssh vps "echo OK"
```

---

## 🛠️ Comandos de Deploy

### Compilar binario para Linux
```bash
make build:linux
# Genera: bin/kiosco (optimizado para VPS)
```

### Backup Local (SEGURO - siempre hacer antes)
```bash
./scripts/deploy.sh db:backup
# Copia: database/database.db → database/database.db.backup.TIMESTAMP
```

### Enviar Binario al VPS (RECOMENDADO)
```bash
./scripts/deploy.sh binary:push
# - Verifica conectividad SSH
- Crea backup del binario anterior en VPS
- Envía nuevo binario
- Lo hace ejecutable
- Pide confirmación
```

### Traer DB del VPS (para desarrollo)
```bash
./scripts/deploy.sh db:pull
# ⚠️  Sobrescribe DB local
# Crea backup automático antes
```

### Enviar DB al VPS (DESTRUCTIVO - usar con cuidado)
```bash
./scripts/deploy.sh db:push
# ⚠️  SOBRESCRIBE DB de producción
# Pide 2x confirmación
# Crea backup en VPS automáticamente
```

### Ver Estado de VPS
```bash
./scripts/deploy.sh status
# Muestra: versión binario, servicio, archivos, DB
```

---

## 📦 Binarios Compilados

### Versión para Desarrollo (local)
```bash
make build:quick
# Rápido, sin validaciones
```

### Versión para Producción (VPS Linux)
```bash
make build:linux
# Compilado con -ldflags "-s -w" (optimizado)
# Arquitectura: amd64 (Intel/AMD 64-bit)
# SO: Linux
```

### Variables de Entorno
```bash
CGO_ENABLED=0      # Sin dependencias C (máxima portabilidad)
GOOS=linux         # Sistema operativo: Linux
GOARCH=amd64       # Arquitectura: 64-bit Intel/AMD
```

---

## ⚠️ Seguridad & Backups

### Backup Strategy
```
Antes de push:
├── ./scripts/deploy.sh db:backup      (local)
└── ssh vps cp /opt/kiosco/database.db /opt/kiosco/database.db.backup

Después de push:
├── VPS automáticamente hace backup    (script genera .backup.TIMESTAMP)
└── Logs en deploy.audit.log           (auditoría local)
```

### Rollback Rápido
```bash
# Si algo falló, revertir binario en VPS:
ssh vps "cp /opt/kiosco/kiosco.backup.* /opt/kiosco/kiosco && systemctl restart kiosco"

# Si DB falló, restaurar:
ssh vps "cp /opt/kiosco/database.db.backup.* /opt/kiosco/database.db"
```

---

## 🔍 Troubleshooting

### "SSH connection refused"
```bash
# Verificar SSH keys
ls -la ~/.ssh/id_ed25519*

# Probar conexión
ssh -v vps "echo OK"
```

### "database.db.backup not found"
```bash
# Crear backup manualmente
./scripts/deploy.sh db:backup
```

### "Binario no encontrado"
```bash
# Compilar primero
make build:linux
```

### Script retorna error
```bash
# Ver audit log
cat deploy.audit.log

# Ver help
./scripts/deploy.sh help
```

---

## 📊 Auditoría

Todos los deploys se registran en `deploy.audit.log`:
```bash
tail -f deploy.audit.log
```

Formato:
```
[2026-05-05 04:35:22] DB pushed to vps (backup creado automáticamente)
[2026-05-05 04:36:11] Binary pushed to vps (backup: kiosco.backup.*)
```

---

## ✅ Checklist Pre-Deploy

- [ ] SSH configurado y funcional (`ssh vps echo OK`)
- [ ] Binario compilado (`make build:linux`)
- [ ] DB sincronizada localmente (`./scripts/deploy.sh db:pull`)
- [ ] Cambios testeados localmente (`make test`)
- [ ] Backup local creado (`./scripts/deploy.sh db:backup`)
- [ ] Cambios en git (`git status`)

---

## 🚀 Deploy en 3 comandos

```bash
make build:linux                 # Compilar
./scripts/deploy.sh binary:push  # Enviar
ssh vps systemctl restart kiosco # Reiniciar
```

---

**Última actualización**: 2026-05-05
**Script de deploy**: scripts/deploy.sh
**Auditoría**: deploy.audit.log