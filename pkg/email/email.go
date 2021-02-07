package email

import (
	"strconv"

	gomail "gopkg.in/gomail.v2"
)

type EmailClient struct {
	from         string
	smtpUser     string
	smtpPassword string
	smtpHost     string
	smtpPort     int
}

func NewEmailClient(from, smtpHost, smtpPort, smtpUser, smtpPassword string) (*EmailClient, error) {
	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return nil, err
	}
	return &EmailClient{
		from:         from,
		smtpUser:     smtpUser,
		smtpPassword: smtpPassword,
		smtpHost:     smtpHost,
		smtpPort:     port,
	}, nil
}

func (c EmailClient) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetBody("text/plain", body)
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(c.from, "admin")},
		"To":      {to},
		"Subject": {subject},
	})

	d := gomail.NewPlainDialer(c.smtpHost, c.smtpPort, c.smtpUser, c.smtpPassword)
	return d.DialAndSend(m)
}
