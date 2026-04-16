package controllers

import (
	"kiosco/internal/auth"
	"kiosco/templates/pages"
	"log"
	"net/http"
	"strings"
	"time"
)

// MostrarLogin muestra la página de login.
// Si ya hay sesión válida redirige al inicio.
func (m *Controlador) MostrarLogin(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(auth.CookieNombre); err == nil {
		if _, _, ok := auth.VerificarToken(cookie.Value); ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	if err := pages.Login("").Render(r.Context(), w); err != nil {
		log.Printf("Error renderizando login: %v", err)
	}
}

// ProcesarLogin valida credenciales contra usuarios existentes en la DB y emite la cookie de sesión.
func (m *Controlador) ProcesarLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	usuario := strings.TrimSpace(r.FormValue("usuario"))
	password := r.FormValue("password")

	u, err := m.servicio.Repo.ObtenerUsuarioPorNombre(usuario)
	if err != nil || !auth.VerificarPassword(u.Contrasenha, password) {
		// Mismo mensaje para no revelar si el usuario existe o no
		pages.Login("Usuario o contraseña incorrectos").Render(r.Context(), w)
		return
	}

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
