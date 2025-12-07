// Package mailer for sending all mails information
package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	gomail "gopkg.in/mail.v2"
)

func sendEmail(params EmailParams) error {
	cfg, ok := templateRegistry[params.Type]
	if !ok {
		return errors.New()
	}

	data := cfg.BuildData(params)

	tmpl, err := template.New("email").Parse(cfg.File)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return errors.New()
	}

	message := gomail.NewMessage()
	message.SetHeader("from", "")
	message.SetHeader("To", params.Recipient)
	message.SetHeader("Subject", params.Title)
	message.SetBody("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api")
	return nil
}


344614ae31d4a3ddfb47fbac536c88c4



package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "https://send.api.mailtrap.io/api/send"
	method := "POST"

	payload := strings.NewReader(`{\"from\":{\"email\":\"hello@demomailtrap.co\",\"name\":\"Mailtrap Test\"},\"to\":[{\"email\":\"odunayoshittu55@gmail.com\"}],\"subject\":\"You are awesome!\",\"text\":\"Congrats for sending test email with Mailtrap!\",\"category\":\"Integration Test\"}`)

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer <YOUR_API_TOKEN>")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
