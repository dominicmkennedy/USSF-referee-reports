package main

import (
    "log"
    "context"
    "time"

    firebase "firebase.google.com/go"
    "google.golang.org/api/option"
    "github.com/gogo/status"
    "google.golang.org/grpc/codes"
)

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

type DBRefereeReport struct {
    Names               map[string]struct{}
    Emails              map[string]struct{}
    PhoneNumbers        map[string]struct{}
    Reports             map[string]interface{}
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
            PlayerID = "000000"
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

    RefereeUSSFID := POST.ReporterUSSFID
    RefereeReport := POST.GetRefereeReport()
    if RefereeUSSFID == "" {
        RefereeUSSFID = "0000000000000000"
    }
    dsnap, err := client.Collection("Referees").Doc(RefereeUSSFID).Get(ctx)
    if status.Code(err) == codes.NotFound {
        _, err = client.Collection("Referees").Doc(RefereeUSSFID).Set(ctx, RefereeReport)
        if err != nil {
            log.Println(err)
        }
    } else if err != nil {
        log.Println(err)
    } else {
        DocData := dsnap.Data()

        if Names, ok := DocData["Names"]; ok {
            for Name := range Names.(map[string]interface{}) {
                RefereeReport.Names[Name] = struct{}{}
            }
        }
        if Emails, ok := DocData["Emails"]; ok {
            for Email := range Emails.(map[string]interface{}) {
                RefereeReport.Emails[Email] = struct{}{}
            }
        }
        if PhoneNumbers, ok := DocData["PhoneNumbers"]; ok {
            for PhoneNumber := range PhoneNumbers.(map[string]interface{}) {
                RefereeReport.PhoneNumbers[PhoneNumber] = struct{}{}
            }
        }

        if Reports, ok := DocData["Reports"]; ok {
            for ReportID, Report := range Reports.(map[string]interface{}) {
                RefereeReport.Reports[ReportID] = Report
            }
        }

        _, err = client.Collection("Referees").Doc(RefereeUSSFID).Set(ctx, RefereeReport)
        if err != nil {
            log.Println(err)
        }

    }
}

func (POST *POSTReport) GetRefereeReport() (DBRefereeReport) {
    return DBRefereeReport{
        Emails:         map[string]struct{}   { POST.ReporterEmail: struct{}{} },
        PhoneNumbers:   map[string]struct{}   { POST.ReporterPhone: struct{}{} },
        Names:          map[string]struct{}   { POST.ReporterName:  struct{}{} },
        Reports:        map[string]interface{}{ POST.ReportID:      struct{
            GameDate        time.Time
            SubmittedDate   time.Time
        } {
            GameDate:       POST.GameDate,
            SubmittedDate:  POST.SubmittedDate,
        },
    },
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
