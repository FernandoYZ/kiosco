package utils

import (
	"encoding/json"
	"fmt"
	"html/template"
	"kiosco/models"
	"strconv"
	"time"
)

func GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatearFecha": func(t time.Time) string {
			dias := []string{"Dom", "Lun", "Mar", "Mié", "Jue", "Vie", "Sáb"}
			return dias[t.Weekday()] + " " + strconv.Itoa(t.Day())
		},
		"formatearFechaCompleta": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"formatearFechaLarga": func(t time.Time) string {
			dias := []string{"Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"}
			meses := []string{"", "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
				"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}
			return fmt.Sprintf("%s %02d de %s", dias[t.Weekday()], t.Day(), meses[t.Month()])
		},
		"formatearMoneda": func(valor float64) string {
			return strconv.FormatFloat(valor, 'f', 2, 64)
		},
		"mayorQueCero": func(valor float64) bool {
			return valor > 0
		},
		"obtenerCantidad": func(consumos map[int]map[string]map[int]int, idEstudiante int, fecha time.Time, idProducto int) int {
			fechaKey := fecha.Format("2006-01-02")
			if consumos[idEstudiante] != nil && consumos[idEstudiante][fechaKey] != nil {
				return consumos[idEstudiante][fechaKey][idProducto]
			}
			return 0
		},
		"calcularTotalDia": func(consumos map[int]map[string]map[int]int, idEstudiante int, fecha time.Time, productos []models.Producto) float64 {
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
		},
		"toJSON": func(v interface{}) template.JS {
			if v == nil {
				return template.JS("{}")
			}
			bytes, err := json.Marshal(v)
			if err != nil {
				return template.JS("{}")
			}
			// Si el resultado es "null", devolver objeto vacío
			if string(bytes) == "null" {
				return template.JS("{}")
			}
			return template.JS(bytes)
		},
		"dividirProductos": func(productos []models.Producto) [][]models.Producto {
			// Dividir productos en grupos de 3
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
		"sub": func(a, b int) int {
			return a - b
		},
		"mod": func(a, b int) bool {
			return a%b == 0
		},
		"nombreMes": func(mes time.Month) string {
			meses := []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
				"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}
			return meses[mes-1]
		},
		"obtenerProducto": func(productos []models.Producto, idProducto int) *models.Producto {
			for _, p := range productos {
				if p.IdProducto == idProducto {
					return &p
				}
			}
			return nil
		},
	}
}
