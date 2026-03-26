package controllers

import "kiosco/internal/services"

// Controlador contiene las dependencias para los controladores
type Controlador struct {
	servicio *services.Servicio
}

// NuevoControlador crea una instancia del controlador con su servicio.
func NuevoControlador() (*Controlador, error) {
	return &Controlador{servicio: services.NuevoServicio()}, nil
}
