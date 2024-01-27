function createGame() {
    nicknames[0] = document.getElementById("nicknameInput").value;
    socket.send(JSON.stringify({
        type: "create-session-command",
        nickname: nicknames[0]
    }));
}

function joinGame() {

}