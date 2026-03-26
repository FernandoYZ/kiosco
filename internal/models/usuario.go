package models

// Usuario representa un usuario del sistema
type Usuario struct {
	IdUsuario   int
	Usuario     string
	Contrasenha string
	PuedeEditar bool
}
