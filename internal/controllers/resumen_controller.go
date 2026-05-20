package controllers

import (
	"kiosco/internal/models"
	"kiosco/templates/pages"
	"log"
	"net/http"
	"strings"
	"time"
)

// ResumenSector — GET /resumen/{sector}
// Vista de lectura: consumos del día agrupados por estudiante
func (m *Controlador) ResumenSector(w http.ResponseWriter, r *http.Request) {
	if !validarAuth(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sector := strings.TrimPrefix(r.URL.Path, "/resumen/")
	if sector != "menor" && sector != "mayor" {
		http.Error(w, "Sector inválido", http.StatusBadRequest)
		return
	}

	fechaStr := r.URL.Query().Get("fecha")
	if fechaStr == "" {
		fechaStr = time.Now().Format("2006-01-02")
	}

	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}

	resumenes, err := m.servicio.Repo.ObtenerResumenDiario(sector, fecha)
	if err != nil {
		log.Printf("Error al obtener resumen diario: %v", err)
		http.Error(w, "Error al cargar resumen", http.StatusInternalServerError)
		return
	}

	fechas := generarFechasSemana(fechaStr)

	// Calcular total de ítems para badge en navbar
	totalItems := 0
	for _, res := range resumenes {
		for _, item := range res.Items {
			totalItems += item.Cantidad
		}
	}

	datos := models.DatosResumenSector{
		Sector:     sector,
		Fecha:      fechaStr,
		Fechas:     fechas,
		Resumenes:  resumenes,
		TotalItems: totalItems,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := pages.ResumenSector(datos).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar resumen: %v", err)
	}
}
