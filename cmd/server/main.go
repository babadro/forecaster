package main

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/babadro/forecaster/internal/infra/restapi"
	"github.com/babadro/forecaster/internal/infra/restapi/operations"
)

func main() {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewPollAPIAPI(swaggerSpec)

	server := restapi.NewServer(api)
	defer func(server *restapi.Server) {
		if err = server.Shutdown(); err != nil {
			log.Printf("error while shutting down server: %v", err)
		}
	}(server)

	server.ConfigureAPI()

	if err = server.Serve(); err != nil {
		log.Printf("error while serving: %v", err)
	}
}
