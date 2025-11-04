package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// ConfigBD contiene la configuración de la base de datos
type ConfigBD struct {
	Usuario    string
	Contraseña string
	Host       string
	Puerto     string
	NombreBD   string
}

// ObtenerConfigBD retorna la configuración desde variables de entorno o valores por defecto.
func ObtenerConfigBD() ConfigBD {
	// Leer variables de entorno
	usuario := os.Getenv("DB_USER")
	contraseña := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	puerto := os.Getenv("DB_PORT")
	nombreBD := os.Getenv("DB_NAME")

	// Usar valores por defecto si las variables de entorno no están definidas
	if usuario == "" {
		usuario = "postgres"
	}
	if contraseña == "" {
		contraseña = "081102"
	}
	if host == "" {
		host = "localhost"
	}
	if puerto == "" {
		puerto = "5432"
	}
	if nombreBD == "" {
		nombreBD = "Kiosco"
	}

	return ConfigBD{
		Usuario:    usuario,
		Contraseña: contraseña,
		Host:       host,
		Puerto:     puerto,
		NombreBD:   nombreBD,
	}
}

// ConectarBD establece la conexión con PostgreSQL
func ConectarBD(config ConfigBD) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Puerto,
		config.Usuario,
		config.Contraseña,
		config.NombreBD,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión: %v", err)
	}

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al conectar con la base de datos: %v", err)
	}

	// Configuración de pool de conexiones (optimizado para 1 usuario)
	db.SetMaxOpenConns(2)    // Máximo 2 conexiones concurrentes
	db.SetMaxIdleConns(1)    // Solo 1 conexión idle
	db.SetConnMaxLifetime(0) // Sin límite de tiempo de vida
	db.SetConnMaxIdleTime(0) // Mantener conexión idle indefinidamente

	log.Println("✓ Conexión exitosa a PostgreSQL")
	return db, nil
}