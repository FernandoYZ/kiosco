package router

import (
	"kiosco/internal/config"
	"kiosco/internal/controllers"
	"kiosco/internal/middleware"
	"net/http"
)

// ConfigurarRutas registra todas las rutas y archivos estáticos, devuelve el mux.
func ConfigurarRutas(controlador *controllers.Controlador) *http.ServeMux {
	mux := http.NewServeMux()

	// Archivos estáticos embebidos (públicos, sin auth)
	config.RegistrarEstaticos(mux)

	// Rutas públicas
	mux.HandleFunc("GET /login", controlador.MostrarLogin)
	mux.HandleFunc("POST /login", controlador.ProcesarLogin)
	mux.HandleFunc("/logout", controlador.Logout)

	// Atajo para proteger HandlerFunc sin repetir middleware.Proteger(...)
	p := middleware.Proteger

	// Rutas principales
	mux.HandleFunc("/", p(controlador.Inicio))
	mux.HandleFunc("/editar-consumos", p(controlador.EditarConsumos))
	mux.HandleFunc("/guardar-consumos-dia", p(controlador.GuardarConsumosDia))
	mux.HandleFunc("/registrar-consumo", p(controlador.RegistrarConsumo))
	mux.HandleFunc("/editar-pagos", p(controlador.EditarPagos))
	mux.HandleFunc("/registrar-pago", p(controlador.RegistrarPago))
	mux.HandleFunc("/eliminar-pago", p(controlador.EliminarPago))
	mux.HandleFunc("/ver-consumo-semanal", p(controlador.VerConsumoSemanal))

	// Configuración de estudiantes
	mux.HandleFunc("/setup", p(controlador.SetupEstudiantes))
	mux.HandleFunc("/setup/estudiante", p(controlador.AgregarEstudiante))
	mux.HandleFunc("/setup/estudiante/actualizar", p(controlador.ActualizarEstudiante))
	mux.HandleFunc("/setup/estudiante/toggle", p(controlador.ToggleEstudiante))

	// Gestión de productos
	mux.HandleFunc("/setup/productos", p(controlador.SetupProductos))
	mux.HandleFunc("/setup/producto", p(controlador.AgregarProducto))
	mux.HandleFunc("/setup/producto/actualizar", p(controlador.ActualizarProducto))
	mux.HandleFunc("/setup/producto/toggle", p(controlador.ToggleProducto))

	return mux
}
