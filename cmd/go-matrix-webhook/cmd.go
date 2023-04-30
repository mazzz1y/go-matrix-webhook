package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mazzz1y/go-matrix-webhook/internal/handler"
	"github.com/mazzz1y/go-matrix-webhook/internal/matrix"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func run(c *cli.Context) error {
	listenAddr := c.String("listen-addr")
	listenPort := c.Int("listen-port")
	listenPath := c.String("listen-path")
	secretHeader := c.String("secret-header")
	mAccessToken := c.String("matrix-access-token")
	mUserId := c.String("matrix-id")
	mUrl := c.String("matrix-url")
	logLevel := c.String("log-level")
	logType := c.String("log-type")

	setLogLevel(logLevel, logType)

	m, err := matrix.NewMatrix(mUrl, mUserId, mAccessToken)
	if err != nil {
		panic(err)
	}
	h := handler.NewHandler(*m, secretHeader)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(listenPath, h).Methods("POST")

	listen := fmt.Sprintf("%s:%d", listenAddr, listenPort)
	log.Info().Err(err).Msgf("listen on: %s", listen)
	return http.ListenAndServe(listen, router)
}
