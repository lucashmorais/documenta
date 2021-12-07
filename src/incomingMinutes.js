import IncomingMinutes from "./IncomingMinutes.svelte";
import "@fortawesome/fontawesome-free/js/all.js";
import "carbon-components-svelte/css/all.css";

const app = new IncomingMinutes({
  target: document.body,
});

export default app;
