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
                playerInfo.style.backgroundColor = "transparent";
            }
        }
    }

    // Local ready button
    if (localIsReady == true) {
        document.getElementById("ready").style.backgroundColor = "green";
    }
    else {
        document.getElementById("ready").style.backgroundColor = "transparent";
    }

    // Show invite link
    let hostUrl = location.protocol + '//' + location.host + location.pathname;
    inviteLink = hostUrl + "?session=" + sessionID;
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
    document.getElementById("joinLink").textContent ="Link Copied!";
    setTimeout(resetLinkBtn, 2000);
}

function resetLinkBtn() {
    document.getElementById("joinLink").textContent = "Copy Invite Link";
}