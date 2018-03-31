package main

import (
	"email-notifier/src/config"
	"github.com/getlantern/systray"
	"email-notifier/src/app"
)

func main() {
	config := config.Load()
	app := app.NewApp(config)
	systray.Run(app.Start, app.Exit)
}