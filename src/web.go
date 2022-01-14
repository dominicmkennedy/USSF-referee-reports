package main

import (
	"encoding/base64"
	"errors"
	"html/template"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/gorilla/schema"
)

type Submission struct {
	PDF           string
	ReporterEmail string
	SendToEmails  []string
	POSTForm      interface{}
}

func WebError(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../static/error.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return
		} else {
			log.Println(err)
		}
	}
}

func SubmitForm(w http.ResponseWriter, r *http.Request) {
	//  parse http POST request
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	//  create a new POSTReport struct
	//  put POST data into struct
	//  then sanitize data
	form := new(POSTReport)
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(time.Now(), DateConverter)
	if err := decoder.Decode(form, r.PostForm); err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	ReportID, err := GetReportID()
	if err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	// I am curentlly of the mind that this function should never toss back an error
	// if internally the form has some malformed data sanitize should first attempt to fix it
	// if that's not posible than the data should be returned to the glang 0 value
	// if referee reports gets a lot of griefers than I may later ammend this statment
	form.SanitizePostData(ReportID)

	//  put POST data into the database
	if err := form.AddToDatabase(); err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	//  create new PDF struct
	//  fill in data from the POSTReport
	//  write the PDF to disk
	//  then store it in cloud
	PDF := new(PDFReport)
	PDF.FillPDF(*form)
	PDFfile, err := PDF.WriteToPDF()
	if err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	if err := SendReport(form, PDFfile); err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	if err := PDF.StorePDF(PDFfile); err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("../static/submitted.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return
		} else {
			log.Println(err)
		}
	}
}

func PostForm(w http.ResponseWriter, r *http.Request) {
	//  parse http POST request
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	//  create a new POSTReport struct
	//  put POST data into struct
	//  then sanitize data
	form := new(POSTReport)
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(time.Now(), DateConverter)
	if err := decoder.Decode(form, r.PostForm); err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}
	form.SanitizePostData("tmp-report")

	//  create new PDF struct
	//  fill in data from the POSTReport
	//  write the PDF to disk
	//  then store it in cloud
	PDF := new(PDFReport)
	PDF.FillPDF(*form)
	PDFfile, err := PDF.WriteToPDF()
	if err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}
	EncodedPDF := base64.StdEncoding.EncodeToString(PDFfile.Bytes())

	SubmissionData := Submission{
		PDF:           EncodedPDF,
		ReporterEmail: form.ReporterEmail,
		SendToEmails:  form.SendToEmail,
		POSTForm:      r.PostForm,
	}

	tmpl := template.Must(template.ParseFiles("../static/review.html"))
	if err := tmpl.Execute(w, SubmissionData); err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return
		} else {
			log.Println(err)
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../static/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return
		} else {
			log.Println(err)
		}
	}
}
