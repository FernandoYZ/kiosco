package repositories

import "database/sql"

// Repositorio maneja todas las operaciones con la base de datos
type Repositorio struct {
	db *sql.DB
}

// NuevoRepositorio crea una nueva instancia del repositorio
func NuevoRepositorio(db *sql.DB) *Repositorio {
	return &Repositorio{db: db}
}
