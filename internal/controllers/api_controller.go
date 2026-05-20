package controllers

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DescargarBD sirve el archivo de base de datos comprimido para descarga en la app mobile.
// Requiere autenticación válida.
func (m *Controlador) DescargarBD(w http.ResponseWriter, r *http.Request) {
	dbPath := "database/database.db"

	// Abrir el archivo
	file, err := os.Open(dbPath)
	if err != nil {
		http.Error(w, "Error al leer base de datos", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Headers para descarga
	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="database_%d.db.gz"`, time.Now().Unix()))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// Comprimir en tiempo real
	gz := gzip.NewWriter(w)
	defer gz.Close()

	// Copiar contenido comprimido
	if _, err := io.Copy(gz, file); err != nil {
		// Ya se envió parte de la respuesta, no podemos cambiar headers
		return
	}

	gz.Flush()
}
