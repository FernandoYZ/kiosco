package controllers

import (
	"encoding/json"
	"kiosco/internal/models"
	"kiosco/templates/pages"
	"log"
	"net/http"
	"time"
)

// RegistroConsumos — GET /registro
// Renderiza la página de registro con selector de sector
func (m *Controlador) RegistroConsumos(w http.ResponseWriter, r *http.Request) {
	if err := pages.RegistroConsumos().Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar registro: %v", err)
	}
}

// ObtenerEstudiantesSector — GET /registro/estudiantes?sector=menor|mayor
// Devuelve JSON con estudiantes activos del sector + productos activos
func (m *Controlador) ObtenerEstudiantesSector(w http.ResponseWriter, r *http.Request) {
	sector := r.URL.Query().Get("sector")
	if sector != "menor" && sector != "mayor" {
		http.Error(w, "Sector inválido", http.StatusBadRequest)
		return
	}

	estudiantes, err := m.servicio.Repo.ObtenerEstudiantesActivosPorSector(sector)
	if err != nil {
		log.Printf("Error al obtener estudiantes del sector %s: %v", sector, err)
		http.Error(w, "Error al cargar estudiantes", http.StatusInternalServerError)
		return
	}

	productos, err := m.servicio.Repo.ObtenerProductosActivos()
	if err != nil {
		log.Printf("Error al obtener productos: %v", err)
		http.Error(w, "Error al cargar productos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"estudiantes": estudiantes,
		"productos":   productos,
	})
}

// GuardarRegistroBatch — POST /registro/guardar
// Body JSON: { "items": [{ "id_estudiante": 1, "id_producto": 5, "cantidad": 2 }], "fecha": "2025-04-11" }
func (m *Controlador) GuardarRegistroBatch(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Items []struct {
			IdEstudiante int `json:"id_estudiante"`
			IdProducto   int `json:"id_producto"`
			Cantidad     int `json:"cantidad"`
		} `json:"items"`
		Fecha string `json:"fecha"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	fecha, err := time.Parse("2006-01-02", payload.Fecha)
	if err != nil {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}

	// Obtener precios de productos activos
	productos, err := m.servicio.Repo.ObtenerProductosActivos()
	if err != nil {
		log.Printf("Error al obtener productos: %v", err)
		http.Error(w, "Error al procesar solicitud", http.StatusInternalServerError)
		return
	}

	precioMap := make(map[int]float64)
	for _, p := range productos {
		precioMap[p.IdProducto] = p.PrecioUnitario
	}

	// Construir lista de consumos a insertar
	consumos := make([]models.Consumo, 0, len(payload.Items))
	for _, item := range payload.Items {
		if item.Cantidad <= 0 {
			continue
		}

		precio, ok := precioMap[item.IdProducto]
		if !ok {
			http.Error(w, "Producto no válido", http.StatusBadRequest)
			return
		}

		consumos = append(consumos, models.Consumo{
			IdEstudiante:        item.IdEstudiante,
			IdProducto:          item.IdProducto,
			Cantidad:            item.Cantidad,
			PrecioUnitarioVenta: precio,
			FechaConsumo:        fecha,
		})
	}

	// Registrar batch con transacción
	if err := m.servicio.Repo.RegistrarConsumosBatch(consumos); err != nil {
		log.Printf("Error al registrar consumos batch: %v", err)
		http.Error(w, "Error al guardar consumos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
