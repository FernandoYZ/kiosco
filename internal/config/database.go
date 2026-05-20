package config

import (
	"database/sql"
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// Se va usar la estructura de la base de datos en SQLite
//go:embed schema.sql
var schemaSQL string

const dbPath = "database/database.db"

var (
	instancia *sql.DB
	una       sync.Once
)

// DB retorna la instancia singleton de la base de datos.
func DB() *sql.DB {
	una.Do(func() {
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			log.Fatal("Error creando directorio de DB:", err)
		}
		dbNueva := !verificarDB()

		var err error
		instancia, err = sql.Open("sqlite", dbPath)
		if err != nil {
			log.Fatal("Error al abrir DB:", err)
		}
		instancia.SetMaxOpenConns(5)
		instancia.SetMaxIdleConns(2)
		instancia.SetConnMaxLifetime(10 * time.Minute)

		pragmas := []string{
			"PRAGMA journal_mode=WAL",
			"PRAGMA synchronous=NORMAL",
			"PRAGMA busy_timeout=5000",
		}
		for _, pragma := range pragmas {
			if _, err := instancia.Exec(pragma); err != nil {
				log.Printf("Warning: PRAGMA %s failed: %v", pragma, err)
			}
		}

		if err := instancia.Ping(); err != nil {
			log.Fatal("No se pudo conectar a la DB:", err)
		}
		if dbNueva {
			inicializarDB(instancia)
		}
	})
	return instancia
}

func verificarDB() bool {
	_, err := os.Stat(dbPath)
	return !os.IsNotExist(err) && err == nil
}

func inicializarDB(db *sql.DB) {
	if _, err := db.Exec(schemaSQL); err != nil {
		log.Fatal("Error al inicializar DB:", err)
	}
	log.Println("Base de datos inicializada correctamente")
}
