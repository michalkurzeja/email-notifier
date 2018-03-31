package checker

import (
	"strconv"
	"log"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap"
)

type Checker struct {
	user string
	password string
	host string
	port uint
}

func NewChecker(user string, password string, host string, port uint) *Checker {
	return &Checker{user, password, host, port}
}

func (c *Checker) CheckUnreadCount() uint {
	cl := c.dial()
	defer cl.Close()

	defer cl.Logout()
	if err := cl.Login(c.user, c.password); err != nil {
		log.Fatal(err)
	}
	
	c.selectInbox(cl)
	return c.getUnseenCount(cl)
}

func (c *Checker) dial() *client.Client {
	cl, err := client.DialTLS(c.getSmtpAddress(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return cl
}

func (c *Checker) getSmtpAddress() string {
	return c.host + ":" + strconv.Itoa(int(c.port))
}

func (c *Checker) selectInbox(cl *client.Client) {
	_, err := cl.Select("INBOX", true)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Checker) getUnseenCount(cl *client.Client) uint {
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	
	unseenIds, err := cl.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	return uint(len(unseenIds))
}