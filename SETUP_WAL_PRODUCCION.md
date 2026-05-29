# Activar WAL Mode en Base de Datos de Producción

## Qué es WAL y por qué es crítico

**WAL = Write-Ahead Logging**

Sin WAL (modo por defecto en SQLite):
- Cuando alguien escribe a la DB, la tabla completa se bloquea
- Otros usuarios que intentan leer/escribir reciben "database is locked"
- Con 3-6 usuarios simultáneos registrando consumos, esto genera errores

Con WAL:
- Múltiples lectores pueden leer MIENTRAS alguien escribe
- Escrituras se colan en una cola sin bloquear lecturas
- Performance: +3x mejor bajo concurrencia

## ⚠️ CRÍTICO: El código de kiosco YA intenta activar WAL

El servidor ahora intenta esto automáticamente en `internal/config/database.go`:

```go
pragmas := []string{
	"PRAGMA journal_mode=WAL",
	"PRAGMA synchronous=NORMAL",
	"PRAGMA busy_timeout=5000",
	"PRAGMA foreign_keys=ON",
}
```

**PERO**: Si la DB de producción está en estado antiguo (sin WAL), necesitas migrar manualmente antes de deployar.

---

## Procedimiento Seguro para Migración en Producción

### Paso 1: Backup de la DB actual

```bash
# En el servidor de producción
cd /ruta/a/datos  # Donde esté database.db

# Hacer backup
cp database.db database.db.backup.$(date +%Y%m%d_%H%M%S)
ls -lh database.db*

# Verificar que el backup es sólido
sqlite3 database.db.backup.* "SELECT COUNT(*) FROM estudiante;" # Debería retornar un número
```

### Paso 2: Verificar estado actual de la DB

```bash
# Conectarse a la DB de producción
sqlite3 database.db

# Dentro de sqlite3:
sqlite> PRAGMA journal_mode;
# Retorna: delete (significa SIN WAL, necesita migración)

sqlite> SELECT COUNT(*) FROM estudiante;
# Verificar que la DB es accesible

sqlite> .quit
```

### Paso 3: Activar WAL (una sola vez)

```bash
sqlite3 database.db << 'EOF'
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
PRAGMA busy_timeout=5000;
PRAGMA foreign_keys=ON;
.quit
EOF
```

**Qué pasa**: 
- Se crean dos archivos nuevos: `database.db-wal` y `database.db-shm`
- No se pierden datos
- Es idempotente (si se ejecuta 2 veces, funciona igual)

### Paso 4: Verificar que WAL está activo

```bash
sqlite3 database.db "PRAGMA journal_mode;"
# Debería retornar: wal
```

### Paso 5: Ver los archivos creados

```bash
ls -lh database.db*
# Deberías ver:
# - database.db         (la DB principal)
# - database.db-wal     (write-ahead log, NUEVO)
# - database.db-shm     (shared memory, NUEVO)
# - database.db.backup.* (tu backup)
```

### Paso 6: Backup de los archivos WAL

```bash
# Hacer backup de los archivos WAL también (son parte de la DB)
cp database.db-wal database.db-wal.backup.$(date +%Y%m%d_%H%M%S) 2>/dev/null || echo "WAL no existe aún"
cp database.db-shm database.db-shm.backup.$(date +%Y%m%d_%H%M%S) 2>/dev/null || echo "SHM no existe aún"
```

---

## Después: Deployar el nuevo código

Una vez que WAL está activado en la DB:

1. Deploy el binario con `resource-optimization` (ya incluye los PRAGMAs)
2. El servidor intentará confirmar WAL al startup (es idempotente, no causa problemas)
3. Múltiples usuarios pueden escribir sin conflictos

---

## Si algo falla: Rollback

### Opción A: Volver a la DB antigua (sin pérdida de datos post-migración)

```bash
# Restaurar el backup
rm database.db database.db-wal database.db-shm
cp database.db.backup.FECHA database.db

# Verificar
sqlite3 database.db "SELECT COUNT(*) FROM estudiante;"
```

### Opción B: Desactivar WAL (mantener los datos, volver a modo antiguo)

```bash
sqlite3 database.db << 'EOF'
PRAGMA journal_mode=DELETE;
PRAGMA synchronous=FULL;
.quit
EOF

# Esperar a que se consolide (puede tardar segundos)
ls -lh database.db*
# Los archivos -wal y -shm deberían desaparecer
```

---

## Checklist de Producción

```
□ Hice backup de la DB actual (database.db.backup.*)
□ Verifiqué que el backup se puede leer (SELECT COUNT(*) FROM estudiante)
□ Ejecuté PRAGMA journal_mode=WAL
□ Verifiqué que PRAGMA journal_mode retorna 'wal'
□ Verifiqué que los archivos -wal y -shm existen
□ Hice backup de los archivos WAL también
□ Deploy del nuevo binario (con code que intenta WAL)
□ Probé que múltiples usuarios simultáneos no ven "database is locked"
□ Monitoreo de memory y performance (debería mejorar)
```

---

## Monitoreo Post-Migración

Después de migrar a WAL, monitorea:

```bash
# Size de la DB (WAL puede crecer temporalmente)
du -h database.db*

# Processes usando la DB
lsof | grep database.db

# Logs del servidor (no debería haber "database is locked")
tail -f /ruta/a/logs/kiosco.log | grep -i "locked\|pragma"
```

---

## FAQ

**P: ¿Es reversible?**
R: Sí, puedes volver a modo DELETE ejecutando `PRAGMA journal_mode=DELETE`. Sin pérdida de datos.

**P: ¿Pierde datos al migrar?**
R: No. WAL es un formato adicional, no destructivo. Los datos existentes se preservan.

**P: ¿Se reintenta el PRAGMA si falló una vez?**
R: Sí. El código intenta `PRAGMA journal_mode=WAL` cada vez que inicia el servidor. Es safe.

**P: ¿Necesito reiniciar el servidor?**
R: No es obligatorio, pero es recomendable deployar el nuevo binario después de migrar la DB. Así el código ya espera WAL.

**P: ¿Cuánto tarda la migración?**
R: Unos segundos (depende del size de la DB). Para kiosco con 3-6 meses de datos: <1 segundo.

---

## Comandos rápidos (copia-pega)

```bash
# 1. Backup
cp database.db database.db.backup.$(date +%Y%m%d_%H%M%S)

# 2. Activar WAL
sqlite3 database.db "PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL; PRAGMA busy_timeout=5000; PRAGMA foreign_keys=ON;"

# 3. Verificar
sqlite3 database.db "PRAGMA journal_mode;"

# 4. Ver archivos
ls -lh database.db*
```

---

## Soporte

Si algo falla:
1. Verifica que la DB no esté bloqueada (`lsof | grep database.db`)
2. Revisa logs del servidor
3. Restore desde backup si es necesario
4. Contacta a desarrollo si persisten errores

---

**Documento generado**: 2026-05-05
**Versión de kiosco**: post resource-optimization
**WAL Mode**: CRÍTICO para producción con 3+ usuarios concurrentes
