package repositories

import (
	"database/sql"
	"kiosco/internal/config"
)

// Repositorio maneja todas las operaciones con la base de datos
type Repositorio struct {
	db *sql.DB
}

// NuevoRepositorio crea una instancia del repositorio usando el singleton de SQLite.
func NuevoRepositorio() *Repositorio {
	return &Repositorio{db: config.DB()}
}
