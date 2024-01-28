const NoSession = -1;
const Open = 1;

let socket;
let sessionID = NoSession;

let serverSideDisconnect = false;

let players = ["", "", "", ""];
let localPlayer = "nickname";
let localIsReady = false;

let currentView = "connecting";

const backgrounds = [];
let selectedBackground = null;

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
  let socketUrl = "ws://192.168.37.247:8090/ws";
  {
    let urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has("local")) {
      socketUrl = "ws://localhost:8080/ws";
    } else if (urlParams.has("host")) {
      socketUrl = "ws://" + urlParams.get("host") + ":8080/ws";
    }
  }

  document.getElementById("connecting").style.display = "flow";
  socket = new WebSocket(socketUrl);

  socket.onerror = function (event) {
    console.log("WebSocket error: ", event);
  };
  socket.onopen = function (event) {
    setView("title");
  };
  socket.onmessage = function (event) {
    onSocketReceive(event);
  };

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

function setView(newView) {
  if (
    newView == GameView.title ||
    newView == GameView.lobby ||
    newView == GameView.gallery ||
    newView == "connecting" ||
    newView == "connection_failed" ||
    newView == "link_required" ||
    newView == "link_invalid"
  ) {
    hideSection(currentView);
    currentView = newView;
    showSection(newView);

    if (newView == GameView.gallery) {
      initGallery();
    }
  } else {
    switch (newView) {
      case GameView.promptselection:
        break;
      case GameView.artstudioGeneric:
        break;
      case GameView.artstudioActive:
        break;
      case GameView.artstudioSticker:
        break;
    }
    newView = "artstudio";
    hideSection(currentView);
    currentView = newView;
    showSection(newView);
  }
}

function onSocketReceive(event) {
  let data = JSON.parse(event.data);
  console.log(data);

  switch (data.type) {
    case EventId.EnterSession:
      sessionID = data.sessionId;
      break;
    case EventId.JoinSessionFailed:
      setView("link_invalid");
      break;
    case EventId.Kicked:
      break;
    case EventId.ChangeGameView:
      if (data.painting.backdrop) {
        setBackground(data.painting.backdrop);
      }
      if (data.painting.graphics) {
        setPainting(data.painting.graphics);
      } else {
        clearPainting();
      }
      setPaintingPrompt(data.painting.prompt);

      if (data.view == GameView.promptselection) {
        setPromptOptions(data.voteOptions);
        setPromptSelectionEnabled(true);
        setVoteOptions([]);
      } else {
        setPromptSelectionEnabled(false);
        setVoteOptions(data.voteOptions);
      }

      if (data.view == GameView.artstudioActive) {
        setPaintingToolsEnabled(true);
      } else {
        setPaintingToolsEnabled(false);
      }

      setView(data.view);
      break;
    case EventId.TimerChanged:
      setTimerSecondsLeft(data.secondsLeft);
      break;
    case EventId.ChangeToolModifier:
      activateChaosEffect(data.modifier, data.duration);
      break;
    case EventId.PaintingChanged:
      setPainting(data.graphics);
      break;
    case EventId.PlayersChanged:
      for (let i = 0; i < data.players.length; i++) {
        players[i] = data.players[i];
      }
      updateLobby();
      break;
    case EventId.PlayerReadyChanged:
      updateLobby(data.players);
      break;
  }
}

function loadBackgrounds() {
  backgrounds["arctic"] = document.getElementById("background-arctic");
  backgrounds["graveyard"] = document.getElementById("background-graveyard");
  backgrounds["pirate_ship"] = document.getElementById("background-pirate_ship");
  backgrounds["theater_stage1"] = document.getElementById("background-theater_stage1");
  backgrounds["desert"] = document.getElementById("background-desert");
}

function setBackground(name) {
  selectedBackground = backgrounds[name];
  drawPainterCanvas();
}

function drawPainting(canvas, paths) {
  const ctx = canvas.getContext("2d");
  if (selectedBackground) {
    ctx.drawImage(selectedBackground, 0, 0, canvas.width, canvas.height);
  } else {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
  }

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
