package main

import (
	"database/sql"
	"time"
)

// Repositorio maneja todas las operaciones con la base de datos
type Repositorio struct {
	db *sql.DB
}

// NuevoRepositorio crea una nueva instancia del repositorio
func NuevoRepositorio(db *sql.DB) *Repositorio {
	return &Repositorio{db: db}
}

// ObtenerEstudiantesActivos retorna todos los estudiantes activos
func (r *Repositorio) ObtenerEstudiantesActivos() ([]Estudiante, error) {
	return r.ObtenerEstudiantesPorGrado(0)
}

// ObtenerEstudiantesPorGrado retorna los estudiantes activos filtrados por grado (0 = todos)
func (r *Repositorio) ObtenerEstudiantesPorGrado(idGrado int) ([]Estudiante, error) {
	query := `
		SELECT e.IdEstudiante, e.Nombres, e.Apellidos, e.IdGrado, e.EstaActivo,
		       g.AnioGrado || ' ' || g.NivelGrado as nombre_grado
		FROM Estudiantes e
		LEFT JOIN Grados g ON e.IdGrado = g.IdGrado
		WHERE e.EstaActivo = true
	`

	var args []interface{}
	if idGrado > 0 {
		query += " AND e.IdGrado = $1"
		args = append(args, idGrado)
	}

	query += " ORDER BY g.NivelGrado, g.NombreGrado, e.Apellidos, e.Nombres"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estudiantes []Estudiante
	for rows.Next() {
		var e Estudiante
		err := rows.Scan(&e.IdEstudiante, &e.Nombres, &e.Apellidos, &e.IdGrado, &e.EstaActivo, &e.NombreGrado)
		if err != nil {
			return nil, err
		}
		estudiantes = append(estudiantes, e)
	}

	return estudiantes, nil
}

// ObtenerProductosActivos retorna todos los productos activos
func (r *Repositorio) ObtenerProductosActivos() ([]Producto, error) {
	query := `
		SELECT IdProducto, Nombre, PrecioUnitario, EstaActivo
		FROM Productos
		WHERE EstaActivo = true
		ORDER BY IdProducto
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []Producto
	for rows.Next() {
		var p Producto
		err := rows.Scan(&p.IdProducto, &p.Nombre, &p.PrecioUnitario, &p.EstaActivo)
		if err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}

	return productos, nil
}

// ObtenerConsumosSemana retorna los consumos de una semana específica
func (r *Repositorio) ObtenerConsumosSemana(fechaInicio, fechaFin time.Time) ([]Consumo, error) {
	query := `
		SELECT IdConsumo, IdEstudiante, IdProducto, Cantidad,
		       PrecioUnitarioVenta, TotalLinea, FechaConsumo
		FROM Consumos
		WHERE FechaConsumo BETWEEN $1 AND $2
		ORDER BY FechaConsumo, IdEstudiante
	`

	rows, err := r.db.Query(query, fechaInicio, fechaFin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumos []Consumo
	for rows.Next() {
		var c Consumo
		err := rows.Scan(&c.IdConsumo, &c.IdEstudiante, &c.IdProducto, &c.Cantidad,
			&c.PrecioUnitarioVenta, &c.TotalLinea, &c.FechaConsumo)
		if err != nil {
			return nil, err
		}
		consumos = append(consumos, c)
	}

	return consumos, nil
}

// ObtenerDeudaAnterior calcula la deuda anterior de un estudiante hasta una fecha
func (r *Repositorio) ObtenerDeudaAnterior(idEstudiante int, fechaLimite time.Time) (float64, error) {
	// Total consumos antes de la fecha límite
	var totalConsumos sql.NullFloat64
	queryConsumos := `
		SELECT SUM(TotalLinea)
		FROM Consumos
		WHERE IdEstudiante = $1 AND FechaConsumo < $2
	`
	err := r.db.QueryRow(queryConsumos, idEstudiante, fechaLimite).Scan(&totalConsumos)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	// Total pagos antes de la fecha límite
	var totalPagos sql.NullFloat64
	queryPagos := `
		SELECT SUM(Monto)
		FROM Pagos
		WHERE IdEstudiante = $1 AND FechaPago < $2
	`
	err = r.db.QueryRow(queryPagos, idEstudiante, fechaLimite).Scan(&totalPagos)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	consumos := 0.0
	if totalConsumos.Valid {
		consumos = totalConsumos.Float64
	}

	pagos := 0.0
	if totalPagos.Valid {
		pagos = totalPagos.Float64
	}

	return consumos - pagos, nil
}

// ObtenerPagosSemana retorna los pagos de una semana específica
func (r *Repositorio) ObtenerPagosSemana(idEstudiante int, fechaInicio, fechaFin time.Time) (float64, error) {
	var total sql.NullFloat64
	query := `
		SELECT SUM(Monto)
		FROM Pagos
		WHERE IdEstudiante = $1 AND FechaPago BETWEEN $2 AND $3
	`
	err := r.db.QueryRow(query, idEstudiante, fechaInicio, fechaFin).Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if total.Valid {
		return total.Float64, nil
	}
	return 0, nil
}

// RegistrarConsumo inserta un nuevo consumo
func (r *Repositorio) RegistrarConsumo(consumo Consumo) error {
	query := `
		INSERT INTO Consumos (IdEstudiante, IdProducto, Cantidad, PrecioUnitarioVenta, TotalLinea, FechaConsumo)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, consumo.IdEstudiante, consumo.IdProducto, consumo.Cantidad,
		consumo.PrecioUnitarioVenta, consumo.TotalLinea, consumo.FechaConsumo)
	return err
}

// ActualizarConsumo actualiza la cantidad de un consumo existente
func (r *Repositorio) ActualizarConsumo(idEstudiante, idProducto int, fecha time.Time, cantidad int, precioUnitario float64) error {
	// Primero verificamos si existe el consumo
	var idConsumo int64
	queryBuscar := `
		SELECT IdConsumo
		FROM Consumos
		WHERE IdEstudiante = $1 AND IdProducto = $2 AND FechaConsumo = $3
		LIMIT 1
	`
	err := r.db.QueryRow(queryBuscar, idEstudiante, idProducto, fecha).Scan(&idConsumo)

	if err == sql.ErrNoRows {
		// No existe, insertamos nuevo
		if cantidad > 0 {
			return r.RegistrarConsumo(Consumo{
				IdEstudiante:        idEstudiante,
				IdProducto:          idProducto,
				Cantidad:            cantidad,
				PrecioUnitarioVenta: precioUnitario,
				TotalLinea:          float64(cantidad) * precioUnitario,
				FechaConsumo:        fecha,
			})
		}
		return nil
	} else if err != nil {
		return err
	}

	// Existe, actualizamos o eliminamos
	if cantidad <= 0 {
		// Eliminar si la cantidad es 0
		queryEliminar := `DELETE FROM Consumos WHERE IdConsumo = $1`
		_, err = r.db.Exec(queryEliminar, idConsumo)
		return err
	}

	// Actualizar
	totalLinea := float64(cantidad) * precioUnitario
	queryActualizar := `
		UPDATE Consumos
		SET Cantidad = $1, PrecioUnitarioVenta = $2, TotalLinea = $3
		WHERE IdConsumo = $4
	`
	_, err = r.db.Exec(queryActualizar, cantidad, precioUnitario, totalLinea, idConsumo)
	return err
}

// RegistrarPago inserta un nuevo pago
func (r *Repositorio) RegistrarPago(pago Pago) error {
	query := `
		INSERT INTO Pagos (IdEstudiante, Monto, FechaPago)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, pago.IdEstudiante, pago.Monto, pago.FechaPago)
	return err
}

// ObtenerPagosSemanaDetalle retorna todos los pagos de un estudiante en una semana
func (r *Repositorio) ObtenerPagosSemanaDetalle(idEstudiante int, fechaInicio, fechaFin time.Time) ([]Pago, error) {
	query := `
		SELECT IdPago, IdEstudiante, Monto, FechaPago
		FROM Pagos
		WHERE IdEstudiante = $1 AND FechaPago BETWEEN $2 AND $3
		ORDER BY FechaPago DESC
	`
	rows, err := r.db.Query(query, idEstudiante, fechaInicio, fechaFin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pagos []Pago
	for rows.Next() {
		var p Pago
		err := rows.Scan(&p.IdPago, &p.IdEstudiante, &p.Monto, &p.FechaPago)
		if err != nil {
			return nil, err
		}
		pagos = append(pagos, p)
	}

	return pagos, nil
}

// EliminarPago elimina un pago específico
func (r *Repositorio) EliminarPago(idPago int) error {
	query := `DELETE FROM Pagos WHERE IdPago = $1`
	_, err := r.db.Exec(query, idPago)
	return err
}

// ObtenerProductoPorId retorna un producto por su ID
func (r *Repositorio) ObtenerProductoPorId(idProducto int) (*Producto, error) {
	query := `
		SELECT IdProducto, Nombre, PrecioUnitario, EstaActivo
		FROM Productos
		WHERE IdProducto = $1
	`
	var p Producto
	err := r.db.QueryRow(query, idProducto).Scan(&p.IdProducto, &p.Nombre, &p.PrecioUnitario, &p.EstaActivo)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ObtenerConsumoExistente verifica si existe un consumo específico y retorna la cantidad
func (r *Repositorio) ObtenerConsumoExistente(idEstudiante, idProducto int, fecha time.Time) (int, error) {
	var cantidad int
	query := `
		SELECT Cantidad
		FROM Consumos
		WHERE IdEstudiante = $1 AND IdProducto = $2 AND FechaConsumo = $3
		LIMIT 1
	`
	err := r.db.QueryRow(query, idEstudiante, idProducto, fecha).Scan(&cantidad)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return cantidad, nil
}

// ObtenerDeudasAnterioresBatch obtiene deudas anteriores para todos los estudiantes de un grado en una sola query
func (r *Repositorio) ObtenerDeudasAnterioresBatch(idGrado int, fechaLimite time.Time) (map[int]float64, error) {
	query := `
		SELECT
			e.IdEstudiante,
			COALESCE(SUM(c.TotalLinea), 0) - COALESCE((
				SELECT SUM(p.Monto)
				FROM Pagos p
				WHERE p.IdEstudiante = e.IdEstudiante AND p.FechaPago < $1
			), 0) as deuda_anterior
		FROM Estudiantes e
		LEFT JOIN Consumos c ON e.IdEstudiante = c.IdEstudiante AND c.FechaConsumo < $2
		WHERE e.EstaActivo = true AND ($3 = 0 OR e.IdGrado = $4)
		GROUP BY e.IdEstudiante
	`

	rows, err := r.db.Query(query, fechaLimite, fechaLimite, idGrado, idGrado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deudas := make(map[int]float64, 50) // Pre-asignar para ~50 estudiantes
	for rows.Next() {
		var idEstudiante int
		var deuda float64
		if err := rows.Scan(&idEstudiante, &deuda); err != nil {
			return nil, err
		}
		deudas[idEstudiante] = deuda
	}

	return deudas, rows.Err()
}

// ObtenerPagosSemanaBatch obtiene pagos de la semana para todos los estudiantes en una sola query
func (r *Repositorio) ObtenerPagosSemanaBatch(idGrado int, fechaInicio, fechaFin time.Time) (map[int]float64, error) {
	query := `
		SELECT
			e.IdEstudiante,
			COALESCE(SUM(p.Monto), 0) as total_pagos
		FROM Estudiantes e
		LEFT JOIN Pagos p ON e.IdEstudiante = p.IdEstudiante
			AND p.FechaPago BETWEEN $1 AND $2
		WHERE e.EstaActivo = true AND ($3 = 0 OR e.IdGrado = $4)
		GROUP BY e.IdEstudiante
	`

	rows, err := r.db.Query(query, fechaInicio, fechaFin, idGrado, idGrado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pagos := make(map[int]float64, 50) // Pre-asignar para ~50 estudiantes
	for rows.Next() {
		var idEstudiante int
		var totalPagos float64
		if err := rows.Scan(&idEstudiante, &totalPagos); err != nil {
			return nil, err
		}
		pagos[idEstudiante] = totalPagos
	}

	return pagos, rows.Err()
}

