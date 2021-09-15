<script>
	import { DataTable } from "carbon-components-svelte";
	import { getNameFromUser, decodeDate, decodeTime, countCompletion } from "./utils.js";
	import Pen16 from "carbon-icons-svelte/lib/Pen16";
	import { Button } from "carbon-components-svelte";
	
	export let processID;
	export let modRightsPromise;

	export function updateUserSequence() {
		let promise = new Promise((resolve, reject) => {
			if (processID != null) {
				fetch("http://localhost:3123/api/v1/user_sequence_with_timestamps?processID=" + processID).
					then((response)=>response.json().
						then(function (raw_response) {
							let seq = raw_response.userSequence
							let times = raw_response.tokenPassingTimestamps
							let numCompletions = times.length
							let rows = seq.Users.map(
								function (user, index) { 
									return {
										id: index,
										name: getNameFromUser(user),
										unix_completion_time: index < numCompletions ? times[index].UnixTimestamp : null 
									}
								}
							)
							let response = {sequence: seq, timestamps: times, user_rows: rows}
							console.log("[updateUserSequence::response]: ", response)
							resolve(response)
						}
					)
				)
			} else {
				resolve (null);
			}
		})
		return promise
	}

	export let sequencePromise = updateUserSequence();
</script>

<style>
	.modificationTokenHolder {
		font-weight: 900!important;		
	}
	.customToolbar {
		display: flex;
		align-items: center;
		justify-content: space-around;
	}
</style>
      
{#await sequencePromise then value}
	<DataTable
		headers={[{ key: 'name', value: 'Nome' }, { key: 'unix_completion_time', value: 'Conclusão' }]}
		rows={value.user_rows}
	>
	<!-- <strong slot="title">Sequência de análise</strong> -->
	<div slot="description" class="customToolbar">
		<span style="font-size: 1rem">
			Apenas o usuário ressaltado pode realizar modificações.
		</span>
		{#await modRightsPromise then canModify}
			<Button disabled={!canModify} iconDescription="Editar" on:click={() => countCompletion(processID)}>Finalizar turno</Button>
		{/await}
	</div>
	<span slot="cell" let:row let:cell>
		{#if !cell.value}
			<div></div>
		{:else if row.id === value.sequence.NumCompletions}
			<div class="modificationTokenHolder">
				{#if cell.key === "unix_completion_time"}
					{decodeDate(cell.value)}
					<br>
					{decodeTime(cell.value)}
				{:else if cell.key === "name"}
					{cell.value}
					<Pen16/>
				{/if}
			</div>
		{:else}
			{#if cell.key === "unix_completion_time"}
				{decodeDate(cell.value)}
				<br>
				{decodeTime(cell.value)}
			{:else}
				{cell.value}
			{/if}
		{/if}
	</span>
	</DataTable>
{/await}