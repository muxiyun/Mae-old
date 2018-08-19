package mail

import (
	"log"
	"time"
	"bytes"
	"io/ioutil"
	"html/template"
	"gopkg.in/gomail.v2"
	"github.com/spf13/viper"
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
	Ms.Ch = make(chan *gomail.Message, viper.GetInt("mail.chanCache"))

	// init Msg
	Ms.Msg = gomail.NewMessage()
	Ms.Msg.SetHeader("From", Ms.Msg.FormatAddress(viper.GetString("mail.username"), viper.GetString("mail.senderNickName")))
	Ms.Msg.SetHeader("Subject", "Notification from Mae")

	notification, err := ioutil.ReadFile("./conf/mailTemplates/notification.tpl")
	if err != nil {
		log.Fatal("error occurred while read from notification.tpl")
	}
	confirm, err := ioutil.ReadFile("./conf/mailTemplates/confirm.tpl")
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


// start the mail service daemon with a goroutine in your main function
// like the following
//
//   go mail.StartMailDaemon()
func StartMailDaemon() {

	d := gomail.NewDialer(viper.GetString("mail.host"), viper.GetInt("mail.port"),
		viper.GetString("mail.username"), viper.GetString("mail.password"))

	var s gomail.SendCloser
	var err error
	open := false
	for {
		select {
		case m, ok := <-Ms.Ch:
			if !ok {
				return
			}
			if !open {
				if s, err = d.Dial(); err != nil {
					panic(err)
				}
				open = true
			}
			if err := gomail.Send(s, m); err != nil {
				log.Println(err.Error())
			}
			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.

		case <-time.After(time.Duration(viper.GetInt("mail.maxFreeTime")) * time.Second):
			if open {
				if err := s.Close(); err != nil {
					panic(err)
				}
				open = false
			}
		}
	}
	close(Ms.Ch)
}