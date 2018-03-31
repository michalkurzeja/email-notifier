package app

import (
	"github.com/faiface/beep"
	"os"
	"github.com/faiface/beep/mp3"
	"time"
	"github.com/faiface/beep/speaker"
	"fmt"
	"email-notifier/src/assets"
	"github.com/getlantern/systray"
	"github.com/gen2brain/beeep"
)

func (a *App) statusUpdater() {
	for {
		a.updateStatus(<-a.unread)
	}
}

func (a *App) updateStatus(unread uint) {
	systray.SetIcon(assets.Get(selectIcon(unread)))
	systray.SetTitle(string(unread))
	systray.SetTooltip(getUnreadTitle(unread))
	a.menu.get("unread").SetTitle(getUnreadTitle(unread))

	if a.shouldNotify(unread) {
		a.notify(unread)
	}

	a.unreadCount = unread	
}

func (a *App) shouldNotify(unread uint) bool {
	return unread > a.unreadCount
}

func (a *App) notify(unread uint) {
	go a.pushNotification(unread)
	go a.playSound()
}

func (a *App) pushNotification(unread uint) {
	beeep.Notify("Nowy email!", getUnreadTitle(unread), assets.GetAbsolutePath(assets.IconUnread))
}

func (a *App) playSound() {
	f, _ := os.Open(assets.GetPath("notification.mp3"))
	defer f.Close()
	
	stream, format, _ := mp3.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	
	play(stream)
}

func play(stream beep.StreamSeekCloser) {
	done := make(chan struct{})
	speaker.Play(beep.Seq(stream, beep.Callback(func() {
        close(done)
	})))
	<-done
}

func selectIcon(unread uint) string {
	if unread > 0 {
		return assets.IconUnread
	}

	return assets.IconAllRead
}

func getUnreadTitle(unread uint) string {
	return fmt.Sprintf("Nieprzeczytane wiadomości: %d", unread)
}