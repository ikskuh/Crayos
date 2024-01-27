let painterTimerNumberElem;

// CONSTANTS
const width = 1920;
const height = 1080;
const lineWidth = 20;
const distanceThreshold = 5; // minimum distance between points to add a new point
const eraserRadius = 40;
const TOOL_PENCIL = "pencil";
const TOOL_ERASER = "eraser";
const EFFECT_COOLDOWN_MS = 10000;
const TIMER_SECONDS = 90;

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

const painterPaths = [];
let mx = -1000;
let my = -1000;

let chaosEffect = null;

function setPaintingToolsEnabled(enabled) {
  document.getElementById('painting-tools').style.display = enabled ? "block" : "none";
  setInputEnabled(enabled);
  if (enabled) {
    initPalette();
    selectTool(TOOL_PENCIL);
  }
}

function setChaosEffectsEnabled(enabled) {
  document.getElementById('chaos-effects').style.display = enabled ? "block" : "none";
}

function setVotingButtonsEnabled(enabled) {
  document.getElementById('voting-buttons').style.display = enabled ? "block" : "none";
}

function onMouseDown(e) {
  mx = e.offsetX;
  my = e.offsetY;
  const point = { x: mx, y: my };
  if (selectedTool == TOOL_PENCIL) {
    pencilBeginPath(point);
  } else if (selectedTool == TOOL_ERASER) {
    eraserDeleteAt(point);
  }
  drawPainterCanvas();
}

function onMouseMove(e) {
  mx = e.offsetX;
  my = e.offsetY;
  if (e.buttons & 1 || chaosEffect == Effect.lock_pencil) {
    const point = { x: mx, y: my };
    if (selectedTool == TOOL_PENCIL) {
      pencilContinuePath(point);
    } else if (selectedTool == TOOL_ERASER) {
      eraserDeleteAt(point);
    }
  }
  drawPainterCanvas();
}

function onMouseUp(e) {
  mx = -1000;
  my = -1000;
  drawPainterCanvas();
  console.log(JSON.stringify(painterPaths));
}

function onMouseEnter(e) {
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
}

function onMouseLeave(e) {
  mx = -1000;
  my = -1000;
  drawPainterCanvas();
}

function setInputEnabled(enabled) {
  const painterCanvas = document.getElementById("painter-canvas");
  if (enabled) {
    painterCanvas.classList.add("input-enabled");
    painterCanvas.addEventListener("mousedown", onMouseDown);
    painterCanvas.addEventListener("mousemove", onMouseMove);
    painterCanvas.addEventListener("mouseup", onMouseUp);
    painterCanvas.addEventListener("mouseenter", onMouseEnter);
    painterCanvas.addEventListener("mouseleave", onMouseLeave);
  } else {
    painterCanvas.classList.remove("input-enabled");
    painterCanvas.removeEventListener("mousedown", onMouseDown);
    painterCanvas.removeEventListener("mousemove", onMouseMove);
    painterCanvas.removeEventListener("mouseup", onMouseUp);
    painterCanvas.removeEventListener("mouseenter", onMouseEnter);
    painterCanvas.removeEventListener("mouseleave", onMouseLeave);
  }
}

function resetPainting() {
  painterPaths.splice(0, painterPaths.length);
}

function initArtstudio() {
  painterTimerNumberElem = document.getElementById("painter-timer-number");

  // countdown
  painterTimerNumberElem.innerText = TIMER_SECONDS;
  const timerInterval = setInterval(() => {
    painterTimerNumberElem.innerText--;
    if (painterTimerNumberElem.innerText == 0) {
      clearInterval(timerInterval);
    }
  }, 1000);

  drawPainterCanvas();
}

function drawPainterCanvas() {
  const painterCanvas = document.getElementById("painter-canvas");
  const ctx = painterCanvas.getContext("2d");
  drawPainting(painterCanvas, painterPaths);

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

  // chaos effect
  if (chaosEffect == Effect.flashlight) {
    ctx.beginPath();
    ctx.rect(0, 0, painterCanvas.width, painterCanvas.height);
    ctx.arc(mx, my, 200, 2 * Math.PI, 0);
    ctx.fillStyle = "#000";
    ctx.fill("evenodd");
  }
}

// CHAOS EFFECTS

function activateChaosEffect(effect) {
  const painterCanvas = document.getElementById("painter-canvas");
  console.log("activate chaos effect: " + effect);
  chaosEffect = effect;
  setTimeout(() => {
    deactivateChaosEffect(effect);
  }, EFFECT_COOLDOWN_MS);

  if (chaosEffect == Effect.flip) {
    painterCanvas.classList.add(Effect.flip);
  } else if (chaosEffect == Effect.drunk) {
    painterCanvas.classList.add(Effect.drunk);
  } else if (chaosEffect == Effect.swap_tool) {
    swapTools();
  }

  drawPainterCanvas();
}

function deactivateChaosEffect(effect) {
  console.log("deactivate chaos effect: " + effect);
  if (effect == Effect.flip) {
    painterCanvas.classList.remove(Effect.flip);
  } else if (effect == Effect.drunk) {
    painterCanvas.classList.remove(Effect.drunk);
  } else if (chaosEffect == Effect.swap_tool) {
    swapTools();
  }

  chaosEffect = null;
  drawPainterCanvas();
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

function swapTools() {
  if (selectedTool == TOOL_PENCIL) {
    selectTool(TOOL_ERASER);
  } else if (selectedTool == TOOL_ERASER) {
    selectTool(TOOL_PENCIL);
  }
}

function pencilBeginPath(point) {
  painterPaths.push({ color: palette[selectedColor], points: [point] });
}

function pencilContinuePath(point) {
  if (painterPaths.length == 0) {
    pencilBeginPath(point);
    return;
  }
  const currentPath = painterPaths[painterPaths.length - 1];
  const lastPoint = currentPath.points[currentPath.points.length - 1];
  if (distanceSquared(lastPoint, point) > distanceThreshold * distanceThreshold) {
    currentPath.points.push(point);
  }
}

function eraserDeleteAt(point) {
  for (let i = 0; i < painterPaths.length; i++) {
    const path = painterPaths[i];
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
  const paletteElem = document.getElementById("palette");

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
  const paletteElem = document.getElementById("palette");
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

function distanceSquared(p1, p2) {
  return Math.pow(p1.x - p2.x, 2) + Math.pow(p1.y - p2.y, 2);
}
