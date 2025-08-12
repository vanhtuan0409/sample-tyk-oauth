package tykloadtest

import (
	"fmt"
	"strings"
)

const (
	APIID           = "myapp_oauth"
	ClientID        = "sample_oauth_app"
	ClientSecret    = "sample_app_secret"
	RedirectURI     = "http://localhost:5050/redirect"
	ListenPath      = "/api"
	TykHTTPEndpoint = "http://api.gondor-local.io:8080"
)

func GetTykURL(p string) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		TykHTTPEndpoint,
		strings.TrimPrefix(ListenPath, "/"),
		strings.TrimPrefix(p, "/"),
	)
}
