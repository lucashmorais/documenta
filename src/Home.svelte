<script>
	import 'carbon-components-svelte/css/all.css';
	import StatusBar from './StatusBar.svelte'
	import ProcessModal from './ProcessModal.svelte'
	// import DataTable from './DataTable.svelte'
	// import DataTable from './DataTable/DataTable.svelte'
	import TrashCan16 from "carbon-icons-svelte/lib/TrashCan16";
	import Edit16 from "carbon-icons-svelte/lib/Edit16";
	import {
		DataTable,
		Toolbar,
		Button,
		ToolbarBatchActions,
		ToolbarContent,
		ToolbarSearch,
		Tag
	} from "carbon-components-svelte";

	let reps = [1,2,3,4,5,6]
	let selectedRowIds = [];
	let editModalIsOpen = false;
	let deleteModalIsOpen;

	let headers=[{ key: 'assunto', value: 'Assunto' }, { key: 'centro', value: 'Centro' }, { key: 'tipo', value: 'Tipo' }, { key: 'pend', value: 'Pendência Atual' }]
	let rows=[]
	// let rows=[{ id: 'a', assunto: 'Assunto 0', centro: 'cs', tipo: 3000, pend: 'Revisão do defensor' }, { id: 'b', assunto: 'Assunto 1', centro: 'brs', tipo: 443, pend: 'Revisão inicial do Secretário' }, { id: 'c', assunto: 'Assunto 2', centro: 'cur', tipo: 80, pend: 'Geração de protocolo' }, { id: 'd', assunto: 'Assunto 6', centro: 'rib', tipo: 3000, pend: 'Round robin' }, { id: 'e', assunto: 'Assunto 4', centro: 'for', tipo: 443, pend: 'Discussão Geral' }, { id: 'f', assunto: 'Assunto 5', centro: 'bh', tipo: 80, pend: 'Encaminhamento final para o ctr' }]

	var processesPromise;
	export function updateProcesses() {
		processesPromise = new Promise((resolve, reject) => {
			fetch("http://localhost:3123/api/v1/processes").
				then((response)=>response.json().
					then(function (processes) {
						console.log(processes)
						rows = []
						let processObj = {}
						console.log("[updateProcesses::processes]:", processes)
						for (const p of processes) {
							processObj = {}
							processObj.id = p.ID
							processObj.assunto = p.Title
							
							//TODO: GET THE FOLLOWING FROM THE DB!
							processObj.centro = "Brasília"
							processObj.pend = "Encaminhamento final"
							processObj.tipo = p.ProcessType.Name

							console.log("[updateProcesses]: Process just built: ", processObj)
							rows.push(processObj)
						}
						resolve(processes)
					}
				)
			)
		})
		return processesPromise
	}
	updateProcesses();
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

<ProcessModal 
	bind:open={editModalIsOpen}
	on:backendModification={updateProcesses}
/>
<StatusBar />
<h1>Documenta</h1>

<div class="content1">
	<h2>Processos esperando ação do usuário</h2>
		<div class="content2">
			<DataTable
				headers={headers}
				rows={rows}
			>
				<Toolbar>
					<ToolbarBatchActions>
					  <Button on:click={() => deleteModalIsOpen = true} icon={TrashCan16}>Eliminar</Button>
					  {#if selectedRowIds.length < 2}
						  <Button on:click={() => editModalIsOpen = true} icon={Edit16}>Editar</Button>
					  {/if}
					</ToolbarBatchActions>
					<ToolbarContent>
					  <!-- <ToolbarSearch /> -->
					  <Button on:click={() => editModalIsOpen = true}>Novo processo</Button>
					</ToolbarContent>
				</Toolbar>
			</DataTable>
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
