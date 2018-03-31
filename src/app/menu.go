package app

import (
	"github.com/getlantern/systray"
)

type menu map[string]*systray.MenuItem

func (m *menu) register(name string, item *systray.MenuItem) {
	(*m)[name] = item
}

func (m *menu) notify(name string) <- chan struct{} {
	return (*m)[name].ClickedCh
}