<script>
	import {
		Modal,
		FluidForm,
		TextArea,
		TextInput,
		Grid,
		Column,
		Row,
		Dropdown,
		Tag,
		Checkbox
	} from "carbon-components-svelte"
	import isEmail from "validator/es/lib/isEmail";
	import isStrongPassword from "validator/es/lib/isStrongPassword";
	
	import { createEventDispatcher } from "svelte";
	import { getAvailableUsers } from "./utils.js";
	import { getEndpointPrefix } from "./config-helper.js"
	import { set } from "js-cookie";
	
	const dispatch = createEventDispatcher();
	
	function signalBackendModification() {
		dispatch("backendModification");
	}
	
	export let open = false;
	export let processPromise;
	export let sequencePromise;
	let sequenceInitializationWasDone = false;
	let textFieldsInitializationWasDone = false;
	
	$: if (processPromise && available_types && available_centers) {
		console.log("[ProcessModal::processPromise]: ", processPromise);		
		if (textFieldsInitializationWasDone == false) {
			processPromise.then((process) => {
					formState.title = process.Title;
					formState.summary = process.Summary;
					formState.selectedCenter = available_centers.findIndex((a) => a.id == process.CenterID);
					formState.selectedType = available_types.findIndex((a) => a.id == process.ProcessTypeID);
					textFieldsInitializationWasDone = true;
					// formState = formState
			})
		}
	} 
	
	$: if (!sequenceInitializationWasDone && sequencePromise) {
		console.log("[ProcessModal::sequencePromise]: ", sequencePromise);		
		sequencePromise.then((sequenceWrapper) => {
			console.log("[ProcessModal::sequenceWrapper.sequence.Users]: ", sequenceWrapper.sequence.Users);
			
			usersPromise.then((allUsers) => {
				if (!sequenceInitializationWasDone) {
					let local_negative_priority_counter = 1;
					for (const user of sequenceWrapper.sequence.Users) {
						let matchedUserIndex = allUsers.findIndex((u) => u.ID == user.ID);
						allUsers[matchedUserIndex].negativePriority = local_negative_priority_counter;
						local_negative_priority_counter++;
					}
					console.log("[ProcessModal::_reactiveIfCallback::allUsers]: ", allUsers);
					selection_sequence_promise = getSelectionSequencePromise(available_users);
					sequenceInitializationWasDone = true;
				}
			});
		})
	} 
	
	let formState = {
		"title": "",
		"summary": "",
		"selectedType": 0,
		"selectedCenter": 0,
	}
	
	let validationIsEnabled = false;
	
	$: coreTitleIsInvalid = formState != null && formState.title != null && formState.title == ""
	$: coreSummaryIsInvalid = formState != null && formState.summary != null && formState.summary == ""

	$: titleIsInvalid = validationIsEnabled && coreTitleIsInvalid
	$: summaryIsInvalid = validationIsEnabled && coreSummaryIsInvalid
	
	$: someInputIsInvalid = coreTitleIsInvalid || coreSummaryIsInvalid
	
	export let purpose = "registering"
	if (processPromise) {
		purpose = "editing";
	}
	$: if (purpose == "registering") { clearForm() }
	
	let purposePromise = null;
	function updatePurposePromise(ignored) {
		purposePromise = new Promise((resolve, reject) => {
			console.log("[purposePromise::_callback::purpose]: ", purpose);
			resolve(purpose)			
		})
	}
	$: updatePurposePromise(purpose)
	
	export let userInfo = null;
	
	function updateFormState(userInfo) {
		if (userInfo != null) {
			formState.email = userInfo.email;
			formState.firstName = userInfo.firstName;
			formState.lastName = userInfo.lastName;
			formState.initials = userInfo.initials;
			formState.roles = userInfo.roles;
			formState = formState
		} else {
			clearForm()
		}
	}
	
	$: console.log("[ProcessModal::userInfo]: ", userInfo)
	$: updateFormState(userInfo)
	
	function enableValidation() {
		validationIsEnabled = true;
	}
	
	function disableValidation() {
		validationIsEnabled = false;
	}
	
	let available_types;
	let available_centers = [{"id": 0, "text": "Brasília"}, {"id": 1, "text": "João Cachoeira"}];
	
	function updateProcessTypes() {
		return new Promise((resolve, reject) => {
			fetch(getEndpointPrefix() + "/api/v1/process_types").
				then((response)=>response.json().
					then(function (types) {
						available_types = []
						console.log("[updateProcessTypes::types]: ", types)
						for (const u of types) {
							let typeObj = {}
							typeObj.id = u.ID
							typeObj.text= u.Name
							typeObj.description = u.Description
							console.log(typeObj)
							available_types.push(typeObj)
						}
						resolve(types)
					})
				)
		})
	}
	updateProcessTypes();
	
	function updateCenters() {
		return new Promise((resolve, reject) => {
			fetch(getEndpointPrefix() + "/api/v1/centers").
				then((response)=>response.json().
					then(function (centers) {
						available_centers = []
						console.log("[updateCenters::centers]: ", centers)
						for (const u of centers) {
							let centerObj = {}
							centerObj.id = u.ID
							centerObj.text= u.Name
							centerObj.shortName= u.ShortName
							centerObj.description = u.Description
							console.log(centerObj)
							available_centers.push(centerObj)
						}
						resolve(centers)
					})
				)
		})
	}
	updateCenters();
	
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
		// if (available_users) {
		// 	for (const user of available_users) {
		// 		user.negativePriority = 0;
		// 	}
		// }
		disableValidation()
	}
	
	function clearFormDelayed(delay) {
		setTimeout(clearForm, delay)
	}
	
	async function submitForm() {
		if (someInputIsInvalid) {
			console.log("[submitForm:someInputIsInvalid:formState]: ", formState)
			console.log("[submitForm]: [titleIsInvalid, summaryIsInvalid] = ", [titleIsInvalid, summaryIsInvalid])
			enableValidation()
			return
		}
		disableValidation()
		
		try {     
			let requestBody;
			let response;
			if (purpose == "registering") {
				requestBody = JSON.stringify({
					"title": formState.title,
					"summary": formState.summary,
					"typeID": available_types[formState.selectedType].id,
					"centerID": available_centers[formState.selectedCenter].id,
					"userSequenceUserIDs": (await selection_sequence_promise).map((user) => user.ID)
				});
				
				console.log("[submitForm:registering:requestBody]: ", requestBody);
				
				response = await fetch(getEndpointPrefix() + "/api/v1/process", {
					method: "post",
					
					body: requestBody,
					headers: {
						"Content-type": "application/json; charset=UTF-8"
					}
				}
				);
			} else if (purpose == "editing") {
				processPromise.then(async function(process) {
					requestBody = JSON.stringify({
						"title": formState.title,
						"summary": formState.summary,
						"typeID": available_types[formState.selectedType].id,
						"centerID": available_centers[formState.selectedCenter].id,
						"userSequenceUserIDs": (await selection_sequence_promise).map((user) => user.ID)
					});
					
					console.log("[submitForm:editing:requestBody]: ", requestBody);
					
					response = await fetch(getEndpointPrefix() + "/api/v1/process/" + process.ID, {
						method: "put",
						
						body: requestBody,
						headers: {
							"Content-type": "application/json; charset=UTF-8"
						}
					}
					);
				})
			} else {
				return;
			}
			
			if (response.status == 200) {
				console.log("[ProcessModal::submitForm]: Successfully performed action: ", purpose);
				open = false;
				clearForm();
				// updateRolesTable();
				signalBackendModification();
				// fireToastNotification("success", {email: formState.userValue});
			} else {
				console.log("[Add role]: Got valid response from server but process registration has failed.")
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
	h4 {
		margin-top: 1.5em;
	}

	h6 {
		margin: 1.3em 0 0.3em 0;
	}
	
	.userTagsWrapper {
		/* margin: 1em 0; */
		margin-left: 0;
	}
	
	.userTagsWrapper > * {
		margin-left: 0!important;
	}
	
	p {
		font-weight: 100;
		color: darkgray;
	}
</style>


<Modal
	bind:open
	modalHeading="Dados básicos"
	primaryButtonText="Confirmar"
	secondaryButtonText="Cancelar"
	on:click:button--secondary={() => (open = false)}
	on:open={() => {
		updateProcessTypes()
		updateCenters()
		updateFormState(userInfo)
		sequenceInitializationWasDone = false;
		textFieldsInitializationWasDone = false;
	}}
	on:close={() => {
		sequenceInitializationWasDone = false;
		textFieldsInitializationWasDone = false;
		disableValidation()
		clearFormDelayed(800)
	}}
	on:submit={() => {
		submitForm()
	}}
>
		{#await purposePromise then p}
			<TextInput bind:value={formState.title} invalidText="Título pequeno demais" invalid={titleIsInvalid} labelText="Título" required />
			<!-- TODO: Ensure race conditions involving `available_types` do not cause any trouble -->
			<Dropdown
				titleText="Tipo de processo"
				bind:selectedIndex={formState.selectedType}
				items={available_types}
			/>
			<Dropdown
				titleText="Centro"
				bind:selectedIndex={formState.selectedCenter}
				items={available_centers}
			/>
			<TextArea bind:value={formState.summary} invalidText="Descrição demasiado curta" invalid={summaryIsInvalid} labelText="Resumo" required />
			{#await usersPromise then users}
				<h4>Sequência de análise</h4>
				<h6>Usuários disponíveis</h6>
				<div class=userTagsWrapper>
					{#await selection_sequence_promise then sequence}
						{#each users as user}
							<Tag type="high-contrast" disabled={user.negativePriority != 0} on:click={getSetNegativePriorityClickCallback(user)}>{user.FirstName} {user.LastName}</Tag>
						{/each}
					{/await}
				</div>
				<h6>Usuários selecionados</h6>
				<div class=userTagsWrapper>
					{#await selection_sequence_promise then sequence}
						{#if sequence.length > 0}
							{#each sequence as user}
								<Tag type="high-contrast" on:click={getSetNegativePriorityClickCallback(user)}>{user.FirstName} {user.LastName}</Tag>
							{/each}
						{:else}
							<p>Nenhum usuário foi selecionado: o processo será gerenciado exclusivamente pelo Secretário.</p>
						{/if}
					{/await}
				</div>
			{/await}
		{/await}
</Modal>