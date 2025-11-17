CREATE TABLE Grados (
   IdGrado SERIAL PRIMARY KEY,
   NivelGrado VARCHAR(20) NOT NULL,
   NombreGrado VARCHAR(20) NOT NULL,
   AnioGrado VARCHAR(3) NOT NULL
);

CREATE TABLE Estudiantes (
   IdEstudiante SERIAL PRIMARY KEY,
   Nombres VARCHAR(255) NOT NULL,
   Apellidos VARCHAR(255) NOT NULL,
   IdGrado INT,
   EstaActivo BOOLEAN DEFAULT true NOT NULL,
   FOREIGN KEY (IdGrado) REFERENCES Grados(IdGrado)
);

CREATE TABLE Productos (
   IdProducto SERIAL PRIMARY KEY,
   Nombre VARCHAR(255) NOT NULL,
   PrecioUnitario DECIMAL(10, 2) NOT NULL,
   EstaActivo BOOLEAN DEFAULT true NOT NULL
);

CREATE TABLE Consumos (
   IdConsumo BIGSERIAL PRIMARY KEY,
   IdEstudiante INT NOT NULL,
   IdProducto INT NOT NULL,
   Cantidad INT NOT NULL DEFAULT 1,
   PrecioUnitarioVenta DECIMAL(10, 2) NOT NULL,
   TotalLinea DECIMAL(10, 2) NOT NULL,
   FechaConsumo DATE NOT NULL,
   FOREIGN KEY (IdEstudiante) REFERENCES Estudiantes(IdEstudiante),
   FOREIGN KEY (IdProducto) REFERENCES Productos(IdProducto)
);

CREATE TABLE Pagos (
   IdPago SERIAL PRIMARY KEY,
   IdEstudiante INT NOT NULL,
   Monto DECIMAL(10, 2) NOT NULL,
   FechaPago TIMESTAMP NOT NULL,
   FOREIGN KEY (IdEstudiante) REFERENCES Estudiantes(IdEstudiante)
);

CREATE INDEX idx_consumos_fecha ON Consumos(FechaConsumo);
CREATE INDEX idx_consumos_estudiante_fecha ON Consumos(IdEstudiante, FechaConsumo);
CREATE INDEX idx_consumos_estudiante_producto_fecha ON Consumos(IdEstudiante, IdProducto, FechaConsumo);

CREATE INDEX idx_pagos_estudiante ON Pagos(IdEstudiante);
CREATE INDEX idx_pagos_estudiante_fecha ON Pagos(IdEstudiante, FechaPago);
CREATE INDEX idx_pagos_fecha ON Pagos(FechaPago);

CREATE INDEX idx_estudiantes_grado_activo ON Estudiantes(IdGrado, EstaActivo);
CREATE INDEX idx_estudiantes_activo ON Estudiantes(EstaActivo) WHERE EstaActivo = true;
CREATE INDEX idx_estudiantes_apellidos_nombres ON Estudiantes(Apellidos, Nombres);

CREATE INDEX idx_grados_nivel_nombre ON Grados(NivelGrado, NombreGrado);

INSERT INTO Productos (Nombre, PrecioUnitario, EstaActivo) VALUES
('Pan', 1.50, TRUE),
('Keke', 1.50, TRUE),
('Postre', 2.00, TRUE),
('Gelatina', 1.50, TRUE),
('Agua', 1.50, TRUE),
('Comida 1', 3.00, TRUE),
('Comida 2', 3.50, TRUE),
('Comida 3', 4.00, TRUE),
('Comida 4', 5.00, TRUE);

INSERT INTO Grados (NivelGrado, NombreGrado, AnioGrado) VALUES
('Primaria', 'Quinto', '5to'),
('Primaria', 'Sexto', '6to'),
('Secundaria', 'Primero', '1ro'),
('Secundaria', 'Segundo', '2do'),
('Secundaria', 'Tercero', '3ro'),
('Secundaria', 'Cuarto', '4to'),
('Secundaria', 'Quinto', '5to');


