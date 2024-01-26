let canvas;
let ctx;

// CONSTANTS
const width = 1600;
const height = 900;

function init() {
  canvas = document.getElementById("canvas");
  ctx = canvas.getContext("2d");

  const resize = (event) => {
    const wrapper = document.getElementById("wrapper");
    const scale = Math.min(window.innerWidth / width, window.innerHeight / height);
    wrapper.style.transform = "translate(-50%, -50%) scale(" + scale + ")";
  };
  resize();
  window.addEventListener("resize", resize);
}
