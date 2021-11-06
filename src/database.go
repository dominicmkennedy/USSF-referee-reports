package main

import (
    "log"
    "context"
    "time"

    firebase "firebase.google.com/go"
    "google.golang.org/api/option"
    "github.com/gogo/status"
    "google.golang.org/grpc/codes"
    "cloud.google.com/go/firestore"
)

type Caution struct {

    GameDate        time.Time
    SubmittedDate   time.Time
    SanctionCode    string
    Reporter        string
    ReportID        string

}

type SendOff struct {

    GameDate        time.Time
    SubmittedDate   time.Time
    SanctionCode    string
    Reporter        string
    ReportID        string

}

type ReportData struct {

    ReportID            string
    ReportDate          time.Time

}

type DBSanction struct {
    GameDate        time.Time
    SubmittedDate   time.Time
    ReportID        string
    Reporter        string
    ReporterUSSFID  string
    SanctionCode    string
    PlayerRole      string
}

type DBPlayerReport struct {
    PlayerName          map[string]struct{}
    Cautions            []interface{}
    SendOffs            []interface{}
}

func (POST *POSTReport) AddToDatabase() {
    ctx := context.Background()
    sa := option.WithCredentialsFile("../creds.json")
    app, err := firebase.NewApp(ctx, nil, sa)
    if err != nil {
        log.Println(err)
    }

    client, err := app.Firestore(ctx)
    if err != nil {
        log.Println(err)
    }
    defer client.Close()

    _, _, err = client.Collection("testing").Add(ctx, POST)
    if err != nil {
        log.Println(err)
    }

    PlayerReports := POST.GetPlayerReports()

    for PlayerID, PlayerReport := range PlayerReports {
        if PlayerID == "" {
            PlayerID = "00000000"
        }
        dsnap, err := client.Collection("Players").Doc(PlayerID).Get(ctx)
        if status.Code(err) == codes.NotFound {
            _, err = client.Collection("Players").Doc(PlayerID).Set(ctx, PlayerReport)
            if err != nil {
                log.Println(err)
            }
        } else if err != nil {
            log.Println(err)
        } else {
            DocData := dsnap.Data()
            if Cautions, ok := DocData["Cautions"]; ok {
                PlayerReport.Cautions = append(PlayerReport.Cautions, Cautions.([]interface{})...)
            }
            if SendOffs, ok := DocData["SendOffs"]; ok {
                PlayerReport.SendOffs = append(PlayerReport.SendOffs, SendOffs.([]interface{})...)
            }
            if PlayerName, ok := DocData["PlayerName"]; ok {
                for Name := range PlayerName.(map[string]interface{}) {
                    PlayerReport.PlayerName[Name] = struct{}{}
                }
            }

            _, err = client.Collection("Players").Doc(PlayerID).Set(ctx, PlayerReport)
            if err != nil {
                log.Println(err)
            }

        }
    }
}

func (POST *POSTReport) GetPlayerReports() (map[string]DBPlayerReport) {

    PlayerReports := make(map[string]DBPlayerReport)

    for _, Caution := range POST.Cautions {

        PlayerID := Caution.PlayerID
        if PlayerReport, ok := PlayerReports[PlayerID]; ok {
            //append
            PlayerReport.Cautions = append(PlayerReports[PlayerID].Cautions, DBSanction{
                SanctionCode:       Caution.Code,
                PlayerRole:         Caution.PlayerRole,
                GameDate:           POST.GameDate,
                SubmittedDate:      POST.SubmittedDate,
                ReportID:           POST.ReportID,
                Reporter:           POST.ReporterName,
                ReporterUSSFID:     POST.ReporterUSSFID,
            })
            PlayerReports[PlayerID] = PlayerReport
        } else {
            //create
            PlayerReports[PlayerID] = DBPlayerReport {

                PlayerName:             map[string]struct{}{
                    Caution.PlayerName: struct{}{},
                },
                Cautions:               []interface{}{DBSanction{
                    SanctionCode:       Caution.Code,
                    PlayerRole:         Caution.PlayerRole,
                    GameDate:           POST.GameDate,
                    SubmittedDate:      POST.SubmittedDate,
                    ReportID:           POST.ReportID,
                    Reporter:           POST.ReporterName,
                    ReporterUSSFID:     POST.ReporterUSSFID,
                }},
                SendOffs:               []interface{}{},
            }
        }

    }

    for _, SendOff := range POST.SendOffs {

        PlayerID := SendOff.PlayerID
        if PlayerReport, ok := PlayerReports[PlayerID]; ok {
            //append
            PlayerReport.SendOffs = append(PlayerReports[PlayerID].SendOffs, DBSanction{
                SanctionCode:       SendOff.Code,
                PlayerRole:         SendOff.PlayerRole,
                GameDate:           POST.GameDate,
                SubmittedDate:      POST.SubmittedDate,
                ReportID:           POST.ReportID,
                Reporter:           POST.ReporterName,
                ReporterUSSFID:     POST.ReporterUSSFID,
            })
            PlayerReports[PlayerID] = PlayerReport
        } else {
            //create
            PlayerReports[PlayerID] = DBPlayerReport {

                PlayerName:             map[string]struct{}{
                    SendOff.PlayerName: struct{}{},
                },
                Cautions:               []interface{}{},
                SendOffs:               []interface{}{DBSanction{
                    SanctionCode:       SendOff.Code,
                    PlayerRole:         SendOff.PlayerRole,
                    GameDate:           POST.GameDate,
                    SubmittedDate:      POST.SubmittedDate,
                    ReportID:           POST.ReportID,
                    Reporter:           POST.ReporterName,
                    ReporterUSSFID:     POST.ReporterUSSFID,
                }},
            }
        }

    }

    return PlayerReports

}

func addtoDB(form *refereeReport) {

    ctx := context.Background()
    sa := option.WithCredentialsFile("../creds.json")
    app, err := firebase.NewApp(ctx, nil, sa)
    if err != nil {
        log.Println(err)
    }

    client, err := app.Firestore(ctx)
    if err != nil {
        log.Println(err)
    }
    defer client.Close()


    DocRef, _, err := client.Collection("reports").Add(ctx, map[string]interface{}{
        //_, err = client.Collection("reports").Doc(form.SubmittedTime.String()).Set(ctx, map[string]interface{}{


        "HomeTeamName":                 form.HomeTeamName,
        "HomeTeamScore":                form.HomeTeamScore,
        "AwayTeamName":                 form.AwayTeamName,
        "AwayTeamScore":                form.AwayTeamScore,

        "GameNumber":                   form.GameNumber,
        "GameDivision":                 form.GameDivisionAgeGroup,
        "GameAssociation":              form.GameAssociationLeague,
        "GameDate":                     form.GameTime,

        "RefereeName":                  form.RefereeName,
        "RefereeGrade":                 form.RefereeGrade,
        "AssistantReferee1Name":        form.AssistantReferee1Name,
        "AssistantReferee1Grade":       form.AssistantReferee1Grade,
        "AssistantReferee2Name":        form.AssistantReferee2Name,
        "AssistantReferee2Grade":       form.AssistantReferee2Grade,
        "FourthOfficialName":           form.FourthOfficialName,
        "FourthOfficialGrade":          form.FourthOfficialGrade,


        "CautionPlayerName":            form.CautionPlayerName,
        "CautionPlayerID":              form.CautionPlayerID,
        "CautionTeam":                  form.CautionTeam,
        "CautionCode":                  form.CautionCode,

        "RedPlayerName":                form.RedPlayerName,
        "RedPlayerID":                  form.RedPlayerID,
        "RedTeam":                      form.RedTeam,
        "RedCode":                      form.RedCode,

        "SupplementalStatement":        form.SupplementalStatement,
        "SupplementalLocationX":        form.SupplementalLocationX,
        "SupplementalLocationY":        form.SupplementalLocationY,
        "SupplementalLocation":         form.SupplementalLocation,

        "SendToEmail":                  form.SendToEmail,
        "Name":                         form.Name,
        "USSFID":                       form.USSFID,
        "ContactNumber":                form.ContactNumber,
        "ContactEmail":                 form.ContactEmail,

        "ipaddr":                       form.ipaddr,
        "SubmittedTime":                form.SubmittedTime,
        "ReportID":                     form.ReportID,

    })
    if err != nil {
        log.Println(err)
    }

    form.ReportID = DocRef.ID
    //form.ReportID = form.SubmittedTime.String()

    if form.USSFID != "" {

        RefereeReports := make([]interface {}, 0)

        dsnap, err := client.Collection("Referees").Doc(form.USSFID).Get(ctx)
        if status.Code(err) == codes.NotFound {
            _, err = client.Collection("Referees").Doc(form.USSFID).Set(ctx, map[string]interface{}{
                "RefereeName":     form.Name,
                "RefereeID":       form.USSFID,
                "Email":           form.ContactEmail,
                "PhoneNumber":     form.ContactNumber,
            })
            if err != nil {
                log.Println(err)
            }
            dsnap, err = client.Collection("Referees").Doc(form.USSFID).Get(ctx)

        }
        if err != nil {
            log.Println(err)
        }

        m := dsnap.Data()

        if m["Reports"] != nil {
            RefereeReports = append(RefereeReports, m["Reports"].([]interface {})...)
        }

        RefereeReport := ReportData{
            ReportID:       form.ReportID,
            ReportDate:     form.SubmittedTime,
        }

        RefereeReports = append(RefereeReports, RefereeReport)


        _, err = client.Collection("Referees").Doc(form.USSFID).Set(ctx, map[string]interface{}{
            "Reports":         RefereeReports,
        }, firestore.MergeAll)
        if err != nil {
            log.Println(err)
        }
    }


    for i := 0; i<len(form.CautionPlayerID); i++ {
        if form.CautionPlayerID[i] == "" {
            continue
        }
        PlayerID        := form.CautionPlayerID[i]
        PlayerName      := form.CautionPlayerName[i]
        PlayerRole      := form.CautionPlayerRole[i]
        ReportCaution   := Caution{

            GameDate:       form.GameTime,
            SubmittedDate:  form.SubmittedTime,
            SanctionCode:   form.CautionCode[i],
            Reporter:       form.Name,
            ReportID:       form.ReportID,

        }
        PlayerCautions := make([]interface {}, 0)


        dsnap, err := client.Collection("players").Doc(PlayerID).Get(ctx)
        if status.Code(err) == codes.NotFound {
            _, err = client.Collection("players").Doc(PlayerID).Set(ctx, map[string]interface{}{
                "Name":         PlayerName,
                "Role":         PlayerRole,
                "PlayerID":     PlayerID,
            })
            if err != nil {
                log.Println(err)
            }
            dsnap, err = client.Collection("players").Doc(PlayerID).Get(ctx)

        }
        if err != nil {
            log.Println(err)
        }

        m := dsnap.Data()

        if m["Cautions"] != nil {
            PlayerCautions = append(PlayerCautions, m["Cautions"].([]interface {})...)
        }
        PlayerCautions = append(PlayerCautions, ReportCaution)

        _, err = client.Collection("players").Doc(PlayerID).Set(ctx, map[string]interface{}{
            "Cautions":         PlayerCautions,
        }, firestore.MergeAll)
        if err != nil {
            log.Println(err)
        }
    }

    for i := 0; i<len(form.RedPlayerID); i++ {
        if form.RedPlayerID[i] == "" {
            continue
        }
        PlayerID        := form.RedPlayerID[i]
        PlayerName      := form.RedPlayerName[i]
        PlayerRole      := form.RedPlayerRole[i]
        ReportSendOff   := SendOff{

            GameDate:       form.GameTime,
            SubmittedDate:  form.SubmittedTime,
            SanctionCode:   form.RedCode[i],
            Reporter:       form.Name,
            ReportID:       form.ReportID,

        }
        PlayerSendOffs := make([]interface {}, 0)

        dsnap, err := client.Collection("players").Doc(PlayerID).Get(ctx)
        if status.Code(err) == codes.NotFound {
            _, err = client.Collection("players").Doc(PlayerID).Set(ctx, map[string]interface{}{
                "Name":         PlayerName,
                "Role":         PlayerRole,
                "PlayerID":     PlayerID,
            })
            if err != nil {
                log.Println(err)
            }
            dsnap, err = client.Collection("players").Doc(PlayerID).Get(ctx)

        }
        if err != nil {
            log.Println(err)
        }

        m := dsnap.Data()

        if m["SendOffs"] != nil {
            PlayerSendOffs = append(PlayerSendOffs, m["SendOffs"].([]interface {})...)
        }
        PlayerSendOffs = append(PlayerSendOffs, ReportSendOff)

        _, err = client.Collection("players").Doc(PlayerID).Set(ctx, map[string]interface{}{
            "SendOffs":         PlayerSendOffs,
        }, firestore.MergeAll)
        if err != nil {
            log.Println(err)
        }
    }

}
