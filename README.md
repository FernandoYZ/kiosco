# Kiosco - Sistema de Control de Consumo Escolar

Sistema web ligero y eficiente para gestionar consumos, pagos y deudas de estudiantes en kioscos escolares. Construido con **Go** y **SQLite** embebido: sin dependencias externas pesadas, sin Docker, un solo binario autocontenido.

---
### Credenciales de Acceso
Para probar el sistema sin configuración previa, utiliza:
| Usuario | Contraseña |
|------|------------|
| prueba | pa$$w0rD |

> [!IMPORTANT]
> Estas credenciales son solo para pruebas. En un entorno real, debes cambiarlas inmediatamente en el archivo `internal/config/schema.sql`.

> [!NOTE]
> El sistema crea automáticamente estos accesos en el primer arranque si no detecta una base de datos existente.

---
## Características principales

- **Vista semanal:** tabla de consumos por estudiante, grado y semana seleccionada
- **Filtrado por grado:** navegación rápida entre Primaria y Secundaria
- **Registro de consumos:** agregar y modificar consumos por producto, estudiante y fecha
- **Edición diaria:** vista dedicada para ajustar todos los productos de un día específico
- **Gestión de pagos:** registro de pagos con historial por estudiante
- **Cálculo de deuda en tiempo real:** deuda anterior + consumos de la semana − pagos
- **Setup de estudiantes:** CRUD completo con habilitación/deshabilitación
- **Setup de productos:** gestión de productos disponibles en el kiosco
- **Autenticación con sesiones firmadas:** cookies HMAC-SHA256, Argon2id para contraseñas
- **CSRF protection:** tokens únicos por sesión, validación en todos los formularios POST
- **Rate limiting:** limitación de intentos de login (5 intentos en 15 minutos)
- **Concurrency management:** límite de 30 conexiones HTTP concurrentes
- **Binario autocontenido:** estáticos y schema SQL embebidos en el binario

> [!TIP]
> Este sistema está pensado para entornos escolares con recursos limitados: instalación simple y sin dependencias externas.

---
## Vista previa del sistema

### Vista general
<p align="center">
  <img src="assets/images/captura_index.png" width="800"/>
</p>

---

### Funcionalidades del sistema

| **Inicio de sesión** | **Productos** |
|-------------------|-------------|
| <img src="assets/images/captura_login.png" width="100%"> | <img src="assets/images/captura_agregar_producto.png" width="100%"> |

| **Estudiantes** | **Consumos** |
|---------------|------------|
| <img src="assets/images/captura_setup_estudiantes.png" width="100%"> | <img src="assets/images/captura_editar_producto.png" width="100%"> |

| **Registro** | **Resumen de consumos** |
|---------------|------------|
| <img src="assets/images/captura_registro_in_roles.png" width="100%"> | <img src="assets/images/captura_resumen_consumo.png" width="100%"> |

---
## Tecnologías utilizadas

| Capa | Tecnología |
|------|------------|
| Backend | Go 1.25+, `net/http` (sin frameworks) |
| Base de datos | SQLite (`modernc.org/sqlite` — pure Go, sin CGO) |
| Templates | `a-h/templ` — templates compilados a Go |
| Autenticación | `golang.org/x/crypto` — Argon2id |
| CSS | TailwindCSS v4 CLI (binario standalone) |
| JS interactividad | Alpine.js 3.x, HTMX 2.x (CDN) |
| Exportación | html-to-image |
| Build frontend | Makefile + curl (sin Node.js) |

---
## Inicio rápido

### Prerrequisitos

- [Go 1.25+](https://go.dev/dl/)
- [templ CLI](https://templ.guide/quick-start/installation) — para generar templates Go
- `curl` — para descargar TailwindCSS CLI

### Desarrollo local

```bash
# 1. Descargar dependencias (TailwindCSS, Alpine.js, HTMX, etc)
make setup

# 2. Iniciar con hot reload
make dev
```

Acceder en: **http://localhost:3200**

### Comandos útiles

```bash
make build          # Compilar para producción
make build-linux    # Compilar para VPS (Linux amd64)
make test           # Ejecutar tests
make fmt            # Formatear código
make clean          # Limpiar artifacts
make help           # Ver todos los comandos
```

Variables de entorno opcionales

| Variable | Por defecto |
|----------|-------------|
| `HOST` | `localhost` |
| `PORT` | `3200` |

> [!NOTE]
> No es obligatorio usar variables de entorno porque vienen por defecto
---
## CSS y JS personalizados

Los assets de entrada están en `assets/`:

- **`assets/main.css`** — punto de entrada de TailwindCSS. Agregar aquí estilos personalizados o directivas `@layer`.
- **`assets/main.js`** — código JavaScript propio (se copia a `public/dist/`).

Al ejecutar `make assets`, estos archivos se compilan (CSS con TailwindCSS) y se descargan las librerías JS (Alpine.js, HTMX) en `public/dist/`. Los archivos en `public/` se embeben en el binario al compilar con `go build`.

---
## Seguridad

### CSRF Protection
Todos los formularios POST están protegidos contra ataques CSRF:
- **Token por sesión:** cada usuario recibe un token único reutilizable
- **Validación dual:** verificación en cookie + campo oculto o header
- **Inyección automática:** tokens se inyectan en contexto para templates Templ

### Rate Limiting
- **Login:** máximo 5 intentos fallidos en 15 minutos por IP
- **Bloqueo automático:** redirección a `/login?error=rate_limit`
- **Reseteo:** contador se limpia al login exitoso

### Concurrency Management
- **Límite global:** máximo 30 conexiones HTTP concurrentes
- **Protección:** devuelve HTTP 503 si se excede el límite
- **Graceful degradation:** previene sobrecarga en SQLite

### WAL Mode (Write-Ahead Logging)
- Base de datos SQLite optimizada para concurrencia: múltiples lectores pueden leer **mientras** alguien escribe
- Mejora performance en multi-usuario: +3x mejor bajo concurrencia
- **Activado automáticamente** por el servidor en startup (ver `internal/config/database.go`)
- Si ves "database is locked" en logs, ejecuta: `make db-verify`

---
## Rutas de la aplicación

### Rutas públicas

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/login` | Formulario de inicio de sesión |
| `POST` | `/login` | Procesar credenciales |
| `GET/POST` | `/logout` | Cerrar sesión |

### Rutas protegidas (requieren sesión)

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/` | Vista principal semanal |
| `GET/POST` | `/editar-consumos` | Editar consumos del día |
| `POST` | `/guardar-consumos-dia` | Guardar cambios de edición diaria |
| `POST` | `/registrar-consumo` | Registrar consumo |
| `GET/POST` | `/editar-pagos` | Gestión de pagos |
| `POST` | `/registrar-pago` | Registrar pago |
| `POST` | `/eliminar-pago` | Eliminar pago |
| `GET` | `/ver-consumo-semanal` | Ver resumen semanal |
| `GET/POST` | `/setup` | Configuración de estudiantes |
| `GET/POST` | `/setup/estudiante` | Crear estudiante |
| `POST` | `/setup/estudiante/actualizar` | Actualizar estudiante |
| `POST` | `/setup/estudiante/toggle` | Habilitar/deshabilitar estudiante |
| `GET/POST` | `/setup/productos` | Configuración de productos |
| `GET/POST` | `/setup/producto` | Crear producto |
| `POST` | `/setup/producto/actualizar` | Actualizar producto |
| `POST` | `/setup/producto/toggle` | Habilitar/deshabilitar producto |

---
## Estructura del proyecto

```
kiosco/
├── cmd/kiosco/main.go
├── internal/
├── templates/
├── assets/
├── public/                       # Estáticos generados (embebidos en el binario)
├── database/database.db          # Creado automáticamente en el primer arranque
├── assets.go                     # embed.FS para public/
├── go.mod
├── go.sum
└── package.json
```
---
## Contacto

**Desarrollador:** Fernando YZ
**GitHub:** [@FernandoYZ](https://github.com/FernandoYZ)
