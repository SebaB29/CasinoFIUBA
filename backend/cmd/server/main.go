package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Servidor backend funcionando (temporal).")
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Servidor backend escuchando en puerto 8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("Error al iniciar el servidor:", err)
    }
}
