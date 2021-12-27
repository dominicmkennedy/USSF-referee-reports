package main

import (
    "log"
    "io/ioutil"

    "gopkg.in/gomail.v2"
)

func SendReport(form *POSTReport) {

    GoogleWorkspacePassword, err := ioutil.ReadFile(PATH_TO_GOOGLE_WORKSPACE_PASSWORD)
    if err != nil {
        log.Panicln(err)
    }

    m := gomail.NewMessage()
    m.SetHeader("From", "automated@referee.report")

    SendTo := make([]string, 0)
    
    if form.ReporterEmail != "" {
        SendTo = append(SendTo, form.ReporterEmail)
    }
    
    for _, Email := range form.SendToEmail {
        if Email != "" {
            SendTo = append(SendTo, Email) 
        }
    }
    
    if len(SendTo) == 0 { return }
    
    m.SetHeader("To", SendTo...)
    m.SetHeader("Subject", "New referee report from " + form.ReporterName)
    m.SetBody("text/html", "A referee report was submitted by " + form.ReporterName + " on " + form.SubmittedDate.Format("2006-01-02 15:04:05") + ". The completed report is attached as a PDF.")
    m.Attach("../reports/" + form.ReportID + ".pdf")

    d := gomail.NewDialer("smtp.gmail.com", 587, "automated@referee.report", string(GoogleWorkspacePassword))

    if err := d.DialAndSend(m); err != nil {
        log.Println(err)
    }
}
