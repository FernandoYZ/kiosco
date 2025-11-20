package config

import "os"

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

	// Usar valores por defecto SOLO para desarrollo local
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
		sslMode = "disable"
	}

	// IMPORTANTE: La contraseña DEBE venir de variable de entorno

	return ConfigBD{
		Usuario:    usuario,
		Contraseña: contraseña,
		Host:       host,
		Puerto:     puerto,
		NombreBD:   nombreBD,
		SSLMode:    sslMode,
	}
}

// ObtenerDireccion retorna la dirección (host:puerto) para el servidor web.
// Lee las variables de entorno HOST y PORT, con valores por defecto para desarrollo.
func ObtenerDireccion() string {
	host := os.Getenv("HOST")
	if host == "" {
		// Para desarrollo, "localhost" es más seguro.
		// Para producción en un VPS (especialmente con Docker), deberías establecer HOST="0.0.0.0"
		host = "localhost"
	}

	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "3200"
	}

	return host + ":" + puerto
}
