function initTitle() {
    extractSessionId();
    document.getElementById("sessionIdInput").value = sessionID;
    document.getElementById("nicknameInput").placeholder = nick_names[Math.floor(Math.random() * nick_names.length)];
}

function createGame() {
    getNickname()
    
    sendCreateSessionCommand(localPlayer);
}

function joinGame() {
    getNickname();
    sessionID = document.getElementById("sessionIdInput").value;
    sendJoinSessionCommand(localPlayer, sessionID);
}

function btnBackToLobby() {
    setView("title");
}

// Extracs session id from current url
function extractSessionId() {
    let urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has("session")) {
        sessionID = urlParams.get("session");
    }
    else {
        sessionID = "";
    }
}

// Checks if the nick is valid, otherwise defaults
function getNickname() {
    const nicknameInput = document.getElementById("nicknameInput");
    localPlayer = nicknameInput.value;
    if (localPlayer == "")
        localPlayer = nicknameInput.placeholder;
}