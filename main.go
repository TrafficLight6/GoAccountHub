package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "account_hub",
		Usage: "account hub example",
		Commands: []*cli.Command{
			StartServerAction(),
			EditPasswordAction(),
			GenerateConfigAction(),
			GenerateTestConfigAction(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
