package main

import (
    "time"
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

    // set the reportID to whatever firestore generates
    r.ReportID = GetReportID()

    // trim slice lengths to reasonable amounts in the case of POST request tampering
    if len(r.Cautions) > 30 { r.Cautions = r.Cautions[:30] }
    if len(r.SendOffs) > 15 { r.SendOffs = r.SendOffs[:15] }
    if len(r.Supplementals) > 5 { r.Supplementals = r.Supplementals[:5] }
    if len(r.SendToEmail) > 30 { r.SendToEmail = r.SendToEmail[:30] }

    /*
    // name shouldn't have an issue except maybe length
    HomeTeamName            string
    //score should convert to an int w/o error
    HomeTeamScore           string
    AwayTeamName            string
    AwayTeamScore           string

    // conform w/ map
    // if else blocks would be quicker than a map
    PlayerSex               string
    // honestly maybe regex idrk
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
    PlayerRole       string
    PlayerName       string
    PlayerID         string
    Team             string
    Code             string

    Supplementals           []Supplemental
    Statement   string
    LocationX   string
    LocationY   string

    ReporterName            string
    ReporterUSSFID          string
    ReporterPhone           string

    */

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

}

func IsStringInMap(str * string, m * map[string]struct{}) {
    if _, found := (*m)[*str]; !found {
        *str = ""
    }
}

func SanitizeEmail(Email * string) {
    _, err := mail.ParseAddress(*Email)
    if err != nil {
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
