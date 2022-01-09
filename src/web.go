package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/schema"
)

type Submission struct {
	PDF           string
	ReporterEmail string
	SendToEmails  []string
}

func PostForm(w http.ResponseWriter, r *http.Request) {
	Start := time.Now()
	ParseStart := time.Now()

	//  parse http POST request
	if err := r.ParseForm(); err != nil {
		log.Println(err)
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
	}
	form.SanitizePostData()

	SanTime := time.Since(SanStart)

	DBStart := time.Now()

	//  put POST data into the database
	form.AddToDatabase()

	DBTime := time.Since(DBStart)
	PDFStart := time.Now()

	//  create new PDF struct
	//  fill in data from the POSTReport
	//  write the PDF to disk
	//  then store it in cloud
	PDF := new(PDFReport)
	PDF.FillPDF(*form)
	PDFfile := PDF.WriteToPDF()
	EncodedPDF := base64.StdEncoding.EncodeToString(PDFfile.Bytes())

	PDFTime := time.Since(PDFStart)
	EmailStart := time.Now()

	SendReport(form, PDFfile)

	EmailTime := time.Since(EmailStart)
	PDFStoreStart := time.Now()

	PDF.StorePDF(PDFfile)

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
		log.Println(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../static/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println(err)
	}
}
