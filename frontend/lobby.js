let inviteLink;

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

    // Show invite link
    let hostUrl = location.protocol + '//' + location.host + location.pathname;
    inviteLink = hostUrl + "?session=" + sessionID;
}

function readyClicked() {
    localIsReady = !localIsReady;
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