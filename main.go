package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strings"
)

// Middleware para verificar el origen de la solicitud
func verifyOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := os.Getenv("DOMAIN_URL")
		validApiKey := os.Getenv("API_KEY_GENBOX")

		// Verifica si la solicitud proviene de un navegador
		if strings.Contains(r.Header.Get("User-Agent"), "Mozilla") {
			// Verifica si el dominio de origen es el permitido
			if r.Header.Get("Origin") != url {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			apiKey := r.Header.Get("X-API-Key")
			if apiKey != validApiKey {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
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
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}

func main() {
	// Manejador del endpoint protegido
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	http.Handle("/protected", verifyOrigin(http.HandlerFunc(protectedEndpoint)))

	// Inicia el servidor en el puerto 8080
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
