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

func addtoDB(form *refereeReport) {

        ctx := context.Background()
        sa := option.WithCredentialsFile("creds.json")
        app, err := firebase.NewApp(ctx, nil, sa)
        if err != nil {
                log.Fatal(err)
        }

        client, err := app.Firestore(ctx)
        if err != nil {
                log.Fatal(err)
        }
        defer client.Close()


        DocRef, _, err := client.Collection("reports").Add(ctx, map[string]interface{}{

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
                log.Fatalf("Failed adding to database: %v", err)
        }

        form.ReportID = DocRef.ID




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
                                log.Fatal(err)
                        } 
                        dsnap, err = client.Collection("players").Doc(PlayerID).Get(ctx)

                }  
                if err != nil {
                        log.Fatal(err)
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
                        log.Fatal(err)
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
                                log.Fatal(err)
                        } 
                        dsnap, err = client.Collection("players").Doc(PlayerID).Get(ctx)

                }  
                if err != nil {
                        log.Fatal(err)
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
                        log.Fatal(err)
                } 
        }

}
