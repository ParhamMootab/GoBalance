package main

import (
    "fmt"
    "log"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from server %s!", r.RemoteAddr)
}

func main() {
    http.HandleFunc("/", handler)
    port := ":8081" // Change the port for each instance
    fmt.Printf("Backend server started at %s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}