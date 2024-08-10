package main

import (
    "fmt"
    "log"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from server %v!", 8083)
}

func healthHandler(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "OK")
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/health", healthHandler)

    port := ":8083" // Change the port for each instance
    fmt.Printf("Backend server started at %s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}