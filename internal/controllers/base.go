package controllers

import (
	"html/template"
	"kiosco/internal/services"
	"kiosco/internal/utils"
	"log"
	"net/http"
)

// Controlador contiene las dependencias para los controladores
type Controlador struct {
	servicio  *services.Servicio
	templates *template.Template
}

// NuevoControlador crea una nueva instancia del controlador
func NuevoControlador(servicio *services.Servicio) (*Controlador, error) {
	// No cargar templates aquí, se cargarán bajo demanda
	return &Controlador{
		servicio:  servicio,
		templates: nil,
	}, nil
}

// cargarTemplate carga un template específico con su layout
func (c *Controlador) cargarTemplate(nombreTemplate string) (*template.Template, error) {
	tmpl := template.New(nombreTemplate).Funcs(utils.GetTemplateFuncs())

	// Cargar layout base
	tmpl, err := tmpl.ParseFiles("internal/views/layouts/base.tmpl")
	if err != nil {
		return nil, err
	}

	// Cargar el template específico
	tmpl, err = tmpl.ParseFiles("internal/views/" + nombreTemplate)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// renderizar es un helper que renderiza un template específico con los datos
func (c *Controlador) renderizar(w http.ResponseWriter, nombreTemplate string, datos interface{}) {
	tmpl, err := c.cargarTemplate(nombreTemplate)
	if err != nil {
		log.Printf("Error al cargar template: %v", err)
		http.Error(w, "Error al cargar la página", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", datos)
	if err != nil {
		log.Printf("Error al renderizar template: %v", err)
		http.Error(w, "Error al renderizar la página", http.StatusInternalServerError)
	}
}
