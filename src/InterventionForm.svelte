<script>
	import {TextArea, ButtonSet, Button, Dropdown} from 'carbon-components-svelte'
	import { createEventDispatcher } from 'svelte'
	let commentContent;
	let selectedIndex = 0;
	export let processID = 0;

	const dispatch = createEventDispatcher();
	
	function handlePostClick() {
		switch(selectedIndex) {
			case 0:	postNewComment()
			break;
			case 1:	postNewMinute()
			break;
		}
	}
	
	// Function that calls the /api/v1/minute endpoint to create a new Minute
	async function postNewMinute() {
		console.log("[postNewMinute]: Entering")
		try {     
			// Here we use the Number function to prevent `processID` from being converted to a string
			let requestBody = JSON.stringify({
						"Content": commentContent,
						"Description": "Texto inicial",
						"ProcessID": Number(processID)
					});

			const response = await fetch('http://localhost:3123/api/v1/minute', {
					method: 'post',

					body: requestBody,

					headers: {
						'Content-type': 'application/json; charset=UTF-8'
					}
				}
			);
			
			console.log("[postNewMinute::requestBody]: ", requestBody);

			response.text().then((text) => {
				console.log(text);
				console.log('[postNewMinute]: Completed!');
			});

		} catch(err) {
			console.error(`Error: ${err}`);
			return;
		}

		clearText();
		dispatch('minuteWasPosted');
	}

	// Function that calls the /api/v1/comment endpoint to create a new Comment
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
		<!-- <Button on:click={postNewComment}>Enviar</Button> -->
		<Button on:click={handlePostClick}>Enviar</Button>
	</ButtonSet>

	<Dropdown
		hideLabel
		direction="top"
		titleText="Contact"
		bind:selectedIndex={selectedIndex}
		items={[{ id: '0', text: 'Novo comentário' }, { id: '1', text: 'Nova minuta' }]}
	/>
</div>