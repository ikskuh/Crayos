let canvas;
let paletteElem;

// CONSTANTS
const width = 1920;
const height = 1080;
const lineWidth = 20;
const distanceThreshold = 10; // minimum distance between points to add a new point

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

const paths = [];

const backgroundUrls = [
  "img/graveyard.png",
];
const backgrounds = [];

function distanceSquared(p1, p2) {
  return Math.pow(p1.x - p2.x, 2) + Math.pow(p1.y - p2.y, 2);
}

function loadBackgrounds() {
  for (let i = 0; i < backgroundUrls.length; i++) {
    const img = new Image();
    img.src = backgroundUrls[i];
    img.onload = () => {
      drawCanvas();
    }
    backgrounds.push(img);
  }
}

function initPainter() {
  canvas = document.getElementById("canvas");

  loadBackgrounds();

  canvas.onmousedown = (e) => {
    paths.push({ color: palette[selectedColor], points: [{ x: e.offsetX, y: e.offsetY }] });
  };

  canvas.onmousemove = (e) => {
    console.log(e.button);
    if (e.buttons & 1) {
      const currentPath = paths[paths.length - 1];
      const lastPoint = currentPath.points[currentPath.points.length - 1];
      const newPoint = { x: e.offsetX, y: e.offsetY };
      if (distanceSquared(lastPoint, newPoint) > distanceThreshold * distanceThreshold) {
        currentPath.points.push(newPoint);
      }
      drawCanvas();
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

function drawCanvas() {
  const ctx = canvas.getContext("2d");
  // ctx.clearRect(0, 0, canvas.width, canvas.height);
  ctx.drawImage(backgrounds[0], 0, 0, canvas.width, canvas.height);
  ctx.lineWidth = lineWidth;
  ctx.lineCap = "round";
  ctx.lineJoin = "round";
  for (let i = 0; i < paths.length; i++) {
    const path = paths[i];
    ctx.beginPath();
    ctx.moveTo(path.points[0].x, path.points[0].y);
    for (let j = 1; j < path.points.length; j++) {
      ctx.lineTo(path.points[j].x, path.points[j].y);
    }
    ctx.strokeStyle = path.color;
    ctx.stroke();
  }
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
