<script>
	import { ClickableTile } from "carbon-components-svelte";
	import DocumentDownload24 from "carbon-icons-svelte/lib/DocumentDownload24";
	import { Button } from "carbon-components-svelte";
	import { FileUploaderButton } from "carbon-components-svelte";
	import { Carousel } from "svelte-images";

	const { Modal, open, close } = Carousel;
	let fileUploader
	let attachmentsPromise;
	let files;

	async function coreSubmit(file) {
		// https://javascript.info/formdata
		let formData = new FormData();
		formData.append("documents", file, file.name);

		const response = await fetch('http://localhost:3123/api/v1/files', {
			method: 'POST',
			body: formData
		});

		let result = await response.json();
		return result.status != "success"
	}

	async function uploadFile() {
		console.log(fileUploader)
		console.log(fileUploader.files)
		for (const file of fileUploader.files) {
			await coreSubmit(file)
		}
		//TODO: ENSURE THIS ONLY RUNS AFTER THE UPLOAD ACTIONS HAVE BEEN COMPLETED
		updateAttachments()
	}

	async function deleteFileWithID(fileID) {
		fetch("http://localhost:3123/api/v1/file/" + fileID, {method: 'DELETE'}).
			then((response)=> {
				updateAttachments()
			})
	}

	export function updateAttachments() {
		attachmentsPromise = new Promise((resolve, reject) => {
			fetch("http://localhost:3123/api/v1/files").
				then((response)=>response.json().
					then(function (attachments) {
						files = attachments
						for (const a of attachments) {
							// console.log(a)
							// console.log(a.Name) 
							// console.log(a.UUID) 
							// console.log(a.ProcessID) 
							a.src = "http://localhost:3123/api/v1/file/" + a.ID;
							console.log(a.src);
						}
						resolve(attachments)
					}
				)
			)
		})
	}
	
	function popModal(idx) {
		console.log("derp1")
		return function() {
			console.log("derp2")
			setTimeout(() => {
				console.log("derp3")
				console.log(files[idx].src)
				open(files, idx);
			}, 0);
		}
	}

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
	{#await attachmentsPromise}
	<!-- <div></div> -->
	{:then files}
		{#each files as file, i}
			<div class="fileTile">
				<!-- <a href={"http://localhost:3123/api/v1/file/" + file.ID} style="text-decoration: none !important;"> -->
					<ClickableTile on:click={popModal(i)}>
						<div style="display: flex; align-items: center; height: 2.5em;">
							<DocumentDownload24 />
							<p class="fileName">{file.Name}</p>
						</div> 
					</ClickableTile>
				<!-- </a> -->
				<div class="deleteHolder">
					<!-- DELETE -->
					<Button kind="danger" class="deleteButton" on:click={deleteFileWithID(file.ID)}>Delete</Button>
				</div>
			</div>
		{/each}
	{/await}
	<!-- <Button kind="secondary">Novo anexo</Button> -->
	<FileUploaderButton multiple disableLabelChanges={true} labelText="Adicionar arquivo" bind:ref={fileUploader} on:change={uploadFile}>Novo anexo</FileUploaderButton>
</div>