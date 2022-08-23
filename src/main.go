package main

import (
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"gopkg.in/ezzarghili/recaptcha-go.v4"
)

/******************************GLOBALS*********************************/

const PATH_TO_FIREBASE_CREDS = "../FirebaseSA.json"

const PAGE_1_TEMPLATE_PATH string = "../templates/pg1.pdf"
const PAGE_2_TEMPLATE_PATH string = "../templates/pg2.pdf"

var PAGE_1_TEMPLATE *os.File
var PAGE_2_TEMPLATE *os.File

var FIREBASE_CLIENT *firestore.Client

const RECAPTCHA_SECRET_KEY = "6LeI5J0hAAAAAGKJE7gHeEQWSvvB-8-mENUjhHlq"

var CAPTCHA recaptcha.ReCAPTCHA

/**********************************************************************/

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
