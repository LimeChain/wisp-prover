package main

import (
	"github.com/LimeChain/crc-prover/pkg/app"
	"github.com/LimeChain/crc-prover/pkg/app/configs"
	"github.com/LimeChain/crc-prover/pkg/app/handlers"
	"github.com/LimeChain/crc-prover/pkg/log"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("cannot read prover config storage", err)
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
