package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/argon2"
)

const (
	CookieNombre  = "kiosco_session"
	tiempoExpiry  = 24 * time.Hour
	argonMemory   = 64 * 1024 // 64 MB
	argonIter     = 1
	argonThreads  = 4
	argonKeyLen   = 32
	argonSaltLen  = 16
)

var (
	llaveSecreta []byte
	unaVez       sync.Once
)

// LlaveEfimera retorna la llave HMAC generada en memoria al iniciar el servidor.
// Se llama una sola vez; las llamadas posteriores devuelven la misma llave.
func LlaveEfimera() []byte {
	unaVez.Do(func() {
		llaveSecreta = make([]byte, 32)
		if _, err := rand.Read(llaveSecreta); err != nil {
			panic("auth: no se pudo generar llave efímera: " + err.Error())
		}
	})
	return llaveSecreta
}

// FirmarToken genera un token firmado con HMAC-SHA256.
// Formato: <base64url(idUsuario:expiry)>.<base64url(hmac)>
func FirmarToken(idUsuario int) string {
	expiry := time.Now().Add(tiempoExpiry).Unix()
	payload := fmt.Sprintf("%d:%d", idUsuario, expiry)
	b64 := base64.RawURLEncoding.EncodeToString([]byte(payload))

	mac := hmac.New(sha256.New, LlaveEfimera())
	mac.Write([]byte(b64))
	firma := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return b64 + "." + firma
}

// VerificarToken valida la firma y la expiración del token.
// Devuelve el idUsuario y true si es válido.
func VerificarToken(token string) (int, bool) {
	partes := strings.SplitN(token, ".", 2)
	if len(partes) != 2 {
		return 0, false
	}
	b64, firmaRecibida := partes[0], partes[1]

	// Verificar HMAC (tiempo constante)
	mac := hmac.New(sha256.New, LlaveEfimera())
	mac.Write([]byte(b64))
	firmaEsperada := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(firmaRecibida), []byte(firmaEsperada)) {
		return 0, false
	}

	// Decodificar payload
	raw, err := base64.RawURLEncoding.DecodeString(b64)
	if err != nil {
		return 0, false
	}
	campos := strings.SplitN(string(raw), ":", 2)
	if len(campos) != 2 {
		return 0, false
	}

	idUsuario, err := strconv.Atoi(campos[0])
	if err != nil {
		return 0, false
	}

	expiry, err := strconv.ParseInt(campos[1], 10, 64)
	if err != nil || time.Now().Unix() > expiry {
		return 0, false
	}

	return idUsuario, true
}

// HashPassword genera un hash Argon2id del password.
// Formato PHC: $argon2id$v=19$m=...,t=...,p=...$<salt>$<hash>
func HashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonIter, argonMemory, argonThreads, argonKeyLen)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory, argonIter, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

// VerificarPassword compara un password contra un hash Argon2id almacenado.
func VerificarPassword(hashAlmacenado, password string) bool {
	partes := strings.Split(hashAlmacenado, "$")
	// ["", "argon2id", "v=19", "m=...,t=...,p=...", "<salt>", "<hash>"]
	if len(partes) != 6 || partes[1] != "argon2id" {
		return false
	}

	var memory uint32
	var iter uint32
	var threads uint8
	if _, err := fmt.Sscanf(partes[3], "m=%d,t=%d,p=%d", &memory, &iter, &threads); err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(partes[4])
	if err != nil {
		return false
	}
	hashEsperado, err := base64.RawStdEncoding.DecodeString(partes[5])
	if err != nil {
		return false
	}

	hashComputado := argon2.IDKey([]byte(password), salt, iter, memory, threads, uint32(len(hashEsperado)))
	return hmac.Equal(hashComputado, hashEsperado)
}
