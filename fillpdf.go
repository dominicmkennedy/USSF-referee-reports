package main

import (
        "os"
        "bytes"
        "strconv"
        "log"
        //"io"

        "github.com/patiek/go-pdftools/pdftk"
        "github.com/patiek/go-pdftools/fdf"
)

//func fillPg1(form refereeReport, page int) (output io.Writer) {
func fillPg1(form refereeReport, page int) {

        var b bytes.Buffer
        err := fdf.Write(&b, fdf.Inputs{

                "HomeTeamName":                 form.HomeTeam,
                "HomeTeamScore":                form.HomeTeamScore,
                "AwayTeamName":                 form.AwayTeamName,
                "AwayTeamScore":                form.AwayTeam,

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

                "CautionPlayerName0":           form.CautionPlayerName[page*10+0],
                "CautionPlayerName1":           form.CautionPlayerName[page*10+1],
                "CautionPlayerName2":           form.CautionPlayerName[page*10+2],
                "CautionPlayerName3":           form.CautionPlayerName[page*10+3],
                "CautionPlayerName4":           form.CautionPlayerName[page*10+4],
                "CautionPlayerName5":           form.CautionPlayerName[page*10+5],
                "CautionPlayerName6":           form.CautionPlayerName[page*10+6],
                "CautionPlayerName7":           form.CautionPlayerName[page*10+7],
                "CautionPlayerName8":           form.CautionPlayerName[page*10+8],
                "CautionPlayerName9":           form.CautionPlayerName[page*10+9],
                "CautionPlayerID0":             form.CautionPlayerID[page*10+0],
                "CautionPlayerID1":             form.CautionPlayerID[page*10+1],
                "CautionPlayerID2":             form.CautionPlayerID[page*10+2],
                "CautionPlayerID3":             form.CautionPlayerID[page*10+3],
                "CautionPlayerID4":             form.CautionPlayerID[page*10+4],
                "CautionPlayerID5":             form.CautionPlayerID[page*10+5],
                "CautionPlayerID6":             form.CautionPlayerID[page*10+6],
                "CautionPlayerID7":             form.CautionPlayerID[page*10+7],
                "CautionPlayerID8":             form.CautionPlayerID[page*10+8],
                "CautionPlayerID9":             form.CautionPlayerID[page*10+9],
                "CautionTeam0":                 form.CautionTeam[page*10+0],
                "CautionTeam1":                 form.CautionTeam[page*10+1],
                "CautionTeam2":                 form.CautionTeam[page*10+2],
                "CautionTeam3":                 form.CautionTeam[page*10+3],
                "CautionTeam4":                 form.CautionTeam[page*10+4],
                "CautionTeam5":                 form.CautionTeam[page*10+5],
                "CautionTeam6":                 form.CautionTeam[page*10+6],
                "CautionTeam7":                 form.CautionTeam[page*10+7],
                "CautionTeam8":                 form.CautionTeam[page*10+8],
                "CautionTeam9":                 form.CautionTeam[page*10+9],
                "CautionCode0":                 form.CautionCode[page*10+0],
                "CautionCode1":                 form.CautionCode[page*10+1],
                "CautionCode2":                 form.CautionCode[page*10+2],
                "CautionCode3":                 form.CautionCode[page*10+3],
                "CautionCode4":                 form.CautionCode[page*10+4],
                "CautionCode5":                 form.CautionCode[page*10+5],
                "CautionCode6":                 form.CautionCode[page*10+6],
                "CautionCode7":                 form.CautionCode[page*10+7],
                "CautionCode8":                 form.CautionCode[page*10+8],
                "CautionCode9":                 form.CautionCode[page*10+9],

                "RedPlayerName0":               form.RedPlayerName[page*5+0],
                "RedPlayerName1":               form.RedPlayerName[page*5+1],
                "RedPlayerName2":               form.RedPlayerName[page*5+2],
                "RedPlayerName3":               form.RedPlayerName[page*5+3],
                "RedPlayerName4":               form.RedPlayerName[page*5+4],
                "RedPlayerID0":                 form.RedPlayerID[page*5+0],
                "RedPlayerID1":                 form.RedPlayerID[page*5+1],
                "RedPlayerID2":                 form.RedPlayerID[page*5+2],
                "RedPlayerID3":                 form.RedPlayerID[page*5+3],
                "RedPlayerID4":                 form.RedPlayerID[page*5+4],
                "RedTeam0":                     form.RedTeam[page*5+0],
                "RedTeam1":                     form.RedTeam[page*5+1],
                "RedTeam2":                     form.RedTeam[page*5+2],
                "RedTeam3":                     form.RedTeam[page*5+3],
                "RedTeam4":                     form.RedTeam[page*5+4],
                "RedCode0":                     form.RedCode[page*5+0],
                "RedCode1":                     form.RedCode[page*5+1],
                "RedCode2":                     form.RedCode[page*5+2],
                "RedCode3":                     form.RedCode[page*5+3],
                "RedCode4":                     form.RedCode[page*5+4],

                "Name":                         form.Name,
                "USSFID":                       form.USSFID,
                "ContactNumber":                form.ContactNumber,
                "ContactEmail":                 form.ContactEmail,
                "SubmittedDate":                form.SubmittedTimeString,
        });
        if err != nil {
                log.Fatal(err)
        } 

        output, err := os.Create("tmp/" + form.ReportID + "-pg1-" + strconv.Itoa(page) + ".pdf")
        if err != nil {
                log.Fatal(err)
                os.Exit(1)
        }
        defer output.Close()

        //output = new(bytes.Buffer)

        err = pdftk.FillForm(output, "templates/pg1.pdf", &b, pdftk.OptionFlatten())
        if err != nil {
                log.Fatal(err)
                os.Exit(1)
        }

        //return output

}

func fillPg2(form refereeReport, page int) {

        var b bytes.Buffer
        err := fdf.Write(&b, fdf.Inputs{

                "HomeTeamName":                 form.HomeTeamName,
                "HomeTeamScore":                form.HomeTeamScore,
                "AwayTeamName":                 form.AwayTeamName,
                "AwayTeamScore":                form.AwayTeamScore,

                "GameDivision":                 form.GameDivisionAgeGroup,
                "GameAssociation":              form.GameAssociationLeague,
                "GameNumber":                   form.GameNumber,
                "GameDate":                     form.GameDate,

                "SupplementalStatement":        form.SupplementalStatement[page],
                "SupplementalLocation":         form.SupplementalLocation[page],

                "Name":                         form.Name,
                "USSFID":                       form.USSFID,
                "ContactNumber":                form.ContactNumber,
                "ContactEmail":                 form.ContactEmail,
                "SubmittedDate":                form.SubmittedTimeString,
        });
        if err != nil {
                log.Fatal(err)
        } 

        output, err := os.Create("tmp/" + form.ReportID + "-pg2-" + strconv.Itoa(page) + ".pdf")
        if err != nil {
                log.Fatal(err)
                os.Exit(1)
        }
        defer output.Close()

        err = pdftk.FillForm(output, "templates/pg2.pdf", &b, pdftk.OptionFlatten())
        if err != nil {
                log.Fatal(err)
                os.Exit(1)
        }
}

func writePDF(form *refereeReport) {

        outfiles := pdftk.NewInputFileMap()

        for i := 0; i<form.pageA; i++ {
                fillPg1(*form, i)
                outfiles[pdftk.InputHandleNameFromInt(i)] = ("tmp/" + form.ReportID + "-pg1-" + strconv.Itoa(i) + ".pdf")
        }

        for i := 0; i<form.pageB; i++ {
                fillPg2(*form, i)
                outfiles[pdftk.InputHandleNameFromInt(i+form.pageA)] = ("tmp/" + form.ReportID + "-pg2-" + strconv.Itoa(i) + ".pdf")
        }

        output, err := os.Create("reports/" + form.ReportID + ".pdf")
        if err != nil {
                log.Fatal(err)
                os.Exit(1)
        }
        defer output.Close()

        err = pdftk.Cat(output, outfiles, []pdftk.PageRange{}, pdftk.OptionFlatten())
        if err != nil {
                log.Fatal(err)
                os.Exit(1)
        }
}
