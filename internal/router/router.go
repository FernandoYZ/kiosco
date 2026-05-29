package router

import (
	"kiosco/internal/config"
	"kiosco/internal/controllers"
	"kiosco/internal/middleware"
	"net/http"
)

// ConfigurarRutas registra todas las rutas y archivos estáticos, devuelve el mux.
func ConfigurarRutas(controlador *controllers.Controlador) http.Handler {
	mux := http.NewServeMux()

	// Archivos estáticos embebidos (públicos, sin auth)
	config.RegistrarEstaticos(mux)

	// Rutas públicas
	mux.HandleFunc("GET /login", middleware.ProtegerLogin(controlador.MostrarLogin))
	mux.HandleFunc("POST /login", middleware.ProtegerLogin(controlador.ProcesarLogin))
	mux.HandleFunc("GET /logout", controlador.Logout)

	// Atajos para proteger HandlerFunc
	proteger := middleware.Proteger             // Solo requiere autenticación
	protegerEdicion := middleware.ProtegerEdicion     // Requiere autenticación + puede_editar = 1

	// Rutas que requieren edición (puede_editar = 1)
	mux.HandleFunc("GET /", protegerEdicion(controlador.Inicio))
	mux.HandleFunc("GET /editar-consumos", proteger(controlador.EditarConsumos))
	mux.HandleFunc("POST /guardar-consumos-dia", proteger(controlador.GuardarConsumosDia))
	mux.HandleFunc("POST /registrar-consumo", protegerEdicion(controlador.RegistrarConsumo))
	mux.HandleFunc("GET /editar-pagos", protegerEdicion(controlador.EditarPagos))
	mux.HandleFunc("POST /registrar-pago", protegerEdicion(controlador.RegistrarPago))
	mux.HandleFunc("POST /eliminar-pago", protegerEdicion(controlador.EliminarPago))
	mux.HandleFunc("GET /ver-consumo-semanal", protegerEdicion(controlador.VerConsumoSemanal))

	// Configuración de estudiantes — requiere edición
	mux.HandleFunc("GET /setup", protegerEdicion(controlador.SetupEstudiantes))
	mux.HandleFunc("POST /setup/estudiante", protegerEdicion(controlador.AgregarEstudiante))
	mux.HandleFunc("POST /setup/estudiante/actualizar", protegerEdicion(controlador.ActualizarEstudiante))
	mux.HandleFunc("POST /setup/estudiante/toggle", protegerEdicion(controlador.ToggleEstudiante))

	// Gestión de productos — solo lectura para usuarios sin edición
	// GET es accesible a todos, POST requiere edición
	mux.HandleFunc("GET /setup/productos", proteger(controlador.SetupProductos))
	mux.HandleFunc("POST /setup/producto", protegerEdicion(controlador.AgregarProducto))
	mux.HandleFunc("POST /setup/producto/actualizar", protegerEdicion(controlador.ActualizarProducto))
	mux.HandleFunc("POST /setup/producto/toggle", protegerEdicion(controlador.ToggleProducto))

	// Registro de consumos por sector — accesible a todos
	mux.HandleFunc("GET /registro", proteger(controlador.RegistroConsumos))
	mux.HandleFunc("GET /registro/menor", proteger(controlador.RegistroSector))
	mux.HandleFunc("GET /registro/mayor", proteger(controlador.RegistroSector))

	// Resumen de consumos por sector — accesible a todos
	mux.HandleFunc("GET /resumen/menor", proteger(controlador.ResumenSector))
	mux.HandleFunc("GET /resumen/mayor", proteger(controlador.ResumenSector))

	return middleware.LimitarConcurrencia(middleware.LimiteConcurrenciaDefault)(mux)
}
