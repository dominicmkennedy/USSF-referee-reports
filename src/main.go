package main

import (
    "fmt"
    "log"
    "time"
    "net/http"
    "html/template"

    "github.com/gorilla/schema"
)

const PATH_TO_FIREBASE_CREDS = "../FirebaseSA.json"
const PATH_TO_GOOGLE_WORKSPACE_PASSWORD = "../GoogleWorkspacePassword.txt"

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
    PDF.WriteToPDF()

    PDFTime := time.Since(PDFStart)
    PDFStoreStart := time.Now()

    PDF.StorePDF()

    PDFStoreTime := time.Since(PDFStoreStart)
    EmailStart := time.Now()

    SendReport(form)

    EmailTime := time.Since(EmailStart)
    Elapsed := time.Since(Start)

    fmt.Fprintf(w, "HTML parsing took %v\n", ParseTime)
    fmt.Fprintf(w, "Form Sanitzation took %v\n", SanTime)
    fmt.Fprintf(w, "Adding to the Database took %v\n", DBTime)
    fmt.Fprintf(w, "Creating the PDF took %v\n", PDFTime)
    fmt.Fprintf(w, "Storing the PDF took %v\n", PDFStoreTime)
    fmt.Fprintf(w, "Emailing the reports took %v\n", EmailTime)
    fmt.Fprintf(w, "Total elapsed time was %v\n", Elapsed)

}

func index(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("../static/index.html"))
    if err := tmpl.Execute(w, nil); err != nil {
        log.Println(err)
    }
}

func main() {

    StartLogger()

    http.Handle("/script.js", http.FileServer(http.Dir("../static")))
    http.HandleFunc("/", index)
    http.HandleFunc("/submit/", PostForm)

    log.Fatal(http.ListenAndServe(":8080", nil))

}
