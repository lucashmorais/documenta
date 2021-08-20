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
		DataTableSkeleton,
		Toolbar,
		Button,
		ToolbarBatchActions,
		ToolbarContent,
		ToolbarSearch,
		Tag
	} from "carbon-components-svelte";

	let selectedRowIds = [];
	let editModalIsOpen = false;
	let deleteModalIsOpen;

	let headers=[{ key: 'assunto', value: 'Assunto' }, { key: 'centro', value: 'Centro' }, { key: 'tipo', value: 'Tipo' }, { key: 'pend', value: 'Pendência Atual' }]
	
	//TODO: Remove the following in favor of promise-based data communication

	var draftProcessesPromise;
	var activeProcessesPromise;
	var finishedOrBlockedProcessesPromise;
	
	async function coreProcessUpdater(resolve, reject, set) {
		let rows = []
		let promises = []
		for (let queryParam of set) {
			promises.push( new Promise((resolve, reject) => {
					fetch("http://localhost:3123/api/v1/processes?statusString=" + queryParam).
						then((response)=>response.json().
							then(function (processes) {
								console.log(processes)
								let processObj = {}
								console.log("[updateProcesses::processes]:", processes)
								for (const p of processes) {
									processObj = {}
									processObj.id = p.ID
									processObj.assunto = p.Title
									processObj.centro = p.Center.Name
									processObj.tipo = p.ProcessType.Name

									//TODO: GET THE FOLLOWING FROM THE DB!
									processObj.pend = "Encaminhamento final"

									console.log("[updateProcesses]: Process just built: ", processObj)
									rows.push(processObj)
								}
								resolve()
							}
						)
					)
				}
			))
		}
		await Promise.all(promises)
		resolve(rows)
	}

	export function updateProcesses() {
		let set = []
		set = ["Rascunho"]
		draftProcessesPromise = new Promise((resolve, reject) => {
			coreProcessUpdater(resolve, reject, set)
		})
		set = ["Ativo"]
		activeProcessesPromise = new Promise((resolve, reject) => {
			coreProcessUpdater(resolve, reject, set)
		})
		set = ["Concluído", "Interditado"]
		finishedOrBlockedProcessesPromise = new Promise((resolve, reject) => {
			coreProcessUpdater(resolve, reject, set)
		})
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
			{#await draftProcessesPromise}
				<DataTableSkeleton />
			{:then draftProcesses}
				<DataTable
					headers={headers}
					rows={draftProcesses}
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
			{/await}
		</div>
	<h2>Processos em andamento</h2>
		<div class="content2">
			{#await activeProcessesPromise}
				<DataTableSkeleton />
			{:then activeProcesses}
				<DataTable
					headers={headers}
					rows={activeProcesses}
				/>
			{/await}
		</div>
	<h2>Processos finalizados</h2>
		<div class="content2">
			{#await finishedOrBlockedProcessesPromise}
				<DataTableSkeleton />
			{:then finishedOrBlockedProcesses}
				<DataTable
					headers={headers}
					rows={finishedOrBlockedProcesses}
				/>
			{/await}
		</div>
</div>
