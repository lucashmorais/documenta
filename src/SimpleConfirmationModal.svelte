<script>
	import {
		Modal
	} from 'carbon-components-svelte'
	
	import { createEventDispatcher } from 'svelte'
	
	const dispatcher = createEventDispatcher();

	export let open = false;
	export let selectedItems = [];
	export let singularString = "";
	export let pluralString = "";
	export let customMessage = "";
</script>

<Modal
	bind:open={open}
	modalHeading="Atenção!"
	primaryButtonText="Confirmar"
	secondaryButtonText="Cancelar"
	on:click:button--secondary={() => (open = false)}
	on:open
	on:close
	on:submit={() => dispatcher('actionConfirmed')}
>
	{#if customMessage != ""} 
		{customMessage}
	{:else}
		{#if selectedItems.length > 1}
			Você tem certeza que gostaria de remover {selectedItems.length} {pluralString}?
		{:else}
			Você tem certeza que gostaria de remover {selectedItems.length} {singularString}?
		{/if}
	{/if}
</Modal>