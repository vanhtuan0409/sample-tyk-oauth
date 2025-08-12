package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"tykloadtest"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing refresh token args")
	}
	rToken := os.Args[1]
	fmt.Println(rToken)

	values := url.Values{}
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", rToken)

	req, err := http.NewRequest(http.MethodPost, tykloadtest.GetTykURL("/oauth/token"), strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatalf("unable to request refresh: +%v", err)
		return
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(tykloadtest.ClientID, tykloadtest.ClientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("unable to request refresh: +%v", err)
		return
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}
