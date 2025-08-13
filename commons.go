package tykloadtest

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

const (
	TykHTTPEndpoint  = "http://localhost:8080"
	TykAdminEndpoint = "http://localhost:8081"
	TykAdminSecret   = "myabcsecret"
)

type OauthApp struct {
	APPID       string         `json:"api_id"`
	ClientID    string         `json:"client_id"`
	RedirectURI string         `json:"redirect_uri"`
	PolicyID    string         `json:"policy_id,omitempty"`
	Secret      string         `json:"secret,omitempty"`
	Metadata    map[string]any `json:"meta_data,omitempty"`
	Description string         `json:"description,omitempty"`

	ListenPath string `json:"-"`
}

var SampleOauthApp = &OauthApp{
	APPID:       "myapp_oauth",
	ClientID:    "sample_oauth_app",
	Secret:      "sample_app_secret",
	RedirectURI: "http://localhost:5050/redirect",
	PolicyID:    "mypolicy",
	Metadata: map[string]any{
		"userid": "user_from_oauth",
	},
	Description: "my sample oauth app",

	ListenPath: "/oauth",
}

func (app OauthApp) GetOauthUrl(p string) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		TykHTTPEndpoint,
		strings.TrimPrefix(app.ListenPath, "/"),
		strings.TrimPrefix(p, "/"),
	)
}

func (app OauthApp) LoginURL(state, scope string) string {
	oauthURL, err := url.Parse(app.GetOauthUrl("/oauth/authorize"))
	if err != nil {
		log.Fatalf("invalid tyk url: %+v", err)
	}

	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("client_id", app.ClientID)
	values.Set("redirect_uri", app.RedirectURI)
	if state != "" {
		values.Set("state", state)
	}
	if scope != "" {
		values.Set("scope", scope)
	}
	oauthURL.RawQuery = values.Encode()
	return oauthURL.String()
}
