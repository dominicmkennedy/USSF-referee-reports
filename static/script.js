function addCaution() {
    var temp = document.getElementById("caution").content;
    var clon = document.importNode(temp, true);
    document.getElementById("cautions").appendChild(clon);
}
function addRed() {
    var temp = document.getElementById("red").content;
    var clon = document.importNode(temp, true);
    document.getElementById("reds").appendChild(clon);
}
function addSupplemental() {
    var temp = document.getElementById("supplemental").content;
    var clon = document.importNode(temp, true);
    document.getElementById("supplementals").appendChild(clon);
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
