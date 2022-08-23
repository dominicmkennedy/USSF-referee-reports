function addCaution() {
    if (typeof addCaution.count == 'undefined') {
        addCaution.count = 0;
    }
    if (addCaution.count >= 30) {
        alert("Maximum of 30 cautions.");
        return;
    }
    var temp = document.getElementById("caution").content;
    var clon = document.importNode(temp, true);
    document.getElementById("cautions").appendChild(clon);
    document.getElementsByName("CautionPlayerName")[0].setAttribute("name", "Cautions."+addCaution.count+".PlayerName");
    document.getElementsByName("CautionPlayerID")[0].setAttribute("name", "Cautions."+addCaution.count+".PlayerID");
    document.getElementsByName("CautionTeam")[0].setAttribute("name", "Cautions."+addCaution.count+".Team");
    
    let homeTeamName = document.getElementById("HomeTeamName").value;
    if (homeTeamName.length === 0) {
        homeTeamName = "Home Team";
    }

    let awayTeamName = document.getElementById("AwayTeamName").value;
    if (awayTeamName.length === 0) {
        awayTeamName = "Away Team";
    }
    
    let cautionTeamBox = document.getElementsByName("Cautions."+addCaution.count+".Team")[0];
    while (cautionTeamBox.firstChild) {
        cautionTeamBox.removeChild(cautionTeamBox.firstChild);
    }
    let el0 = document.createElement("option");
    el0.textContent = "";
    el0.value = "";
    el0.setAttribute("selected", true)
    cautionTeamBox.appendChild(el0)
    let el1 = document.createElement("option");
    el1.textContent = homeTeamName;
    el1.value = homeTeamName;
    cautionTeamBox.appendChild(el1)
    let el2 = document.createElement("option");
    el2.textContent = awayTeamName;
    el2.value = awayTeamName;
    cautionTeamBox.appendChild(el2)

    document.getElementsByName("CautionCode")[0].setAttribute("name", "Cautions."+addCaution.count+".Code");
    addCaution.count++;
}

function addRed() {
    if (typeof addRed.count == 'undefined') {
        addRed.count = 0;
    }
    if (addRed.count >= 15) {
        alert("Maximum of 15 send offs.");
        return;
    }
    var temp = document.getElementById("red").content;
    var clon = document.importNode(temp, true);
    document.getElementById("reds").appendChild(clon);
    document.getElementsByName("RedPlayerName")[0].setAttribute("name", "Sendoffs."+addRed.count+".PlayerName");
    document.getElementsByName("RedPlayerID")[0].setAttribute("name", "Sendoffs."+addRed.count+".PlayerID");
    document.getElementsByName("RedTeam")[0].setAttribute("name", "Sendoffs."+addRed.count+".Team");
    
    let homeTeamName = document.getElementById("HomeTeamName").value;
    if (homeTeamName.length === 0) {
        homeTeamName = "Home Team";
    }

    let awayTeamName = document.getElementById("AwayTeamName").value;
    if (awayTeamName.length === 0) {
        awayTeamName = "Away Team";
    }
    
    let sendoffTeamBox = document.getElementsByName("Sendoffs."+addRed.count+".Team")[0];
    while (sendoffTeamBox.firstChild) {
        sendoffTeamBox.removeChild(sendoffTeamBox.firstChild);
    }
    let el0 = document.createElement("option");
    el0.textContent = "";
    el0.value = "";
    el0.setAttribute("selected", true)
    sendoffTeamBox.appendChild(el0)
    let el1 = document.createElement("option");
    el1.textContent = homeTeamName;
    el1.value = homeTeamName;
    sendoffTeamBox.appendChild(el1)
    let el2 = document.createElement("option");
    el2.textContent = awayTeamName;
    el2.value = awayTeamName;
    sendoffTeamBox.appendChild(el2)
   
    document.getElementsByName("RedCode")[0].setAttribute("name", "Sendoffs."+addRed.count+".Code");
    addRed.count++;
}

function addSupplemental() {
    if (typeof addSupplemental.count == 'undefined') {
        addSupplemental.count = 0;
    }
    if (addSupplemental.count >= 5) {
        alert("Maximum of 5 supplemental reports.");
        return;
    }
    var temp = document.getElementById("supplemental").content;
    var clon = document.importNode(temp, true);
    document.getElementById("supplementals").appendChild(clon);
    document.getElementsByName("SupplementalImg")[0].setAttribute("id", addSupplemental.count*4+0);
    document.getElementsByName("SupplementalImg")[0].removeAttribute("name");
    document.getElementsByName("SupplementalStatement")[0].setAttribute("name", "Supplementals."+addSupplemental.count+".Statement");
    document.getElementsByName("SupplementalLocationX")[0].setAttribute("id", addSupplemental.count*4+1);
    document.getElementsByName("SupplementalLocationX")[0].setAttribute("name", "Supplementals."+addSupplemental.count+".LocationX");
    document.getElementsByName("SupplementalLocationY")[0].setAttribute("id", addSupplemental.count*4+2);
    document.getElementsByName("SupplementalLocationY")[0].setAttribute("name", "Supplementals."+addSupplemental.count+".LocationY");
    document.getElementsByName("SupplementalDiv")[0].setAttribute("id", addSupplemental.count*4+3);
    document.getElementsByName("SupplementalDiv")[0].removeAttribute("name");
    addSupplemental.count++;
}

function addEmail() {
    var temp = document.getElementById("email").content;
    var clon = document.importNode(temp, true);
    document.getElementById("emails").appendChild(clon);
}

function updateTeams() {
    let homeTeamName = document.getElementById("HomeTeamName").value;
    if (homeTeamName.length === 0) {
        homeTeamName = "Home Team";
    }

    let awayTeamName = document.getElementById("AwayTeamName").value;
    if (awayTeamName.length === 0) {
        awayTeamName = "Away Team";
    }

    let cautionTeamBoxes = [];
    for (let i=0; i<addCaution.count; ++i) {
        cautionTeamBoxes.push(document.getElementsByName("Cautions."+i+".Team")[0]);
    }
    cautionTeamBoxes.forEach((elem, _) => {
        while (elem.firstChild) {
            elem.removeChild(elem.firstChild);
        }
        let el0 = document.createElement("option");
        el0.textContent = "";
        el0.value = "";
        el0.setAttribute("selected", true)
        elem.appendChild(el0)
        let el1 = document.createElement("option");
        el1.textContent = homeTeamName;
        el1.value = homeTeamName;
        elem.appendChild(el1)
        let el2 = document.createElement("option");
        el2.textContent = awayTeamName;
        el2.value = awayTeamName;
        elem.appendChild(el2)
    });

    let sendOffTeamBoxes = [];
    for (let i=0; i<addRed.count; ++i) {
        sendOffTeamBoxes.push(document.getElementsByName("Sendoffs."+i+".Team")[0]);
    }
    sendOffTeamBoxes.forEach((elem, _) => {
        while (elem.firstChild) {
            elem.removeChild(elem.firstChild);
        }
        let el0 = document.createElement("option");
        el0.textContent = "";
        el0.value = "";
        el0.setAttribute("selected", true)
        elem.appendChild(el0)
        let el1 = document.createElement("option");
        el1.textContent = homeTeamName;
        el1.value = homeTeamName;
        elem.appendChild(el1)
        let el2 = document.createElement("option");
        el2.textContent = awayTeamName;
        el2.value = awayTeamName;
        elem.appendChild(el2)
    });
}

function FindPosition(oElement) {
    if(typeof( oElement.offsetParent ) != "undefined") {
        for(var posX = 0, posY = 0; oElement; oElement = oElement.offsetParent) {
            posX += oElement.offsetLeft;
            posY += oElement.offsetTop;
        }
        return [ posX, posY ];
    }
    else
    {
        return [ oElement.x, oElement.y ];
    }
}

function GetCoordinates(myImg, e) {
    var num = parseInt(myImg.id);
    var PosX = 0;
    var PosY = 0;
    var ImgPos;
    ImgPos = FindPosition(myImg);
    if (!e) var e = window.event;
    if (e.pageX || e.pageY)
    {
        PosX = e.pageX;
        PosY = e.pageY;
    }
    else if (e.clientX || e.clientY)
    {
        PosX = e.clientX + document.body.scrollLeft
            + document.documentElement.scrollLeft;
        PosY = e.clientY + document.body.scrollTop
            + document.documentElement.scrollTop;
    }
    PosX = PosX - ImgPos[0];
    PosY = PosY - ImgPos[1];
    var p1 = num+1;
    var p2 = num+2;
    var p3 = num+3;
    var width = myImg.offsetWidth;
    var height = myImg.offsetHeight;
    var RelX = (47 * PosX) / width;
    var RelY = (16 * PosY) / height;
    var box = document.getElementById(p3.toString());
    box.setAttribute("style", "position:absolute; top:"+(PosY-8)+"px; left:"+(PosX-8)+"px; height: 15px; width: 15px; background-color:red; border-radius: 50%; display: inline-block;");
    document.getElementById(p1.toString()).value = RelX;
    document.getElementById(p2.toString()).value = RelY;
}
