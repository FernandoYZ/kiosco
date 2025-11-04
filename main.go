package main

import (
	"fmt"
	"kiosco/config"
	"kiosco/handlers"
	"kiosco/repository"
	"kiosco/routes"
	"kiosco/services"
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
	cfg := config.ObtenerConfigBD()
	db, err := config.ConectarBD(cfg)
	if err != nil {
		log.Fatalf("❌ Error fatal al conectar con PostgreSQL: %v", err)
	}
	defer db.Close()

	// Inicializar capas de la aplicación
	repo := repository.NuevoRepositorio(db)
	serv := services.NuevoServicio(repo)
	manejador, err := handlers.NuevoManejador(serv)
	if err != nil {
		log.Fatalf("❌ Error al cargar templates: %v", err)
	}

	// Configurar rutas
	routes.SetupRoutes(manejador)

	// Iniciar servidor
	puerto := config.ObtenerPuerto()
	fmt.Printf("✓ Servidor iniciado en http://localhost%s\n", puerto)
	fmt.Println("✓ Presiona Ctrl+C para detener el servidor")
	fmt.Println()

	if err := http.ListenAndServe(puerto, nil); err != nil {
		log.Fatalf("❌ Error al iniciar servidor: %v", err)
	}
}
