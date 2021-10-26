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

    err := r.ParseForm();
    if err != nil {
        log.Println(err)
    }

    form := new(POSTReport)
    decoder := schema.NewDecoder()
    decoder.RegisterConverter(time.Now(), DateConverter)
    err = decoder.Decode(form, r.PostForm)
    if err != nil {
        log.Println(err)
    }
    form.SanitizePostData()

    PDF := new(PDFReport)
    PDF.FillPDF(*form)
    PDF.WriteToPDF()

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
