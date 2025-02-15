package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var videoExtensions = map[string]bool{
	".mp4":  true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".webm": true,
}

func main() {
	// Servir archivos estáticos
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Servir archivos multimedia
	http.Handle("/media/",
		http.StripPrefix("/media/",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Decodificar nombre de archivo con caracteres especiales
				decodedPath, err := filepath.Abs(filepath.Clean(r.URL.Path))
				if err != nil {
					http.Error(w, "Invalid path", http.StatusBadRequest)
					return
				}
				http.ServeFile(w, r, decodedPath)
			}),
		),
	)

	// Endpoint API para listar películas
	http.HandleFunc("/api/movies", moviesHandler)

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func moviesHandler(w http.ResponseWriter, r *http.Request) {
	movies, err := findVideoFiles(".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func findVideoFiles(root string) ([]string, error) {
	var movies []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorar el directorio static
		if info.IsDir() && info.Name() == "static" {
			return filepath.SkipDir
		}

		if !info.IsDir() && videoExtensions[strings.ToLower(filepath.Ext(path))] {
			relPath, _ := filepath.Rel(root, path)
			movies = append(movies, relPath)
		}
		return nil
	})

	return movies, err
}
