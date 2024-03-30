package impl

import (
	"fmt"
	"go-backend-template/internal/model"
	"net/smtp"
)

type SMTPMailerSvc struct {
	username string
	password string
	host     string
	port     string
	auth     smtp.Auth
}

func NewSMTPMailerSvc(username, password, host, port string) model.MailerSvc {
	auth := smtp.PlainAuth("", username, password, host)

	return &SMTPMailerSvc{
		username: username,
		password: password,
		host:     host,
		port:     port,
		auth:     auth,
	}
}

func (m *SMTPMailerSvc) Send(b *model.MailBuilder) error {
	receivers, content := b.Build()
	fmt.Println(content)
	return smtp.SendMail(m.host+":"+m.port, m.auth, m.username, receivers, []byte(content))
}