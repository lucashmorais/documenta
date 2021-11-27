<script>
	import { getAvailableCenters, getCurrentUserPermissions, setProcessStatus, hasPermission } from "./utils.js"
	import InfoLine from "./InfoLine.svelte"
	import CommentArea from "./CommentArea.svelte"
	import RoutingModal from "./RoutingModal.svelte"
	import Minutes from "./Minutes.svelte"
	import StatusBar from "./StatusBar.svelte"
	import AttachmentsArea from "./AttachmentsArea.svelte"
	import InterventionForm from "./InterventionForm.svelte"
	import SequenceTable from "./SequenceTable.svelte"
	import ProcessModal from './ProcessModal.svelte'
	import ModificationToolbar from './ModificationToolbar.svelte'
	import { getEndpointPrefix } from "./config-helper.js"
	import Cookie from "js-cookie";
	import Edit32 from "carbon-icons-svelte/lib/Edit32";
	import Stamp32 from "carbon-icons-svelte/lib/Stamp32";
	import {
	Button,
	Content,
	Toggle,
	Tile,
	TextInput
	} from "carbon-components-svelte";
	import { constants } from "./constants"

	let isSideNavOpen = false;
	var editModalIsOpen = false;
	var routingModalIsOpen = true;

	var coreRefreshComments;
	var coreRefreshMinutes;
	
	function refreshMinutes() {
		coreRefreshMinutes();
	}
	
	function refreshComments() {
		coreRefreshComments();
	}
	
	function logout() {
		Cookie.remove("documentaLoginToken")
		window.location.href = "/";
	}
	
	let refreshInfoLine;

	let urlParams = new URLSearchParams(window.location.search);
	let processID = urlParams.get("id")
	
	console.log("processID: ", processID)

	export function updateProcess() {
		processPromise = new Promise((resolve, reject) => {
			if (processID != null) {
				fetch(getEndpointPrefix() + "/api/v1/processes?processID=" + processID).
					then((response)=>response.json().
						then(function (wrappedProcess) {
							console.log("[updateProcess]: ", wrappedProcess[0])
							resolve(wrappedProcess[0])
						}
					)
				)
			} else {
				resolve (null);
			}
		})
		return processPromise;
	}
	
	let sequenceChangeEvent;
	function updateSequence(event) {
		console.log("[Document::updateSequence::event]: ", event)
		sequenceChangeEvent = event
	}
	
	let sequencePromise;
	let currentUserPromise;

	let processExaminationState = "analysis"
	
	async function updateProcessExaminationState(_promise) {
		if (_promise) {
			let seqObj = await _promise
			let sequence = seqObj.sequence

			console.log("[Document::setupSequenceCallbacks::__callback::sequence]: ", sequence)

			let kindID = sequence.UserSequenceKindID;

			if (kindID == constants.db.UserSequenceKinds.REVIEW) {
				processExaminationState = "analysis"
			} else if (kindID == constants.db.UserSequenceKinds.APPROVAL) {
				processExaminationState = "approval"
			}
		}		
	}
	$: updateProcessExaminationState(sequencePromise)
	
	function printCurrentUser(p) {
		if (p) {
			p.then((user) => {
				console.log("[Document::printCurrentUser:then]: ", user)			
			})
		}
	}
	
	$: printCurrentUser(currentUserPromise)

	let processPromise;
	updateProcess();

	let availableCenters = [];
	getAvailableCenters().then((centers) => {
		console.log("[getAvailableCentersCallback::centers]: ", centers)
		availableCenters = centers
	})
	
	let inPageModificationHappened = false;
	
	function currentUserHasModificationRights(seqP, userP) {
		if (!seqP || !userP){
			return null;
		}
		return new Promise(async function (resolve, reject) {
			let sequence = await seqP
			userP.then((user) => {
				console.log("[currentUserHasModificationRights::then::then::sequence]: ", sequence)
				console.log("[currentUserHasModificationRights::then::then::user]: ", user)
				let coreSeq = sequence.sequence
				let numCompletions = coreSeq.NumCompletions
				console.log("[currentUserHasModificationRights::then::then::numCompletions]: ", numCompletions)
				
				let result;
				if (numCompletions >= coreSeq.Users.length) {
					result = false
				} else {
					result = coreSeq.Users[numCompletions].ID == user.ID
				}
				console.log("[currentUserHasModificationRights::then::then::result]: ", result)
				resolve(result)
			})
		})				
	}
	
	$: modRightsPromise = currentUserHasModificationRights(sequencePromise, currentUserPromise)
	
	function checkIfPendingConfirmation(_procPromise, _seqPromise, _examinationState) {
		console.log("[checkIfPendingConfirmation]: {_procPromise, _seqPromise, _examinationState}: ",
			{_procPromise, _seqPromise, _examinationState})			

		if (!_procPromise || !_seqPromise || (_examinationState == "")) {
			console.log("[checkIfPendingConfirmation]: Some variable is missing, so we are returning now.")			
			return
		}
		
		_procPromise.then((proc) => {
			_seqPromise.then((seqWrapper) => {
				let seq = seqWrapper.sequence
				if (proc.ProcessStatusID != constants.db.ProcessStatuses.ACTIVE) {
					console.log("[checkIfPendingConfirmation]: Process is not active")			
				} else {
					console.log("[checkIfPendingConfirmation::seq]: ", seq)			
					if (seq.NumCompletions == seq.NumUsers) {
						switch(seq.UserSequenceKindID) {
							case constants.db.UserSequenceKinds.REVIEW:
								console.log("[checkIfPendingConfirmation]: Process is pending review confirmation")			
								setProcessStatus(processID, constants.db.ProcessStatuses.AWAITING_REVIEW_CONFIRMATION)
							break;
							case constants.db.UserSequenceKinds.APPROVAL:
								console.log("[checkIfPendingConfirmation]: Process is pending approval confirmation")			
								setProcessStatus(processID, constants.db.ProcessStatuses.AWAITING_APPROVAL_CONFIRMATION)
							break;
							default:
								console.log("[checkIfPendingConfirmation]: Sequence has bad UserSequenceKindID")			
							break;
						}
					} else {
						console.log("[checkIfPendingConfirmation]: The current sequence's number of " + 
						"confirmations does not match the number of users it describes, so this process " +
						"certainly is not pending review/approval confirmation.")			
					}
				}
			})
		})
	}
	$: checkIfPendingConfirmation(processPromise, sequencePromise, processExaminationState)
	
	let userPermissionsPromise = getCurrentUserPermissions()
	$: if (userPermissionsPromise) {
		userPermissionsPromise.then((value) => {
			console.log("[Document::userPermissions]: ", value)

		})
	}
	
	function handleResponse(__response) {
		if (__response && __response.status == 200) {
		} else {
			console.log(__response)
		}
	}
	
	function getNewTitle() {
		let newTitle = document.getElementById("editableTitle").innerText;
		console.log("[getNewTitle::newTitle]: ", newTitle);
		return newTitle;
	}
	
	function getNewSummary() {
		let newSummary = document.getElementById("editableSummary").innerText;
		console.log("[getNewSummary::newSummary]: ", newSummary);
		return newSummary;
	}
	
	async function submitProcess() {
		try {     
			let requestBody;
			let response;
			processPromise.then(async function(process) {
				requestBody = JSON.stringify({
					"title": getNewTitle(),
					"summary": getNewSummary(),
					"typeID": newProcessTypeID ? newProcessTypeID : process.ProcessTypeID,
					"centerID": newCenterID ? newCenterID : process.CenterID,
					// "typeID": available_types[formState.selectedType].id,
					// "centerID": available_centers[formState.selectedCenter].id,
					// "userSequenceUserIDs": (await sequencePromise).sequence.map((user) => user.ID)
				});

				
				console.log("[submitProcess:editing:requestBody]: ", requestBody);
				
				let asyncResponse = await fetch(getEndpointPrefix() + "/api/v1/process/" + process.ID + "?keep_user_sequence=true", {
					method: "put",
					
					body: requestBody,
					headers: {
						"Content-type": "application/json; charset=UTF-8"
					}
				}
				);
				handleResponse(asyncResponse);
				return;
			});
		} catch(err) {
			console.error(`Error: ${err}`);
			return;
		}
	}
	
	const updatedComments = new Set();
	
	// Function that iterates over the updatedComments set and submits the comments to the server using
	// the '/api/v1/comment' PUT endpoint.
	async function submitComments() {
		if (updatedComments.size == 0) {
			console.log("[submitComments::updatedComments.size]: ", updatedComments.size)
			return;
		}
		
		console.log("[submitComments::updatedComments]: ", updatedComments)
		for (let comment of updatedComments) {
			let requestBody = JSON.stringify(comment);
			console.log("[submitComments::requestBody]: ", requestBody);
			let asyncResponse = await fetch(getEndpointPrefix() + "/api/v1/comment", {
				method: "put",
				
				body: requestBody,
				headers: {
					"Content-type": "application/json; charset=UTF-8"
				}
			});
			handleResponse(asyncResponse);
		}
	}
	
	async function handleCommentModification(event) {
		inPageModificationHappened = true;

		console.log("[handleCommentModification::event]: ", event);
		updatedComments.add(event.detail.payload);
	}
	
	let newCenterID = 0;
	let newProcessTypeID = 0;
	
	function handleInfoLineChange() {
		inPageModificationHappened = true;
	}
</script>

<style>
	h1 {
		margin-top: 1em;
		margin-bottom: 1em;
		text-align: left;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	h2 {
		font-size: 22px;
		margin: 3em 0 1em 0;
		text-align: left;
		display: flex;
		align-items: center;
		justify-content: left;
	}

	.contents {
		margin: 0 15vw;
	}
	
	.records {

	}

	:global(.summary) {
		font-size: 15px;
		line-height: 150%;
	}

	:global(.selectable) {
	}
</style>

<ProcessModal 
	bind:open={editModalIsOpen}
	bind:processPromise={processPromise}
	sequencePromise={sequencePromise}
	on:backendModification={(event) => {
		updateProcess();
		updateSequence(event);
	}}
/>


<StatusBar bind:currentUserPromise/>

<Content>
	<div class="contents">
		<h1>
			{#await processPromise then process}
				<div id="editableTitle" on:input={() => {console.log("[titleInputCallback]: Title was modified"); inPageModificationHappened = true}} contenteditable="true">{process.Title}</div>
			{/await}
{#await Promise.all([processPromise, userPermissionsPromise]) then [proc, userPermissions]}
	{#if
		(
			proc.ProcessStatusID == constants.db.ProcessStatuses.AWAITING_REVIEW_CONFIRMATION
			&& hasPermission(userPermissions, constants.db.Permissions.CONFIRM_PROCESS_REVIEW)
		)
		||
		(
			proc.ProcessStatusID == constants.db.ProcessStatuses.AWAITING_APPROVAL_CONFIRMATION
			&& hasPermission(userPermissions, constants.db.Permissions.CONFIRM_PROCESS_APPROVAL)
		)
		&&
		!editModalIsOpen
	}
		<div style="margin-left: 0.5em">
			<!-- TODO: Add a dynamic description to this button. -->
			<Button kind="secondary" iconDescription="Aprovar ou estender processo" icon={Stamp32} on:click={() => routingModalIsOpen = true}/>
		</div>
		<RoutingModal bind:open={routingModalIsOpen} processPromise={processPromise} processExaminationState={processExaminationState}
			on:close={() => console.log("[RoutingModal::__closeCallback:routingModalIsOpen]: ", routingModalIsOpen)}
			on:processChange={updateProcess}
			on:sequenceChange={updateSequence}
		/>
	{/if}
{/await}
		</h1>
		<InfoLine
			bind:current_center_id={newCenterID}
			bind:current_type_id={newProcessTypeID}
			bind:refreshAndClear={refreshInfoLine}
			on:change={handleInfoLineChange}
			processPromise={processPromise}
		/>
		<h2>Resumo</h2>
		<Tile class="summary">
			{#await processPromise then process}
				{#if process != null}
					<div id="editableSummary" contenteditable="true" on:input={() => {inPageModificationHappened = true; console.log("[summaryInputCallback::process.Summary]: ", process.Summary)} }>
						{process.Summary}
					</div>
				{:else}
					Lorem ipsum dolor sit amet
				{/if}
			{/await}
		</Tile>

		<h2>
			<div>
				Sequência de análise
			</div>
			{#await processPromise then process}
				{#if
					process.ProcessStatusID == constants.db.ProcessStatuses.ACTIVE ||
					process.ProcessStatusID == constants.db.ProcessStatuses.DRAFT
				}
					<div style="margin-left: 0.5em">
						<Button kind="tertiary" iconDescription="Editar processo" icon={Edit32} on:click={() => {editModalIsOpen = true}}/>
					</div>
				{/if}
			{/await}
		</h2>
		<SequenceTable bind:modRightsPromise bind:sequencePromise bind:sequenceChangeEvent processID={processID}/>

		<h2>Anexos</h2>
		<AttachmentsArea bind:modRightsPromise processID={processID}/>

		<h2>Minutas</h2>
		<Minutes bind:modRightsPromise={modRightsPromise} processID={processID} bind:updateMinutes={coreRefreshMinutes} on:minuteWasPosted={refreshMinutes} availableCenters={availableCenters}/>

		{#await modRightsPromise then modRights}
			{#if modRights}
				<h2>Nova intervenção</h2>
				<InterventionForm processID={processID} on:commentWasPosted={refreshComments} on:minuteWasPosted={refreshMinutes} availableCenters={availableCenters}/>
			{/if}
		{/await}

		<h2>Intervenções</h2>
		<CommentArea on:commentModification={handleCommentModification} currentUserPromise={currentUserPromise} processID={processID} bind:updateComments={coreRefreshComments}/>
	</div>
</Content>

{#if inPageModificationHappened}
	<ModificationToolbar
		on:commit={async function() {
			await submitProcess();
			await submitComments();
			updatedComments.clear();
			//await updateProcess();
			await updateProcess();
			// refreshInfoLine()
			inPageModificationHappened = false;
		}}
		on:cancel={async function() {
			await updateProcess();
			refreshComments();
			// refreshInfoLine()
			inPageModificationHappened = false;
		}}
	/>
{/if}