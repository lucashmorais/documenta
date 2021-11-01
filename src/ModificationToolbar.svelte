<script>
	import {
		ButtonSet,
		Button
	} from 'carbon-components-svelte'
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
		<Button size="small" kind="secondary">Cancelar</Button>
		<Button size="small">Aplicar</Button>
	</ButtonSet>
</div>
