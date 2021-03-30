package main

import (
	"bytes"
	"crypto/tls"
	"html/template"

	log "github.com/Sirupsen/logrus"

	"net"
	"net/smtp"
	"strconv"
	"time"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func NewRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail() bool {
	useTls := false
	useStartTls := true
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	conn, err := net.Dial("tcp", config.CheckSsl.Notify.Mail.Server+":"+strconv.Itoa(config.CheckSsl.Notify.Mail.Port))
	if err != nil {
		log.Error(err)
		return false
	}

	// TLS
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         config.CheckSsl.Notify.Mail.Server,
	}

	if useTls {
		conn = tls.Client(conn, tlsconfig)
	}

	client, err := smtp.NewClient(conn, config.CheckSsl.Notify.Mail.Server)
	if err != nil {
		log.Error(err)
		return false
	}

	hasStartTLS, _ := client.Extension("STARTTLS")
	if useStartTls && hasStartTLS {
		if config.CheckSsl.Debug {
			log.Info("STARTTLS ...")
		}
		if err = client.StartTLS(tlsconfig); err != nil {
			log.Error(err)
			return false
		}

		if err = client.Hello(config.CheckSsl.Notify.Mail.Server); err != nil {
		}
		if config.CheckSsl.Debug {
			log.Info("HELO done")
		}
		// Set up authentication information.
		auth := smtp.PlainAuth(
			"",
			config.CheckSsl.Notify.Mail.Auth.Login,
			config.CheckSsl.Notify.Mail.Auth.Password,
			config.CheckSsl.Notify.Mail.Server,
		)

		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				log.Error("Error during AUTH ", err)
				return false
			}
		}
		if config.CheckSsl.Debug {
			log.Info("AUTH done")
		}
		if err := client.Mail(config.CheckSsl.Notify.Mail.Auth.Login); err != nil {
			log.Error("Error: ", err)
			return false
		}
		if config.CheckSsl.Debug {
			log.Info("FROM done")
		}
		if err := client.Rcpt(r.to[0]); err != nil {
			log.Error("Error: ", err)
			return false
		}
		if config.CheckSsl.Debug {
			log.Info("TO done")
		}
		w, err := client.Data()
		if err != nil {
			log.Error("Error: ", err)
			return false
		}

		_, err = w.Write([]byte(body))
		if err != nil {
			log.Error("Error: ", err)
			return false
		}

		err = w.Close()
		if err != nil {
			log.Error("Error: ", err)
			return false
		}

		client.Quit()
	}
	return true
}
func (r *Request) Send(templateName string, items interface{}) {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Error(err)
	}
	if ok := r.sendMail(); ok {
		if config.CheckSsl.Debug {
			log.Info("Email has been sent to ", r.to)
		}
	} else {
		if config.CheckSsl.Debug {
			log.Error("Failed to send the email to ", r.to)
		}
	}
}

func SendMailMsg(subject string, email string) time.Duration {
	start := time.Now()
	receiver := email
	r := NewRequest([]string{receiver}, subject)
	r.Send("templates/mail/mail.html", map[string]string{"username": ""})
	return time.Since(start)
}
