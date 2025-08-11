package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// This handler will be called for all requests not matched by more specific patterns
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requested: %s", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(map[string]bool{
			"ok": true,
		})
	})

	fmt.Println("Server listening on :5050")
	http.ListenAndServe(":5050", nil)
}
