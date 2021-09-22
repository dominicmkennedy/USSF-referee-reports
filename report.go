package main

import (
        "fmt"
        "net/mail"
        "time"
        "crypto/md5"

)

//  maybe break it up into a user report struct, and a metadatastruct
//  either switch to go-sanitize use both go-san and lee-conform
//  maybe this is a first pass of data clean up
//  then I concattonate fields and trim slices etc
type refereeReport struct {

        HomeTeamName            string      `conform:"name"`            //  team names can have nums so may need to fix that
        HomeTeamScore           string      `conform:"num"`             //  used a string bc easy to operate on, if bugs arise, I will change this
        AwayTeamName            string      `cotform:"name"`
        AwayTeamScore           string      `conform:"num"`

        PlayerSex               string      `conform:"name"`            //  maybe a struct or enum to force it
        PlayerAgeOverUnder      string      `conform:"name"`            //  TODO front end redo to add over/under, then I'll concat in san pt2.
        PlayerAgeNumber         string      `conform:"num"`             //  TODO front end redo to add over/under, then I'll concat in san pt2.
        GameLeague              string      `conform:"name"`
        GameDivision            string      `conform:"!js!html"`        
        GameNumber              string      `conform:"!js!html"`
        GameDate                string      `conform:"!js!html"`        //  change frontend, then jam into a go default type
        GameAssociationLeague   string      `schema:"-"`
        GameDivisionAgeGroup    string      `schema:"-"`
        PlayerAge               string      `schema:"-"`                //  TODO front end redo to add over/under, then I'll concat in san pt2.

        RefereeName             string      `conform:"name"`            //  referees struct good for DB?? I'll worry about that when I get to the DB 
        RefereeGrade            string      `conform:"name"`            //  struct or enum??
        AssistantReferee1Name   string      `conform:"name"`
        AssistantReferee1Grade  string      `conform:"name"`
        AssistantReferee2Name   string      `conform:"name"`
        AssistantReferee2Grade  string      `conform:"name"`
        FourthOfficialName      string      `conform:"name"`
        FourthOfficialGrade     string      `conform:"name"`

        CautionPlayerName       []string    `conform:"name"`            //  struct for YC/RC too?? same as with refs we'll see
        CautionPlayerID         []string    `conform:"num"`             //  I belive that coaches and players are just ints but unsure tbh
        CautionTeam             []string    `conform:"name`             //  bool? or js func to set name?
        CautionCode             []string    `conform:"name"`            //  struct or enum??

        RedPlayerName           []string    `conform:"name"`            //  same considerations as before but also maybe bundle supplamentals
        RedPlayerID             []string    `conform:"num"`             //  also consider putting a tag for bench personel vs player
        RedTeam                 []string    `conform:"name"`
        RedCode                 []string    `conform:"name"`

        SupplementalStatement   []string    `conform:"!js,!html"`       //  frontend redo statement maybe char limit&& remove '\n' ++bundle into a struct??
        SupplementalLocationX   []string    `conform:"num"`
        SupplementalLocationY   []string    `conform:"num"`
        SupplementalLocation    []string    `schema:"-"`

        SendToEmail             []string                                //  builtin email struct?? ++best way to sanitize it
        Name                    string      `conform:"name"`
        USSFID                  string      `conform:"num"`             //  frontend redo
        ContactNumber           string      `conform:"num"`             //  frontend redo (int'nl numbers, extensions)
        ContactEmail            string      

        //  data that goes into the report without referee input
        //  ip address
        SubmittedTime           time.Time   `schema:"-"`
        SubmittedTimeString     string      `schema:"-"`
        ReportID                string      `schema:"-"`

        pageA                   int         `schema:"-"`
        pageB                   int         `schema:"-"`
}


func SanitizeSlice(slice *[]string, size int) () {
        if size > len(*slice) {
                padding := make([]string, size-len(*slice))
                *slice = append(*slice, padding...)
        } else
        if size < len(*slice) {
                *slice = (*slice)[:size]
        }
}

//  break this out into many funcs
func (r *refereeReport) SanitizePostData() {
        fmt.Printf("\ndata is being operated on\n")

        r.SubmittedTime = time.Now()
        r.SubmittedTimeString = fmt.Sprintf("%s, %d %d %d:%d %s", 
        r.SubmittedTime.Month(),
        r.SubmittedTime.Day(),
        r.SubmittedTime.Year(),
        r.SubmittedTime.Hour(),
        r.SubmittedTime.Minute(), 
        r.SubmittedTime.Location())


        //  subject to change
        r.ReportID = fmt.Sprintf("%x", md5.Sum([]byte(r.SubmittedTime.String())))

        Max := func (x, y int) int {
                if x < y {
                        return y
                }
                return x
        }
        //  enforce num of reds yellows and supps
        //  and set nil//empty strings as needed
        size := 5
        size = Max(len(r.RedPlayerName), size)
        size = Max(len(r.RedPlayerID), size)
        size = Max(len(r.RedTeam), size)
        size = Max(len(r.RedCode), size)

        size *= 2
        size = Max(len(r.CautionPlayerName), size)
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
        SanitizeSlice(&r.CautionPlayerID, size)                    
        SanitizeSlice(&r.CautionTeam, size)                    
        SanitizeSlice(&r.CautionCode, size)                    

        SanitizeSlice(&r.RedPlayerName, size/2)                    
        SanitizeSlice(&r.RedPlayerID, size/2)                    
        SanitizeSlice(&r.RedTeam, size/2)                    
        SanitizeSlice(&r.RedCode, size/2)                    


        //  do somthing with the date :/
        //  switch san libs
        //  add association??

        //  sanitize email
        e, err := mail.ParseAddress(r.ContactEmail);
        if err != nil { r.ContactEmail = "not a valid email"
} else { r.ContactEmail = e.Address }

for i, _ := range r.SendToEmail {
        e, err := mail.ParseAddress(r.SendToEmail[i]);
        if err != nil { r.SendToEmail[i] = "not a valid email"
} else { r.SendToEmail[i] = e.Address }
    }

    //  pack up age
    r.PlayerAge = r.PlayerAgeOverUnder + " " + r.PlayerAgeNumber


    //  pack up division, leauge, assoc, etc
    r.GameAssociationLeague = r.GameLeague
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


    //  make sure location = satement and the slice size < 25
    //  prfile to make sure that doesnt take too long to load
    size = 0
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
    //  pack up x and y loc and do the math
    r.SupplementalLocation = make([]string, len(r.SupplementalLocationX))
    for i, _ := range r.SupplementalLocation {
            r.SupplementalLocation[i] = r.SupplementalLocationX[i] + " + " + r.SupplementalLocationY[i]
    }
}
