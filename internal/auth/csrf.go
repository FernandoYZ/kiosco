package auth

import (
	"crypto/rand"
	"encoding/base64"
)

const CSRFCookieName = "kiosco_csrf"

// GenerarTokenCSRF genera un token CSRF aleatorio de 32 bytes.
func GenerarTokenCSRF() string {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		panic("auth: no se pudo generar token CSRF: " + err.Error())
	}
	return base64.RawURLEncoding.EncodeToString(token)
}
