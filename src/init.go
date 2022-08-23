package main

import (
	"context"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"github.com/dominicmkennedy/gobrr"
	"google.golang.org/api/option"
	"gopkg.in/ezzarghili/recaptcha-go.v4"
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
	PDFTempalteInit()
	firebaseLoginInit()

	var err error
	CAPTCHA, err = recaptcha.NewReCAPTCHA(RECAPTCHA_SECRET_KEY, recaptcha.V2, time.Second*10)
	if err != nil {
		log.Panicln(err)
	}

	http.HandleFunc("/post/", PostForm)

	http.Handle("/", http.FileServer(http.Dir("../static")))
	log.Print("Listening on :8080")
}
