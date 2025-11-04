package main

import (
	"encoding/json"
	"html/template"
	"log"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Manejador contiene las dependencias para los handlers
type Manejador struct {
	servicio  *Servicio
	templates *template.Template
}

// NuevoManejador crea una nueva instancia del manejador
func NuevoManejador(servicio *Servicio) (*Manejador, error) {
	// Cargar templates con funciones personalizadas
	funcMap := template.FuncMap{
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
		"calcularTotalDia": func(consumos map[int]map[string]map[int]int, idEstudiante int, fecha time.Time, productos []Producto) float64 {
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
		"dividirProductos": func(productos []Producto) [][]Producto {
			// Dividir productos en grupos de 3
			var grupos [][]Producto
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
		"obtenerProducto": func(productos []Producto, idProducto int) *Producto {
			for _, p := range productos {
				if p.IdProducto == idProducto {
					return &p
				}
			}
			return nil
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("templates/*.tmpl")
	if err != nil {
		return nil, err
	}

	return &Manejador{
		servicio:  servicio,
		templates: tmpl,
	}, nil
}

// ManejadorInicio muestra la vista principal con la semana actual
func (m *Manejador) ManejadorInicio(w http.ResponseWriter, r *http.Request) {
	fechaInicio, fechaFin := ObtenerSemanaActual()

	// Verificar si hay parámetros de fecha en la URL
	if fechaParam := r.URL.Query().Get("fecha"); fechaParam != "" {
		if fecha, err := time.Parse("2006-01-02", fechaParam); err == nil {
			fechaInicio, fechaFin = CalcularSemanaDesdeFecha(fecha)
		}
	}

	// Verificar si hay filtro de grado (por defecto primer grado = 1)
	idGrado := 1
	if gradoParam := r.URL.Query().Get("grado"); gradoParam != "" {
		if grado, err := strconv.Atoi(gradoParam); err == nil {
			idGrado = grado
		}
	}

	// Obtener días deshabilitados del parámetro URL
	diasDeshabilitados := r.URL.Query().Get("dias_off")

	datos, err := m.servicio.ObtenerDatosVistaPrincipal(fechaInicio, fechaFin, idGrado, diasDeshabilitados)
	if err != nil {
		log.Printf("Error al obtener datos: %v", err)
		http.Error(w, "Error al cargar datos", http.StatusInternalServerError)
		return
	}

	err = m.templates.ExecuteTemplate(w, "index.tmpl", datos)
	if err != nil {
		log.Printf("Error al renderizar template: %v", err)
	}
}

// ManejadorRegistrarConsumo procesa el formulario de registro de consumo
func (m *Manejador) ManejadorRegistrarConsumo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear datos del formulario (soporta form-urlencoded y multipart)
	err := r.ParseForm()
	if err != nil {
		// Intentar con multipart si falla
		err = r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			log.Printf("Error al parsear formulario: %v", err)
			http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
			return
		}
	}

	idEstudianteStr := r.FormValue("id_estudiante")
	idEstudiante, err := strconv.Atoi(idEstudianteStr)
	if err != nil {
		log.Printf("Error al convertir id_estudiante '%s': %v", idEstudianteStr, err)
		http.Error(w, "ID de estudiante inválido", http.StatusBadRequest)
		return
	}

	idProducto, err := strconv.Atoi(r.FormValue("id_producto"))
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	cantidad, err := strconv.Atoi(r.FormValue("cantidad"))
	if err != nil {
		http.Error(w, "Cantidad inválida", http.StatusBadRequest)
		return
	}

	fechaStr := r.FormValue("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}

	// Registrar el consumo
	err = m.servicio.RegistrarConsumoDesdeFormulario(idEstudiante, idProducto, cantidad, fecha)
	if err != nil {
		log.Printf("Error al registrar consumo: %v", err)
		http.Error(w, "Error al registrar consumo", http.StatusInternalServerError)
		return
	}

	// Obtener el grado para mantenerlo en la redirección
	grado := r.FormValue("grado")
	urlRedireccion := "/?fecha=" + fechaStr
	if grado != "" {
		urlRedireccion += "&grado=" + grado
	}

	// Redireccionar de vuelta a la vista principal
	http.Redirect(w, r, urlRedireccion, http.StatusSeeOther)
}

// ManejadorRegistrarPago procesa el formulario de registro de pago
func (m *Manejador) ManejadorRegistrarPago(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	idEstudiante, err := strconv.Atoi(r.FormValue("id_estudiante"))
	if err != nil {
		http.Error(w, "ID de estudiante inválido", http.StatusBadRequest)
		return
	}

	monto, err := strconv.ParseFloat(r.FormValue("monto"), 64)
	if err != nil {
		http.Error(w, "Monto inválido", http.StatusBadRequest)
		return
	}

	fechaStr := r.FormValue("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		fecha = time.Now() // Si no hay fecha, usar la actual
	}

	// Permitir especificar fecha_pago diferente (para la vista de editar pagos)
	fechaPagoStr := r.FormValue("fecha_pago")
	fechaPago := fecha // Por defecto usar la fecha de la semana
	if fechaPagoStr != "" {
		if fp, err := time.Parse("2006-01-02", fechaPagoStr); err == nil {
			fechaPago = fp
		}
	}

	// Registrar el pago con la fecha especificada
	err = m.servicio.RegistrarPagoDesdeFormulario(idEstudiante, monto, fechaPago)
	if err != nil {
		log.Printf("Error al registrar pago: %v", err)
		http.Error(w, "Error al registrar pago: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Obtener el grado para mantenerlo en la redirección
	grado := r.FormValue("grado")
	redirect := r.FormValue("redirect")

	var urlRedireccion string
	if redirect == "editar-pagos" {
		// Redirigir a la vista de editar pagos
		urlRedireccion = fmt.Sprintf("/editar-pagos?id_estudiante=%d&fecha=%s", idEstudiante, fechaStr)
		if grado != "" {
			urlRedireccion += "&grado=" + grado
		}
	} else {
		// Redirigir a la vista principal
		urlRedireccion = "/?fecha=" + fechaStr
		if grado != "" {
			urlRedireccion += "&grado=" + grado
		}
	}

	// Redireccionar de vuelta
	http.Redirect(w, r, urlRedireccion, http.StatusSeeOther)
}

// ManejadorEditarConsumos muestra la vista para editar consumos de un día
func (m *Manejador) ManejadorEditarConsumos(w http.ResponseWriter, r *http.Request) {
	idEstudiante, err := strconv.Atoi(r.URL.Query().Get("id_estudiante"))
	if err != nil {
		http.Error(w, "ID de estudiante inválido", http.StatusBadRequest)
		return
	}

	fechaStr := r.URL.Query().Get("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}

	idGrado := 0
	if gradoParam := r.URL.Query().Get("grado"); gradoParam != "" {
		if grado, err := strconv.Atoi(gradoParam); err == nil {
			idGrado = grado
		}
	}

	// Obtener datos del estudiante
	estudiantes, err := m.servicio.repo.ObtenerEstudiantesPorGrado(0) // Todos para buscar
	if err != nil {
		http.Error(w, "Error al obtener estudiante", http.StatusInternalServerError)
		return
	}

	var nombreEstudiante string
	for _, est := range estudiantes {
		if est.IdEstudiante == idEstudiante {
			nombreEstudiante = est.Apellidos + ", " + est.Nombres
			break
		}
	}

	// Obtener productos
	productos, err := m.servicio.repo.ObtenerProductosActivos()
	if err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}

	// Obtener consumos existentes
	fechaInicio := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, fecha.Location())
	fechaFin := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 23, 59, 59, 0, fecha.Location())
	consumos, err := m.servicio.repo.ObtenerConsumosSemana(fechaInicio, fechaFin)
	if err != nil {
		http.Error(w, "Error al obtener consumos", http.StatusInternalServerError)
		return
	}

	// Crear mapa de consumos
	consumosPorDia := make(map[int]map[string]map[int]int)
	for _, c := range consumos {
		if consumosPorDia[c.IdEstudiante] == nil {
			consumosPorDia[c.IdEstudiante] = make(map[string]map[int]int)
		}
		fechaKey := c.FechaConsumo.Format("2006-01-02")
		if consumosPorDia[c.IdEstudiante][fechaKey] == nil {
			consumosPorDia[c.IdEstudiante][fechaKey] = make(map[int]int)
		}
		consumosPorDia[c.IdEstudiante][fechaKey][c.IdProducto] = c.Cantidad
	}

	datos := DatosEditarConsumos{
		IdEstudiante:      idEstudiante,
		NombreEstudiante:  nombreEstudiante,
		Fecha:             fecha,
		Productos:         productos,
		Consumos:          consumosPorDia,
		GradoSeleccionado: idGrado,
	}

	err = m.templates.ExecuteTemplate(w, "editar_consumos.tmpl", datos)
	if err != nil {
		log.Printf("Error al renderizar template: %v", err)
	}
}

// ManejadorGuardarConsumosDia guarda todos los consumos de un día
func (m *Manejador) ManejadorGuardarConsumosDia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	idEstudiante, err := strconv.Atoi(r.FormValue("id_estudiante"))
	if err != nil {
		http.Error(w, "ID de estudiante inválido", http.StatusBadRequest)
		return
	}

	fechaStr := r.FormValue("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}

	grado := r.FormValue("grado")

	// Obtener productos para procesar
	productos, err := m.servicio.repo.ObtenerProductosActivos()
	if err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}

	// Procesar cada producto
	for _, producto := range productos {
		cantidadStr := r.FormValue(fmt.Sprintf("cantidad_%d", producto.IdProducto))
		cantidad, err := strconv.Atoi(cantidadStr)
		if err != nil {
			cantidad = 0
		}

		// Actualizar consumo
		err = m.servicio.RegistrarConsumoDesdeFormulario(idEstudiante, producto.IdProducto, cantidad, fecha)
		if err != nil {
			log.Printf("Error al registrar consumo producto %d: %v", producto.IdProducto, err)
		}
	}

	// Redireccionar
	urlRedireccion := "/?fecha=" + fechaStr
	if grado != "" {
		urlRedireccion += "&grado=" + grado
	}

	http.Redirect(w, r, urlRedireccion, http.StatusSeeOther)
}

// ManejadorEditarPagos muestra la vista para editar pagos de una semana
func (m *Manejador) ManejadorEditarPagos(w http.ResponseWriter, r *http.Request) {
	idEstudiante, err := strconv.Atoi(r.URL.Query().Get("id_estudiante"))
	if err != nil {
		http.Error(w, "ID de estudiante inválido", http.StatusBadRequest)
		return
	}

	fechaStr := r.URL.Query().Get("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}

	idGrado := 0
	if gradoParam := r.URL.Query().Get("grado"); gradoParam != "" {
		if grado, err := strconv.Atoi(gradoParam); err == nil {
			idGrado = grado
		}
	}

	// Calcular la semana
	fechaInicio, fechaFin := CalcularSemanaDesdeFecha(fecha)

	// Obtener datos del estudiante
	estudiantes, err := m.servicio.repo.ObtenerEstudiantesPorGrado(0) // Todos para buscar
	if err != nil {
		http.Error(w, "Error al obtener estudiante", http.StatusInternalServerError)
		return
	}

	var nombreEstudiante string
	for _, est := range estudiantes {
		if est.IdEstudiante == idEstudiante {
			nombreEstudiante = est.Apellidos + ", " + est.Nombres
			break
		}
	}

	// Obtener pagos de la semana
	pagos, err := m.servicio.repo.ObtenerPagosSemanaDetalle(idEstudiante, fechaInicio, fechaFin)
	if err != nil {
		http.Error(w, "Error al obtener pagos", http.StatusInternalServerError)
		return
	}

	// Calcular total de pagos
	totalPagos := 0.0
	for _, p := range pagos {
		totalPagos += p.Monto
	}

	datos := DatosEditarPagos{
		IdEstudiante:      idEstudiante,
		NombreEstudiante:  nombreEstudiante,
		FechaInicio:       fechaInicio,
		FechaFin:          fechaFin,
		Pagos:             pagos,
		TotalPagos:        totalPagos,
		GradoSeleccionado: idGrado,
	}

	err = m.templates.ExecuteTemplate(w, "editar_pagos.tmpl", datos)
	if err != nil {
		log.Printf("Error al renderizar template: %v", err)
	}
}

// ManejadorEliminarPago elimina un pago específico
func (m *Manejador) ManejadorEliminarPago(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	idPago, err := strconv.Atoi(r.FormValue("id_pago"))
	if err != nil {
		http.Error(w, "ID de pago inválido", http.StatusBadRequest)
		return
	}

	// Eliminar el pago
	err = m.servicio.repo.EliminarPago(idPago)
	if err != nil {
		log.Printf("Error al eliminar pago: %v", err)
		http.Error(w, "Error al eliminar pago", http.StatusInternalServerError)
		return
	}

	// Redirigir de vuelta a la vista de editar pagos
	idEstudiante := r.FormValue("id_estudiante")
	fechaStr := r.FormValue("fecha")
	grado := r.FormValue("grado")

	urlRedireccion := fmt.Sprintf("/editar-pagos?id_estudiante=%s&fecha=%s", idEstudiante, fechaStr)
	if grado != "" {
		urlRedireccion += "&grado=" + grado
	}

	http.Redirect(w, r, urlRedireccion, http.StatusSeeOther)
}

// ManejadorVerConsumoSemanal muestra el comprobante de consumo semanal de un estudiante
func (m *Manejador) ManejadorVerConsumoSemanal(w http.ResponseWriter, r *http.Request) {
	idEstudiante, err := strconv.Atoi(r.URL.Query().Get("id_estudiante"))
	if err != nil {
		http.Error(w, "ID de estudiante inválido", http.StatusBadRequest)
		return
	}

	fechaStr := r.URL.Query().Get("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}

	idGrado := 0
	if gradoParam := r.URL.Query().Get("grado"); gradoParam != "" {
		if grado, err := strconv.Atoi(gradoParam); err == nil {
			idGrado = grado
		}
	}

	// Calcular la semana
	fechaInicio, fechaFin := CalcularSemanaDesdeFecha(fecha)

	// Obtener datos del estudiante
	estudiantes, err := m.servicio.repo.ObtenerEstudiantesPorGrado(0)
	if err != nil {
		http.Error(w, "Error al obtener estudiante", http.StatusInternalServerError)
		return
	}

	var nombreEstudiante string
	for _, est := range estudiantes {
		if est.IdEstudiante == idEstudiante {
			nombreEstudiante = est.Apellidos + ", " + est.Nombres
			break
		}
	}

	// Obtener consumos de la semana
	consumos, err := m.servicio.repo.ObtenerConsumosSemana(fechaInicio, fechaFin)
	if err != nil {
		http.Error(w, "Error al obtener consumos", http.StatusInternalServerError)
		return
	}

	// Obtener productos
	productos, err := m.servicio.repo.ObtenerProductosActivos()
	if err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}

	// Crear mapa de productos para búsqueda rápida
	productosMap := make(map[int]Producto)
	for _, p := range productos {
		productosMap[p.IdProducto] = p
	}

	// Organizar consumos por día
	diasConsumo := make(map[string][]ConsumoProducto)
	totalesPorDia := make(map[string]float64)

	for _, c := range consumos {
		if c.IdEstudiante != idEstudiante {
			continue
		}

		fechaKey := c.FechaConsumo.Format("2006-01-02")
		producto := productosMap[c.IdProducto]

		consumoProducto := ConsumoProducto{
			Nombre:   producto.Nombre,
			Cantidad: c.Cantidad,
			Precio:   c.PrecioUnitarioVenta,
			Total:    c.TotalLinea,
		}

		diasConsumo[fechaKey] = append(diasConsumo[fechaKey], consumoProducto)
		totalesPorDia[fechaKey] += c.TotalLinea
	}

	// Crear lista de consumos diarios en orden
	var consumosPorDia []ConsumoDiario
	currentDate := fechaInicio
	subTotal := 0.0

	for currentDate.Before(fechaFin) || currentDate.Equal(fechaFin) {
		fechaKey := currentDate.Format("2006-01-02")

		consumoDia := ConsumoDiario{
			Fecha:     currentDate,
			Productos: diasConsumo[fechaKey],
			Total:     totalesPorDia[fechaKey],
		}

		// Solo agregar días con consumo
		if len(consumoDia.Productos) > 0 {
			consumosPorDia = append(consumosPorDia, consumoDia)
			subTotal += consumoDia.Total
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Obtener deuda anterior y pagos
	deudaAnterior, _ := m.servicio.repo.ObtenerDeudaAnterior(idEstudiante, fechaInicio)
	pagos, _ := m.servicio.repo.ObtenerPagosSemana(idEstudiante, fechaInicio, fechaFin)

	total := subTotal + deudaAnterior - pagos

	datos := DatosConsumoSemanal{
		IdEstudiante:      idEstudiante,
		NombreEstudiante:  nombreEstudiante,
		FechaInicio:       fechaInicio,
		FechaFin:          fechaFin,
		ConsumosPorDia:    consumosPorDia,
		SubTotal:          subTotal,
		DeudaAnterior:     deudaAnterior,
		Pagos:             pagos,
		Total:             total,
		GradoSeleccionado: idGrado,
	}

	err = m.templates.ExecuteTemplate(w, "ver_consumo_semanal.tmpl", datos)
	if err != nil {
		log.Printf("Error al renderizar template: %v", err)
	}
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

