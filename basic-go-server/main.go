package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/api/info", infoHandler)
	http.HandleFunc("/api/secret", hiddenHandler) // THIS ROUTE IS "HIDDEN"

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the base /api route. Other routes are /api/status & /api/info\n")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API Status: OK\n")
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API Info: Version 1.0\n")
}

func hiddenHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}
