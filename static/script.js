function addCaution() {
    if (typeof addCaution.count == 'undefined') {
        addCaution.count = 0;
    }
    var temp = document.getElementById("caution").content;
    var clon = document.importNode(temp, true);
    document.getElementById("cautions").appendChild(clon);
    document.getElementsByName("CautionPlayerName")[0].setAttribute("name", "Cautions."+addCaution.count+".PlayerName");
    document.getElementsByName("CautionPlayerRole")[0].setAttribute("name", "Cautions."+addCaution.count+".PlayerRole");
    document.getElementsByName("CautionPlayerID")[0].setAttribute("name", "Cautions."+addCaution.count+".PlayerID");
    document.getElementsByName("CautionTeam")[0].setAttribute("name", "Cautions."+addCaution.count+".Team");
    document.getElementsByName("CautionCode")[0].setAttribute("name", "Cautions."+addCaution.count+".Code");
    addCaution.count++;
}
function addRed() {
    if (typeof addRed.count == 'undefined') {
        addRed.count = 0;
    }
    var temp = document.getElementById("red").content;
    var clon = document.importNode(temp, true);
    document.getElementById("reds").appendChild(clon);
    document.getElementsByName("RedPlayerName")[0].setAttribute("name", "Sendoffs."+addRed.count+".PlayerName");
    document.getElementsByName("RedPlayerRole")[0].setAttribute("name", "Sendoffs."+addRed.count+".PlayerRole");
    document.getElementsByName("RedPlayerID")[0].setAttribute("name", "Sendoffs."+addRed.count+".PlayerID");
    document.getElementsByName("RedTeam")[0].setAttribute("name", "Sendoffs."+addRed.count+".Team");
    document.getElementsByName("RedCode")[0].setAttribute("name", "Sendoffs."+addRed.count+".Code");
    addRed.count++;
}
function addSupplemental() {
    if (typeof addSupplemental.count == 'undefined') {
        addSupplemental.count = 0;
    }
    var temp = document.getElementById("supplemental").content;
    var clon = document.importNode(temp, true);
    document.getElementById("supplementals").appendChild(clon);
    document.getElementsByName("SupplementalStatement")[0].setAttribute("name", "Supplementals."+addSupplemental.count+".Statement");
    document.getElementsByName("SupplementalLocationX")[0].setAttribute("name", "Supplementals."+addSupplemental.count+".LocationX");
    document.getElementsByName("SupplementalLocationY")[0].setAttribute("name", "Supplementals."+addSupplemental.count+".LocationY");
    addSupplemental.count++;
}
function addEmail() {
    var temp = document.getElementById("email").content;
    var clon = document.importNode(temp, true);
    document.getElementById("emails").appendChild(clon);
}
var config = {
    enableTime: true,
    disableMobile: "true",
    altInput: true,
    altFormat: "F j, Y h:i K",
    dateFormat: "U",
}
flatpickr("input[type=datetime-local]", config)
