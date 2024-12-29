package main

import (
	"github.com/asrioth/gator/internal/config"
	"github.com/asrioth/gator/internal/database"
)

type State struct {
	ConfigData *config.Config
	Db         *database.Queries
}
