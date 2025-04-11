package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/api/info", infoHandler)
	fmt.Println("server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	blocked := []string{"<script", "<img", "onerror", "onload", "<svg", "javascript:", "<iframe"}

	isBlocked := false
	for _, pattern := range blocked {
		if strings.Contains(strings.ToLower(query), pattern) {
			isBlocked = true
			break
		}
	}

	allowedBypasses := []string{
		`<scr<script>ipt>alert(1)</scr</script>ipt>`,
		`<svg/onload=alert(1)>`,
		`<body onload=alert(1)>`,
	}

	isBypass := false
	for _, bypass := range allowedBypasses {
		if query == bypass {
			isBypass = true
			break
		}
	}

	if isBlocked && !isBypass {
		http.Error(w, "Bad Request: input filtered", http.StatusBadRequest)
		return
	}

	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head><title>Search Result</title></head>
		<body>
			<h1>Search Results</h1>
			<p>You searched for: %s</p>
		</body>
		</html>
	`, query)

	os.MkdirAll("./xss-results", 0755)
	filename := fmt.Sprintf("xss-results/result_%d.html", time.Now().UnixNano())

	err := os.WriteFile(filename, []byte(htmlContent), 0644)
	if err != nil {
		log.Printf("Failed to write HTML file: %v\n", err)
		http.Error(w, "Could not write file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Saved HTML to: %s\n", filename)
}
