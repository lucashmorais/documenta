import App from './App.svelte';
import '@fortawesome/fontawesome-free/js/all.js';
import 'carbon-components-svelte/css/all.css';

const app = new App({
	target: document.body,
	props: {
		name: 'Main',
	},
});

export default app;
