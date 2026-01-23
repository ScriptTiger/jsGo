//go:build !ecmascript && !wasm

package main

import (
	"net/http"
	"os"
)

func main() {
	os.Stdout.WriteString("Serving at http://localhost:8080...")
	http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))
}
