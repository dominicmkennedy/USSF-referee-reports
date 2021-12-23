package main

import (
    "log"

    "gopkg.in/gomail.v2"
)

func SendReport(form *POSTReport) {

    m := gomail.NewMessage()
    m.SetHeader("From", "automated@referee.report")
    m.SetHeader("To", append(form.SendToEmail, form.ReporterEmail)...)
    m.SetHeader("Subject", "New referee report from " + form.ReporterName)
    m.SetBody("text/html", "A referee report was submitted by " + form.ReporterName + " on " + form.SubmittedDate.Format("2006-01-02 15:04:05")+ ". The completed report is attached as a PDF.")
    m.Attach("../reports/" + form.ReportID + ".pdf")

    d := gomail.NewDialer("smtp.gmail.com", 587, "automated@referee.report", "testing7890!")

    if err := d.DialAndSend(m); err != nil {
        log.Println(err)
    }
}
