package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/TykTechnologies/tyk/apidef/oas"
	"github.com/TykTechnologies/tyk/config"
	"github.com/TykTechnologies/tyk/gateway"
)

const (
	serverConfPath = "./conf/tyk/tyk.json"
	oldAppConfPath = "./conf/tyk/apps"
	newAppConfPath = "./conf/tyk/newapps"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwConfig := config.Config{}
	if err := config.Load([]string{serverConfPath}, &gwConfig); err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	gw := gateway.NewGateway(gwConfig, ctx)

	loader := gateway.APIDefinitionLoader{Gw: gw}
	specs := loader.FromDir(oldAppConfPath)

	for _, s := range specs {
		based, versions, err := oas.MigrateAndFillOAS(s.APIDefinition)
		if err != nil {
			log.Fatalf("unable to migrate app `%s`: %+v", s.APIID, err)
		}
		log.Printf("app `%s` have %d versions", s.APIID, len(versions))

		if err := writeOasFile(based); err != nil {
			log.Fatalf("unable to write oas file for api `%s`", s.APIID)
		}
		if err := writeClassicFile(based); err != nil {
			log.Fatalf("unable to write classic file for api `%s`", s.APIID)
		}
	}
}

func writeOasFile(api oas.APIDef) error {
	f, err := os.OpenFile(
		filepath.Join(newAppConfPath, fmt.Sprintf("%s-oas.json", api.Classic.APIID)),
		os.O_CREATE|os.O_TRUNC|os.O_RDWR,
		0644,
	)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", " ")
	encoder.Encode(api.OAS)
	return nil
}

func writeClassicFile(api oas.APIDef) error {
	f, err := os.OpenFile(
		filepath.Join(newAppConfPath, fmt.Sprintf("%s.json", api.Classic.APIID)),
		os.O_CREATE|os.O_TRUNC|os.O_RDWR,
		0644,
	)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", " ")
	encoder.Encode(api.Classic)
	return nil
}
