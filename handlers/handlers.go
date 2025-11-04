package handlers

import (
	"html/template"
	"kiosco/services"
	"kiosco/utils"
)

// Manejador contiene las dependencias para los handlers
type Manejador struct {
	servicio  *services.Servicio
	templates *template.Template
}

// NuevoManejador crea una nueva instancia del manejador
func NuevoManejador(servicio *services.Servicio) (*Manejador, error) {
	// Cargar templates con funciones personalizadas
	tmpl, err := template.New("").Funcs(utils.GetTemplateFuncs()).ParseGlob("templates/*.tmpl")
	if err != nil {
		return nil, err
	}

	return &Manejador{
		servicio:  servicio,
		templates: tmpl,
	}, nil
}
