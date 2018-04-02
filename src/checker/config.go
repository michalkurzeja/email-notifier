package checker

import (
	"email-notifier/src/config"
)

type Config struct {
	User string
	Password string
	Host string
	Port uint
}

func FromGlobalConfig(globalConfig config.Config) Config {
	return Config{globalConfig.MailUser, globalConfig.MailPassword, globalConfig.SMTPHost, globalConfig.SMTPPort}
}