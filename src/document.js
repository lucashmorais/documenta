import Document from './Document.svelte';
import '@fortawesome/fontawesome-free/js/all.js';
import 'carbon-components-svelte/css/all.css';

const app = new Document({
	target: document.body,
	props: {
		name: 'Main',
	},
});

export default app;
