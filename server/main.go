package main

import (
	"os"

	"zyrouge.me/umi/config"
	"zyrouge.me/umi/server"
	"zyrouge.me/umi/utils"
)

func main() {
	if err := start(); err != nil {
		utils.Logger.Error().Err(err).Msg("failed to start application")
		os.Exit(1)
	}
}

func start() error {
	utils.Logger.Info().Msg("starting...")
	_, err := config.GetConfig()
	if err != nil {
		return err
	}
	return server.Start()
}
