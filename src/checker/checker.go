package checker

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
	"strconv"
)

const DoneFlag = "mk_chk_done"

type Checker struct {
	client   *client.Client
	config   Config
	messages chan<- imap.Message
	unread   chan<- uint
}

func NewChecker(config Config) (*Checker, <-chan imap.Message, <-chan uint) {
	messages := make(chan imap.Message)
	unread := make(chan uint)

	return &Checker{nil, config, messages, unread}, messages, unread
}

func (c *Checker) Check() {
	_, err := c.client.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	c.unread <- uint(len(c.getUnseen()))
	
	notDone := c.getNotDone()
	
	if len(notDone) == 0 {
		return
	}

	c.fetchMessages(notDone)
	c.flagDone(notDone)
}

func (c *Checker) Start() {
	cl, err := client.DialTLS(c.getSmtpAddress(), nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := cl.Login(c.config.User, c.config.Password); err != nil {
		log.Fatal(err)
	}

	c.client = cl
}

func (c *Checker) Stop() {
	c.client.Logout()
	c.client.Close()
}

func (c *Checker) getSmtpAddress() string {
	return c.config.Host + ":" + strconv.Itoa(int(c.config.Port))
}

func (c *Checker) getUnseen() []uint32 {
	return c.getWithoutFlags([]string{imap.SeenFlag})
}

func (c *Checker) getNotDone() []uint32 {
	return c.getWithoutFlags([]string{imap.SeenFlag, DoneFlag})
}

func (c *Checker) getWithoutFlags(flags []string) []uint32 {
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = flags

	unseenIds, err := c.client.UidSearch(criteria)
	if err != nil {
		log.Fatal(err)
	}

	return unseenIds
}

func (c *Checker) fetchMessages(uids []uint32) {
	messages := make(chan *imap.Message, 10)
	done := make(chan error)
	
	set := new(imap.SeqSet)
	set.AddNum(uids...)

	go func() {
		done <- c.client.UidFetch(set, []string{"ENVELOPE"}, messages)
	}()

	go func() {
		for msg := range messages {
			c.messages <- *msg
		}
	}()

	if err := <-done; err != nil {
		log.Fatal(err)
	}
}

// TODO make synchronous
func (c *Checker) flagDone(uids []uint32) {
	done := make(chan error)
	
	go func() {
		set := new(imap.SeqSet)
		set.AddNum(uids...)
		done <- c.client.UidStore(set, "+FLAGS", DoneFlag, nil)
	}()

	if err := <-done; err != nil {
		log.Fatal(err)
	}
}