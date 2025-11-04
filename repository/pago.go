package repository

import (
	"database/sql"
	"kiosco/models"
	"time"
)

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

// RegistrarPago inserta un nuevo pago
func (r *Repositorio) RegistrarPago(pago models.Pago) error {
	query := `
		INSERT INTO Pagos (IdEstudiante, Monto, FechaPago)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, pago.IdEstudiante, pago.Monto, pago.FechaPago)
	return err
}

// ObtenerPagosSemanaDetalle retorna todos los pagos de un estudiante en una semana
func (r *Repositorio) ObtenerPagosSemanaDetalle(idEstudiante int, fechaInicio, fechaFin time.Time) ([]models.Pago, error) {
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

	var pagos []models.Pago
	for rows.Next() {
		var p models.Pago
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
