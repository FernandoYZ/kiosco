package utils

import (
	"encoding/json"
	"fmt"
	"html/template"
	"kiosco/internal/models"
	"strconv"
	"time"
)

// --- Funciones exportadas (usadas en archivos .templ) ---

func FormatearFecha(t time.Time) string {
	dias := []string{"Dom", "Lun", "Mar", "Mié", "Jue", "Vie", "Sáb"}
	return dias[t.Weekday()] + " " + strconv.Itoa(t.Day())
}

func FormatearFechaCompleta(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatearFechaLarga(t time.Time) string {
	dias := []string{"Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"}
	meses := []string{"", "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
		"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}
	return fmt.Sprintf("%s %02d de %s", dias[t.Weekday()], t.Day(), meses[t.Month()])
}

func FormatearMoneda(valor float64) string {
	return strconv.FormatFloat(valor, 'f', 2, 64)
}

func MayorQueCero(valor float64) bool {
	return valor > 0
}

func ObtenerCantidad(consumos map[int]map[string]map[int]int, idEstudiante int, fecha time.Time, idProducto int) int {
	fechaKey := fecha.Format("2006-01-02")
	if consumos[idEstudiante] != nil && consumos[idEstudiante][fechaKey] != nil {
		return consumos[idEstudiante][fechaKey][idProducto]
	}
	return 0
}

func CalcularTotalDia(consumos map[int]map[string]map[int]int, idEstudiante int, fecha time.Time, productos []models.Producto) float64 {
	total := 0.0
	fechaKey := fecha.Format("2006-01-02")
	if consumos[idEstudiante] != nil && consumos[idEstudiante][fechaKey] != nil {
		for _, producto := range productos {
			cantidad := consumos[idEstudiante][fechaKey][producto.IdProducto]
			if cantidad > 0 {
				total += float64(cantidad) * producto.PrecioUnitario
			}
		}
	}
	return total
}

func NombreMes(mes time.Month) string {
	meses := []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
		"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}
	return meses[mes-1]
}

func ObtenerProducto(productos []models.Producto, idProducto int) *models.Producto {
	for i, p := range productos {
		if p.IdProducto == idProducto {
			return &productos[i]
		}
	}
	return nil
}

// GetTemplateFuncs retorna el FuncMap para html/template (vistas .tmpl legacy)
func GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatearFecha":         FormatearFecha,
		"formatearFechaCompleta": FormatearFechaCompleta,
		"formatearFechaLarga":    FormatearFechaLarga,
		"formatearMoneda":        FormatearMoneda,
		"mayorQueCero":           MayorQueCero,
		"obtenerCantidad": func(consumos map[int]map[string]map[int]int, idEstudiante int, fecha time.Time, idProducto int) int {
			return ObtenerCantidad(consumos, idEstudiante, fecha, idProducto)
		},
		"calcularTotalDia": func(consumos map[int]map[string]map[int]int, idEstudiante int, fecha time.Time, productos []models.Producto) float64 {
			return CalcularTotalDia(consumos, idEstudiante, fecha, productos)
		},
		"toJSON": func(v interface{}) template.JS {
			if v == nil {
				return template.JS("{}")
			}
			bytes, err := json.Marshal(v)
			if err != nil || string(bytes) == "null" {
				return template.JS("{}")
			}
			return template.JS(bytes)
		},
		"dividirProductos": func(productos []models.Producto) [][]models.Producto {
			var grupos [][]models.Producto
			for i := 0; i < len(productos); i += 3 {
				end := i + 3
				if end > len(productos) {
					end = len(productos)
				}
				grupos = append(grupos, productos[i:end])
			}
			return grupos
		},
		"sub":             func(a, b int) int { return a - b },
		"mod":             func(a, b int) bool { return a%b == 0 },
		"nombreMes":       func(mes time.Month) string { return NombreMes(mes) },
		"obtenerProducto": func(productos []models.Producto, idProducto int) *models.Producto { return ObtenerProducto(productos, idProducto) },
	}
}
