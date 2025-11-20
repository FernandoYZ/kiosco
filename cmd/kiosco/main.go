package main

import (
	"fmt"
	"kiosco/internal/config"
	"kiosco/internal/controllers"
	"kiosco/internal/db"
	"kiosco/internal/repositories"
	"kiosco/internal/router"
	"kiosco/internal/services"
	"log"
	"net/http"
)

func main() {
	// Banner del sistema
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║    SISTEMA DE CONTROL DE CONSUMO ESCOLAR   ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println()

	// Configurar conexión a la base de datos
	configuracion := config.ObtenerConfigBD()
	baseDatos, err := db.ConectarBD(configuracion)
	if err != nil {
		log.Fatalf("❌ Error fatal al conectar con PostgreSQL: %v", err)
	}
	defer baseDatos.Close()

	// Inicializar capas de la aplicación
	repo := repositories.NuevoRepositorio(baseDatos)
	serv := services.NuevoServicio(repo)
	controlador, err := controllers.NuevoControlador(serv)
	if err != nil {
		log.Fatalf("❌ Error al cargar templates: %v", err)
	}

	// Configurar rutas
	router.ConfigurarRutas(controlador)

	// Iniciar servidor
	direccionServidor := config.ObtenerDireccion()
	fmt.Printf("✓ Servidor escuchando en %s\n", direccionServidor)
	fmt.Println("✓ Presiona Ctrl+C para detener el servidor")
	fmt.Println()

	if err := http.ListenAndServe(direccionServidor, nil); err != nil {
		log.Fatalf("❌ Error al iniciar servidor: %v", err)
	}
}
