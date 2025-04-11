package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api", apiHandler)
	// OOPS!!!! Legacy code pushed to prod :(
	http.HandleFunc("/api/v1/login", legacyLoginHandler) 
	http.HandleFunc("/api/v2/info", infoHandler)
	http.HandleFunc("/api/v2/status", statusHandler)
	http.HandleFunc("/api/v2/login", loginHandler) 

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the base /api route. Other routes are /api/v2/status & /api/v2/info\n")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API Status: OK\n")
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API Info: Version 1.0\n")
}

// LEGACY CODE - Vulnerable back door to admin
func legacyLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt: username=%s, password=%s\n", creds.Username, creds.Password)

	if creds.Username == "admin" && creds.Password == "letmein" {
		fmt.Fprintf(w, "Welcome, admin!\n")
		return
	} else if creds.Username == "admin" {
		http.Error(w, "Incorrect password for admin.\n", http.StatusForbidden)
		return
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}


func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt: username=%s, password=%s\n", creds.Username, creds.Password)
  // Removed in v2 as a patch!!
	// if creds.Username == "admin" && creds.Password == "letmein" {
	// 	fmt.Fprintf(w, "Welcome, admin!\n")
	// 	return
	// }

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}
