<script>
	import { Grid, Row, Column } from "carbon-components-svelte";
	import { DataTable } from "carbon-components-svelte";

	let headers=[{ key: 'centro', value: 'Centro' }, { key: 'tipo', value: 'Tipo' }, { key: 'pend', value: 'Pendência Atual' }]
	let rows=[{ id: 'a', centro: 'sm', tipo: "Consulta", pend: 'Revisão do defensor' }]
</script>

<style>
	div{
		margin-bottom: 5vh;
	}
</style>

<div>
	<DataTable
		size="compact"
		headers={headers}
		rows={rows}
	/>
</div>