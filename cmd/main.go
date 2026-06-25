package main

import (
	"os"

	faust "github.com/crazytaxii/faust/cmd/app"
	"github.com/crazytaxii/faust/pkg/signals"
)

var AppVersion string

func main() {
	app := faust.NewFaustApp(AppVersion)
	if err := app.Run(signals.SetupSignalHandler(), os.Args); err != nil {
		os.Exit(1)
	}
}
