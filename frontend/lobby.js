let inviteLink;

function updateLobby(readyMap = undefined) {
    // Update Nicknames and ready status
    for (let i = 0; i < players.length; i++) {
        let playerInfo = document.getElementById("player" + (i+1));
        playerInfo.value = players[i];
        if (readyMap != undefined) {
            if (readyMap[players[i]] == true) {
                playerInfo.style.backgroundColor = "green";
            }
            else {
                if (players[i] == "")
                    playerInfo.style.backgroundColor = "transparent";
                else
                    playerInfo.style.backgroundColor = "#FD5A46";
            }
        }
    }

    // Local ready button
    if (localIsReady == true) {
        document.getElementById("ready").style.backgroundColor = "green";
    }
    else {
        document.getElementById("ready").style.backgroundColor = "#FD5A46";
    }

    // Show invite link
    let hostUrl = location.protocol + '//' + location.host + location.pathname;
    inviteLink = hostUrl + "?session=" + sessionID;

    // Show ID
    document.getElementById("copyId").value = "ID: " + sessionID;
}

function readyClicked() {
    if (localIsReady) {
        sendUserCommand(UserAction.setNotReady)
        localIsReady = false;
    }
    else {
        sendUserCommand(UserAction.setReady)
        localIsReady = true;
    }
    updateLobby();
}

function btnCopyInvite() {
    navigator.clipboard.writeText(inviteLink);
    tempChangeBtnText("joinLink", "Link Copied!", 2000);
}

function btnCopyId() {
    navigator.clipboard.writeText(sessionID);
    tempChangeBtnText("copyId", "ID Copied!", 2000);
}