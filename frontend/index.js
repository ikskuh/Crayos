const NoSession = -1;
const Open = 1;

let socket;
let sessionID = NoSession;

let serverSideDisconnect = false;   

let players = ["", "", "", ""];
let localPlayer = "nickname";
let localIsReady = false;

let currentGamestate = "connecting";

const backgrounds = [];

function init() {
    loadBackgrounds();

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
        case EventId.ChangeGameView:
            //setView(data.view);
            setView("rating");
            break;
        case EventId.PlayersChanged:
            for (let i = 0; i < data.players.length; i++) {
                players[i] = data.players[i];
            }
            updateLobby();
            break;
        case EventId.EnterSession:
            sessionID = data.sessionId;
            break;
        case EventId.JoinSessionFailed:
            setView("link_invalid");
            break;
        case EventId.Kicked:
            break;
        case EventId.ChangeToolModifier:
            break;
        case EventId.PaintingChanged:
            break;
        case EventId.PlayerReadyChanged:
            updateLobby(data.players);
            break;
    }
}

function loadBackgrounds() {
    backgrounds.push(document.getElementById("background0"));
    backgrounds.push(document.getElementById("background1"));
    backgrounds.push(document.getElementById("background2"));
    backgrounds.push(document.getElementById("background3"));
}

function drawPainting(canvas, paths) {
    const ctx = canvas.getContext("2d");
    ctx.drawImage(backgrounds[selectedBackground], 0, 0, canvas.width, canvas.height);
  
    ctx.lineWidth = lineWidth;
    ctx.lineCap = "round";
    ctx.lineJoin = "round";
    for (let i = 0; i < paths.length; i++) {
      const path = paths[i];
      ctx.beginPath();
      if (path.points.length == 1 && !path.points[0].erased) {
        ctx.arc(path.points[0].x, path.points[0].y, lineWidth / 2, 0, 2 * Math.PI);
        ctx.fillStyle = path.color;
        ctx.fill();
      } else {
        let moved = false;
        for (let j = 0; j < path.points.length; j++) {
          if (path.points[j].erased) {
            moved = false;
            continue;
          }
          if (!moved) {
            ctx.moveTo(path.points[j].x, path.points[j].y);
            moved = true;
          } else {
            ctx.lineTo(path.points[j].x, path.points[j].y);
          }
        }
        ctx.strokeStyle = path.color;
        ctx.stroke();
      }
    }
}