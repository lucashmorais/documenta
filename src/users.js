import Users from './Users.svelte';
import '@fortawesome/fontawesome-free/js/all.js';
import 'carbon-components-svelte/css/all.css';

const app = new Users({
	target: document.body,
});

export default app;
