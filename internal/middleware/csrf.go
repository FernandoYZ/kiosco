package middleware

import (
	"context"
	"crypto/subtle"
	"kiosco/internal/auth"
	"log"
	"net/http"
)

// CSRFTokenContextKey es la clave para almacenar el token CSRF en el context.
type contextKey string

const CSRFTokenContextKey contextKey = "csrf_token"

// InyectarCSRFToken obtiene o genera un token CSRF, lo coloca en cookie y en context.
// Reutiliza el token si ya existe una cookie CSRF válida.
func InyectarCSRFToken(w http.ResponseWriter, r *http.Request) context.Context {
	// Intentar obtener token existente de la cookie
	cookie, err := r.Cookie(auth.CSRFCookieName)
	var token string

	if err == nil && cookie.Value != "" {
		// Reutilizar token existente
		token = cookie.Value
	} else {
		// Generar nuevo token si no existe cookie
		token = auth.GenerarTokenCSRF()
		log.Printf("🔐 CSRF token generated for new session from %s", r.RemoteAddr)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.CSRFCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: false, // No HttpOnly para que JS pueda acceder si es necesario
		Secure:   false, // Cambiar a true en HTTPS producción
		MaxAge:   0,     // Cookie de sesión (expira al cerrar navegador)
		SameSite: http.SameSiteLaxMode,
	})

	return context.WithValue(r.Context(), CSRFTokenContextKey, token)
}

// ValidarCSRF valida que el token en cookie coincida con el token en formulario o header.
// Llamado por RequiereEdicion para validar POSTs antes de ejecutar handlers.
// Registrar todos los fallos (cookie faltante, token faltante, mismatch).
func validarCSRF(r *http.Request) bool {
	// Obtener token de cookie
	cookieToken, err := r.Cookie(auth.CSRFCookieName)
	if err != nil {
		log.Printf("⚠️ CSRF validation failed (missing cookie) from %s on %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		return false
	}

	// Obtener token del formulario (hidden field o query param)
	formToken := r.FormValue("csrf_token")

	// Si no está en formulario, intentar obtener del header (para HTMX)
	if formToken == "" {
		formToken = r.Header.Get("X-CSRF-Token")
	}

	if formToken == "" {
		log.Printf("⚠️ CSRF validation failed (missing token) from %s on %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		return false
	}

	// Comparación en tiempo constante para evitar timing attacks
	valid := subtle.ConstantTimeCompare([]byte(cookieToken.Value), []byte(formToken)) == 1
	if !valid {
		log.Printf("⚠️ CSRF validation failed (token mismatch) from %s on %s %s", r.RemoteAddr, r.Method, r.URL.Path)
	}
	return valid
}
