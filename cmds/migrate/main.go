package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/TykTechnologies/tyk/apidef/oas"
	"github.com/TykTechnologies/tyk/config"
	"github.com/TykTechnologies/tyk/gateway"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwConfig := config.Config{}
	if err := config.Load([]string{"./conf/tyk/tyk.json"}, &gwConfig); err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	gw := gateway.NewGateway(gwConfig, ctx)

	loader := gateway.APIDefinitionLoader{Gw: gw}
	specs := loader.FromDir("./conf/tyk/apps")

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")

	for _, s := range specs {
		_, versions, err := oas.MigrateAndFillOAS(s.APIDefinition)
		if err != nil {
			log.Fatalf("unable to migrate app `%s`: %+v", s.APIID, err)
		}
		log.Printf("app `%s` have %d versions", s.APIID, len(versions))
	}
}
