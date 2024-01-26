let canvas;
let paletteElem;

// CONSTANTS
const width = 1920;
const height = 1080;

const palette = [
  "#FFF",
  "#e42932", // red
  "#ff8652", // orange
  "#552cb7", // purple
  "#00995e", // green**
  "#058cd7", // blue
  "#fff243", // yellow
  "#000",
];
let selectedColor = 7;

function initPainter() {
  canvas = document.getElementById("canvas");
  const ctx = canvas.getContext("2d");

  canvas.onmousemove = (e) => {
    console.log(e.button);
    if (e.buttons & 1) {
      ctx.beginPath();
      ctx.rect(e.offsetX, e.offsetY, 20, 20);
      ctx.fillStyle = palette[selectedColor];
      ctx.fill();
    }
  };

  const resize = (event) => {
    const wrapper = document.getElementById("wrapper");
    const scale = Math.min(window.innerWidth / width, window.innerHeight / height);
    wrapper.style.transform = "translate(-50%, -50%) scale(" + scale + ")";
  };
  resize();
  window.addEventListener("resize", resize);

  initPalette();
}

const colorw = 122;
const colorh = 80;

const makeColorRect = (i) => {
  return {
    x: 2 + (i % 2) * (colorw + 16),
    y: 2 + Math.floor(i / 2) * (colorh + 12),
    w: colorw,
    h: colorh,
  };
};

function initPalette() {
  paletteElem = document.getElementById("palette");

  paletteElem.onclick = (e) => {
    for (let i = 0; i < palette.length; i++) {
      const rect = makeColorRect(i);
      if (e.offsetX >= rect.x && e.offsetY >= rect.y && e.offsetX <= rect.x + rect.w && e.offsetY <= rect.y + rect.h) {
        selectedColor = i;
        drawPalette();
        break;
      }
    }
  };

  drawPalette();
}

function drawPalette() {
  pctx = paletteElem.getContext("2d");
  pctx.clearRect(0, 0, paletteElem.width, paletteElem.height);
  for (let i = 0; i < palette.length; i++) {
    pctx.beginPath();
    const rect = makeColorRect(i);
    pctx.rect(rect.x, rect.y, rect.w, rect.h);
    pctx.fillStyle = palette[i];
    pctx.fill();
    if (i == selectedColor) {
      pctx.strokeStyle = "#000";
      pctx.lineWidth = 4;
      pctx.stroke();
    }
  }
}
