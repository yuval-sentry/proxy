package main

import (
	"fmt"
	"io"
	http "net/http"
)

const Port = ":8080"

// Define a type for the request details
type RequestDetails struct {
    Method string
    URL    string
    Body   []byte
}

// Channel to forward request details
var requestChan = make(chan RequestDetails)

// Number of worker goroutines
const numWorkers = 15

// Worker function to process requests
func worker(id int, requests <-chan RequestDetails) {
    for request := range requests {
        // Process request (e.g., log or handle asynchronously)
        fmt.Printf("Worker %d received request: %s %s\n", id, request.Method, request.URL)
        fmt.Printf("Body: %s\n", request.Body)
    }
}

// Handler function for HTTP requests
func handler(w http.ResponseWriter, r *http.Request) {
    // Read request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusInternalServerError)
        return
    }

    // Extract request details
    request := RequestDetails{
        Method: r.Method,
        URL:    r.URL.String(),
        Body:   body,
    }

    // Send request details through the channel
    requestChan <- request

    // Write response
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintf(w, "Request forwarded to a worker")
}

func main() {
    // Start worker goroutines
    for i := 0; i < numWorkers; i++ {
        go worker(i, requestChan)
    }

    // Register handler function for the root URL pattern "/"
    http.HandleFunc("/", handler)

    // Start the HTTP server on "$Port"
    fmt.Println("Server listening on port 8080" + Port)
    if err := http.ListenAndServe(Port, nil); err != nil {
        fmt.Printf("Failed to start server: %s\n", err)
    }
}