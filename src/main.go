package main

import (
    "fmt"
    "log"
    "time"
    "net/http"
    "html/template"

    "github.com/gorilla/schema"
    //"github.com/leebenson/conform"
)

func PostForm(w http.ResponseWriter, r *http.Request) {

    //  parse http POST request
    err := r.ParseForm();
    if err != nil {
        log.Println(err)
    }

    //  create a new POSTReport struct
    //  put POST data into struct
    //  then sanitize data
    form := new(POSTReport)
    decoder := schema.NewDecoder()
    decoder.RegisterConverter(time.Now(), DateConverter)
    err = decoder.Decode(form, r.PostForm)
    if err != nil {
        log.Println(err)
    }
    form.SanitizePostData()

    //  put POST data into the database
    form.AddToDatabase()

    //  create new PDF struct
    //  fill in data from the POSTReport
    //  write the PDF to disk
    //  then store it in cloud
    PDF := new(PDFReport)
    PDF.FillPDF(*form)
    PDF.WriteToPDF()
    PDF.StorePDF()

    fmt.Fprintf(w, "report: %T\n", form)
    fmt.Fprintf(w, "report: %v\n", form)

}

func index(w http.ResponseWriter, r *http.Request) {
    //http.ServeFile(w, r, "./static/")
    tmpl := template.Must(template.ParseFiles("../static/index.html"))
    err := tmpl.Execute(w, nil)
    if err != nil {
        log.Println(err)
    }
}

func main() {

    StartLogger()

    //http.Handle("/", http.FileServer(http.Dir("./static")))    
    http.Handle("/script.js", http.FileServer(http.Dir("../static")))
    http.HandleFunc("/", index)
    http.HandleFunc("/submit/", PostForm)

    log.Fatal(http.ListenAndServe(":8080", nil))

}
