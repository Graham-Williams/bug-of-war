package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Serve static files from the current directory
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
