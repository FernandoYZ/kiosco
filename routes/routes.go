package routes

import (
	"kiosco/handlers"
	"net/http"
)

func SetupRoutes(manejador *handlers.Manejador) {
	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Configurar rutas
	http.HandleFunc("/", manejador.ManejadorInicio)
	http.HandleFunc("/editar-consumos", manejador.ManejadorEditarConsumos)
	http.HandleFunc("/guardar-consumos-dia", manejador.ManejadorGuardarConsumosDia)
	http.HandleFunc("/registrar-consumo", manejador.ManejadorRegistrarConsumo)
	http.HandleFunc("/editar-pagos", manejador.ManejadorEditarPagos)
	http.HandleFunc("/registrar-pago", manejador.ManejadorRegistrarPago)
	http.HandleFunc("/eliminar-pago", manejador.ManejadorEliminarPago)
	http.HandleFunc("/ver-consumo-semanal", manejador.ManejadorVerConsumoSemanal)
}
