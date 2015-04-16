package mail

import (
	"fmt"
	"github.com/blckur/blckur/settings"
	"net/smtp"
)

func SendMail(to string, subject string, body string) (err error) {
	auth := smtp.PlainAuth("", settings.Email.Username,
		settings.Email.Password, settings.Email.Server)

	msg := fmt.Sprintf("To: %s\r\nSubject%s\r\n\r\n%s", to, subject, body)
	addr := fmt.Sprintf("%s:%d", settings.Email.Server, settings.Email.Port)

	err = smtp.SendMail(addr, auth, settings.Email.From, []string{to},
		[]byte(msg))

	return
}
