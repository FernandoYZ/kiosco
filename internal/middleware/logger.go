package middleware

import (
	"log"
	"net/http"
	"time"
)

// RegistrarRequest es un middleware que registra información de cada request HTTP
func RegistrarRequest(siguiente http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()

		// Llamar al siguiente handler
		siguiente(w, r)

		// Registrar información del request
		duracion := time.Since(inicio)
		log.Printf("[%s] %s %s - %v", r.Method, r.RemoteAddr, r.URL.Path, duracion)
	}
}

// Recuperacion es un middleware que captura panics y evita que el servidor se caiga
func Recuperacion(siguiente http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("❌ PANIC recuperado: %v", err)
				http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			}
		}()

		siguiente(w, r)
	}
}
