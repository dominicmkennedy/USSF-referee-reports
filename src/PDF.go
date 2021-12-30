package main

import (
	"bytes"
	"log"
	"os"
	"strconv"
    "time"

	"github.com/patiek/go-pdftools/fdf"
	"github.com/patiek/go-pdftools/pdftk"
)

type PDFReport struct {
    Pg1Reports  []Pg1Report
    Pg2Reports  []Pg2Report
    ReportID    string
}

type Pg1Report map[string]interface{}
type Pg2Report map[string]interface{}

func (PDF *PDFReport) FillPDF(POST POSTReport) {

    PDF.ReportID = POST.ReportID

    pg1 := 0
    if ((len(POST.Cautions) + 9) / 10) > pg1 { pg1 = ((len(POST.Cautions) + 9) / 10) }
    if ((len(POST.SendOffs) + 4) / 5) > pg1 { pg1 = ((len(POST.SendOffs) + 4) / 5) }
    if pg1 == 0 { pg1 = 1 }

    PDF.Pg1Reports = make([]Pg1Report, pg1)
    for i := 0; i < pg1; i++ {
        PDF.Pg1Reports[i].FillPDF(POST, i)
    }

    pg2 := len(POST.Supplementals)
    PDF.Pg2Reports = make([]Pg2Report, pg2)
    for i := 0; i < pg2; i++ {
        PDF.Pg2Reports[i].FillPDF(POST, i)
    }

}

func (PDF *Pg1Report) FillPDF(POST POSTReport, Page int) {

    *PDF = map[string]interface{}{
        "HomeTeamName":                 POST.HomeTeamState + ": " + POST.HomeTeamName,
        "HomeTeamScore":                POST.HomeTeamScore,
        "AwayTeamName":                 POST.AwayTeamState + ": " + POST.AwayTeamName,
        "AwayTeamScore":                POST.AwayTeamScore,

        "GameDivision":                 "Division: " + POST.GameDivision + " Sex: " + POST.PlayerSex + " Age: " + POST.PlayerAge,
        "GameAssociation":              POST.GameLeague,
        "GameNumber":                   POST.GameNumber,
        "GameDate":                     POST.GameDate.Format(time.RFC1123),

        "RefereeName":                  POST.RefereeName,
        "RefereeGrade":                 POST.RefereeGrade,
        "AssistantReferee1Name":        POST.AssistantReferee1Name,
        "AssistantReferee1Grade":       POST.AssistantReferee1Grade,
        "AssistantReferee2Name":        POST.AssistantReferee2Name,
        "AssistantReferee2Grade":       POST.AssistantReferee2Grade,
        "FourthOfficialName":           POST.FourthOfficialName,
        "FourthOfficialGrade":          POST.FourthOfficialGrade,

        "Name":                         POST.ReporterName,
        "USSFID":                       POST.ReporterUSSFID,
        "ContactNumber":                POST.ReporterPhone,
        "ContactEmail":                 POST.ReporterEmail,
        "SubmittedDate":                POST.SubmittedDate.Format(time.RFC1123),
    }

    for iPDF, iPOST := 0, Page*10; iPDF < 10 && iPOST < len(POST.Cautions); iPDF, iPOST = iPDF+1, iPOST+1 {
        (*PDF)["CautionPlayerName"+strconv.Itoa(iPDF)] = POST.Cautions[iPOST].PlayerRole +
                                                  ": " + POST.Cautions[iPOST].PlayerName
        (*PDF)["CautionPlayerID"+strconv.Itoa(iPDF)]   = POST.Cautions[iPOST].PlayerID
        (*PDF)["CautionTeam"+strconv.Itoa(iPDF)]       = POST.Cautions[iPOST].Team
        (*PDF)["CautionCode"+strconv.Itoa(iPDF)]       = POST.Cautions[iPOST].Code
    }

    for iPDF, iPOST := 0, Page*5; iPDF < 5 && iPOST < len(POST.SendOffs); iPDF, iPOST = iPDF+1, iPOST+1 {
        (*PDF)["RedPlayerName"+strconv.Itoa(iPDF)] = POST.SendOffs[iPOST].PlayerRole +
                                              ": " + POST.SendOffs[iPOST].PlayerName
        (*PDF)["RedPlayerID"+strconv.Itoa(iPDF)]   = POST.SendOffs[iPOST].PlayerID
        (*PDF)["RedTeam"+strconv.Itoa(iPDF)]       = POST.SendOffs[iPOST].Team
        (*PDF)["RedCode"+strconv.Itoa(iPDF)]       = POST.SendOffs[iPOST].Code
    }

}

func (PDF *Pg2Report) FillPDF(POST POSTReport, Page int) {

    *PDF = map[string]interface{}{
        "HomeTeamName":                 POST.HomeTeamState + ": " + POST.HomeTeamName,
        "HomeTeamScore":                POST.HomeTeamScore,
        "AwayTeamName":                 POST.AwayTeamState + ": " + POST.AwayTeamName,
        "AwayTeamScore":                POST.AwayTeamScore,

        "GameDivision":                 "Division: " + POST.GameDivision + " Sex: " + POST.PlayerSex + " Age: " + POST.PlayerAge,
        "GameAssociation":              POST.GameLeague,
        "GameNumber":                   POST.GameNumber,
        "GameDate":                     POST.GameDate.Format(time.RFC1123),

        "SupplementalStatement":        POST.Supplementals[Page].Statement,
        "SupplementalLocation":         POST.Supplementals[Page].LocationX + " somemath " + POST.Supplementals[Page].LocationY,

        "Name":                         POST.ReporterName,
        "USSFID":                       POST.ReporterUSSFID,
        "ContactNumber":                POST.ReporterPhone,
        "ContactEmail":                 POST.ReporterEmail,
        "SubmittedDate":                POST.SubmittedDate.Format(time.RFC1123),
    }

}

func (PDF *PDFReport) WriteToPDF() {

    outfiles := pdftk.NewInputFileMap()

    for i, elem := range (*PDF).Pg1Reports {
        var b bytes.Buffer
        var vals map[string]interface{} = elem
        err := fdf.Write(&b, vals)
        if err != nil {
            log.Println(err)
            continue
        }

        FileName := "../tmp/" + PDF.ReportID + "-pg1-" + strconv.Itoa(i) + ".pdf"
        output, err := os.Create(FileName)
        if err != nil {
            log.Println(err)
        }

        if err := pdftk.FillForm(output, "../templates/pg1.pdf", &b, pdftk.OptionFlatten()); err != nil {
            log.Println(err)
        }

        outfiles[pdftk.InputHandleNameFromInt(i)] = FileName

    }

    for i, elem := range (*PDF).Pg2Reports {
        var b bytes.Buffer
        var vals map[string]interface{} = elem
        err := fdf.Write(&b, vals)
        if err != nil {
            log.Println(err)
            continue
        }

        FileName := "../tmp/" + PDF.ReportID + "-pg2-" + strconv.Itoa(i) + ".pdf"
        output, err := os.Create(FileName)
        if err != nil {
            log.Println(err)
        }

        if err := pdftk.FillForm(output, "../templates/pg2.pdf", &b, pdftk.OptionFlatten()); err != nil {
            log.Println(err)
        }

        outfiles[pdftk.InputHandleNameFromInt(i+len((*PDF).Pg1Reports))] = FileName

    }

    output, err := os.Create("../reports/" + PDF.ReportID + ".pdf")
    if err != nil {
        log.Println(err)
    }
    defer output.Close()

    if err := pdftk.Cat(output, outfiles, []pdftk.PageRange{}, pdftk.OptionFlatten()); err != nil {
        log.Println(err)
    }

}
