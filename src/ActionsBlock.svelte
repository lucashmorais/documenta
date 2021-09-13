<script>
	import { getAttributes, toSVG } from '@carbon/icon-helpers';
	import { Button, Tag } from "carbon-components-svelte";
	import { onMount, createEventDispatcher } from 'svelte'

	import editIcon from '@carbon/icons/es/edit/16';
	import deleteIcon from '@carbon/icons/es/delete/16';
	import Edit16 from "carbon-icons-svelte/lib/Edit16";
	import Delete16 from "carbon-icons-svelte/lib/Delete16";
import RefKindTag from './RefKindTag.svelte';

	export let city;
	export let referenceNumber = 111;
	export let protocolNumber = 222;
	export let isOutReference = false;
	export let editAction = () => {};
	export let modRightsPromise;

	const dispatch = createEventDispatcher();

	let editIconNode;
	let deleteIconNode;

	const editIconSVG = toSVG({
		...editIcon,
		attrs: getAttributes(editIcon.attrs),
	});
	const deleteIconSVG = toSVG({
		...deleteIcon,
		attrs: getAttributes(deleteIcon.attrs),
	});

	function mountIcons() {
		console.log("Appending following nodes:")
		console.log(editIcon)
		console.log(deleteIcon)
		editIconNode.prepend(editIconSVG)
		deleteIconNode.prepend(deleteIconSVG)
	}
	
	function requestDeletion() {
		dispatch("deletionRequested")
	}

	onMount(async () => {
		// mountIcons();
	});
</script>

<style>
	.grid-container {
	display: grid;
	grid-template-columns: 1fr 1fr;
	grid-template-rows: 1fr;
	gap: 1em 1em;
	grid-template-areas:
		". .";
	}

	.cityName {
		margin-bottom: 5em;
		width: 100%;
		flex: 1;
		flex-grow: 1;
		line-height: 2em;
	}

	.flexPlacer {
		display: flex;
		flex-direction: column;
		height: 100%;
	}

	p {
		font-size: 12px;
	}

	.actionsContainer {
		display: flex;
	}

	.actionsContainer * {
		margin-right: 1em;
	}
</style>

<div class="flexPlacer">
	<div class="cityName">
		<div>{city}</div>
		<div>[Prot] {protocolNumber}</div>
		<div>[Ref] {referenceNumber} <RefKindTag /></div>
	</div>
	{#await modRightsPromise then canModify}
		{#if canModify}
			<div class="actionsContainer">
				<div use:editAction>
					<Button iconDescription="Editar" icon={Edit16} kind="tertiary"></Button>
				</div>
				<div>
					<Button iconDescription="Eliminar" icon={Delete16} kind="danger-tertiary" on:click={requestDeletion}></Button>
				</div>
			</div>
		{/if}
	{/await}
</div>