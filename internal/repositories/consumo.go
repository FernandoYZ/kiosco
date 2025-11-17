package repositories

import (
	"database/sql"
	"kiosco/internal/models"
	"time"
)

// ObtenerConsumosSemana retorna los consumos de una semana específica
func (r *Repositorio) ObtenerConsumosSemana(fechaInicio, fechaFin time.Time) ([]models.Consumo, error) {
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

	var consumos []models.Consumo
	for rows.Next() {
		var c models.Consumo
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

// RegistrarConsumo inserta un nuevo consumo
func (r *Repositorio) RegistrarConsumo(consumo models.Consumo) error {
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
			return r.RegistrarConsumo(models.Consumo{
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
