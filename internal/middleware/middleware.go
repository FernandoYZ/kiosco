package middleware

import (
	"kiosco/internal/auth"
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

		next.ServeHTTP(w, r)
	})
}

// RequiereEdicion verifica que el usuario tenga puede_editar = 1.
// Si no tiene permisos, redirige a /registro.
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
			http.Redirect(w, r, "/registro", http.StatusSeeOther)
			return
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
