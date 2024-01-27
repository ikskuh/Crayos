function initGallery() {
    for (let i = 1; i <= 4; i++) {
        let galCanvas = document.getElementById("gallery" + i);
        let ctx = galCanvas.getContext("2d")
        ctx.rect(0, 0, 1284, 736);
        ctx.fill();
        galCanvas.style.width = "692px";
    }
}