package main

import (
	"email-notifier/src/config"
	"github.com/getlantern/systray"
	"email-notifier/src/app"
	"os"
	"log"
)

func main() {
	f := setLogOutput()
	defer f.Close()

	config := config.Load()
	app := app.NewApp(config)

	systray.Run(app.Start, app.Exit)
}

func setLogOutput() *os.File {
	f, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)

	return f
}