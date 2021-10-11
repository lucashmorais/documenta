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
	import { getEndpointPrefix } from "./config-helper.js"
	import Cookie from "js-cookie";
	import {
	Content,
	Tile,
	Dropdown
	} from "carbon-components-svelte";
	import { constants } from "./constants"

	let isSideNavOpen = false;

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
</style>

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
	}
		<RoutingModal open={true} processPromise={processPromise} processExaminationState={processExaminationState}
			on:processChange={updateProcess}
			on:sequenceChange={updateSequence}
		/>
	{/if}
{/await}

<StatusBar bind:currentUserPromise/>

<Content>
	<div class="contents">
		<h1>
		{#await processPromise then process}
			<div>{process.Title}</div>
		{/await}
		<div style="margin-left: 0.5em">
			<Dropdown
				selectedIndex={0}
				items={[{ id: "0", text: "Rascunho" }, { id: "1", text: "Ativo" }, { id: "2", text: "Concluído" }]}
			/>
		</div>
		</h1>
		<InfoLine processPromise={processPromise}/>
		<h2>Resumo</h2>
		<Tile class="summary">
			{#await processPromise then process}
				{#if process != null}
					{process.Summary}
				{:else}
					Lorem ipsum dolor sit amet
				{/if}
			{/await}
		</Tile>

		<h2>Sequência de análise</h2>
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
		<CommentArea processID={processID} bind:updateComments={coreRefreshComments}/>
	</div>
</Content>