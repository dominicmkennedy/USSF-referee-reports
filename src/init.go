package main

import (
	"log"
	"net/http"

	"github.com/dominicmkennedy/gobrr"
)

func PDFTempalteInit() {
	if file, err := gobrr.CopyFilePathToMemfile(Page1TemplatePath); err != nil {
		log.Println(err)
	} else {
		Page1TemplatePath = file.Name()
	}

	if file, err := gobrr.CopyFilePathToMemfile(Page2TemplatePath); err != nil {
		log.Println(err)
	} else {
		Page2TemplatePath = file.Name()
	}
}

func init() {
	StartLogger()
	InitStates()
	PDFTempalteInit()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/review/", PostForm)
	http.HandleFunc("/submit/", SubmitForm)
}

func InitStates() {
	States = map[string]struct{}{
		"AL": {},
		"AK": {},
		"AZ": {},
		"AR": {},
		"CA": {},
		"CO": {},
		"CT": {},
		"DE": {},
		"DC": {},
		"FL": {},
		"GA": {},
		"HI": {},
		"ID": {},
		"IL": {},
		"IN": {},
		"IA": {},
		"KS": {},
		"KY": {},
		"LA": {},
		"ME": {},
		"MD": {},
		"MA": {},
		"MI": {},
		"MN": {},
		"MS": {},
		"MO": {},
		"MT": {},
		"NE": {},
		"NV": {},
		"NH": {},
		"NJ": {},
		"NM": {},
		"NY": {},
		"NC": {},
		"ND": {},
		"OH": {},
		"OK": {},
		"OR": {},
		"PA": {},
		"RI": {},
		"SC": {},
		"SD": {},
		"TN": {},
		"TX": {},
		"UT": {},
		"VT": {},
		"VA": {},
		"WA": {},
		"WV": {},
		"WI": {},
		"WY": {},
	}
}
