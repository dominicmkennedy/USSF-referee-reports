package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "net/http"
    "html/template"

    "github.com/gorilla/schema"
    "github.com/leebenson/conform"
)

func ReadUserIP(r *http.Request) string {
    IPAddress := r.Header.Get("X-Real-Ip")
    if IPAddress == "" {
        IPAddress = r.Header.Get("X-Forwarded-For")
    }
    if IPAddress == "" {
        IPAddress = r.RemoteAddr
    }
    return IPAddress
}

func PostForm(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        fmt.Fprintf(w, "this api endpoint is for POST data only")
        return
    } 

    err := r.ParseForm(); 
    if err != nil {
        log.Println(err)
    }

    form := new(refereeReport)
    decoder := schema.NewDecoder()
    //decoder.IgnoreUnknownKeys(true)
    err = decoder.Decode(form, r.PostForm)
    if err != nil {
        log.Println(err)
    }
    err = conform.Strings(form)
    if err != nil {
        log.Println(err)
    }

    form.SanitizePostData()

    addtoDB(form)
    writePDF(form)
    StorePDF(form)

    ip := ReadUserIP(r)
    fmt.Fprintf(w, "IP: %T\n", ip)
    fmt.Fprintf(w, "IP: %v\n", ip)

    fmt.Fprintf(w, "report: %T\n", form)
    fmt.Fprintf(w, "report: %v\n", form)

}

func index(w http.ResponseWriter, r *http.Request) {
    //http.ServeFile(w, r, "./static/")
    tmpl := template.Must(template.ParseFiles("./static/index.html")) 
    err := tmpl.Execute(w, nil)
    if err != nil {
        log.Println(err)
    }
}

func StartLogger() {
    LogFile, err := os.OpenFile("logs.txt", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
    if err != nil {
        panic(err)
    }

    LogWriter := io.MultiWriter(os.Stdout, LogFile)

    log.SetOutput(LogWriter)

    log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

func main() {

    StartLogger()

    //http.Handle("/", http.FileServer(http.Dir("./static")))    
    http.Handle("/script.js", http.FileServer(http.Dir("./static")))    
    http.HandleFunc("/", index)
    http.HandleFunc("/submit/", PostForm)

    log.Fatal(http.ListenAndServe(":8080", nil))

}
