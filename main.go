package main

import (
	"fmt"
	"log"
	"net/http"

    "github.com/gorilla/schema"
    "github.com/leebenson/conform"
)




//  need to revisit this code
//  probably good to check out some web resources
func site(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "lol 404", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/index.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}


        form := new(refereeReport)
        schema.NewDecoder().Decode(form, r.PostForm)
        conform.Strings(form)

        form.SanitizePostData()

        addtoDB(form)
        writePDF(form)

        fmt.Fprintf(w, "report: %T\n", form)
		fmt.Fprintf(w, "report: %s\n", form)

    default:
		fmt.Fprintf(w, "can only GET or POST")
	}
}

//  probably need to revist this code later
//  probably good to check out some web resources
func main() {


	http.HandleFunc("/", site)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
