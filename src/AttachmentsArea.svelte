<script>
	import { ClickableTile } from "carbon-components-svelte";
	import DocumentDownload24 from "carbon-icons-svelte/lib/DocumentDownload24";
	import { Button } from "carbon-components-svelte";
	import { FileUploaderButton } from "carbon-components-svelte";
	import { Carousel } from "svelte-images";
	import { getEndpointPrefix } from "./config-helper.js"
	import { coreUploadFile } from "./utils.js"

	const { Modal, open, close } = Carousel;
	let fileUploader
	let attachmentsPromise;
	let files;
	let images;

	export let processID;
	export let modRightsPromise;
	let modRights = false;

	async function uploadFile() {
		if (!modRights) {
			return;
		}

		console.log(fileUploader)
		console.log(fileUploader.files)
		for (const file of fileUploader.files) {
			await coreUploadFile(file, processID)
		}
		//TODO: ENSURE THIS ONLY RUNS AFTER THE UPLOAD ACTIONS HAVE BEEN COMPLETED
		updateAttachments()
	}

	async function deleteFileWithID(fileID) {
		if (!modRights) {
			return;
		}

		fetch(getEndpointPrefix() + "/api/v1/file/" + fileID, {method: "DELETE"}).
			then((response)=> {
				updateAttachments()
			})
	}

	export function updateAttachments() {
		attachmentsPromise = new Promise((resolve, reject) => {
			if (processID != "" && processID != "0") {
				fetch(getEndpointPrefix() + "/api/v1/files?processID=" + processID).
					then((response)=>response.json().
						then(function (attachments) {
							files = attachments
							images = []
							let imgCounter = 0
							for (let i = 0; i < attachments.length; i++) {
								let a = attachments[i]
								console.log(a)
								a.src = getEndpointPrefix() + "/api/v1/file/" + a.ID;
								console.log(a.src);
								if (a.ContentType.startsWith("image")) {
									a.imgID = imgCounter
									imgCounter++
									images.push(a)
								}
							}
							resolve(attachments)
						}
					)
				)
			} else {
				resolve([])
			}
		})
	}
	
	function popModal(idx, type = "image/jpeg") {
		console.log("derp1")

		let imageFunction = function() {
			console.log("derp2")
			setTimeout(() => {
				console.log("derp3")
				console.log(files[idx].src)
				open(files, idx);
			}, 0);
		}
		
		function getOnlyImages(fileArray) {
			return fileArray.filter(f => f.ContentType.startsWith("image"))
		}

		imageFunction = function() {
			console.log("derp2")
			setTimeout(() => {
				console.log("derp3")
				console.log(files[idx].src)
				open(images, files[idx].imgID);
			}, 0);
		}

		let pdfFunction = function() {
			// window.open(files[idx].src, "_blank").focus();
			window.open(getEndpointPrefix() + "/PDFVisualizer.html?id=" + String(files[idx].ID), "_blank").focus();
		}
		
		if (type.startsWith("image"))
			return imageFunction;

		switch(type) {
			case "application/pdf":
				return pdfFunction;
			default:
				return downloadFunction;
		}
	}
	
	function updatePlainRightsVar(p) {
		if (p) {
			p.then((canModify) => modRights = canModify)
		}
	}
	
	$: updatePlainRightsVar(modRightsPromise)

	updateAttachments();
</script>

<style>
	.fileGrid {
		display: grid;
		grid-template-columns: 50% 50%;
		grid-template-rows: 1fr;
		gap: 1em 1em;
		grid-template-areas:
			". .";
	}

	.fileTile {
		display: grid;
		grid-template-columns: 1.5fr 0.5fr;
		/* grid-template-rows: 67px; */
		grid-auto-rows: 67px;
		gap: 0px 0px;
		grid-template-areas:
		". .";
		/* max-width: 50%; */
	}

	.fileName {
		padding-left: 1em;
		font-size: 14px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: normal;
		overflow-wrap: break-word;
		word-break: break-word;
	}

	:global(.deleteButton) {
		height: 100%!important;
	}
</style>


<Modal />
<div class="fileGrid">
	{#await attachmentsPromise then files}
		{#each files as file, i}
			<div class="fileTile">
					<ClickableTile on:click={popModal(i, file.ContentType)}>
						<div style="display: flex; align-items: center; height: 2.5em;">
							<DocumentDownload24 />
							<p class="fileName">{file.Name}</p>
						</div> 
					</ClickableTile>

				{#await modRightsPromise then canModify}
					{#if canModify}
						<div class="deleteHolder">
							<Button kind="danger" class="deleteButton" on:click={deleteFileWithID(file.ID)}>Delete</Button>
						</div>
					{/if}
				{/await}
			</div>
		{/each}
	{/await}
	
	{#await modRightsPromise then canModify}
		{#if canModify}
			<FileUploaderButton multiple disableLabelChanges={true} labelText="Adicionar arquivo" bind:ref={fileUploader} on:change={uploadFile}>Novo anexo</FileUploaderButton>
		{/if}
	{/await}
</div>