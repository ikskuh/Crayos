function createGame() {
    localPlayer = document.getElementById("nicknameInput").value;
    socket.send(JSON.stringify({
        type: "create-session-command",
        nickname: localPlayer
    }));
}

function joinGame() {

}