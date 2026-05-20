package controllers

import (
	"fmt"
	"kiosco/internal/auth"
	"kiosco/internal/models"
	"kiosco/templates/pages"
	"log"
	"net/http"
	"strings"
	"time"
)

// RegistroConsumos — GET /registro
// Muestra el selector de sectores
func (m *Controlador) RegistroConsumos(w http.ResponseWriter, r *http.Request) {
	if !validarAuth(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Redirigir /registro/ a /registro (sin trailing slash)
	if r.URL.Path == "/registro/" {
		http.Redirect(w, r, "/registro", http.StatusMovedPermanently)
		return
	}

	if err := pages.RegistroConsumos().Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar selector: %v", err)
	}
}

// RegistroSector — GET /registro/{sector}
// Devuelve página completa con layout y lista de estudiantes como enlaces
func (m *Controlador) RegistroSector(w http.ResponseWriter, r *http.Request) {
	if !validarAuth(r) {
		http.Error(w, "No autenticado", http.StatusUnauthorized)
		return
	}

	sector := strings.TrimPrefix(r.URL.Path, "/registro/")
	fecha := r.URL.Query().Get("fecha")

	if sector != "menor" && sector != "mayor" {
		http.Error(w, "Sector inválido", http.StatusBadRequest)
		return
	}

	if fecha == "" {
		fecha = time.Now().Format("2006-01-02")
	}

	// Cargar estudiantes y productos
	estudiantes, err := m.servicio.Repo.ObtenerEstudiantesActivosPorSector(sector)
	if err != nil {
		log.Printf("Error al obtener estudiantes: %v", err)
		http.Error(w, "Error al cargar estudiantes", http.StatusInternalServerError)
		return
	}

	productos, err := m.servicio.Repo.ObtenerProductosActivos()
	if err != nil {
		log.Printf("Error al obtener productos: %v", err)
		http.Error(w, "Error al cargar productos", http.StatusInternalServerError)
		return
	}

	fechas := generarFechasSemana(fecha)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := pages.RegistroConsumosCon(sector, fechas, fecha, estudiantes, productos).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar página: %v", err)
	}
}

// Helpers

func validarAuth(r *http.Request) bool {
	cookie, err := r.Cookie(auth.CookieNombre)
	if err != nil {
		return false
	}
	_, _, ok := auth.VerificarToken(cookie.Value)
	return ok
}

func generarFechasSemana(fechaStr string) []models.DiaFecha {
	fecha, _ := time.Parse("2006-01-02", fechaStr)

	diaSemana := fecha.Weekday()
	var diferencia int
	if diaSemana == 0 {
		diferencia = -6
	} else {
		diferencia = 1 - int(diaSemana)
	}

	lunes := fecha.AddDate(0, 0, diferencia)
	lunes = time.Date(lunes.Year(), lunes.Month(), lunes.Day(), 0, 0, 0, 0, lunes.Location())

	nombres := []string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"}
	meses := []string{"Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dic"}

	hoy := time.Now().Format("2006-01-02")
	dias := make([]models.DiaFecha, 6)

	for i := 0; i < 6; i++ {
		d := lunes.AddDate(0, 0, i)
		fechaFormato := d.Format("2006-01-02")
		dias[i] = models.DiaFecha{
			Nombre:       nombres[i],
			Fecha:        fechaFormato,
			FechaFormato: fmt.Sprintf("%d %s", d.Day(), meses[int(d.Month())-1]),
			EsHoy:        fechaFormato == hoy,
		}
	}

	return dias
}

