@echo off
REM --- Plantilla de variables de entorno para desarrollo ---
REM Copia este archivo como "env.bat" y completa con tus valores reales
REM NUNCA subas "env.bat" a git (debe estar en .gitignore)

REM Configuración de base de datos
set DB_USER=postgres
set DB_PASSWORD=TU_CONTRASEÑA_AQUI
set DB_HOST=localhost
set DB_PORT=5432
set DB_NAME=Kiosco
set PGSSLMODE=disable

REM --- Confirmar que las variables de entorno se cargaron correctamente ---
echo [env] Variables de entorno cargadas con exito.

REM --- Ejecutar la aplicación Go ---
echo [app] Ejecutando la aplicacion Go...
go run .

pause
