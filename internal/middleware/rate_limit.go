package middleware

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// DatosRateLimit almacena intentos fallidos y timestamp del último intento.
type DatosRateLimit struct {
	Attempts    int
	UltimoIntento time.Time
}

var (
	mapaRateLimit = make(map[string]DatosRateLimit)
	muRateLimit  sync.Mutex

	// Configuración de rate limiting
	MaxIntentosLogin = 5
	VentanaDuracion   = 15 * time.Minute
	IntervaloSweep   = 60 * time.Second
)

// limpiarExpirados elimina entradas expiradas del mapa de rate limit.
// Debe llamarse con muRateLimit tomado.
func limpiarExpirados(ahora time.Time, m map[string]DatosRateLimit) {
	for ip, data := range m {
		if ahora.Sub(data.UltimoIntento) > VentanaDuracion {
			delete(m, ip)
		}
	}
}

// IniciarSweeper arranca la goroutine de limpieza activa del mapaRateLimit.
// Debe llamarse una vez desde main, pasando el context raíz para shutdown limpio.
func IniciarSweeper(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(IntervaloSweep)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				muRateLimit.Lock()
				limpiarExpirados(time.Now(), mapaRateLimit)
				muRateLimit.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// incrementarIntentosLogin incrementa el contador de intentos fallidos para una IP.
func incrementarIntentosLogin(r *http.Request) {
	ip := obtenerIP(r)

	muRateLimit.Lock()
	defer muRateLimit.Unlock()

	ahora := time.Now()

	// Incrementar intento para esta IP
	if data, exists := mapaRateLimit[ip]; exists {
		data.Attempts++
		data.UltimoIntento = ahora
		mapaRateLimit[ip] = data
	} else {
		mapaRateLimit[ip] = DatosRateLimit{Attempts: 1, UltimoIntento: ahora}
	}
}

// verificarRateLimit verifica si la IP ha excedido el límite de intentos.
func verificarRateLimit(r *http.Request) bool {
	ip := obtenerIP(r)

	muRateLimit.Lock()
	defer muRateLimit.Unlock()

	ahora := time.Now()

	// Cleanup pasivo — única llamada por request path
	limpiarExpirados(ahora, mapaRateLimit)

	// Verificar si está bloqueada
	if data, exists := mapaRateLimit[ip]; exists {
		if ahora.Sub(data.UltimoIntento) <= VentanaDuracion && data.Attempts >= MaxIntentosLogin {
			log.Printf("⚠️ Rate limit exceeded for IP %s (attempts: %d)", ip, data.Attempts)
			return false // Bloqueada
		}
	}

	return true // No bloqueada
}

// resetearRateLimitLogin resetea el contador para una IP después de login exitoso.
func resetearRateLimitLogin(r *http.Request) {
	ip := obtenerIP(r)

	muRateLimit.Lock()
	defer muRateLimit.Unlock()

	delete(mapaRateLimit, ip)
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
