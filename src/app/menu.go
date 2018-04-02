package app

import (
	"github.com/getlantern/systray"
)

type menu map[string]*systray.MenuItem

func (m *menu) initialise() {
	m.register("unread", systray.AddMenuItem(getUnreadTitle(0), "Nieprzeczytane"))	
	m.register("refresh", systray.AddMenuItem("Odśwież", "Sprawdź nowe wiadomości"))	
	systray.AddSeparator()
	m.register("quit", systray.AddMenuItem("Wyjdź", "Zakończ aplikację"))
}

func (m *menu) register(name string, item *systray.MenuItem) {
	(*m)[name] = item
}

func (m *menu) get(name string) *systray.MenuItem {
	return (*m)[name]
}

func (m *menu) notify(name string) <- chan struct{} {
	return m.get(name).ClickedCh
}