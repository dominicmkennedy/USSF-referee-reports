package main

import (
    //"errors"
    "fmt"
    "net/smtp"
)

var emailAuth smtp.Auth

/*
func parseTemplate(templateFileName string, data interface{}) (string, error) {
   templatePath, err := filepath.Abs(fmt.Sprintf("gomail/email_templates/%s", templateFileName))
   if err != nil {
      return "", errors.New("invalid template name")
   }
   t, err := template.ParseFiles(templatePath)
   if err != nil {
      return "", err
   }
   buf := new(bytes.Buffer)
   if err = t.Execute(buf, data); err != nil {
      return "", err
   }
   body := buf.String()
   return body, nil
}
*/

func SendEmailSMTP(to []string) (bool, error) {
    emailHost := "smtp.gmail.com"
    emailFrom := "automated@referee.report"
    emailPassword := "testing7890!"
    emailPort := 587
    emailAuth = smtp.PlainAuth("", emailFrom, emailPassword, emailHost)
    //emailBody, err := parseTemplate(template, data)
    emailBody := "aaaaaaaaaaaaaaaaaaaaaa"

    mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
    subject := "Subject: " + "Test Email" + "!\n"
    msg := []byte(subject + mime + "\n" + emailBody)
    addr := fmt.Sprintf("%s:%d", emailHost, emailPort)
    if err := smtp.SendMail(addr, emailAuth, emailFrom, to, msg); err != nil {
        return false, err
    }
    return true, nil
}

func notmain() {

    _, err := SendEmailSMTP([]string{"dominicmkennedy@gmail.com"})
    if err != nil {
        fmt.Println(err)
    }

}
