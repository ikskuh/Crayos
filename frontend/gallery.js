function initGallery()
{
}

function setGalleryCanvases(results)
{
    for (let i = 0; i < results.length; i++)
    {
        let galCanvas = document.getElementById("gallery" + (i+1));
        drawPainting(galCanvas, results[i].graphics.paths, results[i].backdrop);
        //drawFinalStars(galCanvas, 1);

        if (results[i].winner)
            drawWinnerBadge(galCanvas);
    }
}

function drawFinalStars(canvas, points)
{
    let star = new Image;
    let starCnt = round(points);
    star.onload = function() {
        const ctx = canvas.getContext("2d");
        for (let i = 0; i < points; i++)
        {
            ctx.drawImage(star, 50 * i, 500);
        }
    };
    star.src = "img/star.png"
}

function drawWinnerBadge(canvas)
{
    let badge = new Image;
    badge.onload = function() {
        canvas.getContext("2d").drawImage(badge, 1100, 500);
    };
    badge.src = "img/winner_badge.png";
}

function backToLobby() {
    sendUserCommand(UserAction.leaveGallery);
}