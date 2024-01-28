function initTitle() {
    extractSessionId();
    document.getElementById("sessionIdInput").value = sessionID;
}

function createGame() {
    getNickname()
    
    sendCreateSessionCommand(localPlayer);
}

function joinGame() {
    getNickname();
    sessionID = document.getElementById("sessionIdInput").value;

    if (sessionID != "") {
        sendJoinSessionCommand(localPlayer, sessionID);
    }
    else {
        setView("id_required");
    }
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
    localPlayer = document.getElementById("nicknameInput").value;
    if (localPlayer == "")
        localPlayer = "nickname";
}