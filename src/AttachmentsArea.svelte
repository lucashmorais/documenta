<script>
	import { ClickableTile } from "carbon-components-svelte";
	import DocumentDownload24 from "carbon-icons-svelte/lib/DocumentDownload24";
	import { Button } from "carbon-components-svelte";
	import { FileUploaderButton } from "carbon-components-svelte";

	let files = [{name: "Nota relevante"}, {name: "Estudo similar"}, {name: "Consulta relativa"}]
	let fileUploader

	function uploadFile() {
		console.log(fileUploader)
		console.log(fileUploader.files)
		// fileUploader.files[0].arrayBuffer().then((buffer) => console.log(JSON.stringify(buffer)))

		let file = fileUploader.files[0]

		var success = async function ( content ) {
			try {     
					console.log("Starting to send data to server")
					const response = await fetch('http://localhost:3123/api/v1/files', {
							method: 'post',

							body: JSON.stringify({
								fileContents: content,
								name: file.name
							}),

							headers: {
								'Content-type': 'application/json'
							}
					});
					console.log(response)
			} catch(err) {
				console.error(`Error: ${err}`);
				return;
			}
		}

		var fileReader = new FileReader( );
		fileReader.onload = function ( evt ) { success( evt.target.result ) };
		// fileReader.readAsBinaryString( file );
		fileReader.readAsBinaryString( file );
	}
</script>

<style>
	.fileGrid {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1fr;
		grid-template-rows: 1fr;
		gap: 1em 1em;
		grid-template-areas:
			". . . .";
	}
</style>


<div class="fileGrid">
	{#each files as file}
		<div>
			<ClickableTile>
				<div style="display: flex; align-items: center; height: 2.5em;">
					<DocumentDownload24 />
					<p style="padding-left: 1em; word-wrap: break-word; font-size: 14px">{file.name}</p>
				</div> 
			</ClickableTile>
		</div>
	{/each}
	<!-- <Button kind="secondary">Novo anexo</Button> -->
	<FileUploaderButton disableLabelChanges={true} labelText="Adicionar arquivo" bind:ref={fileUploader} on:change={uploadFile}>Novo anexo</FileUploaderButton>
</div>