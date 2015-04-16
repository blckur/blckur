package email

import (
	"crypto/tls"
	"fmt"
	"github.com/blckur/blckur/settings"
	"github.com/dropbox/godropbox/errors"
	"net/mail"
	"net/smtp"
)

type Client struct {
	from   *mail.Address
	client *smtp.Client
}

func (c *Client) Send(to *mail.Address, subject string,
	body string) (err error) {
	headers := map[string]string{}
	headers["From"] = c.from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	msg := ""
	for key, val := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", key, val)
	}
	msg += "\r\n" + body

	err = c.client.Mail(c.from.Address)
	if err != nil {
		err = &SmtpError{
			errors.Wrap(err, "Mail command error"),
		}
		return
	}

	err = c.client.Rcpt(to.Address)
	if err != nil {
		err = &SmtpError{
			errors.Wrap(err, "Rcpt command error"),
		}
		return
	}

	w, err := c.client.Data()
	if err != nil {
		err = &SmtpError{
			errors.Wrap(err, "Data command error"),
		}
		return
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		err = &SmtpError{
			errors.Wrap(err, "Write error"),
		}
		return
	}

	err = w.Close()
	if err != nil {
		err = &SmtpError{
			errors.Wrap(err, "Write close error"),
		}
		return
	}

	return
}

func (c *Client) Close() (err error) {
	err = c.client.Quit()
	if err != nil {
		err = &SmtpError{
			errors.Wrap(err, "Quit command error"),
		}
		return
	}
	return
}

func New() (client *Client, err error) {
	host := settings.Email.Server
	addr := fmt.Sprintf("%s:%d", host, settings.Email.Port)

	auth := smtp.PlainAuth("", settings.Email.Username,
		settings.Email.Password, host)

	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	})
	if err != nil {
		err = &TlsError{
			errors.Wrap(err, "Dial error"),
		}
		return
	}

	smtpClient, err := smtp.NewClient(conn, host)
	if err != nil {
		err = &SmtpError{
			errors.Wrap(err, "Connection error"),
		}
		return
	}

	err = smtpClient.Auth(auth)
	if err != nil {
		err = &SmtpError{
			errors.Wrap(err, "Auth error"),
		}
		return
	}

	client = &Client{
		from:   &mail.Address{settings.Email.Name, settings.Email.From},
		client: smtpClient,
	}

	return
}
