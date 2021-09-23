
package main

import (
        "log"
        "context"

        firebase "firebase.google.com/go"
        "google.golang.org/api/option"
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
        //_, err = client.Collection("reports").Doc(form.ReportID).Set(ctx, map[string]interface{}{
                "HomeTeamName":                 form.HomeTeamName,
                "HomeTeamScore":                form.HomeTeamScore,
                "AwayTeamName":                 form.AwayTeamName,
                "AwayTeamScore":                form.AwayTeamScore,
                "GameNumber":                   form.GameNumber,
                "GameDivision":                 form.GameDivisionAgeGroup,
                "GameAssociation":              form.GameAssociationLeague,
                "GameDate":                     form.GameDate,
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
                
                "SubmittedTime":                form.SubmittedTime,
/*
        RedPlayerName  
        RedPlayerID   
        RedTeam     
        RedCode  

        SupplementalStatement
        SupplementalLocationX 
        SupplementalLocationY
        SupplementalLocation

        SendToEmail   
        Name   
        USSFID   
        ContactNumber 
        ContactEmail    
        ReportID  */  


        })
        if err != nil {
                log.Fatalf("Failed adding to database: %v", err)
        }

        form.ReportID = DocRef.ID

}
