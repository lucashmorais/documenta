<script>
	import {
		Modal,
		Tag,
	} from "carbon-components-svelte"
	import { translatable } from './translatable.js'
	import { ContentSwitcher, Switch } from "carbon-components-svelte";
	
	import { createEventDispatcher } from 'svelte';
	import { getAvailableUsers } from './utils.js';
	import { constants } from './constants'
	
	const dispatch = createEventDispatcher();
	
	function signalBackendModification() {
		dispatch('backendModification');
	}
	
	export let open = false;
	
	let formState = {
		"title": "",
		"summary": "",
		"selectedType": 0,
		"selectedCenter": 0,
	}
	
	function setupSequenceCallbacks(_promise) {
		if (_promise) {
			_promise.then((seqObj) => {
				let sequence = seqObj.sequence
				console.log("[RoutingModal::setupSequenceCallbacks::__callback::sequence]: ", sequence)
				let kindID = sequence.UserSequenceKindID;

				if (kindID == constants.db.UserSequenceKinds.REVIEW) {
					processExaminationState = "analysis"
				} else if (kindID == constants.db.UserSequenceKinds.APPROVAL) {
					processExaminationState = "approval"
				}
			})
		}		
	}
	
	export let processPromise;
	export let sequencePromise;
	
	$: setupSequenceCallbacks(sequencePromise)
	
	export let processExaminationState = "analysis"
	$: if (processExaminationState == 'analysis') { clearForm() }

	console.log("[RoutingModal::_initialization::translatable]: ", translatable)
	
	function getStringSet(_examinationState) {
		let set = translatable.ptBR.routingModal[_examinationState]
		console.log("[RoutingModal::getStringSet::_examinationState]: ", _examinationState)
		console.log("[RoutingModal::getStringSet::set]: ", set)
		return set
	}

	$: stringSet = getStringSet(processExaminationState)
	$: modalHeading = stringSet.modalHeading
	$: ctxSwitchOpt1 = stringSet.ctxSwitchOption1
	$: ctxSwitchOpt2 = stringSet.ctxSwitchOption2
	
	let ctxSelectedIdx = 0;
	
	let purposePromise = null;
	function updatePurposePromise(ignored) {
		purposePromise = new Promise((resolve, reject) => {
			resolve(processExaminationState)			
		})
	}
	$: updatePurposePromise(processExaminationState)
	
	export let userInfo = null;
	
	$: console.log("[ProcessModal::userInfo]: ", userInfo)
	
	let available_types;
	let available_centers = [{"id": 0, "text": "Brasília"}, {"id": 1, "text": "João Cachoeira"}];
	
	function clearForm() {
		const keys = Object.keys(formState)
		for (const key of keys) {
			if (key == "selectedType" || key == "selectedCenter") {
				formState[key] = 0
			}
			else {
				formState[key] = ""
			}
		}
		formState = formState
	}
	
	function clearFormDelayed(delay) {
		setTimeout(clearForm, delay)
	}
	
	async function submitForm() {
		try {     
			let requestBody;
			let response;
			let processID = (await processPromise).ID;

			if (processExaminationState == "analysis") {
				requestBody = JSON.stringify({
					"userSequenceUserIDs": (await selection_sequence_promise).map((user) => user.ID),
					"userSequenceKindID": ctxSelectedIdx + 1,
					"processID": parseInt(processID)
				});
				
				console.log("[submitForm:analysis:requestBody]: ", requestBody);
				
				response = await fetch('http://localhost:3123/api/v1/user_sequence_simple', {
					method: 'post',
					
					body: requestBody,
					headers: {
						'Content-type': 'application/json; charset=UTF-8'
					}
				}
				);
			} else if (processExaminationState == "approval") {
				requestBody = JSON.stringify({
					"title": formState.title,
					"summary": formState.summary,
					"typeID": available_types[formState.selectedType].id,
					"centerID": available_centers[formState.selectedCenter].id,
					"userSequenceUserIDs": (await selection_sequence_promise).map((user) => user.ID)
				});
				
				console.log("[submitForm:approval:requestBody]: ", requestBody);
				
				response = await fetch('http://localhost:3123/api/v1/process', {
					method: 'put',
					
					body: requestBody,
					headers: {
						'Content-type': 'application/json; charset=UTF-8'
					}
				}
				);
			} else {
				return;
			}
			
			if (response.status == 200) {
				console.log('[ProcessModal::submitForm]: Successfully performed action: ', processExaminationState);
				open = false;
				clearForm();
				// updateRolesTable();
				signalBackendModification();
				// fireToastNotification("success", {email: formState.userValue});
			} else {
				console.log('[Add role]: Got valid response from server but process registration has failed.')
				console.log(response)
				// buildErrorToastFromResponse(response)
			}
			
		} catch(err) {
			console.error(`Error: ${err}`);
			return;
		}
	}
	
	let usersPromise = getAvailableUsers((users) => {
		for (const user of users) {
			user.negativePriority = 0;
		}
	});
	
	let available_users = null;
	usersPromise.then((users) => {available_users = users});
	
	function getSelectionSequencePromise(inner_available_users) {
		return new Promise((resolve, reject) => {
			if (inner_available_users == null) {
				return;
			}
			let sequence = inner_available_users.filter((user) => {return user.negativePriority > 0})
			sequence.sort((a, b) => {return a.negativePriority - b.negativePriority})
			console.log("[getSelectionSequencePromise::sequence]: ", sequence)
			resolve(sequence)
		})
	}
	
	let selection_sequence_promise = null;
	$: selection_sequence_promise = getSelectionSequencePromise(available_users)
	
	let negative_priority_counter = 1;
	
	// NOTE: This function was entirely written by Co-pilot.
	function getSetNegativePriorityClickCallback(user) {
		return () => {
			// NOTE: This part was written by me.
			console.log("[SetNegativePriorityClickCallback::user before modification]: ", user)
			if (user.negativePriority == 0) {
				user.negativePriority = negative_priority_counter;
				negative_priority_counter += 1;
			} else {
				user.negativePriority = 0;
			}
			console.log("[SetNegativePriorityClickCallback::user after modification]: ", user)
			available_users = available_users
		}
	}
</script>

<style>
	p {
		margin: 1.8em 0 1em 0;
	}
	
	p:first-child {
		margin-top: 1em
	}

	h4 {
		margin-top: 1.5em;
	}

	h6 {
		margin: 1.3em 0 0.3em 0;
		font-weight: 100;
	}
	
	h6:first-of-type {
		margin-top: 1em;
	}
	
	.userTagsWrapper {
		/* margin: 1em 0; */
		margin-left: 0;
	}
	
	.userTagsWrapper > * {
		margin-left: 0!important;
	}
	
	.userTagsWrapper p {
		margin-top: 0em;
		font-weight: 100;
		font-size: 12px;
		color: darkgray;
	}
</style>


<Modal
	bind:open
	modalHeading={modalHeading}
	primaryButtonText="Confirmar"
	secondaryButtonText="Cancelar"
	on:click:button--secondary={() => (open = false)}
	on:open={() => {
	}}
	on:close={() => {
		clearFormDelayed(800)
	}}
	on:submit={() => {
		submitForm()
	}}
>
		<p>De que maneira se deverá seguir com o processo?</p>
		<ContentSwitcher bind:selectedIndex={ctxSelectedIdx}>
			<Switch text={ctxSwitchOpt1} />
			<Switch text={ctxSwitchOpt2} />
		</ContentSwitcher>
		{#if !(processExaminationState == "approval" && ctxSelectedIdx == 1)}
			{#await purposePromise then p}
				{#await usersPromise then users}
					<p>Qual deverá ser a nova sequência?</p>
					<h6>Usuários disponíveis</h6>
					<div class=userTagsWrapper>
						{#await selection_sequence_promise then sequence}
							{#each users as user}
								<Tag type="gray" disabled={user.negativePriority != 0} on:click={getSetNegativePriorityClickCallback(user)}>{user.FirstName} {user.LastName}</Tag>
							{/each}
						{/await}
					</div>
					<h6>Usuários selecionados</h6>
					<div class=userTagsWrapper>
						{#await selection_sequence_promise then sequence}
							{#if sequence.length > 0}
								{#each sequence as user}
									<Tag type="gray" on:click={getSetNegativePriorityClickCallback(user)}>{user.FirstName} {user.LastName}</Tag>
								{/each}
							{:else}
								<p>Nenhum usuário foi selecionado: o processo será gerenciado exclusivamente pelo Secretário.</p>
							{/if}
						{/await}
					</div>
				{/await}
			{/await}
		{/if}
</Modal>