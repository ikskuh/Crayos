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
let selectedBackgroundName = "";

function init() {
    window.onerror = (event) => {
        console.log(event);
        alert(event);
    };

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
      socketUrl = "ws://" + window.location.hostname + ":8080/ws";
      console.log("socket url: " + socketUrl);
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
  switch (id) {
    case GameView.promptselection:
    case GameView.artstudioGeneric:
    case GameView.artstudioActive:
    case GameView.artstudioSticker:
      id = "artstudio";
      break;
  }

  document.getElementById(id).style.display = "none";
}

function showSection(id) {
  switch (id) {
    case GameView.promptselection:
    case GameView.artstudioGeneric:
    case GameView.artstudioActive:
    case GameView.artstudioSticker:
      id = "artstudio";
      break;
  }
  document.getElementById(id).style.display = "flow";
}

let global_popup_timeout 

function showPopUp(message, duration) {
  let popup = document.getElementById("popup");
  popup.classList.add("visible");
  popup.innerText = message;

  if(global_popup_timeout ) {
    clearTimeout(global_popup_timeout )
  }
  global_popup_timeout  = setTimeout(() => {
     popup.classList.remove("visible");
  }, duration || 1500);
}


function setView(newView) {
  hideSection(currentView);
  currentView = newView;
  showSection(newView);

  if (newView == GameView.gallery) {
      initGallery();
  }
}

function onSocketReceive(event) {
  let data = JSON.parse(event.data);
  
  // hide periodic timer events:
  if (data.type != EventId.TimerChanged) {
    console.log(data);
  }

  switch (data.type) {
    case EventId.EnterSession:
      sessionID = data.sessionId;
      break;
    case EventId.JoinSessionFailed:
        setView("server_error");
      document.getElementById("serverErrorText").textContent = data.reason;
      break;
    case EventId.Kicked:
      alert(data.reason);
      break;
    case EventId.ChangeGameView:
      if (data.painting.backdrop) {
        setBackground(data.painting.backdrop);
      }
      if (data.painting.graphics) {
        setPainting(data.painting);
      } else {
        clearPainting();
      }
      setPaintingPrompt(data.painting.prompt);

      if (data.view == GameView.title) {
        initTitle();
      }

      setView(data.view); // HACK: need currenView before setVoteOptions
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

      if (data.view == GameView.gallery) {
        setGalleryCanvases(data.results)
      }

      if (data.announcer != "") {
        document.getElementById("announcer_text").textContent = data.announcer;
      }
      break;
    case EventId.TimerChanged:
      setTimerSecondsLeft(data.secondsLeft);
      break;
    case EventId.ChangeToolModifier:
      if (data.modifier == "") {
        // special handling: deactivate current effect
        deactivateChaosEffect(chaosEffect);
      } else {
        activateChaosEffect(data.modifier, data.duration);
      }
      break;
    case EventId.PaintingChanged:
      updatePainting(data.graphics);
      break;
    case EventId.PlayersChanged:
      for (let i = 0; i < 4; i++) {
        if (i < data.players.length)
            players[i] = data.players[i];
        else
            players[i] = "";
      }
      updateLobby();
      break;
    case EventId.PlayerReadyChanged:
      updateLobby(data.players);
      break;

    case EventId.DebugMessage:
      let overlay = document.getElementById('debug-overlay');
      overlay.classList.add("visible");
      overlay.innerText = data.message || "&nbsp;";
      break;

    case EventId.PopUp:
      showPopUp(data.message, data.duration);
      break;

    default:
      throw "unhandled message: " + JSON.stringify(data);
  }
}


function loadBackgrounds() {
  backgrounds["arctic"] = document.getElementById("background-arctic");
  backgrounds["graveyard"] = document.getElementById("background-graveyard");
  backgrounds["pirate_ship"] = document.getElementById("background-pirate_ship");
  backgrounds["theater_stage1"] = document.getElementById("background-theater_stage1");
  backgrounds["desert"] = document.getElementById("background-desert");
}

function setBackground(backgroundName) {
  selectedBackgroundName = backgroundName;
  drawPainterCanvas();
}

function drawPainting(canvas, paths, backgroundName) {
  const ctx = canvas.getContext("2d");
  if (backgroundName != "") {
    ctx.drawImage(backgrounds[backgroundName], 0, 0, canvas.width, canvas.height);
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

// Temporarily changes the text of a button for time t
function tempChangeBtnText(btnId, tempText, t) {
    let btn = document.getElementById(btnId);
    let oldText = btn.textContent;
    btn.textContent = tempText;
    setTimeout(resetBtnText, t, btn, oldText);
}
// Belongs to tempChangeBtnText
function resetBtnText(btn, resetText) {
    btn.textContent = resetText;
}