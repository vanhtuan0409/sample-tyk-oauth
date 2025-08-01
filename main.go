package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	TykAddr   = "http://127.0.0.1:8081"
	TykSecret = "myabcsecret"
)

func createTykAPIKey() error {
	url := fmt.Sprintf("%s/tyk/keys", TykAddr)

	body := []byte(`{
		"org_id": "1",
		"apply_policies": ["mypolicy"]
	}`)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-tyk-authorization", TykSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return nil
}

func main() {
	if err := createTykAPIKey(); err != nil {
		log.Fatalf("unable to create tyk api key: %v", err)
	}

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
