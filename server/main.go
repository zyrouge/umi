package main

import (
	"os"

	"zyrouge.me/umi/application"
	"zyrouge.me/umi/database"
	"zyrouge.me/umi/events_live"
	"zyrouge.me/umi/migrations"
	"zyrouge.me/umi/scheduler"
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
	utils.Logger.Info().Msg("starting application...")
	if err := application.LoadConfig(); err != nil {
		return err
	}
	if err := database.LoadConnection(); err != nil {
		return err
	}
	if err := migrations.MigrateDatabase(); err != nil {
		return err
	}
	if err := scheduler.StartSchedulers(); err != nil {
		return err
	}
	if err := events_live.StartWebsocketManager(); err != nil {
		return err
	}
	return server.StartServer()
}
