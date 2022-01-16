package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

type DBSanction struct {
	GameDate       time.Time
	MisconductCode string
	ReporterID     string
	Reporter       string
	ReportID       string
	Team           string
	Role           string
}

type DBPlayerReport struct {
	PlayerName map[string]struct{}
	Cautions   []interface{}
	SendOffs   []interface{}
}

type DBRefereeReport struct {
	Names        map[string]struct{}
	Emails       map[string]struct{}
	PhoneNumbers map[string]struct{}
	Reports      map[string]interface{}
}

func GetReportID() (string, error) {
	DocRef, _, err := FIREBASE_CLIENT.Collection("Reports").Add(context.Background(), struct{}{})
	if err != nil {
		return "", fmt.Errorf("Error creating a report: %v", err)
	}

	return DocRef.ID, nil
}

func (POST *POSTReport) AddToDatabase() error {
	ctx := context.Background()

    if _, err := FIREBASE_CLIENT.Collection("Reports").Doc(POST.ReportID).Set(ctx, POST); err != nil {
		return fmt.Errorf("Error uploading report data: %v", err)
	}

	PlayerReports := POST.GetPlayerReports()

	for PlayerID, PlayerReport := range PlayerReports {
		if PlayerID == "" {
			PlayerID = "000000"
		}
		dsnap, err := FIREBASE_CLIENT.Collection("Players").Doc(PlayerID).Get(ctx)
		if status.Code(err) == codes.NotFound {
			if _, err = FIREBASE_CLIENT.Collection("Players").Doc(PlayerID).Set(ctx, PlayerReport); err != nil {
				return fmt.Errorf("Error uploading new player report: %v", err)
			}
		} else if err != nil {
			return fmt.Errorf("Error retriving player report: %v", err)
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

			if _, err = FIREBASE_CLIENT.Collection("Players").Doc(PlayerID).Set(ctx, PlayerReport); err != nil {
				return fmt.Errorf("Error uploading player report: %v", err)
			}
		}
	}

	RefereeUSSFID := POST.ReporterUSSFID
	RefereeReport := POST.GetRefereeReport()
	if RefereeUSSFID == "" {
		RefereeUSSFID = "0000000000000000"
	}
	dsnap, err := FIREBASE_CLIENT.Collection("Referees").Doc(RefereeUSSFID).Get(ctx)
	if status.Code(err) == codes.NotFound {
		if _, err = FIREBASE_CLIENT.Collection("Referees").Doc(RefereeUSSFID).Set(ctx, RefereeReport); err != nil {
			return fmt.Errorf("Error uploading new referee report: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("Error retriving referee report: %v", err)
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

		if _, err = FIREBASE_CLIENT.Collection("Referees").Doc(RefereeUSSFID).Set(ctx, RefereeReport); err != nil {
			return fmt.Errorf("Error uploading referee report: %v", err)
		}
	}

	return nil
}

func (POST *POSTReport) GetRefereeReport() DBRefereeReport {
	var Email map[string]struct{}
	if len(POST.ReporterEmail) == 0 {
		Email = make(map[string]struct{})
	} else {
		Email = map[string]struct{}{POST.ReporterEmail: {}}
	}

	var PhoneNumber map[string]struct{}
	if len(POST.ReporterPhone) == 0 {
		PhoneNumber = make(map[string]struct{})
	} else {
		PhoneNumber = map[string]struct{}{POST.ReporterPhone: {}}
	}

	var Name map[string]struct{}
	if len(POST.ReporterName) == 0 {
		Name = make(map[string]struct{})
	} else {
		Name = map[string]struct{}{POST.ReporterName: {}}
	}

	var Report map[string]interface{}
	if len(POST.ReportID) == 0 {
		Report = make(map[string]interface{})
	} else {
		Report = map[string]interface{}{POST.ReportID: POST.SubmittedDate}
	}

	return DBRefereeReport{
		Emails:       Email,
		PhoneNumbers: PhoneNumber,
		Names:        Name,
		Reports:      Report,
	}
}

func (POST *POSTReport) GetPlayerReports() map[string]DBPlayerReport {
	PlayerReports := make(map[string]DBPlayerReport)

	for _, Caution := range POST.Cautions {
		PlayerID := Caution.PlayerID
		if PlayerReport, ok := PlayerReports[PlayerID]; ok {
			//append
			PlayerReport.PlayerName[Caution.PlayerName] = struct{}{}
			PlayerReport.Cautions = append(PlayerReports[PlayerID].Cautions, DBSanction{
				GameDate:       POST.GameDate,
				MisconductCode: Caution.Code,
				ReporterID:     POST.ReporterUSSFID,
				Reporter:       POST.ReporterName,
				ReportID:       POST.ReportID,
				Team:           Caution.Team,
				Role:           Caution.PlayerRole,
			})
			PlayerReports[PlayerID] = PlayerReport
		} else {
			//create
			PlayerReports[PlayerID] = DBPlayerReport{

				PlayerName: map[string]struct{}{
					Caution.PlayerName: {},
				},
				Cautions: []interface{}{DBSanction{
					GameDate:       POST.GameDate,
					MisconductCode: Caution.Code,
					ReporterID:     POST.ReporterUSSFID,
					Reporter:       POST.ReporterName,
					ReportID:       POST.ReportID,
					Team:           Caution.Team,
					Role:           Caution.PlayerRole,
				}},
				SendOffs: []interface{}{},
			}
		}
	}

	for _, SendOff := range POST.SendOffs {
		PlayerID := SendOff.PlayerID
		if PlayerReport, ok := PlayerReports[PlayerID]; ok {
			//append
			PlayerReport.PlayerName[SendOff.PlayerName] = struct{}{}
			PlayerReport.SendOffs = append(PlayerReports[PlayerID].SendOffs, DBSanction{
				GameDate:       POST.GameDate,
				MisconductCode: SendOff.Code,
				ReporterID:     POST.ReporterUSSFID,
				Reporter:       POST.ReporterName,
				ReportID:       POST.ReportID,
				Team:           SendOff.Team,
				Role:           SendOff.PlayerRole,
			})
			PlayerReports[PlayerID] = PlayerReport
		} else {
			//create
			PlayerReports[PlayerID] = DBPlayerReport{

				PlayerName: map[string]struct{}{
					SendOff.PlayerName: {},
				},
				Cautions: []interface{}{},
				SendOffs: []interface{}{DBSanction{
					GameDate:       POST.GameDate,
					MisconductCode: SendOff.Code,
					ReporterID:     POST.ReporterUSSFID,
					Reporter:       POST.ReporterName,
					ReportID:       POST.ReportID,
					Team:           SendOff.Team,
					Role:           SendOff.PlayerRole,
				}},
			}
		}
	}

	return PlayerReports
}
