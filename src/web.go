package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/schema"
)

func PostForm(w http.ResponseWriter, r *http.Request) {
	//  parse http POST request
	if err := r.ParseForm(); err != nil {
		log.Printf("error parsing http POST request: %v", err)
		http.Redirect(w, r, "/error.html", http.StatusFound)
		return
	}

	//  create a new POSTReport struct
	//  put POST data into struct
	//  then sanitize data
	form := new(POSTReport)
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(time.Now(), DateConverter)
	if err := decoder.Decode(form, r.PostForm); err != nil {
		log.Printf("error decoding POST data into form struct: %v", err)
		http.Redirect(w, r, "/error.html", http.StatusFound)
		return
	}

	if err := CAPTCHA.Verify(form.RecaptchaResponse); err != nil {
		log.Printf("error verifying captcha: %v", err)
		http.Redirect(w, r, "/error.html", http.StatusFound)
		return
	}

	// I am curentlly of the mind that this function should never toss back an error
	// if internally the form has some malformed data sanitize should first attempt to fix it
	// if that's not posible than the data should be returned to the glang 0 value
	// if referee reports gets a lot of griefers than I may later ammend this statment

	// located in report.go
	form.SanitizePostData()

	//  put POST data into the database
	// located in database.go
	if err := form.AddToDatabase(); err != nil {
		log.Printf("error adding form to database: %v", err)
		http.Redirect(w, r, "/error.html", http.StatusFound)
		return
	}

	//  create new PDF struct
	//  fill in data from the POSTReport
	//  write the PDF to RAM
	//  then store it in cloud
	//  located in pdf.go
	PDF := new(PDFReport)
	PDF.FillPDF(*form)
	PDFfile, err := PDF.WriteToPDF()
	if err != nil {
		log.Printf("error generating PDF file: %v", err)
		http.Redirect(w, r, "/error.html", http.StatusFound)
		return
	}

	// located in storage.go
	if err := PDF.StorePDF(PDFfile); err != nil {
		log.Printf("error storing report: %v", err)
		http.Redirect(w, r, "/error.html", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/success.html", http.StatusFound)
}
