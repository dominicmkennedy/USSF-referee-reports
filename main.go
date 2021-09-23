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
        schema.NewDecoder().Decode(form, r.PostForm)
        conform.Strings(form)

        form.SanitizePostData()

        addtoDB(form)
        writePDF(form)
        StorePDF(form)

        ip := r.Header.Get("X-FORWARDED-FOR")
        fmt.Fprintf(w, "IP: %T\n", ip)
        fmt.Fprintf(w, "IP: %s\n", ip)

        fmt.Fprintf(w, "report: %T\n", form)
        fmt.Fprintf(w, "report: %s\n", form)

}

func index(w http.ResponseWriter, r *http.Request) {
        //http.ServeFile(w, r, "./static/")
        tmpl := template.Must(template.ParseFiles("./static/index.html")) 
        tmpl.Execute(w, nil)
}

func main() {

        //http.Handle("/", http.FileServer(http.Dir("./static")))    
        http.Handle("/script.js", http.FileServer(http.Dir("./static")))    
        http.HandleFunc("/", index)
        http.HandleFunc("/submit/", PostForm)

        log.Fatal(http.ListenAndServe(":8080", nil))

}
