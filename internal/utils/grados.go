package utils

import (
	"encoding/json"
	"kiosco/internal/models"
)

// GradosNombres returns a slice with "Todos" as the first element followed by
// each grade name from the static grade list, in their original order.
func GradosNombres(grados []models.InfoGrado) []string {
	nombres := make([]string, 0, len(grados)+1)
	nombres = append(nombres, "Todos")
	for _, g := range grados {
		nombres = append(nombres, g.Nombre)
	}
	return nombres
}

// GradosJSON encodes a string slice to a JSON array string suitable for
// embedding in an Alpine x-data attribute. Returns "[]" on marshal failure.
func GradosJSON(grados []string) string {
	b, err := json.Marshal(grados)
	if err != nil {
		return "[]"
	}
	return string(b)
}
