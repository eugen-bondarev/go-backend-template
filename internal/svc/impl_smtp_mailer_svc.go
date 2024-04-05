package svc

import (
	"fmt"
	"net/smtp"
)

type SMTPMailerSvc struct {
	username string
	password string
	host     string
	port     string
	auth     smtp.Auth
}

func NewSMTPMailer(username, password, host, port string) IMailer {
	auth := smtp.PlainAuth("", username, password, host)

	return &SMTPMailerSvc{
		username: username,
		password: password,
		host:     host,
		port:     port,
		auth:     auth,
	}
}

func (m *SMTPMailerSvc) Send(b *MailBuilder) error {
	receivers, content := b.Build()
	fmt.Println(content)
	return smtp.SendMail(m.host+":"+m.port, m.auth, m.username, receivers, []byte(content))
}
