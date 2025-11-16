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
	SSLMode    string
}

// ObtenerConfigBD retorna la configuración desde variables de entorno o valores por defecto.
func ObtenerConfigBD() ConfigBD {
	// Leer variables de entorno
	usuario := os.Getenv("DB_USER")
	contraseña := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	puerto := os.Getenv("DB_PORT")
	nombreBD := os.Getenv("DB_NAME")
	sslMode := os.Getenv("PGSSLMODE")

	// Usar valores por defecto SOLO para desarrollo local (sin credenciales sensibles)
	if usuario == "" {
		usuario = "postgres"
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
	if sslMode == "" {
		sslMode = "disable" // Valor por defecto para desarrollo local
	}

	// IMPORTANTE: La contraseña DEBE venir de variable de entorno
	// No usar valores por defecto para credenciales sensibles

	return ConfigBD{
		Usuario:    usuario,
		Contraseña: contraseña,
		Host:       host,
		Puerto:     puerto,
		NombreBD:   nombreBD,
		SSLMode:    sslMode,
	}
}

// ConectarBD establece la conexión con PostgreSQL
func ConectarBD(config ConfigBD) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Puerto,
		config.Usuario,
		config.Contraseña,
		config.NombreBD,
		config.SSLMode,
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

// ObtenerPuerto retorna la dirección del servidor desde variable de entorno o usa el valor por defecto
func ObtenerPuerto() string {
	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "3200"
	}
	// En producción (con PORT definido), escuchar en todas las interfaces (0.0.0.0)
	// En desarrollo (sin PORT), escuchar en localhost
	if os.Getenv("PORT") != "" {
		return "0.0.0.0:" + puerto
	}
	return ":" + puerto
}
