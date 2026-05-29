package middleware

import "net/http"

// LimiteConcurrenciaDefault es el número máximo de requests HTTP concurrentes
// permitidos antes de retornar HTTP 503. Ajustar según capacidad del hardware.
const LimiteConcurrenciaDefault = 50

// LimitarConcurrencia crea un middleware que limita el número de conexiones HTTP
// concurrentes usando un canal buffered (semáforo). Las requests que excedan el
// límite reciben HTTP 503 Service Unavailable inmediatamente.
func LimitarConcurrencia(max int) func(http.Handler) http.Handler {
	sem := make(chan struct{}, max)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
				next.ServeHTTP(w, r)
			default:
				http.Error(w, "Servidor ocupado (límite de conexiones alcanzado)", http.StatusServiceUnavailable)
			}
		})
	}
}
