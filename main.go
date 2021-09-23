package main

import (
        "fmt"
        "log"
        "net/http"

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

        fmt.Fprintf(w, "report: %T\n", form)
        fmt.Fprintf(w, "report: %s\n", form)

}

func main() {

        http.Handle("/", http.FileServer(http.Dir("./static")))    
        http.HandleFunc("/submit/", PostForm)

        log.Fatal(http.ListenAndServe(":8080", nil))

}
