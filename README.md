# Kiosco - Sistema de Control de Consumo Escolar

Kiosco es una aplicación web desarrollada en Go para gestionar los consumos de productos por parte de estudiantes en un quiosco escolar. Permite llevar un registro detallado de las compras diarias, calcular deudas, registrar pagos y visualizar un resumen semanal por grado.

## Características Principales

- **Vista Semanal:** Muestra un resumen de los consumos de todos los estudiantes para la semana actual o una semana seleccionada.
- **Filtrado por Grado:** Permite filtrar la vista principal para mostrar solo los estudiantes de un grado específico.
- **Registro de Consumos:** Interfaz rápida para añadir o modificar la cantidad de productos consumidos por un estudiante en un día específico.
- **Registro de Pagos:** Permite registrar los pagos (o "descuentos") realizados por los estudiantes para saldar sus deudas.
- **Cálculo de Deuda:** Calcula automáticamente la deuda total de cada estudiante, incluyendo:
    - **Deuda Anterior:** Deuda acumulada de semanas previas.
    - **Subtotal Semanal:** Suma de todos los consumos en la semana visible.
    - **Pagos/Descuentos:** Montos pagados durante la semana.
- **Edición de Consumos Diarios:** Vista dedicada para editar todos los consumos de un estudiante para un día concreto.
- **Días no Hábiles:** Permite deshabilitar días específicos (feriados, etc.) para que no se cuenten en los cálculos ni se muestren en la interfaz.
- **Interfaz Optimizada:** La interfaz está diseñada para ser rápida y eficiente, ideal para un entorno de ventas ágil como un quiosco.

## Tecnologías Utilizadas

- **Backend:** Go (Lenguaje de Programación)
- **Base de Datos:** PostgreSQL
- **Frontend:** HTML, CSS (TailwindCSS), JavaScript

## Cómo Empezar

Sigue estos pasos para configurar y ejecutar el proyecto en tu entorno local.

### Prerrequisitos

- **Go:** Asegúrate de tener Go instalado (versión 1.21 o superior).
- **PostgreSQL:** Necesitas una instancia de PostgreSQL en ejecución.

### Instalación

1.  **Clonar el Repositorio:**
    ```bash
    git clone <URL-DEL-REPOSITORIO>
    cd kiosco
    ```

2.  **Configurar la Base de Datos:**
    - Conéctate a tu servidor PostgreSQL.
    - Crea la base de datos `kiosco`.
    - Ejecuta el script `db.sql` para crear las tablas e insertar los datos iniciales.
      ```sql
      -- Ejemplo usando psql
      psql -U tu_usuario
      CREATE DATABASE kiosco;
      \c kiosco
      \i db.sql
      ```

3.  **Configurar la Conexión:**
    - El proyecto se configura a través de variables de entorno. Puedes crear un archivo `.env` en la raíz del proyecto con el siguiente contenido:
      ```
      DB_USER=postgres
      DB_PASSWORD=tu_contraseña
      DB_HOST=localhost
      DB_PORT=5432
      DB_NAME=kiosco
      ```

4.  **Ejecutar la Aplicación:**
    - Abre una terminal en la raíz del proyecto.
    - Ejecuta el siguiente comando para descargar las dependencias y arrancar el servidor:
      ```bash
      go run .
      ```
    - El servidor se iniciará en `http://localhost:3200`.

## Estructura del Proyecto

- `main.go`: Punto de entrada de la aplicación.
- `config/`: Configuración de la base de datos.
- `db.sql`: Script SQL para la base de datos.
- `models/`: Estructuras de datos (structs).
- `repository/`: Capa de acceso a datos.
- `services/`: Capa de lógica de negocio.
- `handlers/`: Capa de presentación (manejadores HTTP).
- `routes/`: Definición de las rutas de la aplicación.
- `templates/`: Plantillas HTML.
- `static/`: Archivos estáticos (CSS, JS).
- `utils/`: Funciones de utilidad.

## Endpoints de la API

La aplicación funciona con un esquema de renderizado del lado del servidor. Las rutas principales son:

- `GET /`: Muestra la vista principal con la tabla de consumos de la semana.
  - **Parámetros (Query):**
    - `fecha`: (YYYY-MM-DD) Para navegar a una semana específica.
    - `grado`: (ID numérico) Para filtrar por grado.
    - `dias_off`: (YYYY-MM-DD,YYYY-MM-DD) Fechas separadas por comas para deshabilitar días.
- `POST /registrar-consumo`: Registra o actualiza el consumo de un producto para un estudiante en una fecha.
- `POST /registrar-pago`: Registra un pago para un estudiante.
- `GET /editar-consumos`: Muestra una página para editar todos los consumos de un estudiante en un día específico.
- `POST /guardar-consumos-dia`: Guarda todos los cambios realizados en la página de edición de consumos.