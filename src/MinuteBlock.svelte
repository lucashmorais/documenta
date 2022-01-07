<script>
	import "carbon-components-svelte/css/all.css";
	import {
		getNameFromUser
	} from "./utils.js"
	import {
		Tile,
		ClickableTile,
		Button
	} from "carbon-components-svelte";
	import DocumentAdd16 from "carbon-icons-svelte/lib/DocumentAdd16";
	import Add32 from "carbon-icons-svelte/lib/Add32";
	import WatsonHealth3DCurveAutoColon16 from "carbon-icons-svelte/lib/WatsonHealth3DCurveAutoColon16";
	
	export let createMinuteAssignmentModalOpeningHandler;
	export let createProcessCreationModalOpeningHandler;
	export let creationModalHandler;
	export let minutesPromise;
	export let disableEditing = false;
	export let disableTitleWrapper = false;
	export let enableProcessView = false;
	export let enableMinutePlaceholder = true;
	export let blockTitle = "Minutas";
</script>

<style>
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
	
	.content2 {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1fr;
		grid-template-rows: 1fr;
		gap: 2em 2em;
		grid-auto-flow: row;
		grid-template-areas: ". . . .";
		margin-top: 2em;
	}
								
	.actionSet {
		margin-top: 1em;
	}
</style>

{#if !disableTitleWrapper}
<h2>{blockTitle}</h2>
{/if}
	<div class="content2">
		{#if !disableEditing}
			<div class="element">
				<ClickableTile
					style="height: 100%; display: flex; align-items: center; justify-content: center; flex-direction: column"
					on:click={creationModalHandler}
				>
					<Add32 />
					<p>Adicionar minuta</p>
				</ClickableTile>
			</div>
		{/if}
		{#await minutesPromise then minutes}
			{#if enableMinutePlaceholder && minutes.length == 0}
				<div class="element">
					<Tile>
						Nenhuma minuta foi encontrada
					</Tile>
				</div>
			{/if}
			{#each minutes as minute}
				<div class="element">
					<Tile>
						<div class="">
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
								Título
							</h3>
							<p>
								{minute.Description ? minute.Description : "Sem título"}
							</p>
								
							{#if true}
								<h3>
									ID
								</h3>
								<p>
									{minute.ID}
								</p>
								
							{/if}
							<div class="actionSet">
								{#if enableProcessView}
									<Button kind="tertiary" iconDescription="Visualizar processo" on:click={() => window.open('document.html?id=' + minute.ProcessID, '_blank').focus()}>
										Visualizar processo
									</Button>
								{:else if !disableEditing}
									<Button kind="secondary" icon={WatsonHealth3DCurveAutoColon16} iconDescription="Atribuir" on:click={createMinuteAssignmentModalOpeningHandler(minute.ID)} />
									<Button kind="tertiary" icon={DocumentAdd16} iconDescription="Gerar processo" on:click={createProcessCreationModalOpeningHandler(minute.ID)} />
								{/if}
							</div>
							<Button style="margin-top: 0.5em;" kind="tertiary" iconDescription="Visualizar anexo" on:click={() => window.open('api/v1/file/' + minute.UserFileID, '_blank').focus()}>
								Visualizar anexo
							</Button>
						</div>
					</Tile>
				</div>
			{/each}
		{/await}
	</div>
