package checker

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"strconv"
	"log"
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

func (c *Checker) Check() error {
	if err := c.client.Check(); err != nil {
		log.Println(err)
		log.Println("Connection not healthy. Restarting...")

		if err := c.Start(); err != nil {
			log.Println(err)
		}
	}

	unseen, err := c.getUnseen()
	if err != nil {
		return err
	}

	c.unread <- uint(len(unseen))

	var notDone []uint32
	notDone, err = c.getNotDone()
	if err != nil {
		return err
	}
	
	if len(notDone) == 0 {
		return nil
	}

	if err := c.fetchMessages(notDone); err != nil {
		return err
	}

	if err := c.flagDone(notDone); err != nil {
		return err
	}

	return nil
}

func (c *Checker) Start() error {
	c.Stop()

	cl, err := client.DialTLS(c.getSmtpAddress(), nil)
	if err != nil {
		return err
	}

	if err := cl.Login(c.config.User, c.config.Password); err != nil {
		return err
	}

	_, err = cl.Select("INBOX", false)

	if err != nil {
		return err
	}

	c.client = cl

	return nil
}

func (c *Checker) Stop() {
	if c.client == nil {
		return
	}

	c.client.Logout()
	c.client.Close()
}

func (c *Checker) getSmtpAddress() string {
	return c.config.Host + ":" + strconv.Itoa(int(c.config.Port))
}

func (c *Checker) getUnseen() ([]uint32, error) {
	return c.getWithoutFlags([]string{imap.SeenFlag})
}

func (c *Checker) getNotDone() ([]uint32, error) {
	return c.getWithoutFlags([]string{imap.SeenFlag, DoneFlag})
}

func (c *Checker) getWithoutFlags(flags []string) ([]uint32, error) {
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = flags

	unseenIds, err := c.client.UidSearch(criteria)
	if err != nil {
		return []uint32{}, nil
	}

	return unseenIds, nil
}

func (c *Checker) fetchMessages(uids []uint32) error {
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
		return err
	}

	return nil
}

func (c *Checker) flagDone(uids []uint32) error {
	done := make(chan error)
	
	go func() {
		set := new(imap.SeqSet)
		set.AddNum(uids...)
		done <- c.client.UidStore(set, "+FLAGS", DoneFlag, nil)
	}()

	if err := <-done; err != nil {
		return err
	}

	return nil
}