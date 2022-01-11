<script>
	import { getAttributes, toSVG } from '@carbon/icon-helpers';
	import { Button, Tag } from "carbon-components-svelte";
	import { Modal } from "carbon-components-svelte";	
	import { onMount, createEventDispatcher } from 'svelte'

	import editIcon from '@carbon/icons/es/edit/16';
	import deleteIcon from '@carbon/icons/es/delete/16';
	import Edit16 from "carbon-icons-svelte/lib/Edit16";
	import Delete16 from "carbon-icons-svelte/lib/Delete16";
	import Download16 from "carbon-icons-svelte/lib/Download16";

	// import RefKindTag from './RefKindTag.svelte';
	import ProtocolSelectionModal from './ProtocolSelectionModal.svelte';

	export let city;
	export let referenceNumber = 111;
	export let isOutReference = false;
	export let editAction = () => {};
	export let modRightsPromise;
	export let minuteID;
	export let minuteOutboundProtocol;

	const dispatch = createEventDispatcher();

	let editIconNode;
	let deleteIconNode;
	let protocolSelectionModalIsOpen = false;

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
	
	function requestDownload() {
		dispatch("downloadRequested")
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
		margin-top: 1em;
		margin-right: 1em;
	}
	
	.protocolNumber:hover {
		cursor: pointer;
		color: gray
	}
</style>

<ProtocolSelectionModal
	minuteID={minuteID}
	bind:open={protocolSelectionModalIsOpen}
	on:protocolChange
/>

<div class="flexPlacer">
	<div class="cityName">
		<div>{city}</div>
		[Prot] <span class=protocolNumber on:click={() => protocolSelectionModalIsOpen = true}>{minuteOutboundProtocol}/22</span>
		<!-- <div>[Ref] {referenceNumber} <RefKindTag /></div> -->
	</div>
	{#await modRightsPromise then canModify}
		{#if canModify}
			<div class="topButton" use:editAction>
				<Button iconDescription="Editar" icon={Edit16} kind="tertiary"></Button>
			</div>
			<div class="actionsContainer">
				<div>
					<Button iconDescription="Eliminar" icon={Delete16} kind="danger-tertiary" on:click={requestDeletion}></Button>
				</div>
				<div class="topButton">
					<Button iconDescription="Download" icon={Download16} kind="ghost" isSelected={true} on:click={requestDownload}></Button>
				</div>
			</div>
		{/if}
	{/await}
</div>