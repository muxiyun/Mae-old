package mail


import (
	"bytes"
	"io/ioutil"
	"html/template"
	"gopkg.in/gomail.v2"
)


type Event struct {
	Level string
	UserName string
	Who string
	Action string
	What string
	When string
}

type MailService struct{
	Ch chan *gomail.Message
	Msg *gomail.Message
	Tmpl *template.Template

}

var Ms MailService

func  Setup(){
	// init Ch
	Ms.Ch=make(chan *gomail.Message,20)

	// init Msg
	Ms.Msg=gomail.NewMessage()
	Ms.Msg.SetHeader("From", Ms.Msg.FormatAddress("3480437308@qq.com","Mae Notification Robot"))
	Ms.Msg.SetAddressHeader("Cc", "3480437308@qq.com", "Andrewpqc")
	Ms.Msg.SetHeader("Subject", "Notification from Mae")

	content,err:=ioutil.ReadFile("./pkg/mail/notification.tpl")
	if err!=nil{

	}
	// init Tmpl
	Ms.Tmpl,_= template.New("notification").Parse(string(content))

}

func SendEmail(e Event,recipients []string) {
	var tpl bytes.Buffer
	Ms.Tmpl.Execute(&tpl, e)


	Ms.Msg.SetHeaders(map[string][]string{
		"To":      recipients,
	})
	Ms.Msg.SetBody("text/html", tpl.String())

	Ms.Ch <- Ms.Msg
}



