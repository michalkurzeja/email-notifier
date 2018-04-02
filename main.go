package main

import (
	// "os"
	// "log"
	"email-notifier/src/config"
	"github.com/getlantern/systray"
	"email-notifier/src/app"
)

func main() {
	// f, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()

	// log.SetOutput(f)

	config := config.Load()
	app := app.NewApp(config)
	systray.Run(app.Start, app.Exit)
}