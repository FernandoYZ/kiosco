package router

import (
	"kiosco/internal/controllers"
	"net/http"
)

// ConfigurarRutas configura todas las rutas de la aplicación
func ConfigurarRutas(controlador *controllers.Controlador) {
	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Configurar rutas de la aplicación
	http.HandleFunc("/", controlador.Inicio)
	http.HandleFunc("/editar-consumos", controlador.EditarConsumos)
	http.HandleFunc("/guardar-consumos-dia", controlador.GuardarConsumosDia)
	http.HandleFunc("/registrar-consumo", controlador.RegistrarConsumo)
	http.HandleFunc("/editar-pagos", controlador.EditarPagos)
	http.HandleFunc("/registrar-pago", controlador.RegistrarPago)
	http.HandleFunc("/eliminar-pago", controlador.EliminarPago)
	http.HandleFunc("/ver-consumo-semanal", controlador.VerConsumoSemanal)
}
