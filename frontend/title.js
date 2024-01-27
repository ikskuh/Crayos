function createGame() {
    localPlayer = document.getElementById("nicknameInput").value;
    socket.send(JSON.stringify({
        type: "create-session-command",
        nickname: localPlayer
    }));
}

function joinGame() {
    let urlParams = new URLSearchParams(queryString);
    if (urlParams.has("session")) {
        sessionID = urlParams.get("session");
        
    }
    else {
        setView("link_required");
    }
}

function btnBackToLobby() {
    setView("title");
}