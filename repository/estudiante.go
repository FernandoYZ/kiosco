package repository

import (
	"kiosco/models"
	"time"
)

// ObtenerEstudiantesActivos retorna todos los estudiantes activos
func (r *Repositorio) ObtenerEstudiantesActivos() ([]models.Estudiante, error) {
	return r.ObtenerEstudiantesPorGrado(0)
}

// ObtenerEstudiantesPorGrado retorna los estudiantes activos filtrados por grado (0 = todos)
func (r *Repositorio) ObtenerEstudiantesPorGrado(idGrado int) ([]models.Estudiante, error) {
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

	var estudiantes []models.Estudiante
	for rows.Next() {
		var e models.Estudiante
		err := rows.Scan(&e.IdEstudiante, &e.Nombres, &e.Apellidos, &e.IdGrado, &e.EstaActivo, &e.NombreGrado)
		if err != nil {
			return nil, err
		}
		estudiantes = append(estudiantes, e)
	}

	return estudiantes, nil
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
