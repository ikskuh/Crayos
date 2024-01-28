function initTitle() {
    extractSessionId();
    document.getElementById("sessionIdInput").value = sessionID;
}

function createGame() {
    localPlayer = document.getElementById("nicknameInput").value;
    if (localPlayer == "")
        localPlayer = "nickname";
        
    sendCreateSessionCommand(localPlayer);
}

function joinGame() {
    localPlayer = document.getElementById("nicknameInput").value;
    if (localPlayer == "")
        localPlayer = "nickname";

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