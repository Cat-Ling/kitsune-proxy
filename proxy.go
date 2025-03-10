package main

import (
	"fmt"
	"io"
	"net/http"
)

// Base API URL
const baseAPI = "https://api.kitsunee.me"

// ProxyHandler forwards all requests dynamically and fixes CORS
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Construct the full API URL
	targetURL := baseAPI + r.URL.Path

	// Forward the request to the target API
	resp, err := http.Get(targetURL)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	// ✅ Manually add CORS headers so browsers accept the response
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Copy response headers from Kitsune (if needed)
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Write response status code
	w.WriteHeader(resp.StatusCode)

	// Stream the response body to the client
	io.Copy(w, resp.Body)
}

func main() {
	// Handle all requests through the proxy
	http.HandleFunc("/", proxyHandler)

	fmt.Println("Proxy server running on http://localhost:54878")
	err := http.ListenAndServe(":54878", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
