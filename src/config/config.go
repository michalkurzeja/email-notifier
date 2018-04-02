package config

import (
	"github.com/vrischmann/jsonutil"
)

type Config struct {
	CheckInterval jsonutil.Duration
	SMTPHost      string
	SMTPPort      uint
	MailUser      string
	MailPassword  string
}