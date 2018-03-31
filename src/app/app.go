package app

import (
	"email-notifier/src/config"
	"email-notifier/src/checker"
	"github.com/getlantern/systray"
	"time"
	"github.com/skratchdot/open-golang/open"
)

type App struct {
	config		config.Config
	menu        menu
	checker     *checker.Checker
	unread      chan uint
	unreadCount	uint
}

func NewApp(config config.Config) App {
	checker := checker.NewChecker(config.MailUser, config.MailPassword, config.SmtpHost, config.SmtpPort)

	return App{config, make(menu), checker, make(chan uint), 0}
}

func (a *App) Start() {
	a.menu.initialise()
	a.updateStatus(0)

	go a.menuHandler()
	go a.checkHandler()
	go a.statusUpdater()
}

func (a *App) Exit() {}

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
	for {
		<-time.Tick(a.config.CheckInterval.Duration)
		a.unread <- a.checker.CheckUnreadCount()
	}
}