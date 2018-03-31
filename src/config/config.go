package config

import (
	"github.com/vrischmann/jsonutil"
)

type Config struct {
	CheckInterval jsonutil.Duration
	SmtpHost string
	SmtpPort uint
	MailUser string
	MailPassword string
}