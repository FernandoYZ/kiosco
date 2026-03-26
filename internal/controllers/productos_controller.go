package controllers

import (
	"kiosco/templates/pages"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// SetupProductos muestra la página de gestión de productos
func (m *Controlador) SetupProductos(w http.ResponseWriter, r *http.Request) {
	productos, err := m.servicio.Repo.ObtenerTodosProductos()
	if err != nil {
		log.Printf("Error al obtener productos: %v", err)
		productos = nil
	}

	if err := pages.SetupProductos(productos).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar setup productos: %v", err)
	}
}

// AgregarProducto inserta un producto vía formulario y responde con fragmento HTMX
func (m *Controlador) AgregarProducto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(r.FormValue("nombre"))
	precio, err := strconv.ParseFloat(r.FormValue("precio_unitario"), 64)
	if err != nil || nombre == "" || precio <= 0 {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	prod, err := m.servicio.Repo.InsertarProducto(nombre, precio)
	if err != nil {
		log.Printf("Error al insertar producto: %v", err)
		http.Error(w, "Error al agregar producto", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		if err := pages.FilaProducto(prod).Render(r.Context(), w); err != nil {
			log.Printf("Error al renderizar fila producto: %v", err)
		}
		return
	}

	http.Redirect(w, r, "/setup/productos", http.StatusSeeOther)
}

// ActualizarProducto modifica nombre y precio de un producto existente
func (m *Controlador) ActualizarProducto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	idProducto, err := strconv.Atoi(r.FormValue("id_producto"))
	nombre := strings.TrimSpace(r.FormValue("nombre"))
	precio, errPrecio := strconv.ParseFloat(r.FormValue("precio_unitario"), 64)
	if err != nil || errPrecio != nil || nombre == "" || precio <= 0 {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if err := m.servicio.Repo.ActualizarProducto(idProducto, nombre, precio); err != nil {
		log.Printf("Error al actualizar producto %d: %v", idProducto, err)
		http.Error(w, "Error al actualizar producto", http.StatusInternalServerError)
		return
	}

	prod, err := m.servicio.Repo.ObtenerProductoPorId(idProducto)
	if err != nil {
		log.Printf("Error al obtener producto %d: %v", idProducto, err)
		http.Error(w, "Error al obtener producto", http.StatusInternalServerError)
		return
	}

	if err := pages.FilaProducto(*prod).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar fila producto: %v", err)
	}
}

// ToggleProducto habilita o deshabilita un producto
func (m *Controlador) ToggleProducto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	idProducto, err := strconv.Atoi(r.FormValue("id_producto"))
	if err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	activo := r.FormValue("esta_activo") == "1"
	if err := m.servicio.Repo.CambiarEstadoProducto(idProducto, activo); err != nil {
		log.Printf("Error al cambiar estado de producto %d: %v", idProducto, err)
		http.Error(w, "Error al cambiar estado", http.StatusInternalServerError)
		return
	}

	prod, err := m.servicio.Repo.ObtenerProductoPorId(idProducto)
	if err != nil {
		log.Printf("Error al obtener producto %d: %v", idProducto, err)
		http.Error(w, "Error al obtener producto", http.StatusInternalServerError)
		return
	}

	if err := pages.FilaProducto(*prod).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar fila producto: %v", err)
	}
}
