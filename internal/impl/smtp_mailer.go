package impl

import (
	"fmt"
	"go-backend-template/internal/model"
	"net/smtp"
)

type SMTPMailer struct {
	username string
	password string
	host     string
	port     string
	auth     smtp.Auth
}

func NewSMTPMailer(username, password, host, port string) model.Mailer {
	auth := smtp.PlainAuth("", username, password, host)

	return &SMTPMailer{
		username: username,
		password: password,
		host:     host,
		port:     port,
		auth:     auth,
	}
}

func (m *SMTPMailer) Send(b *model.MailBuilder) error {
	receivers, content := b.Build()
	fmt.Println(content)
	return smtp.SendMail(m.host+":"+m.port, m.auth, m.username, receivers, []byte(content))
}
