import Home from './Home.svelte';
import '@fortawesome/fontawesome-free/js/all.js';
import 'carbon-components-svelte/css/all.css';

const app = new Home({
	target: document.body,
	props: {
		name: 'Main',
	},
});

export default app;
