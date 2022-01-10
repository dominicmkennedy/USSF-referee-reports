package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gorilla/schema"
)

type Submission struct {
	PDF           string
	ReporterEmail string
	SendToEmails  []string
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

func PostForm(w http.ResponseWriter, r *http.Request) {
	Start := time.Now()
	ParseStart := time.Now()

	//  parse http POST request
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	ParseTime := time.Since(ParseStart)
	SanStart := time.Now()

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

	SanTime := time.Since(SanStart)
	DBStart := time.Now()

	//  put POST data into the database
	if err := form.AddToDatabase(); err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	DBTime := time.Since(DBStart)
	PDFStart := time.Now()

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

	PDFTime := time.Since(PDFStart)
	EmailStart := time.Now()

	if err := SendReport(form, PDFfile); err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	EmailTime := time.Since(EmailStart)
	PDFStoreStart := time.Now()

	if err := PDF.StorePDF(PDFfile); err != nil {
		log.Println(err)
		WebError(w, r)
		return
	}

	PDFStoreTime := time.Since(PDFStoreStart)
	Elapsed := time.Since(Start)

	fmt.Fprintf(os.Stdout, "HTML parsing took %v\n", ParseTime)
	fmt.Fprintf(os.Stdout, "Form Sanitzation took %v\n", SanTime)
	fmt.Fprintf(os.Stdout, "Adding to the Database took %v\n", DBTime)
	fmt.Fprintf(os.Stdout, "Creating the PDF took %v\n", PDFTime)
	fmt.Fprintf(os.Stdout, "Storing the PDF took %v\n", PDFStoreTime)
	fmt.Fprintf(os.Stdout, "Emailing the reports took %v\n", EmailTime)
	fmt.Fprintf(os.Stdout, "Total elapsed time was %v\n", Elapsed)

	SubmissionData := Submission{
		PDF:           EncodedPDF,
		ReporterEmail: form.ReporterEmail,
		SendToEmails:  form.SendToEmail,
	}

	tmpl := template.Must(template.ParseFiles("../static/submitted.html"))
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
