package email

import (
	"crypto/tls"
	"fmt"
	"github.com/blckur/blckur/settings"
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
		return
	}

	err = c.client.Rcpt(to.Address)
	if err != nil {
		return
	}

	w, err := c.client.Data()
	if err != nil {
		return
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return
	}

	err = w.Close()
	if err != nil {
		return
	}

	return
}

func (c *Client) Close() (err error) {
	err = c.client.Quit()
	return
}

func New() (client *Client, err error) {
	client = &Client{
		from: &mail.Address{settings.Email.Name, settings.Email.From},
	}

	host := settings.Email.Server
	addr := fmt.Sprintf("%s:%d", host, settings.Email.Port)

	auth := smtp.PlainAuth("", settings.Email.Username,
		settings.Email.Password, host)

	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	})
	if err != nil {
		return
	}

	client.client, err = smtp.NewClient(conn, host)
	if err != nil {
		return
	}

	err = client.client.Auth(auth)
	if err != nil {
		return
	}

	return
}
