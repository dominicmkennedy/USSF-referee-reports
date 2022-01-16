package main

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/dominicmkennedy/gobrr"
	"google.golang.org/api/option"
)

func PDFTempalteInit() {
	if file, err := gobrr.CopyFilePathToMemfile(PAGE_1_TEMPLATE_PATH); err != nil {
        log.Panicln(err)
	} else {
        PAGE_1_TEMPLATE = file
	}

	if file, err := gobrr.CopyFilePathToMemfile(PAGE_2_TEMPLATE_PATH); err != nil {
        log.Panicln(err)
	} else {
        PAGE_2_TEMPLATE = file
	}
}

func firebaseLoginInit() {
	ctx := context.Background()
	sa := option.WithCredentialsFile(PATH_TO_FIREBASE_CREDS)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
        log.Panicln(err)
	}

	FIREBASE_CLIENT, err = app.Firestore(ctx)
	if err != nil {
        log.Panicln(err)
	}
}

func init() {
	StartLogger()
	InitStates()
    PDFTempalteInit()
    firebaseLoginInit()

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
