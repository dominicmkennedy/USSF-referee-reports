package main

import (
    "time"
    "strings"
    "reflect"
    "strconv"
    "net/mail"
)

type Sanction struct {

    PlayerRole       string
    PlayerName       string
    PlayerID         string
    Team             string
    Code             string

}

type Supplemental struct {

    Statement   string
    LocationX   string
    LocationY   string

}

type POSTReport struct {

    HomeTeamState           string
    HomeTeamName            string
    HomeTeamScore           string
    AwayTeamState           string
    AwayTeamName            string
    AwayTeamScore           string

    PlayerSex               string
    PlayerAge               string

    GameAssociation         string
    GameDivision            string
    GameLeague              string
    GameNumber              string
    GameDate                time.Time
    SubmittedDate           time.Time

    RefereeName             string
    RefereeGrade            string
    AssistantReferee1Name   string
    AssistantReferee1Grade  string
    AssistantReferee2Name   string
    AssistantReferee2Grade  string
    FourthOfficialName      string
    FourthOfficialGrade     string

    Cautions                []Sanction
    SendOffs                []Sanction
    Supplementals           []Supplemental

    SendToEmail             []string

    ReporterName            string
    ReporterUSSFID          string
    ReporterPhone           string
    ReporterEmail           string

    ReportID                string

}

func (r *POSTReport) SanitizePostData() {

    //TODO consider string lengths

    // set the reportID to whatever firestore generates
    r.ReportID = GetReportID()

    // trim slice lengths to reasonable amounts in the case of POST request tampering
    if len(r.Cautions) > 30 { r.Cautions = r.Cautions[:30] }
    if len(r.SendOffs) > 15 { r.SendOffs = r.SendOffs[:15] }
    if len(r.Supplementals) > 5 { r.Supplementals = r.Supplementals[:5] }
    if len(r.SendToEmail) > 30 { r.SendToEmail = r.SendToEmail[:30] }

    r.HomeTeamName = strings.ToUpper(r.HomeTeamName)
    r.AwayTeamName = strings.ToUpper(r.AwayTeamName)

    // makes sure the user has actually put in ints
    // even though this is internally stored as a string
    SanitizeInt(&r.HomeTeamScore)
    SanitizeInt(&r.AwayTeamScore)

    if r.PlayerSex != "Men" && r.PlayerSex != "Women" && r.PlayerSex != "Co-ed" { r.PlayerSex = "" }

    SanitizeRefGrade(&r.RefereeGrade)
    SanitizeRefGrade(&r.AssistantReferee1Grade)
    SanitizeRefGrade(&r.AssistantReferee2Grade)
    SanitizeRefGrade(&r.FourthOfficialGrade)

    r.RefereeName = strings.ToUpper(r.RefereeName)
    r.AssistantReferee1Name = strings.ToUpper(r.AssistantReferee1Name)
    r.AssistantReferee2Name = strings.ToUpper(r.AssistantReferee2Name)
    r.FourthOfficialName = strings.ToUpper(r.FourthOfficialName)

    for i := range r.Cautions {
        SanitizeSanction(&r.Cautions[i])
    }
    for i := range r.SendOffs {
        SanitizeSanction(&r.SendOffs[i])
    }

    SanitizeEmail(&r.ReporterEmail)

    for i := range r.SendToEmail {
        SanitizeEmail(&r.SendToEmail[i])
    }


    // TODO scope this as both global and constant
    // so it is declared once at runtime
    // and is imutable
    States := map[string]struct{} {
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

    IsStringInMap(&r.AwayTeamState, &States)
    IsStringInMap(&r.HomeTeamState, &States)
    
    /*

    // honestly maybe regex idrk
    PlayerAge               string

    GameAssociation         string
    GameDivision            string
    GameLeague              string
    GameNumber              string

    Cautions                []Sanction
    SendOffs                []Sanction
    Team             string
    Code             string

    Supplementals           []Supplemental
    Statement   string

    */
    

    r.ReporterName = strings.ToUpper(r.ReporterName)

    SanitizeInt(&r.ReporterUSSFID)
    SanitizeInt(&r.ReporterPhone)

}

func SanitizeSupplemental(S *Supplemental) () {

    //TODO sanitize Supplemental statment

    //TODO ranges for these numbers based on math
    SanitizeInt(&S.LocationX)
    SanitizeInt(&S.LocationY)

}

func SanitizeInt(s *string) () {
    
    if _, err := strconv.Atoi(*s); err != nil {
        *s = ""
    }

}

func SanitizeSanction(S *Sanction) () {

    S.PlayerName = strings.ToUpper(S.PlayerName)
    
    if S.PlayerRole != "Player" && S.PlayerRole != "Bench Personnoel" {
        S.PlayerRole = ""
    }
   
    SanitizeInt(&S.PlayerID)

    //TODO sanitize the sanctions misconduct code
    //TODO sanitize the sanctions Team

}

func SanitizeRefGrade(Grade *string) () {

    if
    *Grade != "Grassroots" &&
    *Grade != "Regional" &&
    *Grade != "Regional Emeritus" &&
    *Grade != "National" &&
    *Grade != "National Emeritus" &&
    *Grade != "PRO" &&
    *Grade != "FIFA" {
        *Grade = ""
    }

}

func IsStringInMap(str * string, m * map[string]struct{}) {
    if _, found := (*m)[*str]; !found {
        *str = ""
    }
}

func SanitizeEmail(Email * string) {
    if _, err := mail.ParseAddress(*Email); err != nil {
        *Email = ""
    }
}

func DateConverter(POSTString string) (reflect.Value) {
    IntTime, err := strconv.ParseInt(POSTString, 10, 64)
    if err != nil {
        return reflect.ValueOf(time.Now())
    }
    return reflect.ValueOf(time.Unix(IntTime, 0))
}
