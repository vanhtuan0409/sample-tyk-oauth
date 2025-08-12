package main

import (
	"bytes"
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
		"apply_policies": ["mypolicy"],
		"meta_data": {
			"userid": "user_from_key"
		}
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
}
