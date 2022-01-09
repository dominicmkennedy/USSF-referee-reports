package main

import (
	"log"
	"net/http"
)

/******************************GLOBALS*********************************/

const PATH_TO_FIREBASE_CREDS = "../FirebaseSA.json"
const PATH_TO_GOOGLE_WORKSPACE_PASSWORD = "../GoogleWorkspacePassword.txt"

var Page1TemplatePath string = "../templates/pg1.pdf"
var Page2TemplatePath string = "../templates/pg2.pdf"

var States map[string]struct{}

/**********************************************************************/

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
