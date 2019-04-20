package main

import (
	"net/http"

	"github.com/rakyll/statik/fs"
	log "github.com/sirupsen/logrus"
)

func swaggerUI() (http.Handler, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	return http.FileServer(statikFS), nil
}

func serveSwaggerUI(port int) {
	log.Infof("Starting SwaggerUI on port %d\n", port)
	sh, err := swaggerUI()
	if err != nil {
		log.Error(err)
		return
	}

	http.Handle("/swaggerui/", http.StripPrefix("/swaggerui/", sh))
	http.ListenAndServe(":9090", nil)
}
