<script>
	import {TextArea, ButtonSet, Button, Dropdown} from 'carbon-components-svelte'
	import { createEventDispatcher } from 'svelte'
	let commentContent;

	const dispatch = createEventDispatcher();

	async function postNewComment() {
		console.log("[postNewComment]: Entering")
		try {     
			const response = await fetch('http://localhost:3123/api/v1/comment', {
					method: 'post',

					body: JSON.stringify({
						content: commentContent,
					}),

					headers: {
						'Content-type': 'application/json; charset=UTF-8'
					}
				}
			);

			response.text().then((text) => {
				console.log(text);
				console.log('[postNewComment]: Completed!');
			});

		} catch(err) {
			console.error(`Error: ${err}`);
			return;
		}

		clearText();
		dispatch('commentWasPosted');
	}

	function clearText() {
		commentContent = "";
	}

</script>

<style>
	.interventionButtonSet {
		margin-top: 0.5em;
		display: flex;
		justify-content: space-between;
	}
</style>

<TextArea
	rows={10}
	placeholder="Comentário por adicionar..."
	bind:value={commentContent}
/>

<div class="interventionButtonSet">
	<ButtonSet>
		<Button on:click={clearText} kind="secondary">Limpar</Button>
		<Button on:click={postNewComment}>Enviar</Button>
	</ButtonSet>

	<Dropdown
		hideLabel
		direction="top"
		titleText="Contact"
		selectedIndex={0}
		items={[{ id: '0', text: 'Novo comentário' }, { id: '1', text: 'Nova minuta' }]}
	/>
</div>