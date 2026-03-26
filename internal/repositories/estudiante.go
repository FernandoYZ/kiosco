package repositories

import (
	"kiosco/internal/models"
	"time"
)

// ObtenerEstudiantesActivos retorna todos los estudiantes activos
func (r *Repositorio) ObtenerEstudiantesActivos() ([]models.Estudiante, error) {
	return r.ObtenerEstudiantesPorGrado(0)
}

// InsertarEstudiante agrega un nuevo estudiante activo
func (r *Repositorio) InsertarEstudiante(nombres, apellidos string, idGrado int) (models.Estudiante, error) {
	result, err := r.db.Exec(`
		INSERT INTO estudiantes (nombres, apellidos, id_grado, esta_activo)
		VALUES (?, ?, ?, 1)
	`, nombres, apellidos, idGrado)
	if err != nil {
		return models.Estudiante{}, err
	}

	id, _ := result.LastInsertId()
	return models.Estudiante{
		IdEstudiante: int(id),
		Nombres:      nombres,
		Apellidos:    apellidos,
		IdGrado:      idGrado,
		EstaActivo:   true,
	}, nil
}

// ObtenerEstudiantesPorGrado retorna los estudiantes activos filtrados por grado (0 = todos)
func (r *Repositorio) ObtenerEstudiantesPorGrado(idGrado int) ([]models.Estudiante, error) {
	query := `
		SELECT e.id_estudiante, e.nombres, e.apellidos, e.id_grado, e.esta_activo,
		       g.anio_grado || ' ' || g.nivel_grado as nombre_grado
		FROM estudiantes e
		LEFT JOIN grados g ON e.id_grado = g.id_grado
		WHERE e.esta_activo = 1
	`

	var args []interface{}
	if idGrado > 0 {
		query += " AND e.id_grado = ?"
		args = append(args, idGrado)
	}

	query += " ORDER BY g.nivel_grado, g.nombre_grado, e.apellidos, e.nombres"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estudiantes []models.Estudiante
	for rows.Next() {
		var e models.Estudiante
		if err := rows.Scan(&e.IdEstudiante, &e.Nombres, &e.Apellidos, &e.IdGrado, &e.EstaActivo, &e.NombreGrado); err != nil {
			return nil, err
		}
		estudiantes = append(estudiantes, e)
	}

	return estudiantes, rows.Err()
}

// ObtenerTodosEstudiantes retorna todos los estudiantes (activos e inactivos)
func (r *Repositorio) ObtenerTodosEstudiantes() ([]models.Estudiante, error) {
	rows, err := r.db.Query(`
		SELECT e.id_estudiante, e.nombres, e.apellidos, e.id_grado, e.esta_activo,
		       g.anio_grado || ' ' || g.nivel_grado as nombre_grado
		FROM estudiantes e
		LEFT JOIN grados g ON e.id_grado = g.id_grado
		ORDER BY e.esta_activo DESC, g.nivel_grado, g.nombre_grado, e.apellidos, e.nombres
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estudiantes []models.Estudiante
	for rows.Next() {
		var e models.Estudiante
		if err := rows.Scan(&e.IdEstudiante, &e.Nombres, &e.Apellidos, &e.IdGrado, &e.EstaActivo, &e.NombreGrado); err != nil {
			return nil, err
		}
		estudiantes = append(estudiantes, e)
	}
	return estudiantes, rows.Err()
}

// ObtenerEstudiantePorId retorna un estudiante por su ID
func (r *Repositorio) ObtenerEstudiantePorId(id int) (models.Estudiante, error) {
	var e models.Estudiante
	err := r.db.QueryRow(`
		SELECT e.id_estudiante, e.nombres, e.apellidos, e.id_grado, e.esta_activo,
		       g.anio_grado || ' ' || g.nivel_grado as nombre_grado
		FROM estudiantes e
		LEFT JOIN grados g ON e.id_grado = g.id_grado
		WHERE e.id_estudiante = ?
	`, id).Scan(&e.IdEstudiante, &e.Nombres, &e.Apellidos, &e.IdGrado, &e.EstaActivo, &e.NombreGrado)
	return e, err
}

// ActualizarEstudiante modifica los datos de un estudiante
func (r *Repositorio) ActualizarEstudiante(id int, nombres, apellidos string, idGrado int) error {
	_, err := r.db.Exec(`
		UPDATE estudiantes SET nombres = ?, apellidos = ?, id_grado = ?
		WHERE id_estudiante = ?
	`, nombres, apellidos, idGrado, id)
	return err
}

// CambiarEstadoEstudiante habilita o deshabilita un estudiante
func (r *Repositorio) CambiarEstadoEstudiante(id int, activo bool) error {
	estado := 0
	if activo {
		estado = 1
	}
	_, err := r.db.Exec(`
		UPDATE estudiantes SET esta_activo = ? WHERE id_estudiante = ?
	`, estado, id)
	return err
}

// ObtenerDeudasAnterioresBatch obtiene deudas anteriores para todos los estudiantes de un grado
func (r *Repositorio) ObtenerDeudasAnterioresBatch(idGrado int, fechaLimite time.Time) (map[int]float64, error) {
	fechaLimiteStr := fechaLimite.Format("2006-01-02")

	rows, err := r.db.Query(`
		SELECT
			e.id_estudiante,
			COALESCE(SUM(c.total_linea), 0) - COALESCE((
				SELECT SUM(p.monto) FROM pagos p
				WHERE p.id_estudiante = e.id_estudiante AND p.fecha_pago < ?
			), 0) as deuda_anterior
		FROM estudiantes e
		LEFT JOIN consumos c ON e.id_estudiante = c.id_estudiante AND c.fecha_consumo < ?
		WHERE e.esta_activo = 1 AND (? = 0 OR e.id_grado = ?)
		GROUP BY e.id_estudiante
	`, fechaLimiteStr, fechaLimiteStr, idGrado, idGrado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deudas := make(map[int]float64)
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
