<script>
	import 'carbon-components-svelte/css/all.css';
	import StatusBar from './StatusBar.svelte'
	// import DataTable from './DataTable.svelte'
	import DataTable from './DataTable/DataTable.svelte'

	let reps = [1,2,3,4,5,6]

	let headers=[{ key: 'assunto', value: 'Assunto' }, { key: 'centro', value: 'Centro' }, { key: 'tipo', value: 'Tipo' }, { key: 'pend', value: 'Pendência Atual' }]
	let rows=[{ id: 'a', assunto: 'Assunto 0', centro: 'cs', tipo: 3000, pend: 'Revisão do defensor' }, { id: 'b', assunto: 'Assunto 1', centro: 'brs', tipo: 443, pend: 'Revisão inicial do Secretário' }, { id: 'c', assunto: 'Assunto 2', centro: 'cur', tipo: 80, pend: 'Geração de protocolo' }, { id: 'd', assunto: 'Assunto 6', centro: 'rib', tipo: 3000, pend: 'Round robin' }, { id: 'e', assunto: 'Assunto 4', centro: 'for', tipo: 443, pend: 'Discussão Geral' }, { id: 'f', assunto: 'Assunto 5', centro: 'bh', tipo: 80, pend: 'Encaminhamento final para o ctr' }]
	let rowsShort=[{ id: 'f', assunto: 'Assunto 5', centro: 'bh', tipo: 80, pend: 'Encaminhamento final para o ctr' }]
</script>

<style>
	h1 {
		text-align: center;
		margin: 2em 0 1em 0;
	}

	h1:focus-visible {
		outline-style: none;
	}

	h2 {
		font-size: 16px;
	}

	.content1 {
		margin: 0 5vw;
	}

	.content2 {
		/* margin: 1vh 5vw 4vh 5vw; */
		margin-top: 1em;
		margin-bottom: 2em;
	}

</style>

<StatusBar />
<h1>Documenta</h1>

<div class="content1">
	<h2>Processos esperando ação do usuário</h2>
		<div class="content2">
			<DataTable
				headers={headers}
				rows={rowsShort}
			/>
		</div>
	<h2>Processos em andamento</h2>
		<div class="content2">
			<DataTable
				headers={headers}
				rows={rows}
			/>
		</div>
	<h2>Processos finalizados</h2>
		<div class="content2">
			<DataTable
				headers={headers}
				rows={rows}
			/>
		</div>
</div>
