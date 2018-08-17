package mail

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
	"io/ioutil"
	"log"
)

type NotificationEvent struct {
	Level    string
	UserName string
	Who      string
	Action   string
	What     string
	When     string
}

type ConfirmEvent struct {
	UserName    string
	ConfirmLink string
}

type MailService struct {
	Ch               chan *gomail.Message
	Msg              *gomail.Message
	NotificationTmpl *template.Template
	ConfirmTmpl      *template.Template
}

var Ms MailService

func Setup() {
	// init Ch
	Ms.Ch = make(chan *gomail.Message, 20)

	// init Msg
	Ms.Msg = gomail.NewMessage()
	Ms.Msg.SetHeader("From", Ms.Msg.FormatAddress("3480437308@qq.com", "Mae Notification Robot"))
	//Ms.Msg.SetAddressHeader("Cc", "3480437308@qq.com", "Andrewpqc")
	Ms.Msg.SetHeader("Subject", "Notification from Mae")

	notification, err := ioutil.ReadFile("./pkg/mail/notification.tpl")
	if err != nil {
		log.Fatal("error occurred while read from notification.tpl")
	}
	confirm, err := ioutil.ReadFile("./pkg/mail/confirm.tpl")
	if err != nil {
		log.Fatal("error occurred while read from confirm.tpl")
	}

	// init Tmpl
	Ms.NotificationTmpl, err = template.New("notification").Parse(string(notification))
	if err != nil {
		log.Fatal("error occurred while parse notification template")
	}
	Ms.ConfirmTmpl, err = template.New("confirm").Parse(string(confirm))
	if err != nil {
		log.Fatal("error occurred while parse confirm template")
	}

}

// send notification emails to all admin users
func SendNotificationEmail(e NotificationEvent, recipients []string) {
	var tpl bytes.Buffer
	Ms.NotificationTmpl.Execute(&tpl, e)

	Ms.Msg.SetHeaders(map[string][]string{
		"To": recipients,
	})
	Ms.Msg.SetBody("text/html", tpl.String())

	Ms.Ch <- Ms.Msg
}

// send confirm email to register user
func SendConfirmEmail(ce ConfirmEvent, recipient string) {
	var tpl bytes.Buffer
	Ms.ConfirmTmpl.Execute(&tpl, ce)

	Ms.Msg.SetHeader("To", recipient)

	Ms.Msg.SetBody("text/html", tpl.String())

	Ms.Ch <- Ms.Msg
}
