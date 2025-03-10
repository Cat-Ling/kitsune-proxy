package main

import (
	"fmt"
	"io"
	"net/http"
)

// Base API URL
const baseAPI = "https://api.kitsunee.me"

// ProxyHandler forwards requests and fixes CORS
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// ✅ Handle CORS preflight requests
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Construct full API URL
	targetURL := baseAPI + r.URL.Path

	// Forward request
	resp, err := http.Get(targetURL)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Write response status
	w.WriteHeader(resp.StatusCode)

	// Stream response body
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/", proxyHandler)

	fmt.Println("Proxy server running on http://localhost:54878")
	err := http.ListenAndServe(":54878", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
