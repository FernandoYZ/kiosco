package controllers

import (
	"kiosco/internal/utils"
	"kiosco/templates/pages"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// SetupEstudiantes muestra la página de configuración inicial de estudiantes
func (m *Controlador) SetupEstudiantes(w http.ResponseWriter, r *http.Request) {
	grados := utils.ObtenerGradosEstaticos()

	estudiantes, err := m.servicio.Repo.ObtenerTodosEstudiantes()
	if err != nil {
		log.Printf("Error al obtener estudiantes: %v", err)
		estudiantes = nil
	}

	if err := pages.SetupEstudiantes(grados, estudiantes).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar setup: %v", err)
	}
}

// AgregarEstudiante agrega un estudiante vía formulario y responde con fragmento HTMX
func (m *Controlador) AgregarEstudiante(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	apellidos := strings.TrimSpace(r.FormValue("apellidos"))
	nombres := strings.TrimSpace(r.FormValue("nombres"))
	idGrado, err := strconv.Atoi(r.FormValue("id_grado"))
	if err != nil || apellidos == "" || nombres == "" {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	est, err := m.servicio.Repo.InsertarEstudiante(nombres, apellidos, idGrado)
	if err != nil {
		log.Printf("Error al insertar estudiante: %v", err)
		http.Error(w, "Error al agregar estudiante", http.StatusInternalServerError)
		return
	}

	est.NombreGrado = utils.NombreGrado(idGrado)

	grados := utils.ObtenerGradosEstaticos()

	if r.Header.Get("HX-Request") == "true" {
		if err := pages.FilaEstudiante(est, grados).Render(r.Context(), w); err != nil {
			log.Printf("Error al renderizar fila estudiante: %v", err)
		}
		return
	}

	http.Redirect(w, r, "/setup", http.StatusSeeOther)
}

// ActualizarEstudiante modifica los datos de un estudiante existente
func (m *Controlador) ActualizarEstudiante(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	idEstudiante, err := strconv.Atoi(r.FormValue("id_estudiante"))
	apellidos := strings.TrimSpace(r.FormValue("apellidos"))
	nombres := strings.TrimSpace(r.FormValue("nombres"))
	idGrado, errGrado := strconv.Atoi(r.FormValue("id_grado"))
	if err != nil || errGrado != nil || apellidos == "" || nombres == "" {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if err := m.servicio.Repo.ActualizarEstudiante(idEstudiante, nombres, apellidos, idGrado); err != nil {
		log.Printf("Error al actualizar estudiante %d: %v", idEstudiante, err)
		http.Error(w, "Error al actualizar estudiante", http.StatusInternalServerError)
		return
	}

	est, err := m.servicio.Repo.ObtenerEstudiantePorId(idEstudiante)
	if err != nil {
		log.Printf("Error al obtener estudiante %d: %v", idEstudiante, err)
		http.Error(w, "Error al obtener estudiante", http.StatusInternalServerError)
		return
	}

	grados := utils.ObtenerGradosEstaticos()

	if err := pages.FilaEstudiante(est, grados).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar fila estudiante: %v", err)
	}
}

// ToggleEstudiante habilita o deshabilita un estudiante
func (m *Controlador) ToggleEstudiante(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	idEstudiante, err := strconv.Atoi(r.FormValue("id_estudiante"))
	estaActivoVal := r.FormValue("esta_activo")
	if err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	activo := estaActivoVal == "1"
	if err := m.servicio.Repo.CambiarEstadoEstudiante(idEstudiante, activo); err != nil {
		log.Printf("Error al cambiar estado de estudiante %d: %v", idEstudiante, err)
		http.Error(w, "Error al cambiar estado", http.StatusInternalServerError)
		return
	}

	est, err := m.servicio.Repo.ObtenerEstudiantePorId(idEstudiante)
	if err != nil {
		log.Printf("Error al obtener estudiante %d: %v", idEstudiante, err)
		http.Error(w, "Error al obtener estudiante", http.StatusInternalServerError)
		return
	}

	grados := utils.ObtenerGradosEstaticos()

	if err := pages.FilaEstudiante(est, grados).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar fila estudiante: %v", err)
	}
}
