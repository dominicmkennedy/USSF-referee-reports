package main

import (
        "fmt"
        "net/mail"
        "time"
        "net"
)

type refereeReport struct {

        HomeTeamName            string                                  //  needs san
        HomeTeamState           string                                  //  needs san
        HomeTeamScore           string      `conform:"num"`             
        AwayTeamName            string                                  //  needs san
        AwayTeamState           string                                  //  needs san
        AwayTeamScore           string      `conform:"num"`

        GameAssociationLeague   string      `schema:"-"`
        GameDivisionAgeGroup    string      `schema:"-"`
        GameNumber              string                                  //  needs san
        GameDate                string      `schema:"-"`                //  frontend redo
        
        //  ass/league
        GameAssociation         string      `conform:"name"`
        GameLeague              string      `conform:"name"`
        
        //  division/age
        GameDivision            string                                  //  needs san        
        PlayerSex               string      `conform:"alpha"`
        PlayerAge               string      `schema:"-"`                //  front end redo
        PlayerAgeOverUnder      string      `conform:"name"`            //  front end redo
        PlayerAgeNumber         string      `conform:"num"`             //  front end redo
        
        //  date
        GameTime                time.Time   `schema:"-"`                

        RefereeName             string      `conform:"name"`            
        RefereeGrade            string      `conform:"name"`
        AssistantReferee1Name   string      `conform:"name"`
        AssistantReferee1Grade  string      `conform:"name"`
        AssistantReferee2Name   string      `conform:"name"`
        AssistantReferee2Grade  string      `conform:"name"`
        FourthOfficialName      string      `conform:"name"`
        FourthOfficialGrade     string      `conform:"name"`

        CautionPlayerRole       []string    `conform:"name"`            //  struct for YC/RC too?? same as with refs we'll see
        CautionPlayerName       []string    `conform:"name"`            //  struct for YC/RC too?? same as with refs we'll see
        CautionPlayerID         []string    `conform:"num"`             //  I belive that coaches and players are just ints but unsure tbh
        CautionTeam             []string    `conform:"name`             //  frontend redo
        CautionCode             []string    `conform:"name"`            

        RedPlayerRole           []string    `conform:"name"`            //  link to supps??
        RedPlayerName           []string    `conform:"name"`            //  link to supps??
        RedPlayerID             []string    `conform:"num"`             //  tag for bench personel vs player
        RedTeam                 []string    `conform:"name"`
        RedCode                 []string    `conform:"name"`

        SupplementalStatement   []string                                //  needs san
        SupplementalLocationX   []string    `conform:"num"`
        SupplementalLocationY   []string    `conform:"num"`
        SupplementalLocation    []string    `schema:"-"`

        SendToEmail             []string                                
        Name                    string      `conform:"name"`
        USSFID                  string      `conform:"num"`             //  frontend redo
        ContactNumber           string      `conform:"num"`             //  frontend redo 
        ContactEmail            string      

        ipaddr                  net.IP      `schema:"-"`          
        SubmittedTime           time.Time   `schema:"-"`
        SubmittedTimeString     string      `schema:"-"`
        ReportID                string      `schema:"-"`

        pageA                   int         `schema:"-"`
        pageB                   int         `schema:"-"`
}

func (r *refereeReport) SanitizePostData() {
        
        //  sanitize with regexp
        //HomeTeamName            string                                  //  needs san
        //AwayTeamName            string                                  //  needs san
        //GameDivision            string                                  //  needs san        
        //GameNumber              string                                  //  needs san
        //SupplementalStatement   []string                                //  needs san
        
        SanitizeSupplementals(r)
        SanitizeSanctions(r)
        FormatSubmittedTime(r)
        FormatPlayerAge(r)
        SanitizeContactEmail(r)
        SanitizeSendToEmail(r)
        FormatAssociationLeague(r)
        FormatDivisionAgeGroup(r)

        //  GameDate       
}

func SanitizeSanctions (r *refereeReport) {

        Max := func (x, y int) int {
                if x < y {
                        return y
                }
                return x
        }

        SanitizeSlice := func (slice *[]string, size int) {
                if size > len(*slice) {
                        padding := make([]string, size-len(*slice))
                        *slice = append(*slice, padding...)
                } else
                if size < len(*slice) {
                        *slice = (*slice)[:size]
                }
        }

        //  enforce num of reds yellows and supps
        //  and set nil//empty strings as needed
        size := 5
        size = Max(len(r.RedPlayerName), size)
        size = Max(len(r.RedPlayerRole), size)
        size = Max(len(r.RedPlayerID), size)
        size = Max(len(r.RedTeam), size)
        size = Max(len(r.RedCode), size)

        size *= 2
        size = Max(len(r.CautionPlayerName), size)
        size = Max(len(r.CautionPlayerRole), size)
        size = Max(len(r.CautionPlayerID), size)
        size = Max(len(r.CautionTeam), size)
        size = Max(len(r.CautionCode), size)

        if size % 10 != 0 {
                remainder := (size % 10)
                remainder = 10 - remainder 
                size = remainder + size 
        }

        if size > 50 {
                size = 50
        }

        r.pageA = size / 10

        SanitizeSlice(&r.CautionPlayerName, size)
        SanitizeSlice(&r.CautionPlayerRole, size)
        SanitizeSlice(&r.CautionPlayerID, size)                    
        SanitizeSlice(&r.CautionTeam, size)                    
        SanitizeSlice(&r.CautionCode, size)                    

        SanitizeSlice(&r.RedPlayerName, size/2)                    
        SanitizeSlice(&r.RedPlayerRole, size/2)                    
        SanitizeSlice(&r.RedPlayerID, size/2)                    
        SanitizeSlice(&r.RedTeam, size/2)                    
        SanitizeSlice(&r.RedCode, size/2)                    

}

func SanitizeSupplementals (r *refereeReport) {

        Max := func (x, y int) int {
                if x < y {
                        return y
                }
                return x
        }

        SanitizeSlice := func (slice *[]string, size int) {
                if size > len(*slice) {
                        padding := make([]string, size-len(*slice))
                        *slice = append(*slice, padding...)
                } else
                if size < len(*slice) {
                        *slice = (*slice)[:size]
                }
        }

        //  make sure location = satement and the slice size < 25
        //  prfile to make sure that doesnt take too long to load
        size := 0
        size = Max(len(r.SupplementalStatement), size)
        size = Max(len(r.SupplementalLocationX), size)
        size = Max(len(r.SupplementalLocationY), size)
        if size > 25 {
                size = 25
        }

        r.pageB = size

        SanitizeSlice(&r.SupplementalStatement, size)
        SanitizeSlice(&r.SupplementalLocationX, size)
        SanitizeSlice(&r.SupplementalLocationY, size)

        r.SupplementalLocation = make([]string, len(r.SupplementalLocationX))
        for i, _ := range r.SupplementalLocation {
                r.SupplementalLocation[i] = r.SupplementalLocationX[i] + " + " + r.SupplementalLocationY[i]
        }
}

func FormatSubmittedTime (r *refereeReport) {
        r.SubmittedTime = time.Now()
        r.SubmittedTimeString = fmt.Sprintf("%s, %d %d %d:%d %s", 
        r.SubmittedTime.Month(),
        r.SubmittedTime.Day(),
        r.SubmittedTime.Year(),
        r.SubmittedTime.Hour(),
        r.SubmittedTime.Minute(), 
        r.SubmittedTime.Location())
}

func FormatPlayerAge (r *refereeReport) {
        r.PlayerAge = r.PlayerAgeOverUnder + " " + r.PlayerAgeNumber
}

func SanitizeContactEmail (r *refereeReport) {
        e, err := mail.ParseAddress(r.ContactEmail);
        if err != nil { 
                r.ContactEmail = ""
        } else { 
                r.ContactEmail = e.Address 
        }
}

func SanitizeSendToEmail (r *refereeReport) {
        
        SanitizeSlice := func (slice *[]string, size int) {
                if size < len(*slice) {
                        *slice = (*slice)[:size]
                }
        }

        SanitizeSlice(&r.SendToEmail, 25)

        for i, _ := range r.SendToEmail {
                e, err := mail.ParseAddress(r.SendToEmail[i]);
                if err != nil { 
                        r.SendToEmail[i] = ""
                } else { 
                        r.SendToEmail[i] = e.Address 
                }
        }
}

func FormatAssociationLeague(r *refereeReport) {
        r.GameAssociationLeague = r.GameLeague
}

func FormatDivisionAgeGroup(r *refereeReport) {
        if len(r.GameDivision) != 0 {
                r.GameDivision = "Division: " + r.GameDivision + "   "
        }
        if len(r.PlayerSex) != 0 {
                r.PlayerSex = "Sex: " + r.PlayerSex + "   "
        }
        if len(r.PlayerAge) != 0 {
                r.PlayerAge = "Age: " + r.PlayerAge
        }
        r.GameDivisionAgeGroup = r.GameDivision + r.PlayerSex + r.PlayerAge 
}
