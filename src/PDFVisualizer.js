import PDFVisualizer from './PDFVisualizer.svelte';
// import '@fortawesome/fontawesome-free/js/all.js';
// import 'carbon-components-svelte/css/all.css';

const app = new PDFVisualizer({
	target: document.body,
	props: {
		name: 'Main',
	},
});

export default app;
