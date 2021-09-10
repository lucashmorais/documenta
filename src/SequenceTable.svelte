<script>
	import { DataTable } from "carbon-components-svelte";
	import Pen16 from "carbon-icons-svelte/lib/Pen16";
	
	export let processID;

	export function updateUserSequence() {
		let promise = new Promise((resolve, reject) => {
			if (processID != null) {
				fetch("http://localhost:3123/api/v1/user_sequences?processID=" + processID).
					then((response)=>response.json().
						then(function (sequence) {
							let lastIdx = sequence.length - 1
							let seq = sequence[lastIdx]
							let rows = seq.Users.map(function (user, index) { return {id: index, name: user.FirstName + " " + user.LastName}})
							let response = {sequence: seq, user_rows: rows}
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
		headers={[{ key: 'name', value: 'Nome' }]}
		rows={value.user_rows}
	>
	<span slot="cell" let:row let:cell>
		{#if row.id === value.sequence.NumCompletions}
			<div class="modificationTokenHolder">
				{cell.value}
				<Pen16/>
			</div>
		{:else}
			{cell.value}
		{/if}
      </span>
	</DataTable>
{/await}