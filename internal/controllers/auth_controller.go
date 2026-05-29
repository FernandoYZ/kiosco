package controllers

import (
	"kiosco/internal/auth"
	"kiosco/internal/middleware"
	"kiosco/templates/components"
	"kiosco/templates/pages"
	"log"
	"net/http"
	"strings"
	"time"
)

// MostrarLogin muestra la página de login.
// Si ya hay sesión válida redirige al inicio.
// También maneja errores de seguridad (CSRF, rate limit).
func (m *Controlador) MostrarLogin(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(auth.CookieNombre); err == nil {
		if _, _, ok := auth.VerificarToken(cookie.Value); ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	// Verificar si hay error de seguridad
	errorType := r.URL.Query().Get("error")
	if errorType == "rate_limit" {
		boton1 := components.BotonError("Volver a intentar", "/login", true)
		boton2 := components.BotonError("Ir al inicio", "/", false)
		if err := components.ErrorPage(
			"Demasiados intentos",
			"Has realizado demasiados intentos de login. Intenta de nuevo en 15 minutos.",
			boton1, boton2,
		).Render(r.Context(), w); err != nil {
			log.Printf("Error renderizando error page: %v", err)
		}
		return
	}

	if errorType == "csrf" {
		boton1 := components.BotonError("Recargar", "/login", true)
		boton2 := components.BotonError("Ir al inicio", "/", false)
		if err := components.ErrorPage(
			"Token de seguridad expirado",
			"El token de seguridad ha expirado. Recarga la página e intenta de nuevo.",
			boton1, boton2,
		).Render(r.Context(), w); err != nil {
			log.Printf("Error renderizando error page: %v", err)
		}
		return
	}

	// Mostrar login normal
	// El token CSRF ya viene inyectado por ProtegerLogin
	if err := pages.Login("").Render(r.Context(), w); err != nil {
		log.Printf("Error renderizando login: %v", err)
	}
}

// ProcesarLogin valida credenciales contra usuarios existentes en la DB y emite la cookie de sesión.
func (m *Controlador) ProcesarLogin(w http.ResponseWriter, r *http.Request) {
	usuario := strings.TrimSpace(r.FormValue("usuario"))
	password := r.FormValue("password")

	u, err := m.servicio.Repo.ObtenerUsuarioPorNombre(usuario)
	if err != nil || !auth.VerificarPassword(u.Contrasenha, password) {
		// Incrementar contador de intentos fallidos
		middleware.IncrementarIntentosLogin(r)
		// Inyectar token CSRF para re-renderizar formulario
		ctx := middleware.InyectarCSRFToken(w, r)
		// Mismo mensaje para no revelar si el usuario existe o no
		if err := pages.Login("Usuario o contraseña incorrectos").Render(ctx, w); err != nil {
			log.Printf("Error renderizando login: %v", err)
		}
		return
	}

	// Reset rate limit al login exitoso
	middleware.ResetearRateLimitLogin(r)

	token := auth.FirmarToken(u.IdUsuario, u.PuedeEditar)
	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieNombre,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // cambiar a true cuando se use HTTPS en producción
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(24 * time.Hour / time.Second),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout borra la cookie y redirige al login.
func (m *Controlador) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieNombre,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
