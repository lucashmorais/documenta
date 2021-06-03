<script>
	import { ClickableTile } from "carbon-components-svelte";
	import DocumentDownload24 from "carbon-icons-svelte/lib/DocumentDownload24";
	import { Button } from "carbon-components-svelte";
	import { FileUploaderButton } from "carbon-components-svelte";

	let files = [{Name: "Nota relevante"}, {Name: "Estudo similar"}, {Name: "Consulta relativa"}]
	let fileUploader

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
			coreSubmit(file)
		}
	}

	export function updateAttachments() {
		let attachmentsPromise = new Promise((resolve, reject) => {
			fetch("http://localhost:3123/api/v1/files").
				then((response)=>response.json().
					then(function (attachments) {
						for (const a of attachments) {
							console.log(a)
							console.log(a.Name) 
							console.log(a.UUID) 
							console.log(a.ProcessID) 
						}
						files = attachments
						resolve(attachments)
					}
				)
			)
		})
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

	.fileName {
		padding-left: 1em;
		font-size: 14px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: normal;
	}
</style>


<div class="fileGrid">
	{#each files as file}
		<div>
			<a href={"http://localhost:3123/api/v1/file/" + file.ID} style="text-decoration: none !important;">
				<ClickableTile>
					<div style="display: flex; align-items: center; height: 2.5em;">
						<DocumentDownload24 />
						<p class="fileName">{file.Name}</p>
					</div> 
				</ClickableTile>
			</a>
		</div>
	{/each}
	<!-- <Button kind="secondary">Novo anexo</Button> -->
	<FileUploaderButton multiple disableLabelChanges={true} labelText="Adicionar arquivo" bind:ref={fileUploader} on:change={uploadFile}>Novo anexo</FileUploaderButton>
</div>