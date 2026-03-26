package repositories

import (
	"database/sql"
	"kiosco/internal/models"
	"time"
)

// ObtenerPagosSemana retorna el total de pagos de un estudiante en una semana
func (r *Repositorio) ObtenerPagosSemana(idEstudiante int, fechaInicio, fechaFin time.Time) (float64, error) {
	fechaInicioStr := fechaInicio.Format("2006-01-02")
	fechaFinStr := fechaFin.Format("2006-01-02")

	var total sql.NullFloat64
	err := r.db.QueryRow(`
		SELECT SUM(monto) FROM pagos
		WHERE id_estudiante = ? AND fecha_pago BETWEEN ? AND ?
	`, idEstudiante, fechaInicioStr, fechaFinStr).Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if total.Valid {
		return total.Float64, nil
	}
	return 0, nil
}

// RegistrarPago inserta un nuevo pago
func (r *Repositorio) RegistrarPago(pago models.Pago) error {
	fechaStr := pago.FechaPago.Format("2006-01-02")
	_, err := r.db.Exec(`
		INSERT INTO pagos (id_estudiante, monto, fecha_pago)
		VALUES (?, ?, ?)
	`, pago.IdEstudiante, pago.Monto, fechaStr)
	return err
}

// ObtenerPagosSemanaDetalle retorna todos los pagos de un estudiante en una semana
func (r *Repositorio) ObtenerPagosSemanaDetalle(idEstudiante int, fechaInicio, fechaFin time.Time) ([]models.Pago, error) {
	fechaInicioStr := fechaInicio.Format("2006-01-02")
	fechaFinStr := fechaFin.Format("2006-01-02")

	rows, err := r.db.Query(`
		SELECT id_pago, id_estudiante, monto, fecha_pago
		FROM pagos
		WHERE id_estudiante = ? AND fecha_pago BETWEEN ? AND ?
		ORDER BY fecha_pago DESC
	`, idEstudiante, fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pagos []models.Pago
	for rows.Next() {
		var p models.Pago
		if err := rows.Scan(&p.IdPago, &p.IdEstudiante, &p.Monto, &p.FechaPago); err != nil {
			return nil, err
		}
		pagos = append(pagos, p)
	}

	return pagos, rows.Err()
}

// EliminarPago elimina un pago específico
func (r *Repositorio) EliminarPago(idPago int) error {
	_, err := r.db.Exec(`DELETE FROM pagos WHERE id_pago = ?`, idPago)
	return err
}

// ObtenerPagosSemanaBatch obtiene pagos de la semana para todos los estudiantes en una sola query
func (r *Repositorio) ObtenerPagosSemanaBatch(idGrado int, fechaInicio, fechaFin time.Time) (map[int]float64, error) {
	fechaInicioStr := fechaInicio.Format("2006-01-02")
	fechaFinStr := fechaFin.Format("2006-01-02")

	rows, err := r.db.Query(`
		SELECT
			e.id_estudiante,
			COALESCE(SUM(p.monto), 0) as total_pagos
		FROM estudiantes e
		LEFT JOIN pagos p ON e.id_estudiante = p.id_estudiante
			AND p.fecha_pago BETWEEN ? AND ?
		WHERE e.esta_activo = 1 AND (? = 0 OR e.id_grado = ?)
		GROUP BY e.id_estudiante
	`, fechaInicioStr, fechaFinStr, idGrado, idGrado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pagos := make(map[int]float64)
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
