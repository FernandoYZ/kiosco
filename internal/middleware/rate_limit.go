package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimitData almacena intentos fallidos y timestamp del último intento.
type RateLimitData struct {
	Attempts    int
	LastAttempt time.Time
}

var (
	rateLimitMap = make(map[string]RateLimitData)
	rateLimitMu  sync.Mutex

	// Configuración de rate limiting
	MaxAttemptsLogin = 5
	WindowDuration   = 15 * time.Minute
)

// limpiarExpirados elimina entradas expiradas del mapa de rate limit.
// Debe llamarse con rateLimitMu tomado.
func limpiarExpirados(ahora time.Time, m map[string]RateLimitData) {
	for ip, data := range m {
		if ahora.Sub(data.LastAttempt) > WindowDuration {
			delete(m, ip)
		}
	}
}

// incrementarIntentosLogin incrementa el contador de intentos fallidos para una IP.
func incrementarIntentosLogin(r *http.Request) {
	ip := obtenerIP(r)

	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()

	ahora := time.Now()

	// Incrementar intento para esta IP
	if data, exists := rateLimitMap[ip]; exists {
		data.Attempts++
		data.LastAttempt = ahora
		rateLimitMap[ip] = data
	} else {
		rateLimitMap[ip] = RateLimitData{Attempts: 1, LastAttempt: ahora}
	}
}

// verificarRateLimit verifica si la IP ha excedido el límite de intentos.
func verificarRateLimit(r *http.Request) bool {
	ip := obtenerIP(r)

	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()

	ahora := time.Now()

	// Cleanup pasivo — única llamada por request path
	limpiarExpirados(ahora, rateLimitMap)

	// Verificar si está bloqueada
	if data, exists := rateLimitMap[ip]; exists {
		if ahora.Sub(data.LastAttempt) <= WindowDuration && data.Attempts >= MaxAttemptsLogin {
			return false // Bloqueada
		}
	}

	return true // No bloqueada
}

// resetearRateLimitLogin resetea el contador para una IP después de login exitoso.
func resetearRateLimitLogin(r *http.Request) {
	ip := obtenerIP(r)

	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()

	delete(rateLimitMap, ip)
}

// obtenerIP extrae la dirección IP del cliente del request.
func obtenerIP(r *http.Request) string {
	// Intentar obtener de X-Forwarded-For (detrás de proxy)
	// Tomar solo la primera IP (más a la izquierda) y descartar posibles encabezados añadidos por proxies intermedios
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.TrimSpace(strings.SplitN(forwarded, ",", 2)[0])
	}

	// Obtener de RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
