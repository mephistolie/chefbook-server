package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
)

const (
	signature = "<hr>Sented by Broccy the Broccoli from ChefBook 🥦"
)

type MailService struct {
	host      string
	port      string
	sender    string
	password  string
	tlsConfig *tls.Config
}

func NewMailService(host, port, sender, password string) *MailService {
	return &MailService{
		host:     host,
		port:     port,
		sender:   sender,
		password: password,
		tlsConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		},
	}
}

func (m *MailService) SendEmailVerificationCode(code int, recipient string) error {

	fmt.Print("1")

	from := mail.Address{Name: "Broccy from ChefBook", Address: m.sender}
	to := mail.Address{Name: "", Address: recipient}
	subject := "ChefBook Verification Code"
	body := "<html><body><h1 style=\"text-align: center;\">" + fmt.Sprintf("%d", code) + "</h1>" +
		"<p style=\"text-align: center;\">is your ChefBook Verification Code</p></body></html>"

	fmt.Print("2")
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\";"
	headers["Subject"] = subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body + "\r\n" + signature

	auth := smtp.PlainAuth("", from.Address, m.password, m.host)
	fmt.Print("3")

	conn, err := tls.Dial("tcp", m.host+":"+m.port, m.tlsConfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, m.host)
	if err != nil {
		return err
	}

	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return err
}
