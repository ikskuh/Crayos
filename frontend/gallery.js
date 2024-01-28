function initGallery()
{
    for (let i = 1; i <= 4; i++)
    {
        let galCanvas = document.getElementById("gallery" + i);
        let ctx = galCanvas.getContext("2d")
        ctx.rect(0, 0, 1384, 736);
        ctx.fill();
    }
}

function setGalleryCanvases(results)
{
    let winnerCanvas;
    let highestScore = 0;
    for (let i = 0; i < results.length; i++)
    {
        let galCanvas = document.getElementById("gallery" + (i+1));
        drawPainting(galCanvas, results[i].graphics, results[i].backdrop);
        //drawFinalStars(galCanvas, 1);

        if (results[i].winner)
        {
            //highestScore = results[i].points[i];
            winnerCanvas = galCanvas;
        }
    }

    //drawWinnerBadge(winnerCanvas);
}

function drawFinalStars(canvas, points)
{
    let star = new Image;
    let starCnt = round(points);
    star.onload = function() {
        let ctx = canvas.GetContext("2d");
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
        canvas.GetContext("2d").drawImage(badge, 300, 300);
    };
    badge.src = "img/winner_badge.png";
}

function backToLobby() {
    sendUserCommand(UserAction.leaveGallery);
}