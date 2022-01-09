package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ainsleyclark/go-mail/drivers"
	"github.com/ainsleyclark/go-mail/mail"
)

func SendReport(form *POSTReport, PDFfile *bytes.Buffer) error {
	GoogleWorkspacePassword, err := ioutil.ReadFile(PATH_TO_GOOGLE_WORKSPACE_PASSWORD)
	if err != nil {
		return fmt.Errorf("error could not log into email account: %v", err)
	}

	cfg := mail.Config{
		URL:         "smtp.gmail.com",
		FromAddress: "automated@referee.report",
		FromName:    "automated",
		Password:    string(GoogleWorkspacePassword),
		Port:        587,
	}

	mailer, err := drivers.NewSMTP(cfg)
	if err != nil {
		return fmt.Errorf("error creating a mailer: %v", err)
	}

	SendTo := make([]string, 0)
	if form.ReporterEmail != "" {
		SendTo = append(SendTo, form.ReporterEmail)
	}

	for _, Email := range form.SendToEmail {
		if Email != "" {
			SendTo = append(SendTo, Email)
		}
	}

	if len(SendTo) == 0 {
		log.Printf("report was submitted but no emails were sent\n")
		return nil
	}

	Subject := "New referee report from " + form.ReporterName

	Body := "A referee report was submitted by " + form.ReporterName + " on " + form.SubmittedDate.Format("2006-01-02 15:04:05") + ". The completed report is attached as a PDF."

	tx := &mail.Transmission{
		Recipients: SendTo,
		Subject:    Subject,
		PlainText:  Body,
		Attachments: []mail.Attachment{
			{
				Filename: "report.pdf",
				Bytes:    PDFfile.Bytes(),
			},
		},
	}

	_, err = mailer.Send(tx)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}
