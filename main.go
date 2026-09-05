package main

import (
	"fmt"
	"os"

	action "github.com/TrafficLight6/GoAccountHub/cliAction"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "account_hub",
		Usage: "account hub example",
		Commands: []*cli.Command{
			action.StartServerAction(),
			action.EditPasswordAction(),
			action.GenerateConfigAction(),
			action.GenerateTestConfigAction(),

			action.AddApplicationKeyAction(),
			action.DeleteApplicationKeyAction(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
