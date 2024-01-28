// CONSTANTS
const width = 1920;
const height = 1080;
const lineWidth = 20;
const distanceThreshold = 5; // minimum distance between points to add a new point
const eraserRadius = 40;
const TOOL_PENCIL = "pencil";
const TOOL_ERASER = "eraser";
const TIMER_SECONDS = 90;

let selectedTool = null;

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

let painterPaths = [];
let mx = -1000;
let my = -1000;

let chaosEffect = null;

let paintingSenderInterval = null;
function startPaintingSender() {
  stopPaintingSender();
  paintingSenderInterval = setInterval(() => {
    sendPainting();
  }, 100);
}
function stopPaintingSender() {
  if (paintingSenderInterval) {
    clearInterval(paintingSenderInterval);
    paintingSenderInterval = null;
  }
}

function setPaintingToolsEnabled(enabled) {
  document.getElementById("painting-tools").style.display = enabled ? "block" : "none";
  setInputEnabled(enabled);
  if (enabled) {
    initPalette();
    selectTool(TOOL_PENCIL);
    startPaintingSender();
  } else {
    selectTool(null);
    stopPaintingSender();
  }
}

function setPromptSelectionEnabled(enabled) {
  document.getElementById("prompt-selection").style.display = enabled ? "block" : "none";
}

function setPromptOptions(prompts) {
  setPaintingPrompt("Vote for a prompt!");
  for (let i = 0; i < 3; i++) {
    const button = document.getElementById("prompt" + i);
    button.innerText = prompts[i];
    button.onclick = () => sendVoteCommand(prompts[i]);
  }
}

function setVoteOptions(voteOptions) {
  for (let i = 0; i < 5; i++) {
    const button = document.getElementById("vote" + i);
    if (voteOptions[i]) {
      button.style.display = "block";
      button.style.backgroundImage = "url('img/" + voteOptions[i] + ".png')";
      button.onclick = () => sendVoteCommand(voteOptions[i]);
    } else {
      button.style.display = "none";
    }
  }
}

function setPaintingPrompt(prompt) {
  document.getElementById("painter-prompt-text").innerText = prompt;
}

function setTimerSecondsLeft(secondsLeft) {
  if (secondsLeft < 0) {
    document.getElementById("timer-text").style.display = "none";
  } else {
    document.getElementById("timer-text").style.display = "block";
    document.getElementById("timer-number").innerText = secondsLeft;
  }
}

function setPainting(graphics) {
  painterPaths = graphics.paths;
  mx = graphics.mx;
  my = graphics.my;
  drawPainterCanvas();
}

function clearPainting() {
  painterPaths.splice(0, painterPaths.length);
  mx = -1000;
  my = -1000;
  drawPainterCanvas();
}

function sendPainting() {
  sendSetPaintingCommand({
    paths: painterPaths,
    mx: mx,
    my: my,
  });
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
  drawPainterCanvas();
  sendPainting();
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

function drawPainterCanvas() {
  const painterCanvas = document.getElementById("painter-canvas");
  const ctx = painterCanvas.getContext("2d");
  drawPainting(painterCanvas, painterPaths, selectedBackgroundName);

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
    console.log("flashlight");
    ctx.beginPath();
    ctx.rect(0, 0, painterCanvas.width, painterCanvas.height);
    ctx.arc(mx, my, 200, 2 * Math.PI, 0);
    ctx.fillStyle = "#000";
    ctx.fill("evenodd");
  }
}

// CHAOS EFFECTS

function activateChaosEffect(effect, duration) {
  const painterCanvas = document.getElementById("painter-canvas");
  console.log("activate chaos effect: " + effect);
  chaosEffect = effect;
  setTimeout(() => {
    deactivateChaosEffect(effect);
  }, duration);

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
  const painterCanvas = document.getElementById("painter-canvas");
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
  document.getElementById(TOOL_PENCIL).classList.remove("selected");
  document.getElementById(TOOL_ERASER).classList.remove("selected");
  if (tool) document.getElementById(tool).classList.add("selected");
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

const makeColorRect = (i) => {
  const w = 122;
  const h = 80;
  return {
    x: 2 + (i % 2) * (w + 16),
    y: 2 + Math.floor(i / 2) * (h + 12),
    w,
    h,
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
  const ctx = paletteElem.getContext("2d");
  ctx.clearRect(0, 0, paletteElem.width, paletteElem.height);
  for (let i = 0; i < palette.length; i++) {
    ctx.beginPath();
    const rect = makeColorRect(i);
    ctx.rect(rect.x, rect.y, rect.w, rect.h);
    ctx.fillStyle = palette[i];
    ctx.fill();
    if (i == selectedColor) {
      ctx.strokeStyle = "#000";
      ctx.lineWidth = 4;
      ctx.stroke();
    }
  }
}

function distanceSquared(p1, p2) {
  return Math.pow(p1.x - p2.x, 2) + Math.pow(p1.y - p2.y, 2);
}
