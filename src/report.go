package main

import (
    "fmt"
    "net/mail"
    "time"
    "net"
    "strings"
    "reflect"
    "strconv"
)

type Sanction struct {
    PlayerRole       string
    PlayerName       string
    PlayerID         string
    Team             string
    Code             string
}

type Supplemental struct {
    Statement   string
    LocationX   string
    LocationY   string
}

type POSTReport struct {

    HomeTeamState           string
    HomeTeamName            string
    HomeTeamScore           string
    AwayTeamState           string
    AwayTeamName            string
    AwayTeamScore           string

    PlayerSex               string
    PlayerAge               string

    GameAssociation         string
    GameDivision            string
    GameLeague              string
    GameNumber              string
    GameDate                time.Time
    SubmittedDate           time.Time

    RefereeName             string
    RefereeGrade            string
    AssistantReferee1Name   string
    AssistantReferee1Grade  string
    AssistantReferee2Name   string
    AssistantReferee2Grade  string
    FourthOfficialName      string
    FourthOfficialGrade     string

    Cautions                []Sanction
    SendOffs                []Sanction
    Supplementals           []Supplemental

    SendToEmail             []string

    ReporterName            string
    ReporterUSSFID          string
    ReporterPhone           string
    ReporterEmail           string

    ReportID                string

}

//TODO the rest of this function
func (r *POSTReport) SanitizePostData() {

    if  len(r.Cautions) > 30 {
        r.Cautions = r.Cautions[:30]
    }

    if  len(r.SendOffs) > 15 {
        r.SendOffs = r.SendOffs[:15]
    }

    if  len(r.Supplementals) > 5 {
        r.Supplementals = r.Supplementals[:5]
    }

    r.ReportID = "testing7890"

}

func DateConverter(POSTString string) (reflect.Value) {
    IntTime, err := strconv.ParseInt(POSTString, 10, 64)
    if err != nil {
        return reflect.ValueOf(time.Now())
    }
    return reflect.ValueOf(time.Unix(IntTime, 0))
}









type refereeReport struct {

    HomeTeamName            string      `conform:"upper,trim"`
    HomeTeamState           string
    HomeTeamScore           string      `conform:"num"`             
    HomeTeam                string      `schema:"-"`
    AwayTeamName            string      `conform:"upper,trim"`      
    AwayTeamState           string
    AwayTeamScore           string      `conform:"num"`
    AwayTeam                string      `schema:"-"`

    /**************************NEEDS ATTENTION**********************************/
    GameAssociationLeague   string      `schema:"-"`
    GameDivisionAgeGroup    string      `schema:"-"`
    GameNumber              string      `conform:"upper,trim"`      //  TODO san
    GameAssociation         string      `conform:"name"`
    GameLeague              string      `conform:"name"`
    GameDivision            string                                  //  TODO san        
    PlayerSex               string
    PlayerAge               string      `schema:"-"`                //  front end redo
    PlayerAgeOverUnder      string      `conform:"name"`            //  front end redo
    PlayerAgeNumber         string      `conform:"num"`             //  front end redo

    GameTimeString          string      //`schema:"-"`              //  frontend redo
    GameTime                time.Time   `schema:"-"`                
    /**************************NEEDS ATTENTION**********************************/

    RefereeName             string      `conform:"name"`            
    RefereeGrade            string
    AssistantReferee1Name   string      `conform:"name"`
    AssistantReferee1Grade  string
    AssistantReferee2Name   string      `conform:"name"`
    AssistantReferee2Grade  string
    FourthOfficialName      string      `conform:"name"`
    FourthOfficialGrade     string

    CautionPlayerRole       []string                        
    CautionPlayerName       []string    `conform:"name"`
    CautionPlayerID         []string    `conform:"num"`             //  I belive that coaches and players are just ints but unsure tbh
    CautionTeam             []string                                //  frontend redo
    CautionCode             []string

    RedPlayerRole           []string                                //  link to supps??
    RedPlayerName           []string    `conform:"name"`
    RedPlayerID             []string    `conform:"num"`             //  tag for bench personel vs player
    RedTeam                 []string
    RedCode                 []string

    /**************************NEEDS ATTENTION**********************************/
    SupplementalStatement   []string                                //  needs san
    SupplementalLocationX   []string    `conform:"num"`
    SupplementalLocationY   []string    `conform:"num"`
    SupplementalLocation    []string    `schema:"-"`
    /**************************NEEDS ATTENTION**********************************/

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

    SanitizeSupplementalSlices(r)
    SanitizeSanctionSlices(r)
    SanitizeSendToEmailSlice(r)

    SanitizeStatement (r)
    SanitizeTeamStates(r)
    SanitizeTeamNames(r)
    SanitizeContactEmail(r)
    SanitizeSendToEmailAddress(r)
    SanitizeRefereeGrades (r)
    SanitizePlayerRoles(r)
    SanitizeCautionCodes(r)
    SanitizeRedCodes(r)
    SanitizeSanctionTeam(r)
    SanitizeSex(r)
    SanitizeGameNumber(r)

    FormatSubmittedTime(r)
    FormatPlayerAge(r)
    FormatAssociationLeague(r)
    FormatDivisionAgeGroup(r)
    FormatTeams(r)

}

func SanitizeGameNumber (r *refereeReport) {
    whitelist := "0123456789 ABCDEFGHIJKLMNOPQRSTUVWXYZ-"

    InWhitelist := func (w string) func(rune) rune {
        return func(r rune) rune {
            if strings.ContainsRune(w, r) {
                return r
            }
            return -1
        }
    }

    r.GameNumber = strings.Map(InWhitelist(whitelist), r.GameNumber)
}

func SanitizeSex (r *refereeReport) {
    if r.PlayerSex == "Male"   ||
    r.PlayerSex == "Female" ||
    r.PlayerSex == "Co-ed"  {
        r.PlayerSex = ""
    }
}

func SanitizeSanctionTeam (r *refereeReport) {

    Teams := make(map[string]struct{})
    var Exists struct{}

    Teams[r.HomeTeamName] = Exists
    Teams[r.AwayTeamName] = Exists

    for i := range r.CautionTeam {
        _, in := Teams[r.CautionTeam[i]]
        if !in {
            r.CautionTeam[i] = ""
        }
    }

    for i := range r.RedTeam {
        _, in := Teams[r.RedTeam[i]]
        if !in {
            r.RedTeam[i] = ""
        }
    }

}

func SanitizeCautionCodes (r *refereeReport) {

    Codes := make(map[string]struct{})
    var Exists struct{}

    Codes["Bench Personnel"]                        = Exists
    Codes["Dissent"]                                = Exists
    Codes["Persistant Offenses"]                    = Exists
    Codes["Delaying the Restart"]                   = Exists
    Codes["Failure to Respect Distance"]            = Exists
    Codes["Entering Re-entering or Leaving"]        = Exists
    Codes["Entering the referee review area"]       = Exists
    Codes["Excessively using the 'review' signal"]  = Exists

    for i := range r.CautionCode {
        _, in := Codes[r.CautionCode[i]]
        if !in {
            r.CautionCode[i] = ""
        }
    }

}

func SanitizeRedCodes (r *refereeReport) {

    Codes := make(map[string]struct{})
    var Exists struct{}

    Codes["Serious Foul Play"]                  = Exists
    Codes["Violent Conduct"]                    = Exists
    Codes["Spitting or Biting"]                 = Exists
    Codes["Denying a goal by handling"]         = Exists
    Codes["DOGSO"]                              = Exists
    Codes["Abusive Language/Gestures"]          = Exists
    Codes["Receiving a Second Caution"]         = Exists
    Codes["Entering the video operation room"]  = Exists

    for i := range r.RedCode {
        _, in := Codes[r.RedCode[i]]
        if !in {
            r.RedCode[i] = ""
        }
    }

}

func SanitizePlayerRoles (r *refereeReport) {

    Roles := make(map[string]struct{})
    var Exists struct{}

    Roles["Bench Personnel"] = Exists
    Roles["Player"] = Exists

    for i := range r.CautionPlayerRole {
        _, in := Roles[r.CautionPlayerRole[i]]
        if !in {
            r.CautionPlayerRole[i] = ""
        }
    }

    for i := range r.RedPlayerRole {
        _, in := Roles[r.RedPlayerRole[i]]
        if !in {
            r.RedPlayerRole[i] = ""
        }
    }

}

func SanitizeRefereeGrades (r *refereeReport) {

    grades := make(map[string]struct{})
    var exists struct{}

    grades["Grassroots"] = exists
    grades["Regional"] = exists
    grades["Regional Emeritus"] = exists
    grades["National"] = exists
    grades["National Emeritus"] = exists
    grades["PRO"] = exists
    grades["FIFA"] = exists

    _, in := grades[r.RefereeGrade]
    if !in {
        r.RefereeGrade = ""
    }
    _, in = grades[r.AssistantReferee1Grade]
    if !in {
        r.AssistantReferee1Grade = ""
    }
    _, in = grades[r.AssistantReferee2Grade]
    if !in {
        r.AssistantReferee2Grade = ""
    }
    _, in = grades[r.FourthOfficialGrade]
    if !in {
        r.FourthOfficialGrade = ""
    }
}

func SanitizeStatement (r *refereeReport) {

    whitelist := "0123456789 ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz,.\"'`!?*-()"
    //TODO does puncuation need to be cleaned??

    InWhitelist := func (w string) func(rune) rune {
        return func(r rune) rune {
            if strings.ContainsRune(w, r) {
                return r
            }
            return -1
        }
    }

    for i := range r.SupplementalStatement {
        r.SupplementalStatement[i] = strings.Map(InWhitelist(whitelist), r.SupplementalStatement[i])
    }

}

func FormatTeams (r *refereeReport) {

    r.HomeTeam = r.HomeTeamState + ": " + r.HomeTeamName
    r.AwayTeam = r.AwayTeamState + ": " + r.AwayTeamName

}

func SanitizeTeamStates (r *refereeReport) {

    states := make(map[string]struct{})
    var exists struct{}

    states["AL"] = exists
    states["AK"] = exists
    states["AZ"] = exists
    states["AR"] = exists
    states["CA"] = exists
    states["CO"] = exists
    states["CT"] = exists
    states["DE"] = exists
    states["FL"] = exists
    states["GA"] = exists
    states["HI"] = exists
    states["ID"] = exists
    states["IL"] = exists
    states["IN"] = exists
    states["IA"] = exists
    states["KS"] = exists
    states["KY"] = exists
    states["LA"] = exists
    states["ME"] = exists
    states["MD"] = exists
    states["MA"] = exists
    states["MI"] = exists
    states["MN"] = exists
    states["MS"] = exists
    states["MO"] = exists
    states["MT"] = exists
    states["NE"] = exists
    states["NV"] = exists
    states["NH"] = exists
    states["NJ"] = exists
    states["NM"] = exists
    states["NY"] = exists
    states["NC"] = exists
    states["ND"] = exists
    states["OH"] = exists
    states["OK"] = exists
    states["OR"] = exists
    states["PA"] = exists
    states["RI"] = exists
    states["SC"] = exists
    states["SD"] = exists
    states["TN"] = exists
    states["TX"] = exists
    states["UT"] = exists
    states["VT"] = exists
    states["VA"] = exists
    states["WA"] = exists
    states["WV"] = exists
    states["WI"] = exists
    states["WY"] = exists

    _, in := states[r.HomeTeamState]
    if !in {
        r.HomeTeamState = ""
    }
    _, in  = states[r.AwayTeamState];
    if !in {
        r.AwayTeamState = ""
    }

}

func SanitizeTeamNames (r *refereeReport) {

    whitelist := "0123456789 ABCDEFGHIJKLMNOPQRSTUVWXYZ"

    InWhitelist := func (w string) func(rune) rune {
        return func(r rune) rune {
            if strings.ContainsRune(w, r) {
                return r
            }
            return -1
        }
    }

    r.HomeTeamName = strings.Map(InWhitelist(whitelist), r.HomeTeamName)
    r.AwayTeamName = strings.Map(InWhitelist(whitelist), r.AwayTeamName)

}

func SanitizeSanctionSlices (r *refereeReport) {

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

func SanitizeSupplementalSlices (r *refereeReport) {

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
    for i := range r.SupplementalLocation {
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

func SanitizeSendToEmailSlice (r *refereeReport) {

    SanitizeSlice := func (slice *[]string, size int) {
        if size < len(*slice) {
            *slice = (*slice)[:size]
        }
    }

    SanitizeSlice(&r.SendToEmail, 25)

}

func SanitizeSendToEmailAddress (r *refereeReport) {

    for i := range r.SendToEmail {
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
