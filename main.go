package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Middleware para verificar el origen de la solicitud
func verifyOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica si la solicitud proviene de un navegador
		if strings.Contains(r.Header.Get("User-Agent"), "Mozilla") {
			// Verifica si el dominio de origen es el permitido
			if r.Header.Get("Origin") == "http://dominio-permitido.com" {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Si no cumple con los criterios, devolver un error de acceso prohibido
		http.Error(w, "Acceso prohibido", http.StatusForbidden)
	})
}

type Response struct {
	Message string `json:"message"`
}

// Controlador del endpoint protegido
func protectedEndpoint(w http.ResponseWriter, r *http.Request) {

	response := Response{
		Message: "Protected endpoint reached successfully!",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	// Manejador del endpoint protegido
	http.Handle("/protected", verifyOrigin(http.HandlerFunc(protectedEndpoint)))

	// Inicia el servidor en el puerto 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
