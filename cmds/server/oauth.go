package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"tykloadtest"
)

func renderLoginHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
	<title>Login</title>
</head>
<body>
	<form method="post" action="/login">
		<input type="hidden" name="client_id" value="%s" />
		<input type="hidden" name="redirect_uri" value="%s" />
		<input type="hidden" name="state" value="%s" />
		<input type="hidden" name="scope" value="%s" />
		<input type="submit" value="Submit"/>
	</form>
</body>
</html>`,
		queries.Get("client_id"),
		queries.Get("redirect_uri"),
		queries.Get("state"),
		queries.Get("scope"),
	)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		responseError(w, err)
		return
	}

	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("client_id", r.FormValue("client_id"))
	values.Set("redirect_uri", r.FormValue("redirect_uri"))
	values.Set("state", r.FormValue("state"))
	values.Set("scope", r.FormValue("scope"))
	values.Set("meta_data", `{"meta_data": {"foo": "bar"}}`)

	req, err := http.NewRequest(http.MethodPost, tykloadtest.SampleOauthApp.GetOauthUrl("/tyk/oauth/authorize-client"), strings.NewReader(values.Encode()))
	if err != nil {
		responseError(w, err)
		return
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("x-tyk-authorization", "myabcsecret")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		responseError(w, err)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(io.TeeReader(resp.Body, os.Stdout))
	var data struct {
		Code       string `json:"code"`
		RedirectTo string `json:"redirect_to"`
		State      string `json:"state"`
	}
	if err := decoder.Decode(&data); err != nil {
		responseError(w, err)
		return
	}
	if data.RedirectTo == "" {
		responseError(w, errors.New("empty redirect to"))
		return
	}

	w.Header().Set("location", data.RedirectTo)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	code := queries.Get("code")
	log.Printf("received code %s", code)

	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)

	req, err := http.NewRequest(http.MethodPost, tykloadtest.SampleOauthApp.GetOauthUrl("/oauth/token"), strings.NewReader(values.Encode()))
	if err != nil {
		responseError(w, err)
		return
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(tykloadtest.SampleOauthApp.ClientID, tykloadtest.SampleOauthApp.Secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		responseError(w, err)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}

func createOauthApp() error {
	createUrl := fmt.Sprintf("%s/tyk/oauth/clients/create", tykloadtest.TykAdminEndpoint)

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(tykloadtest.SampleOauthApp); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, createUrl, &body)
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-tyk-authorization", tykloadtest.TykAdminSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return nil
}

func responseError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("api error: %v", err)
}
