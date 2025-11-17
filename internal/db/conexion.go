package db

import (
	"database/sql"
	"fmt"
	"kiosco/internal/config"
	"log"

	_ "github.com/lib/pq"
)

// ConectarBD establece la conexión con PostgreSQL
func ConectarBD(configuracion config.ConfigBD) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		configuracion.Host,
		configuracion.Puerto,
		configuracion.Usuario,
		configuracion.Contraseña,
		configuracion.NombreBD,
		configuracion.SSLMode,
	)

	bd, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión: %v", err)
	}

	// Verificar la conexión
	if err := bd.Ping(); err != nil {
		return nil, fmt.Errorf("error al conectar con la base de datos: %v", err)
	}

	// Configuración de pool de conexiones (optimizado para 1 usuario)
	bd.SetMaxOpenConns(2)    // Máximo 2 conexiones concurrentes
	bd.SetMaxIdleConns(1)    // Solo 1 conexión idle
	bd.SetConnMaxLifetime(0) // Sin límite de tiempo de vida
	bd.SetConnMaxIdleTime(0) // Mantener conexión idle indefinidamente

	log.Println("✓ Conexión exitosa a PostgreSQL")
	return bd, nil
}
