package config

import "os"

// ObtenerDireccion es para usar las variables de entorno HOST y PORT por
// defecto usa 127.0.0.1:3200 para facilidad de desarrollo y despliegue local
func ObtenerDireccion() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "3200"
	}

	return host + ":" + puerto
}
