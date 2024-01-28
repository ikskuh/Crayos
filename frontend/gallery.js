function initGallery()
{
}

function setGalleryCanvases(results)
{
    for (let i = 0; i < results.length; i++)
    {
        let galCanvas = document.getElementById("gallery" + (i+1));
        drawPainting(galCanvas, results[i].graphics.paths, results[i].backdrop);
        drawFinalPoints(galCanvas, results[i].score);

        if (results[i].winner)
            drawWinnerBadge(galCanvas);
    }
}

function drawFinalPoints(canvas, points)
{
    // Draw Star
    let star = new Image;
    let pointsDisplayed = points.toFixed(1);
    const ctx = canvas.getContext("2d");
    star.onload = function() {
        ctx.drawImage(star, 37, 500, 200, 200);
    };
    star.src = "img/star.png"

    // Draw Number
    ctx.font = "50px serif";
    ctx.textAlign = "center";
    ctx.fillText(pointsDisplayed, 137, 500);
}

function drawWinnerBadge(canvas)
{
    let badge = new Image;
    badge.onload = function() {
        canvas.getContext("2d").drawImage(badge, 1100, 500, 200, 200);
    };
    badge.src = "img/winner_badge.png";
}

function backToLobby() {
    sendUserCommand(UserAction.leaveGallery);
}