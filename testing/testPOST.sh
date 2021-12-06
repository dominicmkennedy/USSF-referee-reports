#!/bin/sh

curl -d "HomeTeamState=&HomeTeamName=&HomeTeamScore=&AwayTeamState=&AwayTeamName=&AwayTeamScore=&PlayerSex=&PlayerAge=&GameAssociation=&GameLeague=&GameDivision=&GameNumber=&GameDate=&SubmittedDate=NULL&RefereeName=&RefereeGrade=&AssistantReferee1Name=&AssistantReferee1Grade=&AssistantReferee2Name=&AssistantReferee2Grade=&FourthOfficialName=&FourthOfficialGrade=&ReporterName=&ReporterUSSFID=&ReporterPhone=&ReporterEmail=" -X POST http://localhost:8080/submit
