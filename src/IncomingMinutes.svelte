<script>
	import "carbon-components-svelte/css/all.css";
	import StatusBar from "./StatusBar.svelte"
	import { writable } from 'svelte/store';
	import SimpleConfirmationModal from "./SimpleConfirmationModal.svelte"
	import {
		getNameFromUser,
		coreProcessUpdater,
		postNewMinute,
		coreUploadFile
	} from "./utils.js"
	import {
		Dropdown,
		TextInput,
		Tile,
		ClickableTile,
		Modal,
		DataTable,
		DataTableSkeleton,
		FileUploader,
		ToastNotification,
		Button
	} from "carbon-components-svelte";
	import DocumentAdd16 from "carbon-icons-svelte/lib/DocumentAdd16";
	import Add32 from "carbon-icons-svelte/lib/Add32";
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
	let creationModalIsOpen = false;
	let processCreationModalIsOpen = false;
	let processesPromise;
	let selectedRowIds;
	let selectedMinuteID;
	let currentAttachment;

	let headers=[{ key: 'assunto', value: 'Assunto' }, { key: 'centro', value: 'Centro' }, { key: 'tipo', value: 'Tipo' }, {key: 'estado', value: 'Estado'}, {key: 'autor', value: 'Autor'}]
	
	let formState = {
		"selectedCenter": 0,
		"protocol": "",
		"description": ""
	}
	
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
	
	function handleMinuteCreation() {
		if (!currentAttachment) {
			fireToastNotification("missing_attachment");			
			return;
		}
		let fileID = 0;
		coreUploadFile(currentAttachment, 0).then((response) => response.json().then(
			function (fileIDs) {
				console.log(fileIDs)
				fileID = fileIDs[0];
				postNewMinute(
					"",
					formState.description,
					0,
					available_centers[formState.selectedCenter].id,
					fileID,
					true,
					formState.protocol,
					() => fireToastNotification("error")
				)
				creationModalIsOpen = false;
				minutesPromise = getMinutes()
			}
		))
	}
	
	let available_centers = []
	function updateCenters() {
		return new Promise((resolve, reject) => {
			fetch(getEndpointPrefix() + "/api/v1/centers").
				then((response)=>response.json().
					then(function (centers) {
						available_centers = []
						console.log("[updateCenters::centers]: ", centers)
						for (const u of centers) {
							let centerObj = {}
							centerObj.id = u.ID
							centerObj.text= u.Name
							centerObj.shortName= u.ShortName
							centerObj.description = u.Description
							console.log(centerObj)
							available_centers.push(centerObj)
						}
						resolve(centers)
					})
				)
		})
	}
	updateCenters();
	
	let notifications = writable({});
	$notifications = [];
	
	function fireToastNotification(kind, extras = null) {
		const s = new Date().toLocaleString()
		const l = $notifications.length;	// get our current items list count

		switch(kind) {
			case "success":
				$notifications[l] = {
					kind: "success",
					title: "Sucesso",
					subtitle: "A minuta foi criada com sucesso",
					caption: s,
					iconDescription: "Fechar notificação"
				}
			break;
			case "missing_attachment":
				$notifications[l] = {
					kind: "warning",
					title: "Anexo faltante",
					subtitle: "A minuta não pôde ser criada porque nenhum anexo foi indicado.",
					caption: s,
					iconDescription: "Fechar notificação"
				}
			break;
			default:
				$notifications[l] = {
					kind: "error",
					title: "Erro",
					subtitle: "A minuta não pôde ser criada. Revise os dados de entrada ou contate o administrador do sistema.",
					caption: s,
					iconDescription: "Fechar notificação"
				}
		}
	}
	
	// setInterval(() => {
	// 	fireToastNotification("error")
	// }, 2000);
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

	.absoluteToastWrapper {
		position: absolute;
		bottom: 1em;	
		right: 1em;
	}

	.stickyToast {
		position: sticky;
		z-index: 1000;
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

<Modal
  bind:open={creationModalIsOpen}
  modalHeading="Adicionar minuta"
  primaryButtonText="Confirmar"
  secondaryButtonText="Cancelar"
  on:click:button--secondary={() => (creationModalIsOpen = false)}
  on:open
  on:close
  on:submit={handleMinuteCreation}
>
	<TextInput labelText="Protocolo de entrada" bind:value={formState.protocol} />
	<TextInput labelText="Descrição breve" bind:value={formState.description} />
	<Dropdown
		titleText="Centro"
		bind:selectedIndex={formState.selectedCenter}
		items={available_centers}
	/>
	<div class="actionSet">
		<FileUploader
		  on:change={function (event) { 
			  currentAttachment = event.target.files[0]
			  console.log(currentAttachment)
		}}
		  labelTitle="Conteúdo da minuta"
		  buttonLabel="Adicionar arquivo"
		  accept={['.jpg', '.jpeg', '.png', '.pdf', '.PDF']}
		  status="complete"
		/>
	</div>
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
	<div class="absoluteToastWrapper">
		<div class="stickyToast">
		      {#each $notifications as toast}
				<svelte:component this={ToastNotification} {...toast}/>
		      {/each}
		</div>
	</div>
	<h2>Minutas não alocadas</h2>
		<div class="content2">
			<div class="element">
				<ClickableTile
					style="height: 100%; display: flex; align-items: center; justify-content: center; flex-direction: column"
					on:click={() => creationModalIsOpen = true}
				>
					<Add32 />
					<p>Adicionar minuta</p>
				</ClickableTile>
			</div>
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
