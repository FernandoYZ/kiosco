package models

// Grado representa un nivel educativo
type Grado struct {
	IdGrado     int
	NivelGrado  string // Primaria, Secundaria
	NombreGrado string // Primero, Segundo, etc.
	AnioGrado   string // 1ro, 2do, etc.
}

// InfoGrado contiene información básica de un grado
type InfoGrado struct {
	IdGrado int
	Nombre  string
}
