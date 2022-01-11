<script>
  import {
	Tile,
	Accordion,
	AccordionItem,
  } from "carbon-components-svelte";
	import ActionsBlock from "./ActionsBlock.svelte";
	import InterventionForm from "./InterventionForm.svelte";
	import SimpleConfirmationModal from "./SimpleConfirmationModal.svelte";
	import { onMount } from "svelte";
	import { decodeDate, decodeTime } from "./utils.js"
	import { jsPDF } from "jspdf/dist/jspdf.es";
	import { getEndpointPrefix } from "./config-helper.js"
	
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
			fetch(getEndpointPrefix() + "/api/v1/minute/" + minuteID, {
					method: "delete",
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
			// console.log("[Minutes::updateMinutes]: Just entering minutesPromise"s inner function")
			if (processID == "" || processID == "0") {
				resolve([])
			}
			fetch(getEndpointPrefix() + "/api/v1/minutes?processID=" + processID).
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
			fetch(getEndpointPrefix() + "/api/v1/minute_versions?minuteID=" + minuteID).
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
	
	function coreDownloadMinute(textContents) {
	        var pdf = new jsPDF({
			format: "a6",
			orientation: "landscape",
		});
		// console.log("[generateExamplePDF::getFontSize()]: ", pdf.getFontSize())
		// console.log("[generateExamplePDF::getFontList()]: ", pdf.getFontList())
		pdf.setFont("Times", "Roman", "normal");
		pdf.setFontSize(13);
		pdf.text("sjc 15/20", 130, 10, {align: "right"});
		pdf.text(textContents, 18, 35, {maxWidth: 112});
		pdf.text(130, 95, "São Paulo, 26-8-20", {align: "right"});
		pdf.save("hello_world.pdf");	
	}

	function generateExamplePDF() {
		coreDownloadMinute("Informamos que João XX, professor de ZZ, estará em Lisboa e em Madrid de 10 a 15 de maio. " +
		"Agradeceríamos saber se poderia residir em alguma residência de professores nestas cidades.");
	}
	// generateExamplePDF();
	
	function downloadMinute(minute) {
		getMinuteVersions(minute.ID).then((versions) => {
			let numVersions = versions.length;
			if (numVersions > 0) {
				coreDownloadMinute(versions[numVersions - 1].Content)
			} else {
				coreDownloadMinute(minute.Content)
			}
		})
	}
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

<SimpleConfirmationModal bind:open={deleteModalOpen} on:actionConfirmed={deleteMinute(deleteeIdentifier)} selectedItems={[0]} singularString="minuta" />

<div class="interventions">
	{#await minutesPromise then minutes}
		{#each minutes as minute}
			<Tile>
				<ActionsBlock
					bind:modRightsPromise
					city={minute.Center.Name}
					editAction={setupHidingCallback}
					minuteID={minute.ID}
					minuteOutboundProtocol={`${minute.ProtocolPrefix}${minute.OutboundProtocolNumber}`}
					on:deletionRequested={() => {deleteeIdentifier = minute.ID; deleteModalOpen = true}}
					on:downloadRequested={() => {downloadMinute(minute)}}
					on:protocolChange={updateMinutes}
				/>
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