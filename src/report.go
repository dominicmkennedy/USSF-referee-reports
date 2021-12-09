package main

import (
    "time"
    "reflect"
    "strconv"
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

//TODO the rest of this function
func (r *POSTReport) SanitizePostData() {

    r.ReportID = GetReportID()

    if  len(r.Cautions) > 30 {
        r.Cautions = r.Cautions[:30]
    }

    if  len(r.SendOffs) > 15 {
        r.SendOffs = r.SendOffs[:15]
    }

    if  len(r.Supplementals) > 5 {
        r.Supplementals = r.Supplementals[:5]
    }

}

func DateConverter(POSTString string) (reflect.Value) {
    IntTime, err := strconv.ParseInt(POSTString, 10, 64)
    if err != nil {
        return reflect.ValueOf(time.Now())
    }
    return reflect.ValueOf(time.Unix(IntTime, 0))
}
