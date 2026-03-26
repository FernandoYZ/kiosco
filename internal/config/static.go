package config

import (
	"io/fs"
	"net/http"

	kiosco "kiosco"
)

// RegistrarEstaticos monta los archivos estáticos embebidos en el mux dado.
func RegistrarEstaticos(mux *http.ServeMux) {
	sub, err := fs.Sub(kiosco.StaticFS, "public")
	if err != nil {
		panic("embed: no se pudo crear sub-FS de public: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(sub))

	mux.Handle("GET /dist/", fileServer)
	mux.Handle("GET /fonts/", fileServer)
	mux.Handle("GET /images/", fileServer)

	mux.HandleFunc("GET /favicon.webp", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, kiosco.StaticFS, "public/favicon.webp")
	})
	mux.HandleFunc("GET /manifest.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/manifest+json")
		http.ServeFileFS(w, r, kiosco.StaticFS, "public/manifest.json")
	})
	mux.HandleFunc("GET /sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Service-Worker-Allowed", "/")
		http.ServeFileFS(w, r, kiosco.StaticFS, "public/sw.js")
	})
}
