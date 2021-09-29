package main

import (
        "log"
        "context"

        firebase "firebase.google.com/go"
        "google.golang.org/api/option"
        "github.com/gogo/status"
        "google.golang.org/grpc/codes"
        "cloud.google.com/go/firestore"
)

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

        PlayerID        := form.CautionPlayerID[0]
        PlayerName      := form.CautionPlayerName[0]
        PlayerCautions  := make([]interface {}, 0)

        
        dsnap, err := client.Collection("players").Doc(PlayerID).Get(ctx)
        if status.Code(err) == codes.NotFound {
                log.Printf("doc does not yet exist creating doc\n")
                _, err = client.Collection("players").Doc(PlayerID).Set(ctx, map[string]interface{}{
                        "Name":     PlayerName,
                        "Cautions": PlayerCautions,
                })
                dsnap, err = client.Collection("players").Doc(PlayerID).Get(ctx)

        }  
        if err != nil {
                log.Fatal(err)
        } 

        m := dsnap.Data()
        log.Printf("Player data: %#v\n", m)
        
        PlayerCautions = append(PlayerCautions, m["Cautions"].([]interface {})...) 
        PlayerCautions = append(PlayerCautions, form.CautionCode[0]) 
       
       _, err = client.Collection("players").Doc(PlayerID).Set(ctx, map[string]interface{}{
                "Cautions":         PlayerCautions, 
        }, firestore.MergeAll)

}
