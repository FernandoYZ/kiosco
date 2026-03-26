package main

import (
	"fmt"
	"kiosco/internal/auth"
	"kiosco/internal/config"
	"kiosco/internal/controllers"
	"kiosco/internal/router"
	"log"
	"net/http"
)

func main() {
	fmt.Println("==== SISTEMA DE CONTROL DE CONSUMO ESCOLAR ====")
	fmt.Println()

	// Generar llave efímera HMAC en memoria (se invalida al reiniciar el servidor)
	auth.LlaveEfimera()
	fmt.Println("✓ Llave de sesión generada en memoria")

	// Inicializar controlador (SQLite, repositorio y servicio se inicializan internamente)
	controlador, err := controllers.NuevoControlador()
	if err != nil {
		log.Fatalf("❌ Error al iniciar controlador: %v", err)
	}

	// Configurar rutas y archivos estáticos
	mux := router.ConfigurarRutas(controlador)

	// Iniciar servidor
	direccionServidor := config.ObtenerDireccion()
	fmt.Printf("✓ Servidor escuchando en %s\n", direccionServidor)
	fmt.Println("✓ Presiona Ctrl+C para detener el servidor")
	fmt.Println()

	if err := http.ListenAndServe(direccionServidor, mux); err != nil {
		log.Fatalf("❌ Error al iniciar servidor: %v", err)
	}
}
