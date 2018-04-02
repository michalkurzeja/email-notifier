package app

import (
	"email-notifier/src/assets"
	"email-notifier/src/checker"
	"email-notifier/src/config"
	"log"
	"time"
	"github.com/michalkurzeja/notificator"
	"github.com/emersion/go-imap"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

type App struct {
	config   config.Config
	menu     menu
	checker  *checker.Checker
	messages <-chan imap.Message
	unread   <-chan uint
	notifier *notificator.Notificator
}

func NewApp(config config.Config) App {
	checker, messages, unread := checker.NewChecker(checker.FromGlobalConfig(config))
	notifier := notificator.New(notificator.Options{
		DefaultIcon: assets.GetPath(assets.IconUnread),
		AppName:     "Powiadamiacz",
	})

	return App{config, make(menu), checker, messages, unread, notifier}
}

func (a *App) Start() {
	log.Println("Starting up!")

	a.menu.initialise()
	a.updateStatus(0)

	go a.menuHandler()
	go a.checkHandler()
	go a.statusUpdater()
}

func (a *App) Exit() {
	log.Println("Stopping!")

	a.checker.Stop()
}

func (a *App) menuHandler() {
	for {
		select {
		case <-a.menu.notify("unread"):
			open.Start("https://mail.google.com")
		case <-a.menu.notify("quit"):
			systray.Quit()
		}
	}
}

func (a *App) checkHandler() {
	a.checker.Start()

	for {
		<-time.Tick(a.config.CheckInterval.Duration)
		a.checker.Check()
	}
}
