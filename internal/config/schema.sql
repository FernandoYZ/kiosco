-- 1. Tabla de Grados
CREATE TABLE grados (
    id_grado INTEGER PRIMARY KEY AUTOINCREMENT,
    nivel_grado TEXT NOT NULL, -- Primaria, Secundaria
    nombre_grado TEXT NOT NULL, -- Quinto, Sexto...
    anio_grado TEXT NOT NULL    -- 5to, 6to...
);

-- 2. Tabla de Estudiantes
CREATE TABLE estudiantes (
    id_estudiante INTEGER PRIMARY KEY AUTOINCREMENT,
    nombres TEXT NOT NULL,
    apellidos TEXT NOT NULL,
    id_grado INTEGER,
    esta_activo INTEGER DEFAULT 1 NOT NULL, -- 1 para True, 0 para False
    FOREIGN KEY (id_grado) REFERENCES grados(id_grado)
);

-- 3. Tabla de Productos
CREATE TABLE productos (
    id_producto INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    precio_unitario NUMERIC(10, 2) NOT NULL,
    esta_activo INTEGER DEFAULT 1 NOT NULL
);

-- 4. Tabla de Consumos con Columna Generada
CREATE TABLE consumos (
    id_consumo INTEGER PRIMARY KEY AUTOINCREMENT,
    id_estudiante INTEGER NOT NULL,
    id_producto INTEGER NOT NULL,
    cantidad INTEGER NOT NULL DEFAULT 1,
    precio_unitario_venta NUMERIC(10, 2) NOT NULL,
    -- Columna generada automáticamente
    total_linea NUMERIC(10, 2) GENERATED ALWAYS AS (cantidad * precio_unitario_venta) STORED,
    fecha_consumo DATE NOT NULL,
    FOREIGN KEY (id_estudiante) REFERENCES estudiantes(id_estudiante),
    FOREIGN KEY (id_producto) REFERENCES productos(id_producto)
);

-- 5. Tabla de Pagos
CREATE TABLE pagos (
    id_pago INTEGER PRIMARY KEY AUTOINCREMENT,
    id_estudiante INTEGER NOT NULL,
    monto NUMERIC(10, 2) NOT NULL,
    fecha_pago DATE NOT NULL,
    FOREIGN KEY (id_estudiante) REFERENCES estudiantes(id_estudiante)
);

-- 6. Tabla de Usuarios (Nueva)
CREATE TABLE usuarios (
    id_usuario INTEGER PRIMARY KEY AUTOINCREMENT,
    usuario TEXT NOT NULL UNIQUE,
    contrasenha TEXT NOT NULL,
    puede_editar INTEGER NOT NULL DEFAULT 0
);

---
--- ÍNDICES
---
CREATE INDEX idx_consumos_estudiante_fecha ON consumos(id_estudiante, fecha_consumo);
CREATE INDEX idx_pagos_estudiante_fecha ON pagos(id_estudiante, fecha_pago);
CREATE INDEX idx_estudiantes_activo ON estudiantes(esta_activo) WHERE esta_activo = 1;
CREATE INDEX idx_estudiantes_apellido_nombre ON estudiantes(apellidos, nombres);

---
--- DATOS POR DEFECTO
---

INSERT INTO grados (nivel_grado, nombre_grado, anio_grado) VALUES
('Primaria', 'Quinto', '5to'),
('Primaria', 'Sexto', '6to'),
('Secundaria', 'Primero', '1ro'),
('Secundaria', 'Segundo', '2do'),
('Secundaria', 'Tercero', '3ro'),
('Secundaria', 'Cuarto', '4to'),
('Secundaria', 'Quinto', '5to');

INSERT INTO productos (nombre, precio_unitario) VALUES
('Gelatina', 1.50),
('Comida 1', 5.00),
('Keke', 1.50),
('Postre', 2.00),
('Chaufa', 3.50),
('Comida 2', 3.00);

-- Agregar usuario de prueba
INSERT INTO usuarios (usuario, contrasenha) VALUES
('prueba', '$argon2id$v=19$m=16,t=3,p=1$bG1tMWowTEZMVTVSVU5VYg$9Ib8jWAgKLD1PSaZoIdBBA');
-- contraseña: pa$$w0rD

-- database: ..\database\database.db 
-- Directorio al ejecutar el binario
