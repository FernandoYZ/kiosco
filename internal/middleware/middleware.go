package middleware

import (
	"kiosco/internal/auth"
	"log"
	"net/http"
)

// cookieInvalida borra la cookie y redirige al login
func cookieInvalida(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieNombre,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// RequiereAuth verifica la cookie firmada antes de pasar al handler.
// Redirige a /login si la cookie no existe, es inválida o está expirada.
// También inyecta token CSRF para GETs.
func RequiereAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(auth.CookieNombre)
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if _, _, ok := auth.VerificarToken(cookie.Value); !ok {
			cookieInvalida(w, r)
			return
		}

		// Inyectar token CSRF en context para GETs
		if r.Method == "GET" {
			ctx := InyectarCSRFToken(w, r)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// RequiereEdicion verifica que el usuario tenga puede_editar = 1.
// Si no tiene permisos, redirige a /registro.
// SECURITY: Validates CSRF on all POSTs (logs failures), injects token in context for both
// GET and POST (so templates can include csrf_token field in re-rendered forms from HTMX responses).
// Context key: middleware.CSRFTokenContextKey
func RequiereEdicion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(auth.CookieNombre)
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, puedeEditar, ok := auth.VerificarToken(cookie.Value)
		if !ok {
			cookieInvalida(w, r)
			return
		}

		if !puedeEditar {
			log.Printf("⚠️ Permission denied (no edit permission) from %s on %s %s", r.RemoteAddr, r.Method, r.URL.Path)
			http.Redirect(w, r, "/registro", http.StatusSeeOther)
			return
		}

		// Validar CSRF en POSTs
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				log.Printf("⚠️ ParseForm error from %s on %s %s: %v", r.RemoteAddr, r.Method, r.URL.Path, err)
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}
			if !validarCSRF(r) {
				log.Printf("⚠️ CSRF validation failed on POST from %s on %s", r.RemoteAddr, r.URL.Path)
				http.Error(w, "CSRF token invalid", http.StatusForbidden)
				return
			}
			// Inyectar token en contexto para que templates re-renderizadas puedan accederlo
			ctx := InyectarCSRFToken(w, r)
			r = r.WithContext(ctx)
		}

		// Inyectar token CSRF en context para GETs
		if r.Method == "GET" {
			ctx := InyectarCSRFToken(w, r)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// ProtegerEdicion es un helper para adaptar HandlerFunc con RequiereEdicion.
func ProtegerEdicion(h http.HandlerFunc) http.HandlerFunc {
	return RequiereEdicion(h).ServeHTTP
}

// Proteger es un helper para adaptar HandlerFunc directamente.
func Proteger(h http.HandlerFunc) http.HandlerFunc {
	return RequiereAuth(h).ServeHTTP
}

// ProtegerLogin aplica validación CSRF + rate limiting para POST /login
func ProtegerLogin(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GET: inyectar token CSRF
		if r.Method == "GET" {
			ctx := InyectarCSRFToken(w, r)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
			return
		}

		// POST: validar CSRF + rate limit
		if r.Method == "POST" {
			// Verificar rate limit ANTES de parsear formulario
			if !verificarRateLimit(r) {
				http.Redirect(w, r, "/login?error=rate_limit", http.StatusSeeOther)
				return
			}

			// Parsear form
			if err := r.ParseForm(); err != nil {
				log.Printf("⚠️ ParseForm error on POST /login from %s: %v", r.RemoteAddr, err)
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}

			// Validar CSRF
			if !validarCSRF(r) {
				http.Redirect(w, r, "/login?error=csrf", http.StatusSeeOther)
				return
			}

			// El handler determinará si login fue exitoso
			// Si falla, llamar a incrementarIntentosLogin
			// Si exitoso, llamar a resetearRateLimitLogin
			// Almacenar funciones en request para que el handler las use
			h.ServeHTTP(w, r)
			return
		}
	})
}

// IncrementarIntentosLogin es llamado por el handler de login cuando falla autenticación
func IncrementarIntentosLogin(r *http.Request) {
	incrementarIntentosLogin(r)
}

// ResetearRateLimitLogin es llamado por el handler de login cuando es exitoso
func ResetearRateLimitLogin(r *http.Request) {
	resetearRateLimitLogin(r)
}
