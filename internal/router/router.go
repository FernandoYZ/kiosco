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

	// Atajos para proteger HandlerFunc
	p := middleware.Proteger             // Solo requiere autenticación
	pe := middleware.ProtegerEdicion     // Requiere autenticación + puede_editar = 1

	// Rutas que requieren edición (puede_editar = 1)
	mux.HandleFunc("/", pe(controlador.Inicio))
	mux.HandleFunc("/editar-consumos", pe(controlador.EditarConsumos))
	mux.HandleFunc("/guardar-consumos-dia", pe(controlador.GuardarConsumosDia))
	mux.HandleFunc("/registrar-consumo", pe(controlador.RegistrarConsumo))
	mux.HandleFunc("/editar-pagos", pe(controlador.EditarPagos))
	mux.HandleFunc("/registrar-pago", pe(controlador.RegistrarPago))
	mux.HandleFunc("/eliminar-pago", pe(controlador.EliminarPago))
	mux.HandleFunc("/ver-consumo-semanal", pe(controlador.VerConsumoSemanal))

	// Configuración de estudiantes — requiere edición
	mux.HandleFunc("GET /setup", pe(controlador.SetupEstudiantes))
	mux.HandleFunc("POST /setup/estudiante", pe(controlador.AgregarEstudiante))
	mux.HandleFunc("POST /setup/estudiante/actualizar", pe(controlador.ActualizarEstudiante))
	mux.HandleFunc("POST /setup/estudiante/toggle", pe(controlador.ToggleEstudiante))

	// Gestión de productos — solo lectura para usuarios sin edición
	// GET es accesible a todos, POST requiere edición
	mux.HandleFunc("GET /setup/productos", p(controlador.SetupProductos))
	mux.HandleFunc("POST /setup/producto", p(controlador.AgregarProducto))
	mux.HandleFunc("POST /setup/producto/actualizar", p(controlador.ActualizarProducto))
	mux.HandleFunc("POST /setup/producto/toggle", p(controlador.ToggleProducto))

	// Registro de consumos por sector — accesible a todos (solo lectura y envío de datos)
	mux.HandleFunc("GET /registro", p(controlador.RegistroConsumos))
	mux.HandleFunc("GET /registro/estudiantes", p(controlador.ObtenerEstudiantesSector))
	mux.HandleFunc("POST /registro/guardar", p(controlador.GuardarRegistroBatch))

	return mux
}
