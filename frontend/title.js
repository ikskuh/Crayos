function createGame() {
    localPlayer = document.getElementById("nicknameInput").value;
    socket.send(JSON.stringify({
        type: "create-session-command",
        nickname: localPlayer
    }));

    
}

function joinGame() {
    localPlayer = document.getElementById("nicknameInput").value;
    let urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has("session")) {
        sessionID = urlParams.get("session");
        socket.send(JSON.stringify({
            type: "join-session-command",
            nickname: localPlayer,
            sessionid: sessionID
        }));
    }
    else {
        setView("link_required");
    }
}

function btnBackToLobby() {
    setView("title");
}