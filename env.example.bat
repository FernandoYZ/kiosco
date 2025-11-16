@echo off
REM --- Variables sensibles para desarrollo ---

REM Cargar las variables de entorno
set DB_USER=postgres
set DB_PASSWORD=contrasena_segura_db
set DB_HOST=localhost
set DB_PORT=5432
set DB_NAME=Kiosco
set PGSSLMODE=disable
set PGCHANNELBINDING=disable

REM --- Confirmar que las variables de entorno se cargaron correctamente ---
echo [env] Variables de entorno cargadas con éxito.

REM --- Ejecutar la aplicación Go ---
echo [app] Ejecutando la aplicación Go...
go run .

pause
