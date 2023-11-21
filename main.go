package main

import (
	"github.com/syobonaction/fur_lange/tools"
	"github.com/turret-io/go-menu/menu"
)

func main() {
	commandOptions := []menu.CommandOption{
		{Command: "Server", Description: "Runs the mongo server.", Function: tools.RunServer},
		{Command: "Migrate", Description: "Migrate data to pgsql", Function: tools.MigratePgsql},
		{Command: "Collect", Description: "Collect data from APN", Function: tools.Collect},
	}

	menuOptions := menu.NewMenuOptions("'menu' for help > ", 0)

	menu := menu.NewMenu(commandOptions, menuOptions)
	menu.Start()
}
