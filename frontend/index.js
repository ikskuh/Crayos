const NoSession = -1;
const Open = 1;

let socket;
let sessionID = NoSession;

let serverSideDisconnect = false;   // If this is true the client won't try to reconnect

let players = ["", "", "", ""];
let localPlayer = "nickname"
let localIsReady = false;

const gamestates = ["title", "connecting", "connection_failed", "lobby", "painer", "viewer", "troll", "voting", "winner"];
let currentGamestate = "connecting";

function init() {
    initSocket();

    const resize = (event) => {
      const wrapper = document.getElementById("wrapper");
      const scale = Math.min(window.innerWidth / width, window.innerHeight / height);
      wrapper.style.transform = "translate(-50%, -50%) scale(" + scale + ")";
    };
    resize();
    window.addEventListener("resize", resize);
}

function initSocket() {
    // TODO: remove this hack
    if (!document.getElementById("connecting")) {
        // We're in chaospaint.html, not index.html
        initPainter();
        return;
    }
    document.getElementById("connecting").style.display = "flow";
    socket = new WebSocket("ws://192.168.37.247:8090/ws");

    socket.onerror = function(event){console.log("WebSocket error: ", event);}
    socket.onopen = function(event){setView("title")};
    socket.onmessage = function(event){onSocketReceive(event)};

    setInterval(timeoutCheck, 3000);
}

function timeoutCheck() {
    if (socket.readyState != Open && serverSideDisconnect == false) {
        setView("connection_failed");
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

function setView(newState) {
    hideSection(currentGamestate);
    currentGamestate = newState;
    showSection(newState);
}

function onSocketReceive(event) {
    let data = JSON.parse(event.data);
    console.log(data);

    switch (data.type) {
        case "change-game-view-event":
            setView(data.view);
            break;
        case "players-changed-event":
            for (let i = 0; i < data.players.length; i++) {
                players[i] = data.players[i];
            }
            updateLobby();
            break;
        case "enter-session-event":
            sessionID = data.session;
            break;
        

    }
}