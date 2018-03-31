package app

import (
	"fmt"
	"time"
	"github.com/getlantern/systray"
	"email-notifier/src/assets"
	"email-notifier/src/checker"
)

type App struct {
	menu menu
	checker *checker.Checker
	unread chan int
}

func NewApp() App {
	return App{make(menu), checker.NewChecker(), make(chan int)}
}

func (a *App) Start() {
	a.setUp()
}

func (a *App) Exit() {}

func (a *App) setUp() {
	systray.SetIcon(assets.Get("gmail-red.ico"))
	setUnread(0)

	a.menu.register("quit", systray.AddMenuItem("Wyjdź", "Zakończ aplikację"))

	go a.menuHandler()
	go a.checkHandler()
	go a.statusUpdater()
}

func (a *App) menuHandler() {
	for {
		select {
		case <- a.menu.notify("quit"):
			systray.Quit()
		}
	}
}

func (a *App) checkHandler() {
	for {
		<- time.Tick(5 * time.Second)
		a.unread <- a.checker.CheckUnreadCount()
	}
}

func (a *App) statusUpdater() {
	unread := <- a.unread


}

func setUnread(unread uint) {
	systray.SetTitle(string(unread))	
	systray.SetTooltip(fmt.Sprintf("Nieprzeczytane: %d", unread))
}