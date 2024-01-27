let viewerCanvas;
let viewerTimerNumberElem;

function initViewer() {
  viewerCanvas = document.getElementById("viewer-canvas");
  viewerTimerNumberElem = document.getElementById("viewer-timer-number");

  // countdown
  viewerTimerNumberElem.innerText = TIMER_SECONDS;

  drawPainting(viewerCanvas, []);
}