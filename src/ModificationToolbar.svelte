<script>
	import {
		ButtonSet,
		Button
	} from 'carbon-components-svelte'

	import { createEventDispatcher } from "svelte";
	
	const dispatch = createEventDispatcher();
	
	function signalCommit() {
		dispatch("commit");
	}
	
	function signalCancellation() {
		dispatch("cancel");
	}
</script>
<style>
	.floater {
		position: sticky;
		bottom: 5%;
		width: fit-content;
		margin-left: auto;
		margin-right: 5%;
	}
</style>

<div class=floater>
	<ButtonSet stacked={true}>
		<Button on:click={signalCancellation} size="small" kind="secondary">Limpar</Button>
		<Button on:click={signalCommit} size="small">Aplicar</Button>
	</ButtonSet>
</div>
