package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tykloadtest"
)

func main() {
	if err := createOauthApp(); err != nil {
		log.Fatalf("unable to create oauth app: %+v", err)
	}

	log.Printf("oauth start url: %s", tykloadtest.SampleOauthApp.LoginURL("mystate", "myscope"))
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			renderLoginHandler(w, r)
		} else {
			loginHandler(w, r)
		}
	})
	http.HandleFunc("/redirect", redirectHandler)

	// This handler will be called for all requests not matched by more specific patterns
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userid := r.Header.Get("x-anduin-userid")
		log.Printf("Requested: %s - userid: %s", r.URL.Path, userid)

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
