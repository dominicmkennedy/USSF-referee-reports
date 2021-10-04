package main

import (
        "fmt"
        "log"
        "net/http"
        "html/template"

        "github.com/gorilla/schema"
        "github.com/leebenson/conform"
)

func PostForm(w http.ResponseWriter, r *http.Request) {

        if r.Method != "POST" {
                fmt.Fprintf(w, "this api endpoint is for POST data only")
                return
        } 
        
        err := r.ParseForm(); 
        if err != nil {
                log.Fatal(err)
        }

        form := new(refereeReport)
        err = schema.NewDecoder().Decode(form, r.PostForm)
        if err != nil {
                log.Fatal(err)
        }
        err = conform.Strings(form)
        if err != nil {
                log.Fatal(err)
        }

        form.SanitizePostData()

        addtoDB(form)
        writePDF(form)
        StorePDF(form)

        ip := r.Header.Get("X-FORWARDED-FOR")
        fmt.Fprintf(w, "IP: %T\n", ip)
        fmt.Fprintf(w, "IP: %s\n", ip)

        fmt.Fprintf(w, "report: %T\n", form)
        fmt.Fprintf(w, "report: %v\n", form)

}

func index(w http.ResponseWriter, r *http.Request) {
        //http.ServeFile(w, r, "./static/")
        tmpl := template.Must(template.ParseFiles("./static/index.html")) 
        err := tmpl.Execute(w, nil)
        if err != nil {
                log.Fatal(err)
        }
}

func main() {

        //http.Handle("/", http.FileServer(http.Dir("./static")))    
        http.Handle("/script.js", http.FileServer(http.Dir("./static")))    
        http.HandleFunc("/", index)
        http.HandleFunc("/submit/", PostForm)

        log.Fatal(http.ListenAndServe(":8080", nil))

}
