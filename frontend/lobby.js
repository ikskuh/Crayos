function updateLobby() {
    // Update Nicknames and ready status
    for (let i = 0; i < players.length; i++) {
        document.getElementById("player" + (i+1)).value = players[i];
    }

    // Local ready button
    if (localIsReady == true) {
        document.getElementById("ready").style.backgroundColor = "green";
    }
    else {
        document.getElementById("ready").style.backgroundColor = "transparent";
    }
}

function readyClicked() {
    localIsReady = !localIsReady;
    updateLobby();
}