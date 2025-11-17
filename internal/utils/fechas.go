package utils

import (
	"fmt"
	"kiosco/internal/models"
	"time"
)

// ObtenerGradosEstaticos retorna la lista estática de grados (sin "Todos")
func ObtenerGradosEstaticos() []models.InfoGrado {
	return []models.InfoGrado{
		{IdGrado: 1, Nombre: "5to Primaria"},
		{IdGrado: 2, Nombre: "6to Primaria"},
		{IdGrado: 3, Nombre: "1ro Secundaria"},
		{IdGrado: 4, Nombre: "2do Secundaria"},
		{IdGrado: 5, Nombre: "3ro Secundaria"},
		{IdGrado: 6, Nombre: "4to Secundaria"},
		{IdGrado: 7, Nombre: "5to Secundaria"},
	}
}

// ObtenerSemanaActual retorna el lunes y viernes de la semana actual
func ObtenerSemanaActual() (time.Time, time.Time) {
	ahora := time.Now()

	// Obtener el día de la semana (0 = Domingo, 1 = Lunes, ...)
	diaSemana := int(ahora.Weekday())

	// Ajustar para que Lunes sea 0
	if diaSemana == 0 {
		diaSemana = 7
	}
	diaSemana-- // Ahora Lunes = 0, Martes = 1, ..., Domingo = 6

	// Calcular el lunes de esta semana
	lunes := ahora.AddDate(0, 0, -diaSemana)
	lunes = time.Date(lunes.Year(), lunes.Month(), lunes.Day(), 0, 0, 0, 0, lunes.Location())

	// Calcular el sábado
	sabado := lunes.AddDate(0, 0, 5)
	sabado = time.Date(sabado.Year(), sabado.Month(), sabado.Day(), 23, 59, 59, 0, sabado.Location())

	return lunes, sabado
}

// ObtenerDiasHabiles retorna los días de lunes a sábado
func ObtenerDiasHabiles(inicio time.Time) []time.Time {
	dias := make([]time.Time, 6)
	for i := 0; i < 6; i++ {
		dias[i] = inicio.AddDate(0, 0, i)
	}
	return dias
}

// FormatearSemana retorna el texto de la semana (ej: "27 AL 31 DE OCTUBRE")
func FormatearSemana(inicio, fin time.Time) string {
	meses := []string{
		"ENERO", "FEBRERO", "MARZO", "ABRIL", "MAYO", "JUNIO",
		"JULIO", "AGOSTO", "SEPTIEMBRE", "OCTUBRE", "NOVIEMBRE", "DICIEMBRE",
	}

	mes := meses[inicio.Month()-1]

	// Si ambos están en el mismo mes
	if inicio.Month() == fin.Month() {
		return fmt.Sprintf("%d AL %d DE %s", inicio.Day(), fin.Day(), mes)
	}

	// Si están en meses diferentes
	mesInicio := meses[inicio.Month()-1]
	mesFin := meses[fin.Month()-1]
	return fmt.Sprintf("%d DE %s AL %d DE %s", inicio.Day(), mesInicio, fin.Day(), mesFin)
}

// CalcularSemanaDesdeFecha calcula el lunes y sábado de la semana que contiene una fecha
func CalcularSemanaDesdeFecha(fecha time.Time) (time.Time, time.Time) {
	diaSemana := int(fecha.Weekday())
	if diaSemana == 0 {
		diaSemana = 7
	}
	diaSemana--

	lunes := fecha.AddDate(0, 0, -diaSemana)
	lunes = time.Date(lunes.Year(), lunes.Month(), lunes.Day(), 0, 0, 0, 0, lunes.Location())

	sabado := lunes.AddDate(0, 0, 5)
	sabado = time.Date(sabado.Year(), sabado.Month(), sabado.Day(), 23, 59, 59, 0, sabado.Location())

	return lunes, sabado
}
