package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dominicmkennedy/gobrr"
	"github.com/patiek/go-pdftools/fdf"
	"github.com/patiek/go-pdftools/pdftk"
)

type PDFReport struct {
	Pg1Reports []fdf.Inputs
	Pg2Reports []fdf.Inputs
	ReportID   string
}

func (PDF *PDFReport) FillPDF(POST POSTReport) {
	PDF.ReportID = POST.ReportID

	pg1 := 0
	if ((len(POST.Cautions) + 9) / 10) > pg1 {
		pg1 = ((len(POST.Cautions) + 9) / 10)
	}
	if ((len(POST.SendOffs) + 4) / 5) > pg1 {
		pg1 = ((len(POST.SendOffs) + 4) / 5)
	}
	if pg1 == 0 {
		pg1 = 1
	}

	PDF.Pg1Reports = make([]fdf.Inputs, pg1)
	for i := 0; i < pg1; i++ {
		PDF.FillPg1PDF(POST, i)
	}

	pg2 := len(POST.Supplementals)
	PDF.Pg2Reports = make([]fdf.Inputs, pg2)
	for i := 0; i < pg2; i++ {
		PDF.FillPg2PDF(POST, i)
	}
}

func (PDF PDFReport) FillPg1PDF(POST POSTReport, Page int) {
	PDF.Pg1Reports[Page] = map[string]interface{}{
		"HomeTeamName":  POST.HomeTeamState + ": " + POST.HomeTeamName,
		"HomeTeamScore": POST.HomeTeamScore,
		"AwayTeamName":  POST.AwayTeamState + ": " + POST.AwayTeamName,
		"AwayTeamScore": POST.AwayTeamScore,

		"GameDivision":    "Division: " + POST.GameDivision + " Sex: " + POST.PlayerSex + " Age: " + POST.PlayerAge,
		"GameAssociation": POST.GameLeague,
		"GameNumber":      POST.GameNumber,
		"GameDate":        POST.GameDate.Format(time.RFC1123),

		"RefereeName":            POST.RefereeName,
		"RefereeGrade":           POST.RefereeGrade,
		"AssistantReferee1Name":  POST.AssistantReferee1Name,
		"AssistantReferee1Grade": POST.AssistantReferee1Grade,
		"AssistantReferee2Name":  POST.AssistantReferee2Name,
		"AssistantReferee2Grade": POST.AssistantReferee2Grade,
		"FourthOfficialName":     POST.FourthOfficialName,
		"FourthOfficialGrade":    POST.FourthOfficialGrade,

		"Name":          POST.ReporterName,
		"USSFID":        POST.ReporterUSSFID,
		"ContactNumber": POST.ReporterPhone,
		"ContactEmail":  POST.ReporterEmail,
		"SubmittedDate": POST.SubmittedDate.Format(time.RFC1123),
	}

	for PDFIndex, POSTIndex := 0, Page*10; PDFIndex < 10 && POSTIndex < len(POST.Cautions); PDFIndex, POSTIndex = PDFIndex+1, POSTIndex+1 {
		PDF.Pg1Reports[Page]["CautionPlayerName"+strconv.Itoa(PDFIndex)] = POST.Cautions[POSTIndex].PlayerName
		PDF.Pg1Reports[Page]["CautionPlayerID"+strconv.Itoa(PDFIndex)] = POST.Cautions[POSTIndex].PlayerID
		PDF.Pg1Reports[Page]["CautionTeam"+strconv.Itoa(PDFIndex)] = POST.Cautions[POSTIndex].Team
		PDF.Pg1Reports[Page]["CautionCode"+strconv.Itoa(PDFIndex)] = POST.Cautions[POSTIndex].Code
	}

	for PDFIndex, POSTIndex := 0, Page*5; PDFIndex < 5 && POSTIndex < len(POST.SendOffs); PDFIndex, POSTIndex = PDFIndex+1, POSTIndex+1 {
		PDF.Pg1Reports[Page]["RedPlayerName"+strconv.Itoa(PDFIndex)] = POST.SendOffs[POSTIndex].PlayerName
		PDF.Pg1Reports[Page]["RedPlayerID"+strconv.Itoa(PDFIndex)] = POST.SendOffs[POSTIndex].PlayerID
		PDF.Pg1Reports[Page]["RedTeam"+strconv.Itoa(PDFIndex)] = POST.SendOffs[POSTIndex].Team
		PDF.Pg1Reports[Page]["RedCode"+strconv.Itoa(PDFIndex)] = POST.SendOffs[POSTIndex].Code
	}
}

func (PDF PDFReport) FillPg2PDF(POST POSTReport, Page int) {
	var SupplementalLocation string
	Marker := "x"

	// takes values 0-15 for width
	SupplementalLocationY, err := strconv.ParseFloat(POST.Supplementals[Page].LocationY, 64)
	if err != nil {
		Marker = " "
	}

	// takes int values 0-46 for height
	SupplementalLocationX, err := strconv.ParseFloat(POST.Supplementals[Page].LocationX, 64)
	if err != nil {
		Marker = " "
	}

	for i := 0; i < int(SupplementalLocationY); i++ {
		SupplementalLocation += "\n"
	}
	log.Println((int(SupplementalLocationY)))
	log.Println((int(SupplementalLocationX)))

	for i := 0; i < int(SupplementalLocationX); i++ {
		SupplementalLocation += " "
	}

	SupplementalLocation += Marker

	PDF.Pg2Reports[Page] = map[string]interface{}{
		"HomeTeamName":  POST.HomeTeamState + ": " + POST.HomeTeamName,
		"HomeTeamScore": POST.HomeTeamScore,
		"AwayTeamName":  POST.AwayTeamState + ": " + POST.AwayTeamName,
		"AwayTeamScore": POST.AwayTeamScore,

		"GameDivision":    "Division: " + POST.GameDivision + " Sex: " + POST.PlayerSex + " Age: " + POST.PlayerAge,
		"GameAssociation": POST.GameLeague,
		"GameNumber":      POST.GameNumber,
		"GameDate":        POST.GameDate.Format(time.RFC1123),

		"SupplementalStatement": POST.Supplementals[Page].Statement,
		"SupplementalLocation":  SupplementalLocation,

		"Name":          POST.ReporterName,
		"USSFID":        POST.ReporterUSSFID,
		"ContactNumber": POST.ReporterPhone,
		"ContactEmail":  POST.ReporterEmail,
		"SubmittedDate": POST.SubmittedDate.Format(time.RFC1123),
	}
}

func writePDFPage(pageData fdf.Inputs, pageTemplatePath string) (*os.File, error) {
	fdfData := bytes.NewBuffer(nil)
	if err := fdf.Write(fdfData, pageData); err != nil {
		return nil, fmt.Errorf("error filling fdf: %v", err)
	}

	file, err := gobrr.CreateEmptyMemfile()
	if err != nil {
		return nil, fmt.Errorf("error creating memfile: %v", err)
	}

	if err := pdftk.FillForm(file, pageTemplatePath, fdfData, pdftk.OptionFlatten()); err != nil {
		return nil, fmt.Errorf("error filling pdf: %v", err)
	}

	return file, nil
}

func (PDF PDFReport) WriteToPDF() (*bytes.Buffer, error) {
	outfiles := pdftk.NewInputFileMap()

	for i, pg1Data := range PDF.Pg1Reports {
		file, err := writePDFPage(pg1Data, PAGE_1_TEMPLATE.Name())
		if err != nil {
			return nil, fmt.Errorf("error generating pdf pg1: %v", err)
		}
		outfiles[pdftk.InputHandleNameFromInt(i)] = file.Name()
		defer file.Close()
	}

	for i, pg2Data := range PDF.Pg2Reports {
		file, err := writePDFPage(pg2Data, PAGE_2_TEMPLATE.Name())
		if err != nil {
			return nil, fmt.Errorf("error generating pdf pg1: %v", err)
		}
		outfiles[pdftk.InputHandleNameFromInt(i+len(PDF.Pg1Reports))] = file.Name()
		defer file.Close()
	}

	Output := bytes.NewBuffer(nil)

	if err := pdftk.Cat(Output, outfiles, []pdftk.PageRange{}, pdftk.OptionFlatten()); err != nil {
		return nil, fmt.Errorf("error concatonating final pdf: %v", err)
	}

	return Output, nil
}
