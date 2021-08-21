<script>
	import { Grid, Row, Column } from "carbon-components-svelte";
	import { DataTable } from "carbon-components-svelte";
	import { Button } from "carbon-components-svelte";

	let headers=[{ key: 'centro', value: 'Centro' }, { key: 'tipo', value: 'Tipo' }, { key: 'pend', value: 'Pendência Atual' }]
	let rows=[{ id: 'a', centro: 'sm', tipo: "Consulta", pend: 'Revisão do defensor' }]
	
	export let processPromise;

	$: processPromise.then((process) => {
		rows=[{ id: 'a', centro: process.Center.Name, tipo: process.ProcessType.Name, pend: 'Revisão do defensor' }]
	})
</script>

<style>
	.mainWrapper {
		margin-bottom: 5vh;
		width: 90%;
		margin: auto;
		justify-content: center;
	}
</style>

<div class="mainWrapper">
	<DataTable
		size="short"
		headers={headers}
		rows={rows}
	/>
</div>