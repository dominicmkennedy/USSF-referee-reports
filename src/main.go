package main

import (
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
)

/******************************GLOBALS*********************************/

const PATH_TO_FIREBASE_CREDS = "../FirebaseSA.json"
const PATH_TO_GOOGLE_WORKSPACE_PASSWORD = "../GoogleWorkspacePassword.txt"

const PAGE_1_TEMPLATE_PATH string = "../templates/pg1.pdf"
const PAGE_2_TEMPLATE_PATH string = "../templates/pg2.pdf"

var PAGE_1_TEMPLATE *os.File
var PAGE_2_TEMPLATE *os.File

var FIREBASE_CLIENT *firestore.Client

var States map[string]struct{}

/**********************************************************************/

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
