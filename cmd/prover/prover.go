package main

import (
	"os"

	"github.com/LimeChain/crc-prover/pkg/app"
	"github.com/LimeChain/crc-prover/pkg/app/configs"
	"github.com/LimeChain/crc-prover/pkg/app/handlers"
	"github.com/LimeChain/crc-prover/pkg/log"
)

func main() {

	config, err := configs.ReadConfigFromFile("prover")
	if err != nil {
		log.Errorw("cannot read prover config storage", err)
		os.Exit(1)
	}

	log.SetLevelStr(config.Log.Level)
	// init handlers for router

	var appHandlers = app.Handlers{
		ZKHandler: handlers.NewZKHandler(config.Prover),
	}
	router := appHandlers.Routes()

	server := app.NewServer(router)

	// start the server
	server.Run(config.Server.Port)

}
