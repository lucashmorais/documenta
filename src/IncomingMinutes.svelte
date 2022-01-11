<script>
	import "carbon-components-svelte/css/all.css";
	import StatusBar from "./StatusBar.svelte"
	import { writable } from 'svelte/store';
	import SimpleConfirmationModal from "./SimpleConfirmationModal.svelte"
	import MinuteBlock from "./MinuteBlock.svelte"
	import {
		coreProcessUpdater,
		postNewMinute,
		coreUploadFile,
		getUnassignedMinutes,
		getAssignedMinutes
	} from "./utils.js"
	import {
		Dropdown,
		TextInput,
		Modal,
		DataTable,
		DataTableSkeleton,
		FileUploader,
		ToastNotification
	} from "carbon-components-svelte";
	import { getEndpointPrefix } from "./config-helper.js"
	
// 	import { marked } from 'marked';
// 	import { jsPDF } from "jspdf";
// 	let source = `
// # H1 heading

// ## H2 heading

// ### H3 heading

// --------

// **bold text**

// *italicized text*

// --------

// 1. First item
// 2. Second item
// 3. Third item

// - First item
// - Second item
// - Third item

// [Svelte](https://svelte.dev/)
// `;
	// let htmlSource = marked(source);

	// var doc = new jsPDF();
	// doc.setFontSize(5);
	
	// doc.html(
	// 	htmlSource,
	// 	{
	// 	callback: function (doc) {
	// 		doc.save('a4.pdf');
	// 	}
	// });

	let unassignedMinutesPromise = getUnassignedMinutes()
	let assignedMinutesPromise = getAssignedMinutes()
	
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
						unassignedMinutesPromise = getUnassignedMinutes()
						assignedMinutesPromise = getAssignedMinutes()
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
					() => fireToastNotification("error"),
					null,
					() => unassignedMinutesPromise = getUnassignedMinutes()
				)
				creationModalIsOpen = false;
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

	.content1 {
		margin: 0 5vw;
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

	.container{
		background: #ff3e00d6;
		padding:10px 30px;
	}
	.header {
		height: 10vh;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.header-title {
		margin: 0;
		color:#fff;
	}
	.markdown-editor {
		width: 100%;
		display: flex;
		align-items:flex-start;
		justify-content: space-evenly;
	}
	.left-panel, .right-panel {
		width: 50%;
		border: solid 1px black;
		height: 85vh;
		background: #ffffff;
	}
	.right-panel {
		overflow: auto;
	}
	.source {
		border: none;
		width: 100%;
		height: 100%;
		background: #001628;
		color: #83ba52;
	}
	.source:focus {
		outline: none;
	}
	.output {
		width: 100%;
		padding: 0 2em;
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
	<MinuteBlock
		createMinuteAssignmentModalOpeningHandler={createMinuteAssignmentModalOpeningHandler}
		createProcessCreationModalOpeningHandler={createProcessCreationModalOpeningHandler}
		creationModalHandler={() => creationModalIsOpen = true}
		minutesPromise={unassignedMinutesPromise}
		blockTitle="Minutas esperando atribuição"
		enableMinutePlaceholder={false}
	/>
	<MinuteBlock
		minutesPromise={assignedMinutesPromise}
		blockTitle="Minutas atribuídas"
		disableEditing={true}
		enableProcessView={true}
	/>
</div>

<!-- <main class="container">
	<header class="header">
	    <h1 class="header-title">Svelte markdown editor</h1>
	</header>
	
	<div class="markdown-editor">
	    <div class="left-panel">
		<textarea bind:value={source} class="source"></textarea>
	    </div>
	
	    <div class="right-panel">
		<div class="output">{@html htmlSource}</div>
	    </div>
	</div>
</main> -->