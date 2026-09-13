package main

import (
	"fmt"
	"io/fs"
	"os"

	action "github.com/TrafficLight6/GoAccountHub/cliAction"
	"github.com/urfave/cli/v2"
)

func main() {
	//Frontend assets are compiled into the binary (see frontend.go)
	frontendFS, err := fs.Sub(embeddedFrontend, "GAHFrontend/dist")
	if err != nil {
		fmt.Println("⚠️ Error when Read Frontend Assets:", err)
		frontendFS = nil
	}

	app := &cli.App{
		Name:  "account_hub",
		Usage: "account hub example",
		Commands: []*cli.Command{
			action.StartServerAction(frontendFS),
			action.EditPasswordAction(),
			action.GenerateConfigAction(),
			action.GenerateTestConfigAction(),
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
