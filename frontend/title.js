function createGame() {
    localPlayer = document.getElementById("nicknameInput").value;
    sendCreateSessionCommand(localPlayer);
}

function joinGame() {
    localPlayer = document.getElementById("nicknameInput").value;
    let urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has("session")) {
        sessionID = urlParams.get("session");
        sendJoinSessionCommand(localPlayer, sessionID);
    }
    else {
        setView("link_required");
    }
}

function btnBackToLobby() {
    setView("title");
}