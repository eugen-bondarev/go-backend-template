package svc

import "fmt"

type IMailer interface {
	Send(builder *MailBuilder) error
}

type MailBuilder struct {
	subject   string
	sender    string
	content   string
	receivers []string
}

func NewMailBuilder(receiver string, content string) *MailBuilder {
	return &MailBuilder{
		receivers: []string{receiver},
		content:   content,
	}
}

func (b *MailBuilder) WithSubject(subject string) *MailBuilder {
	b.subject = subject
	return b
}

func (b *MailBuilder) WithSender(name string, email string) *MailBuilder {
	b.sender = fmt.Sprintf("%s: <%s>", name, email)
	return b
}

func (b *MailBuilder) Build() ([]string, string) {
	content := b.content
	if len(b.sender) > 0 {
		content = fmt.Sprintf("From: %s\n%s", b.sender, content)
	}
	if len(b.subject) > 0 {
		content = fmt.Sprintf("Subject: %s\n%s", b.subject, content)
	}
	return b.receivers, content
}
