package main

import (
	"net/mail"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Sanction struct {
	PlayerName string
	PlayerID   string
	Team       string
	Code       string
}

type Supplemental struct {
	Statement string
	LocationX string
	LocationY string
}

type POSTReport struct {
	HomeTeamState string
	HomeTeamName  string
	HomeTeamScore string
	AwayTeamState string
	AwayTeamName  string
	AwayTeamScore string

	PlayerSex       string
	PlayerAge       string
	GameAssociation string
	GameDivision    string
	GameLeague      string
	GameNumber      string
	GameDate        time.Time
	SubmittedDate   time.Time

	RefereeName            string
	RefereeGrade           string
	AssistantReferee1Name  string
	AssistantReferee1Grade string
	AssistantReferee2Name  string
	AssistantReferee2Grade string
	FourthOfficialName     string
	FourthOfficialGrade    string

	Cautions      []Sanction
	SendOffs      []Sanction
	Supplementals []Supplemental

	//SendToEmail []string

	ReporterName   string
	ReporterUSSFID string
	ReporterPhone  string
	ReporterEmail  string

	ReportID string

	Reviewed          bool
	RecaptchaResponse string `schema:"g-recaptcha-response"`
}

func (r *POSTReport) SanitizePostData() {
	r.ReportID = uuid.NewV4().String()

	intMatcher := regexp.MustCompile(`\D`)

	r.HomeTeamName = strings.ToUpper(r.HomeTeamName)
	r.AwayTeamName = strings.ToUpper(r.AwayTeamName)
	r.HomeTeamScore = intMatcher.ReplaceAllString(r.HomeTeamScore, "")
	r.AwayTeamScore = intMatcher.ReplaceAllString(r.AwayTeamScore, "")

	r.PlayerAge = strings.ToUpper(r.PlayerAge)
	r.GameDivision = strings.ToUpper(r.GameDivision)
	r.GameNumber = r.ReportID

	r.RefereeName = strings.ToUpper(r.RefereeName)
	r.AssistantReferee1Name = strings.ToUpper(r.AssistantReferee1Name)
	r.AssistantReferee2Name = strings.ToUpper(r.AssistantReferee2Name)
	r.FourthOfficialName = strings.ToUpper(r.FourthOfficialName)
	SanitizeRefGrade(&r.RefereeGrade)
	SanitizeRefGrade(&r.AssistantReferee1Grade)
	SanitizeRefGrade(&r.AssistantReferee2Grade)
	SanitizeRefGrade(&r.FourthOfficialGrade)

	if len(r.Cautions) > 30 {
		r.Cautions = r.Cautions[:30]
	}
	if len(r.SendOffs) > 15 {
		r.SendOffs = r.SendOffs[:15]
	}
	if len(r.Supplementals) > 5 {
		r.Supplementals = r.Supplementals[:5]
	}
	//if len(r.SendToEmail) > 30 {
	//r.SendToEmail = r.SendToEmail[:30]
	//}

	for i := range r.Cautions {
		r.Cautions[i] = r.SanitizeSanction(r.Cautions[i])
	}
	for i := range r.SendOffs {
		r.SendOffs[i] = r.SanitizeSanction(r.SendOffs[i])
	}
	for i := range r.Supplementals {
		r.Supplementals[i] = SanitizeSupplemental(r.Supplementals[i])
	}
	//for i := range r.SendToEmail {
	//SanitizeEmail(&r.SendToEmail[i])
	//}

	r.ReporterName = strings.ToUpper(r.ReporterName)
	r.ReporterUSSFID = intMatcher.ReplaceAllString(r.ReporterUSSFID, "")
	r.ReporterPhone = intMatcher.ReplaceAllString(r.ReporterPhone, "")

	SanitizeEmail(&r.ReporterEmail)

	r.Reviewed = false
}

func SanitizeSupplemental(S Supplemental) Supplemental {
	var newSupplemental Supplemental

	newSupplemental.Statement = regexp.MustCompile(`[\t\n\r]+`).ReplaceAllString(S.Statement, " ")

	if len(newSupplemental.Statement) > 1500 {
		newSupplemental.Statement = newSupplemental.Statement[:1500]
	}

	newSupplemental.LocationX = S.LocationX
	newSupplemental.LocationY = S.LocationY

	return newSupplemental
}

func (r *POSTReport) SanitizeSanction(S Sanction) Sanction {
	var newSanction Sanction

	newSanction.PlayerName = strings.ToUpper(S.PlayerName)

	if S.Team == "Home Team" {
		S.Team = r.HomeTeamName
	} else if S.Team == "Away Team" {
		S.Team = r.AwayTeamName
	}
	newSanction.Team = strings.ToUpper(S.Team)

	newSanction.PlayerID = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(S.PlayerID, "")
	newSanction.PlayerID = strings.ToUpper(newSanction.PlayerID)

	newSanction.Code = S.Code

	return newSanction
}

func SanitizeRefGrade(Grade *string) {
	if *Grade != "Grassroots" &&
		*Grade != "Regional" &&
		*Grade != "Regional Emeritus" &&
		*Grade != "National" &&
		*Grade != "National Emeritus" &&
		*Grade != "PRO" &&
		*Grade != "FIFA" {
		*Grade = ""
	}
}

func SanitizeEmail(Email *string) {
	if _, err := mail.ParseAddress(*Email); err != nil {
		*Email = ""
	}

	*Email = strings.ToLower(*Email)
}

func DateConverter(POSTString string) reflect.Value {
	IntTime, err := strconv.ParseInt(POSTString, 10, 64)
	if err != nil {
		return reflect.ValueOf(time.Now())
	}
	return reflect.ValueOf(time.Unix(IntTime, 0))
}
