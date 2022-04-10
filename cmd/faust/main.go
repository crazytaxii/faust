package main

import (
	"os"

	faust "github.com/crazytaxii/faust/cmd/faust/app"
)

var AppVersion string

func main() {
	app := faust.NewFaustApp(AppVersion)
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
