package app

import (
	// "io/ioutil"
	// "log"
	"github.com/emersion/go-imap"
	"github.com/faiface/beep"
	"os"
	"github.com/faiface/beep/mp3"
	"time"
	"github.com/faiface/beep/speaker"
	"fmt"
	"email-notifier/src/assets"
	"github.com/getlantern/systray"
	"github.com/michalkurzeja/notificator"
)

func (a *App) statusUpdater() {
	for {
		select {
		case unread := <- a.unread:
			a.updateStatus(unread)
		case message := <-a.messages:
			a.notify(message) 
		}
	}
}

func (a *App) updateStatus(unread uint) {
	systray.SetIcon(assets.Get(selectIcon(unread)))
	systray.SetTitle(string(unread))
	systray.SetTooltip(getUnreadTitle(unread))
	a.menu.get("unread").SetTitle(getUnreadTitle(unread))
}

func (a *App) notify(message imap.Message) {
	// go a.pushNotification(message)
	go a.playSound()
}

func (a *App) pushNotification(message imap.Message) {
	// sectionName, _ := imap.NewBodySectionName("BODY[1.1]")
	// bodyText := message.Body[sectionName]
	// body := make([]byte, 50)

	// body, err := ioutil.ReadAll(bodyText)
	// // read, err := bodyText.Read(body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(body)

	a.notifier.Push(message.Envelope.Subject, message.Envelope.Subject, "", notificator.UR_NORMAL)
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
	if unread == 0 {
		return "Brak wiadomości"
	}
	return fmt.Sprintf("Masz wiadomość: %d", unread)
}