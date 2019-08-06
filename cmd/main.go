package main

import (
	"fmt"
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/ycornamusaz/summon/internal/command"
	"github.com/ycornamusaz/summon/pkg/summon"
)

var (
	CLIArgs   []string  = os.Args
	CLIWriter io.Writer = os.Stdout
)

/* RunCLI starts defines and runs summon's command line interface
 */
func RunCLI() error {
	app := cli.NewApp()
	app.Name = "summon"
	app.Usage = "Parse secrets.yml and export environment variables"
	app.Version = summon.VERSION
	app.Writer = CLIWriter
	app.Flags = command.Flags
	app.Action = command.Action

	return app.Run(CLIArgs)
}

func main() {
	if err := RunCLI(); err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
