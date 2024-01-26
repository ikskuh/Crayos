let canvas;
let ctx;

// CONSTANTS
const width = 1920;
const height = 1080;

const palette = [
    "#FFF",
    "#000",
];

function init() {
    initPainter();
}

function initPainter() {
  canvas = document.getElementById("canvas");
  ctx = canvas.getContext("2d");

  canvas.onmousemove = (e) => {
    ctx.beginPath();
    ctx.rect(e.offsetX, e.offsetY, 20, 20);
    ctx.fillStyle = 'black;'
    ctx.fill();
  };

  const resize = (event) => {
    const wrapper = document.getElementById("wrapper");
    const scale = Math.min(window.innerWidth / width, window.innerHeight / height);
    wrapper.style.transform = "translate(-50%, -50%) scale(" + scale + ")";
  };
  resize();
  window.addEventListener("resize", resize);
}

function drawPalette() {
    palCanvas = document.getElementById("canvas");
    palCtx = canvas.getContext("2d");

    const w = 122;
    const h = 80;
}
