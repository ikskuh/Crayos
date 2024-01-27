const NoSession = -1;
const Open = 1;

let socket;
let sessionID = NoSession;

let ServerSideDisconnect = false;

function init() {
    initSocket();
    initPainter();
}

function initSocket() {
    document.getElementById("connecting").style.display = "flow";
    socket = new WebSocket("ws://192.168.0.100:8090");

    socket.onopen = (event) => {
        initTitle();
    };
    socket.addEventListener("error", (event) => {
        console.log("WebSocket error: ", event);
    });

    setInterval(timeoutCheck, 3000);
}

function timeoutCheck() {
    if (socket.readyState != Open & ServerSideDisconnect == false) {
        document.getElementById("connection_failed").style.display = "flow";
        hideSection("connecting");
    }
}

function btnReconnect() {
    location.reload();
}

function hideSection(id) {
    document.getElementById(id).style.display = "none";
}

function showSection(id) {
    document.getElementById(id).style.display = "flow";
}