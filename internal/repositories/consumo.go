package repositories

import (
	"database/sql"
	"kiosco/internal/models"
	"time"
)

// ObtenerConsumosSemana retorna los consumos de una semana específica
func (r *Repositorio) ObtenerConsumosSemana(fechaInicio, fechaFin time.Time) ([]models.Consumo, error) {
	fechaInicioStr := fechaInicio.Format("2006-01-02")
	fechaFinStr := fechaFin.Format("2006-01-02")

	query := `
		SELECT id_consumo, id_estudiante, id_producto, cantidad,
		       precio_unitario_venta, total_linea, fecha_consumo
		FROM consumos
		WHERE fecha_consumo BETWEEN ? AND ?
		ORDER BY fecha_consumo, id_estudiante
	`

	rows, err := r.db.Query(query, fechaInicioStr, fechaFinStr)
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

	return consumos, rows.Err()
}

// ObtenerDeudaAnterior calcula la deuda anterior de un estudiante hasta una fecha
func (r *Repositorio) ObtenerDeudaAnterior(idEstudiante int, fechaLimite time.Time) (float64, error) {
	fechaLimiteStr := fechaLimite.Format("2006-01-02")

	var totalConsumos sql.NullFloat64
	err := r.db.QueryRow(`
		SELECT SUM(total_linea) FROM consumos
		WHERE id_estudiante = ? AND fecha_consumo < ?
	`, idEstudiante, fechaLimiteStr).Scan(&totalConsumos)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	var totalPagos sql.NullFloat64
	err = r.db.QueryRow(`
		SELECT SUM(monto) FROM pagos
		WHERE id_estudiante = ? AND fecha_pago < ?
	`, idEstudiante, fechaLimiteStr).Scan(&totalPagos)
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

// RegistrarConsumo inserta un nuevo consumo (total_linea es GENERATED, no se inserta)
func (r *Repositorio) RegistrarConsumo(consumo models.Consumo) error {
	fechaStr := consumo.FechaConsumo.Format("2006-01-02")
	_, err := r.db.Exec(`
		INSERT INTO consumos (id_estudiante, id_producto, cantidad, precio_unitario_venta, fecha_consumo)
		VALUES (?, ?, ?, ?, ?)
	`, consumo.IdEstudiante, consumo.IdProducto, consumo.Cantidad,
		consumo.PrecioUnitarioVenta, fechaStr)
	return err
}

// ActualizarConsumo actualiza, inserta o elimina un consumo según la cantidad
func (r *Repositorio) ActualizarConsumo(idEstudiante, idProducto int, fecha time.Time, cantidad int, precioUnitario float64) error {
	fechaStr := fecha.Format("2006-01-02")

	var idConsumo int64
	err := r.db.QueryRow(`
		SELECT id_consumo FROM consumos
		WHERE id_estudiante = ? AND id_producto = ? AND fecha_consumo = ?
		LIMIT 1
	`, idEstudiante, idProducto, fechaStr).Scan(&idConsumo)

	if err == sql.ErrNoRows {
		if cantidad > 0 {
			return r.RegistrarConsumo(models.Consumo{
				IdEstudiante:        idEstudiante,
				IdProducto:          idProducto,
				Cantidad:            cantidad,
				PrecioUnitarioVenta: precioUnitario,
				FechaConsumo:        fecha,
			})
		}
		return nil
	} else if err != nil {
		return err
	}

	if cantidad <= 0 {
		_, err = r.db.Exec(`DELETE FROM consumos WHERE id_consumo = ?`, idConsumo)
		return err
	}

	// total_linea es GENERATED, solo actualizamos cantidad y precio
	_, err = r.db.Exec(`
		UPDATE consumos SET cantidad = ?, precio_unitario_venta = ?
		WHERE id_consumo = ?
	`, cantidad, precioUnitario, idConsumo)
	return err
}

// ObtenerConsumoExistente verifica si existe un consumo y retorna la cantidad
func (r *Repositorio) ObtenerConsumoExistente(idEstudiante, idProducto int, fecha time.Time) (int, error) {
	fechaStr := fecha.Format("2006-01-02")

	var cantidad int
	err := r.db.QueryRow(`
		SELECT cantidad FROM consumos
		WHERE id_estudiante = ? AND id_producto = ? AND fecha_consumo = ?
		LIMIT 1
	`, idEstudiante, idProducto, fechaStr).Scan(&cantidad)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return cantidad, err
}
