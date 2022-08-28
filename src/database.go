package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
)

type DBSanction struct {
	GameDate       time.Time
	MisconductCode string
	ReporterID     string
	ReportID       string
	Reporter       string
	Team           string
	Index          int
}

// This struct is cursed, but I can't think of a better way rn
type DBPlayerReport struct {
	PlayerName []interface{}
	Cautions   []interface{}
	SendOffs   []interface{}
}

func (POST *POSTReport) AddToDatabase() error {
	ctx := context.Background()

	// add the sanitized report struct to the reports table
	if _, err := FIREBASE_CLIENT.Collection("Reports").Doc(POST.ReportID).Set(ctx, POST); err != nil {
		return fmt.Errorf("Error uploading report data: %v", err)
	}

	// add the reporter in the database
	RefereeUSSFID := POST.ReporterUSSFID
	if RefereeUSSFID == "" {
		RefereeUSSFID = "0000000000000000"
	}
	_, err := FIREBASE_CLIENT.Collection("Referees").Doc(RefereeUSSFID).Set(ctx, map[string]interface{}{
		"Names":        firestore.ArrayUnion(POST.ReporterName),
		"Emails":       firestore.ArrayUnion(POST.ReporterEmail),
		"PhoneNumbers": firestore.ArrayUnion(POST.ReporterPhone),
		"Reports":      map[string]time.Time{POST.ReportID: POST.SubmittedDate},
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("Error updating referee field: %v", err)
	}

	// for each player referenced in the report add that player's caution/send off into the AddToDatabaseNewTesting
	PlayerReports := POST.GetPlayerReports()

	for PlayerID, PlayerReport := range PlayerReports {
		if PlayerID == "" {
			PlayerID = "00000000000"
		}

		_, err = FIREBASE_CLIENT.Collection("Players").Doc(PlayerID).Set(ctx, map[string]interface{}{
			"PlayerName": firestore.ArrayUnion(PlayerReport.PlayerName...),
			"Cautions":   firestore.ArrayUnion(PlayerReport.Cautions...),
			"SendOffs":   firestore.ArrayUnion(PlayerReport.SendOffs...),
		}, firestore.MergeAll)
		if err != nil {
			return fmt.Errorf("Error updating player field: %v", err)
		}

		_, err := FIREBASE_CLIENT.Collection("Players").Doc(PlayerID).Update(ctx, []firestore.Update{
			{
				Path:  "NumCautions",
				Value: firestore.FieldTransformIncrement(len(PlayerReport.Cautions)),
			},
			{
				Path:  "NumSendOffs",
				Value: firestore.FieldTransformIncrement(len(PlayerReport.SendOffs)),
			},
			{
				Path: "NumPoints",
				Value: firestore.FieldTransformIncrement(
					len(PlayerReport.Cautions) +
						(2 * len(PlayerReport.SendOffs))),
			},
		})
		if err != nil {
			return fmt.Errorf("Error updating num sanctions: %v", err)
		}
	}

	return nil
}

// this could possibly be shorter
// consider rewriting
func (POST *POSTReport) GetPlayerReports() map[string]DBPlayerReport {
	PlayerReports := make(map[string]DBPlayerReport)

	for i, Caution := range POST.Cautions {
		PlayerID := Caution.PlayerID
		if PlayerReport, ok := PlayerReports[PlayerID]; ok {
			//append
			PlayerReport.PlayerName = append(PlayerReport.PlayerName, Caution.PlayerName)
			PlayerReport.Cautions = append(PlayerReports[PlayerID].Cautions, DBSanction{
				GameDate:       POST.GameDate,
				MisconductCode: Caution.Code,
				ReporterID:     POST.ReporterUSSFID,
				ReportID:       POST.ReportID,
				Reporter:       POST.ReporterName,
				Team:           Caution.Team,
				Index:          i,
			})
			PlayerReports[PlayerID] = PlayerReport
		} else {
			//create
			PlayerReports[PlayerID] = DBPlayerReport{

				PlayerName: []interface{}{Caution.PlayerName},
				Cautions: []interface{}{DBSanction{
					GameDate:       POST.GameDate,
					MisconductCode: Caution.Code,
					ReporterID:     POST.ReporterUSSFID,
					ReportID:       POST.ReportID,
					Reporter:       POST.ReporterName,
					Team:           Caution.Team,
					Index:          i,
				}},
				SendOffs: []interface{}{},
			}
		}
	}

	for i, SendOff := range POST.SendOffs {
		PlayerID := SendOff.PlayerID
		if PlayerReport, ok := PlayerReports[PlayerID]; ok {
			//append
			PlayerReport.PlayerName = append(PlayerReport.PlayerName, SendOff.PlayerName)
			PlayerReport.SendOffs = append(PlayerReports[PlayerID].SendOffs, DBSanction{
				GameDate:       POST.GameDate,
				MisconductCode: SendOff.Code,
				ReporterID:     POST.ReporterUSSFID,
				ReportID:       POST.ReportID,
				Reporter:       POST.ReporterName,
				Team:           SendOff.Team,
				Index:          i,
			})
			PlayerReports[PlayerID] = PlayerReport
		} else {
			//create
			PlayerReports[PlayerID] = DBPlayerReport{

				PlayerName: []interface{}{SendOff.PlayerName},
				Cautions:   []interface{}{},
				SendOffs: []interface{}{DBSanction{
					GameDate:       POST.GameDate,
					MisconductCode: SendOff.Code,
					ReporterID:     POST.ReporterUSSFID,
					ReportID:       POST.ReportID,
					Reporter:       POST.ReporterName,
					Team:           SendOff.Team,
					Index:          i,
				}},
			}
		}
	}

	return PlayerReports
}
