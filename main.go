package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Banner del sistema
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║    SISTEMA DE CONTROL DE CONSUMO ESCOLAR   ║")
	fmt.Println("║               Kiosco v1.0                  ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println()

	// Configurar conexión a la base de datos
	config := ObtenerConfigBD()
	db, err := ConectarBD(config)
	if err != nil {
		log.Fatalf("❌ Error fatal al conectar con PostgreSQL: %v", err)
	}
	defer db.Close()

	// Inicializar capas de la aplicación
	repo := NuevoRepositorio(db)
	servicio := NuevoServicio(repo)
	manejador, err := NuevoManejador(servicio)
	if err != nil {
		log.Fatalf("❌ Error al cargar templates: %v", err)
	}

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

	// Iniciar servidor
	puerto := ":3200"
	fmt.Printf("✓ Servidor iniciado en http://localhost%s\n", puerto)
	fmt.Println("✓ Presiona Ctrl+C para detener el servidor")
	fmt.Println()

	if err := http.ListenAndServe(puerto, nil); err != nil {
		log.Fatalf("❌ Error al iniciar servidor: %v", err)
	}
}
