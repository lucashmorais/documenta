<script>
  import {
	Tile,
	Accordion,
	AccordionItem,
  } from "carbon-components-svelte";
	import ActionsBlock from './ActionsBlock.svelte';
	import InterventionForm from './InterventionForm.svelte';
	import SimpleDeleteModal from './SimpleDeleteModal.svelte';
	import { onMount } from 'svelte';
	import { decodeDate, decodeTime } from './utils.js'
	
	export let availableCenters = [];
	export let modRightsPromise;
	
	var deleteModalOpen = false;
	var deleteeIdentifier;

	function pushNodeReferenceToArray(node) {
		node.is_hidden = true;
		newVersionBlocks.push(node)
		console.log(newVersionBlocks.length)
	}

	function hideNodeAndPushReference(node) {
		node.style.display = "none";
		pushNodeReferenceToArray(node);
	}

	function removeOuterPadding(node) {
		node.parentNode.style["padding-right"] = "0!important";
		node.style["background-color"] = "red!important";
		node.parentNode.style["background-color"] = "blue!important";
	}

	function setupHidingCallback(node) {
		// console.log("We are now going to add an event listener to the following node:");
		console.log(node);
		// console.log("We would like this node to be able to toggle the display of the following text input node:");
		node.callbackID = numTextHidingCallbacksAdded;
		console.log(newVersionBlocks[node.callbackID]);
		numTextHidingCallbacksAdded++;
		node.addEventListener("click", () => {
			if (newVersionBlocks[node.callbackID].is_hidden) {
				newVersionBlocks[node.callbackID].style.display = "";
				newVersionBlocks[node.callbackID].is_hidden = false;
			} else {
				newVersionBlocks[node.callbackID].style.display = "none";
				newVersionBlocks[node.callbackID].is_hidden = true;
			}
		});
	}

	onMount(() => {
		document.querySelectorAll(".removeOuterPadding").forEach((node) => node.parentNode.style["padding-right"] = "0");
	})

	export let processID = 0;
	
	function deleteMinute(minuteID) {
		deleteModalOpen = false;
		console.log("[deleteMinute]: Entering")
		try {     
			fetch('http://localhost:3123/api/v1/minute/' + minuteID, {
					method: 'delete',
				}
			).then((response) => updateMinutes())
		} catch(err) {
			console.error(`Error: ${err}`);
			return;
		}
	}
	
	let minutesPromise;
	export function updateMinutes() {
		minutesPromise = new Promise((resolve, reject) => {
			// console.log("[Minutes::updateMinutes]: Just entering minutesPromise's inner function")
			if (processID == "" || processID == "0") {
				resolve([])
			}
			fetch("http://localhost:3123/api/v1/minutes?processID=" + processID).
				then((response)=>response.json().
					then(function (minutes) {
						// console.log("[Minutes::updateMinutes::minutes]: ", minutes)
						for (let i = 0; i < minutes.length; i++) {
							let a = minutes[i]
							console.log(a)
						}
						resolve(minutes)
					}
				)
			)
		})
	}

	function getMinuteVersions(minuteID) {
		return new Promise((resolve, reject) => {
			console.log("[Minutes::getMinuteVersions]: Just entering minutesPromise's inner function")
			if (processID == "" || processID == "0") {
				resolve([])
			}
			fetch("http://localhost:3123/api/v1/minute_versions?minuteID=" + minuteID).
				then((response)=>response.json().
					then(function (versions) {
						console.log("[Minutes::getMinuteVersions::versions]: ", versions)
						for (let i = 0; i < versions.length; i++) {
							let a = versions[i]
							console.log(a)
						}
						resolve(versions)
					}
				)
			)
		})
	}

	let newVersionBlocks = [];
	let numTextHidingCallbacksAdded = 0;

	updateMinutes();
</script>

<style>
	.interventions {
		display: grid;
		grid-template-columns: minmax(150px, 1fr) 4fr;
		grid-template-rows: 1fr;
		gap: 0px 0px;
		grid-template-areas:
			". .";
		 gap: 0.5em 0.5em;
	}
</style>

<SimpleDeleteModal bind:open={deleteModalOpen} on:deletionConfirmed={deleteMinute(deleteeIdentifier)} selectedItems={[0]} singularString="minuta" />

<div class="interventions">
	{#await minutesPromise then minutes}
		{#each minutes as minute}
			<Tile>
				<ActionsBlock bind:modRightsPromise city={minute.Center.Name} editAction={setupHidingCallback} on:deletionRequested={() => {deleteeIdentifier = minute.ID; deleteModalOpen = true}}/>
			</Tile>
			<Tile>
				<Accordion>
					{#await modRightsPromise then canModify}
						{#if canModify}
							<div use:hideNodeAndPushReference>
								<AccordionItem open title="Nova versão">
									<div class="removeOuterPadding">
										<InterventionForm processID={processID} minuteID={minute.ID} minuteOnly={true} availableCenters={availableCenters} on:minuteWasPosted />
									</div>
								</AccordionItem>
							</div>
						{/if}
					{/await}
					{#await getMinuteVersions(minute.ID) then versions}
						{#each versions.reverse() as version}
							<AccordionItem open={version.ID == versions[0].ID} title="{decodeDate(version.UnixCreatedAt)}, {decodeTime(version.UnixCreatedAt)}, {version.User.FirstName} {version.User.LastName}: {version.Description}">
								<p>
								{version.Content}
								</p>
							</AccordionItem>
						{/each}
						<AccordionItem open={versions.length == 0} title="{decodeDate(minute.UnixCreatedAt)}, {decodeTime(minute.UnixCreatedAt)}, {minute.User.FirstName} {minute.User.LastName}: {minute.Description}">
							<p>
							{minute.Content}
							</p>
						</AccordionItem>
					{/await}
				</Accordion>
			</Tile>
		{/each}
	{/await}
</div>