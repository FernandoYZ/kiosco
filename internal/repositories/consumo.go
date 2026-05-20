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

// ActualizarConsumo actualiza, inserta o elimina un consumo según la cantidad.
// UPSERT: safe for idempotent resubmission — SELECT → INSERT (qty>0) | UPDATE (row exists, qty>0) | DELETE (qty<=0) | noop (no row, qty<=0).
// Two identical submissions always produce exactly 1 row; qty=0 deletes the row.
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

// ObtenerResumenDiario retorna los consumos del día agrupados por estudiante para un sector.
// sector: "menor" (id_grado IN 1,2,3,4) | "mayor" (id_grado IN 5,6,7)
func (r *Repositorio) ObtenerResumenDiario(sector string, fecha time.Time) ([]models.ResumenEstudiante, error) {
	var gradoList string
	if sector == "menor" {
		gradoList = "1,2,3,4"
	} else {
		gradoList = "5,6,7"
	}

	query := `
		SELECT
			e.id_estudiante,
			e.nombres,
			e.apellidos,
			g.anio_grado || ' ' || g.nivel_grado AS nombre_grado,
			p.nombre AS nombre_producto,
			c.cantidad
		FROM consumos c
		JOIN estudiantes e ON c.id_estudiante = e.id_estudiante
		JOIN grados g ON e.id_grado = g.id_grado
		JOIN productos p ON c.id_producto = p.id_producto
		WHERE c.fecha_consumo = ?
		  AND e.id_grado IN (` + gradoList + `)
		ORDER BY e.apellidos, e.nombres, p.nombre
	`

	rows, err := r.db.Query(query, fecha.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Agrupación en memoria: slice para preservar orden de inserción
	var orden []int
	index := make(map[int]*models.ResumenEstudiante)

	for rows.Next() {
		var idEst int
		var nombres, apellidos, nombreGrado, nombreProducto string
		var cantidad int

		if err := rows.Scan(&idEst, &nombres, &apellidos, &nombreGrado, &nombreProducto, &cantidad); err != nil {
			return nil, err
		}

		if _, ok := index[idEst]; !ok {
			index[idEst] = &models.ResumenEstudiante{
				IdEstudiante: idEst,
				Nombres:      nombres,
				Apellidos:    apellidos,
				NombreGrado:  nombreGrado,
			}
			orden = append(orden, idEst)
		}

		index[idEst].Items = append(index[idEst].Items, models.ItemConsumo{
			NombreProducto: nombreProducto,
			Cantidad:       cantidad,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	resumenes := make([]models.ResumenEstudiante, 0, len(orden))
	for _, id := range orden {
		resumenes = append(resumenes, *index[id])
	}
	return resumenes, nil
}

// RegistrarConsumosBatch inserta múltiples consumos en una transacción atómica
func (r *Repositorio) RegistrarConsumosBatch(consumos []models.Consumo) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO consumos (id_estudiante, id_producto, cantidad, precio_unitario_venta, fecha_consumo)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range consumos {
		fechaStr := c.FechaConsumo.Format("2006-01-02")
		if _, err := stmt.Exec(c.IdEstudiante, c.IdProducto, c.Cantidad, c.PrecioUnitarioVenta, fechaStr); err != nil {
			return err
		}
	}

	return tx.Commit()
}
