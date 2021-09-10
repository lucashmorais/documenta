<script>
	import { DataTable } from "carbon-components-svelte";
	
	export let processID;

	export function updateUserSequence() {
		let promise = new Promise((resolve, reject) => {
			if (processID != null) {
				fetch("http://localhost:3123/api/v1/user_sequences?processID=" + processID).
					then((response)=>response.json().
						then(function (sequence) {
							let lastIdx = sequence.length - 1
							let seq = sequence[lastIdx]
							let rows = seq.Users.map(function (user) { return {id: user.ID, name: user.FirstName + " " + user.LastName}})
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
      
{#await sequencePromise}
{:then value}
	<DataTable
		useStaticWidth
		headers={[{ key: 'name', value: 'Nome' }]}
		rows={value.user_rows}
	/>
{/await}