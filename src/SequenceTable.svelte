<script>
	import { DataTable } from "carbon-components-svelte";
	import { decodeDate, decodeTime } from './utils.js'
	import Pen16 from "carbon-icons-svelte/lib/Pen16";
	
	export let processID;

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
										name: user.FirstName + " " + user.LastName,
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

	let sequencePromise = updateUserSequence();
</script>

<style>
	.modificationTokenHolder{
		font-weight: 900!important;		
	}
</style>
      
{#await sequencePromise}-
{:then value}
	<DataTable
		useStaticWidth
		headers={[{ key: 'name', value: 'Nome' }, { key: 'unix_completion_time', value: 'Conclusão' }]}
		rows={value.user_rows}
	>
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