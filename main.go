package main

import (
	"github.com/getlantern/systray"
	"email-notifier/src/app"
)

func main() {
	app := app.NewApp()
	systray.Run(app.Start, app.Exit)
}