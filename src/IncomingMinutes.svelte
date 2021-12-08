<script>
	import "carbon-components-svelte/css/all.css";
	import StatusBar from "./StatusBar.svelte"
	import SimpleConfirmationModal from "./SimpleConfirmationModal.svelte"
	import { getNameFromUser, coreProcessUpdater } from "./utils.js"
	import {
		Tile,
		Modal,
		DataTable,
		DataTableSkeleton,
		Button
	} from "carbon-components-svelte";
	import DocumentAdd16 from "carbon-icons-svelte/lib/DocumentAdd16";
	import WatsonHealth3DCurveAutoColon16 from "carbon-icons-svelte/lib/WatsonHealth3DCurveAutoColon16";
	import { getEndpointPrefix } from "./config-helper.js"

	export function getMinutes() {
		return new Promise((resolve) => {
			fetch(getEndpointPrefix() + "/api/v1/minutes?unassigned=true&incoming=true").
				then((response)=>response.json().
					then(function (minutes) {
						// console.log("[Minutes::updateMinutes::minutes]: ", minutes)
						for (let i = 0; i < minutes.length; i++) {
							let a = minutes[i]
							console.log(a)
						}
						resolve(minutes)
					}
				)
			)
		})
	}
	let minutesPromise = getMinutes()
	
	let assignmentModalIsOpen = false;
	let processCreationModalIsOpen = false;
	let processesPromise;
	let selectedRowIds;
	let selectedMinuteID;

	let headers=[{ key: 'assunto', value: 'Assunto' }, { key: 'centro', value: 'Centro' }, { key: 'tipo', value: 'Tipo' }, {key: 'estado', value: 'Estado'}, {key: 'autor', value: 'Autor'}]
	
	export function updateProcesses() {
		let set = ["Ativo", "Rascunho"]
		processesPromise = new Promise((resolve, reject) => {
			coreProcessUpdater(resolve, reject, set, true)
		})
	}
	updateProcesses();
	
	$: if(selectedRowIds) {
		console.log("selectedRowIds: ", selectedRowIds)
	}
	
	function createMinuteAssignmentModalOpeningHandler(newMinuteID) {
		return () => {
			selectedMinuteID = newMinuteID;
			assignmentModalIsOpen = true;
		}
	}
	
	function createProcessCreationModalOpeningHandler(newMinuteID) {
		return () => {
			selectedMinuteID = newMinuteID;
			processCreationModalIsOpen = true;
		}
	}
	
	function coreHandleMinuteAssignment(targetProcessID) {
		return new Promise((resolve) => {
			fetch(getEndpointPrefix() + `/api/v1/minute/${selectedMinuteID}?processID=${targetProcessID}`, {method: 'PATCH'}).
				then((response)=>response.json().
					then((resp) => {
						console.log("[handleMinuteAssignment::response]: ", resp)
						assignmentModalIsOpen = false;
						minutesPromise = getMinutes()
						resolve(resp)
					})
				)
		})
	}
	
	function handleMinuteAssignment() {
		coreHandleMinuteAssignment(selectedRowIds[0])
	}
	
	function createProcess() {
		let requestBody = JSON.stringify({
			"title": "processo_minuta",
			"summary": "Processo criado a partir de uma minuta.",
			"typeID": 1,
			"centerID": 2,
		});
		
		console.log("[submitForm:registering:requestBody]: ", requestBody);

		return new Promise((resolve) => {
			fetch(getEndpointPrefix() + "/api/v1/process", {
				method: "post",
				
				body: requestBody,
				headers: {
					"Content-type": "application/json; charset=UTF-8"
				}
			}).then((response)=>response.json().
				then(function (createdProcess) {
					resolve(createdProcess)
				})
			)
		})
	}
</script>

<style>
	h1 {
		text-align: center;
		margin: 2em 0 1em 0;
	}

	h1:focus-visible {
		outline-style: none;
	}

	
	h2:not(.realh2) {
		font-size: 20px;
	}
	
	h3 {
		font-size: 16px;
		font-weight: bold;
		margin-top: 1em;
	}
	
	h3:first-child {
		margin-top: 0;
	}

	.content1 {
		margin: 0 5vw;
	}
	
	.content2 {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1fr;
		grid-template-rows: 1fr 1fr 1fr;
		gap: 2em 2em;
		grid-auto-flow: row;
		grid-template-areas: ". . . ." ". . . ." ". . . .";
		margin-top: 2em;
	}
	
	.element {
		/* max-width: 200px; */
	}
								
	.actionSet {
		margin-top: 1em;
	}
	
	.actionSet * {
		margin-right: 3em
	}
</style>

<Modal
  bind:open={assignmentModalIsOpen}
  modalHeading="Atribuir minuta a um processo"
  primaryButtonText="Confirmar"
  secondaryButtonText="Cancelar"
  on:click:button--secondary={() => (assignmentModalIsOpen = false)}
  on:open
  on:close
  on:submit={handleMinuteAssignment}
>
	{#await processesPromise}
		<DataTableSkeleton />
	{:then processes}
		<DataTable
			headers={headers}
			rows={processes}
			radio
			bind:selectedRowIds
		/>
	{/await}
</Modal>

<SimpleConfirmationModal
	customMessage="Tem certeza de que gostaria de criar um processo incluindo esta minuta?"
	bind:open={processCreationModalIsOpen}
	on:actionConfirmed={async function() {
		let createdProcess = await createProcess()
		console.log("[SimpleConfirmationModal::__submitCallback::createdProcess]: ", createdProcess)
		processCreationModalIsOpen = false;
		coreHandleMinuteAssignment(createdProcess.ID)
	}}
/>

<StatusBar />
<h1>Documenta</h1>

<div class="content1">
	<h2>Minutas não alocadas</h2>
		<div class="content2">
			{#await minutesPromise then minutes}
				{#each minutes as minute}
					<div class="element">
						<Tile>
							<h3>
								Autor
							</h3>
							<p>
								{getNameFromUser(minute.User)}
							</p>
							<h3>
								Centro
							</h3>
							<p>
								{minute.Center.Name}
							</p>
							<h3>
								Resumo breve
							</h3>
							<p>
								{minute.Content}
							</p>
							
							<div class="actionSet">
								<Button kind="secondary" icon={WatsonHealth3DCurveAutoColon16} iconDescription="Atribuir" on:click={createMinuteAssignmentModalOpeningHandler(minute.ID)} />
								<Button kind="tertiary" icon={DocumentAdd16} iconDescription="Gerar processo" on:click={createProcessCreationModalOpeningHandler(minute.ID)} />
							</div>
						</Tile>
					</div>
				{/each}
			{/await}
		</div>
</div>
