let canvas;
let paletteElem;
let timerNumberElem;

// CONSTANTS
const width = 1920;
const height = 1080;
const lineWidth = 20;
const distanceThreshold = 5; // minimum distance between points to add a new point
const eraserRadius = 40;
const TOOL_PENCIL = "pencil";
const TOOL_ERASER = "eraser";
const EFFECT_COOLDOWN_MS = 10000;

let selectedTool = TOOL_PENCIL;

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
let mx = -100;
let my = -100;

const backgroundUrls = [
  "img/arctic.png",
  "img/graveyard.png",
  "img/pirate_ship.png",
  "img/theater_stage1.png"
];
const backgrounds = [];
let selectedBackground = 0;

let chaosEffect = "flashlight";

function initPainter() {
  document.getElementById("painter").style.display = "block";
  canvas = document.getElementById("canvas");
  timerNumberElem = document.getElementById("timer-number");

  loadBackgrounds();

  canvas.onmousedown = (e) => {
    mx = e.offsetX;
    my = e.offsetY;
    const point = { x: mx, y: my };
    if (selectedTool == TOOL_PENCIL) {
      pencilBeginPath(point);
    } else if (selectedTool == TOOL_ERASER) {
      eraserDeleteAt(point);
    }
    drawCanvas();
  };
  canvas.onmousemove = (e) => {
    mx = e.offsetX;
    my = e.offsetY;
    if (e.buttons & 1) {
      const point = { x: mx, y: my };
      if (selectedTool == TOOL_PENCIL) {
        pencilContinuePath(point);
      } else if (selectedTool == TOOL_ERASER) {
        eraserDeleteAt(point);
      }
    }
    drawCanvas();
  };
  canvas.onmouseup = (e) => {
    // TODO: send paths to server
    console.log(JSON.stringify(paths));
  };
  canvas.onmouseenter = (e) => {
    if (e.buttons & 1) {
      mx = e.offsetX;
      my = e.offsetY;
      const point = { x: mx, y: my };
      if (selectedTool == TOOL_PENCIL) {
        pencilBeginPath(point);
      } else if (selectedTool == TOOL_ERASER) {
        eraserDeleteAt(point);
      }
    }
  };
  canvas.onmouseleave = (e) => {
    mx = -100;
    my = -100;
    drawCanvas();
  };

  initPalette();
  selectTool(TOOL_PENCIL);
  paths.splice(0, paths.length);

  selectedBackground = Math.floor(Math.random() * backgrounds.length);

  timerNumberElem.innerText = 60;
  const timerInterval = setInterval(() => {
    timerNumberElem.innerText--;
    if (timerNumberElem.innerText == 0) {
      clearInterval(timerInterval);
    }
  }, 1000);

}

function drawCanvas() {
  const ctx = canvas.getContext("2d");
  ctx.drawImage(backgrounds[selectedBackground], 0, 0, canvas.width, canvas.height);

  ctx.lineWidth = lineWidth;
  ctx.lineCap = "round";
  ctx.lineJoin = "round";
  for (let i = 0; i < paths.length; i++) {
    const path = paths[i];
    ctx.beginPath();
    if (path.points.length == 1 && !path.points[0].erased) {
      ctx.arc(path.points[0].x, path.points[0].y, lineWidth / 2, 0, 2 * Math.PI);
      ctx.fillStyle = path.color;
      ctx.fill();
    } else {
      let moved = false;
      for (let j = 0; j < path.points.length; j++) {
        if (path.points[j].erased) {
          moved = false;
          continue;
        }
        if (!moved) {
          ctx.moveTo(path.points[j].x, path.points[j].y);
          moved = true;
        } else {
          ctx.lineTo(path.points[j].x, path.points[j].y);
        }
      }
      ctx.strokeStyle = path.color;
      ctx.stroke();
    }
  }

  // preview tool
  if (selectedTool == TOOL_PENCIL) {
    ctx.fillStyle = palette[selectedColor];
    ctx.beginPath();
    ctx.arc(mx, my, lineWidth / 2, 0, 2 * Math.PI);
    ctx.fill();
  } else if (selectedTool == TOOL_ERASER) {
    ctx.strokeStyle = "#000";
    ctx.lineWidth = 2;
    ctx.beginPath();
    ctx.arc(mx, my, eraserRadius, 0, 2 * Math.PI);
    ctx.stroke();
  }
}

// CHAOS EFFECTS

function onChaosEffect(effect) {
  chaosEffect = effect;
  setTimeout(() => {
    deactivateChaosEffect();
  }, EFFECT_COOLDOWN_MS);

  if (chaosEffect == Effect.flip) {
    canvas.classList.add(Effect.flip);
  } else if (chaosEffect == Effect.drunk) {
    canvas.classList.add(Effect.drunk);
  }
}

function deactivateChaosEffect() {
  if (chaosEffect == Effect.flip) {
    canvas.classList.remove(Effect.flip);
  } else if (chaosEffect == Effect.drunk) {
    canvas.classList.remove(Effect.drunk);
  }
  chaosEffect = null;
}

// TOOLS

function selectTool(tool) {
  selectedTool = tool;
  if (tool == TOOL_PENCIL) {
    document.getElementById(TOOL_PENCIL).classList.add("selected");
    document.getElementById("eraser").classList.remove("selected");
  } else if (tool == "eraser") {
    document.getElementById("eraser").classList.add("selected");
    document.getElementById(TOOL_PENCIL).classList.remove("selected");
  }
}

function pencilBeginPath(point) {
  paths.push({ color: palette[selectedColor], points: [point] });
}

function pencilContinuePath(point) {
  const currentPath = paths[paths.length - 1];
  const lastPoint = currentPath.points[currentPath.points.length - 1];
  if (distanceSquared(lastPoint, point) > distanceThreshold * distanceThreshold) {
    currentPath.points.push(point);
  }
}

function eraserDeleteAt(point) {
  for (let i = 0; i < paths.length; i++) {
    const path = paths[i];
    for (let j = 0; j < path.points.length; j++) {
      const p = path.points[j];
      if (distanceSquared(p, point) < eraserRadius * eraserRadius) {
        p.erased = true;
      }
    }
  }
}

// PALETTE

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

function loadBackgrounds() {
  for (let i = 0; i < backgroundUrls.length; i++) {
    const img = new Image();
    img.src = backgroundUrls[i];
    img.onload = () => {
      drawCanvas();
    };
    backgrounds.push(img);
  }
}

function distanceSquared(p1, p2) {
  return Math.pow(p1.x - p2.x, 2) + Math.pow(p1.y - p2.y, 2);
}
