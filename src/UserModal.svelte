<script>
	import {
		Modal,
		FluidForm,
		TextArea,
		TextInput,
		Grid,
		Column,
		Row,
		Checkbox
	} from "carbon-components-svelte"

	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	function signalBackendModification() {
		dispatch('backendModification');
	}

	export let open = false;
	
	let formState = {
		"email": "",
		"first_name": "",
		"last_name": "",
		"initiais": "",
	}
	
	let validationIsEnabled = false;
	let emailIsInvalid = false;
	let firstNameIsInvalid = false;
	let lastNameIsInvalid = false;
	let initiaisAreInvalid = false;
	
	let purpose = "editing"

	export let userInfo = null;
	
	function updateFormState(userInfo) {
		if (userInfo != null) {
			//TODO: REMOVE THE FOLLOWING ONCE WE START RETRIEVING SUCH INFORMATION FROM THE DB!
			userInfo.roles = []

			formState.email = userInfo.email;
			formState.firstName = userInfo.firstName;
			formState.lastName = userInfo.lastName;
			formState.initials = userInfo.initials;
			formState.roles = userInfo.initials;
		} else {
			clearForm()
		}
	}
	
	$: console.log("[UserModal::userInfo]: ", userInfo)
	$: updateFormState(userInfo)
	
	function enableValidation() {
		validationIsEnabled = true;
	}
	
	function disableValidation() {
		validationIsEnabled = false;
	}
	
	function someInputIsInvalid() {
		//TODO: ACTUALLY IMPLEMENT THIS
		return false;
	}
	
	let available_roles = []
	$: splitRolesPromise = getSplitRoles(3, userInfo);
	
	export function updateRoles() {
		return new Promise((resolve, reject) => {
			fetch("http://localhost:3123/api/v1/roles").
				then((response)=>response.json().
					then(function (roles) {
						available_roles = []
						let roleObj = {}
						for (const u of roles) {
							roleObj = {}
							// console.log(u)
							roleObj.id = u.ID
							roleObj.name= u.Name
							roleObj.selected = false
							console.log(roleObj)
							available_roles.push(roleObj)
						}
						resolve(roles)
					}
				)
			)
		})
	}
	
	// TODO: ENSURE THIS WORKS FOR ALL POSSIBLE REMAINDER VALUES
	function getSplitRoles(numBlocks, ignoredParam) {
		return new Promise(async function(resolve) {
			await updateRoles()
			let numRoles = available_roles.length
			let numPerBlock = numRoles / numBlocks
			
			let splitRoles = []
			for (let i = 0; i < numRoles; i += numPerBlock) {
				splitRoles.push(available_roles.slice(i, i + numPerBlock))
			}
			resolve(splitRoles)
		});
	}
	
	function getSelectedRoles() {
		return available_roles.filter(p => p.selected)
	}
	
	function getSelectedRoleIds() {
		return getSelectedRoles().map(p => p.id)
	}
	
	function unselectAllRoles() {
		available_roles.map(p => p.selected = false)
	}

	function toggleRoleSelection(role) {
		role.selected = !role.selected
		console.log(getSelectedRoles())
		console.log(getSelectedRoleIds())
	}
	
	function roleIsActive(role) {
		// console.log("[UserModal::roleIsActive::role]: ", userInfo)
		// console.log("[UserModal::roleIsActive::userInfo]: ", userInfo)
		const isActive = userInfo != null && userInfo.roleNames.filter((p) => p == role.name).length > 0
		if (isActive) {
			toggleRoleSelection(role)
		}
		return isActive;
	}
	
	function clearForm() {
		const keys = Object.keys(formState)
		for (const key of keys) {
			formState[key] = ""
		}
		formState = formState
		disableValidation()
		unselectAllRoles()
		setTimeout(() => splitRolesPromise = getSplitRoles(3), 700)
		// splitRolesPromise = null
	}
	
	function clearFormDelayed(delay) {
		setTimeout(clearForm, delay)
	}

	async function submitForm() {
		if (someInputIsInvalid()) {
			console.log("[submitForm:someInputIsInvalid:formState]: ", formState)
			enableValidation()
			return
		}
		disableValidation()

		try {     
			let requestBody;
			let response;
			if (purpose == "registering") {
				requestBody = JSON.stringify({
							"email": formState.email,
							"firstName": formState.firstName,
							"lastName": formState.lastName,
							"initials": formState.initials,
							"roles": getSelectedRoleIds(),
						});

				console.log("[submitForm:registering:requestBody]: ", requestBody);

				response = await fetch('http://localhost:3123/api/v1/role', {
						method: 'post',

						body: requestBody,
						headers: {
							'Content-type': 'application/json; charset=UTF-8'
						}
					}
				);
			} else if (purpose == "editing") {
				requestBody = JSON.stringify({
							"id": userInfo.id,
							"email": formState.email,
							"firstName": formState.firstName,
							"lastName": formState.lastName,
							"initials": formState.initials,
							"roles": getSelectedRoleIds(),
						});

				console.log("[submitForm:editing:requestBody]: ", requestBody);

				response = await fetch('http://localhost:3123/api/v1/user', {
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
				console.log('[UserModal::submitForm]: Successfully performed action: ', purpose);
				open = false;
				clearForm();
				// updateRolesTable();
				signalBackendModification();
				// fireToastNotification("success", {email: formState.userValue});
			} else {
				console.log('[Add role]: Got valid response from server but role registration has failed.')
				console.log(response)
				// buildErrorToastFromResponse(response)
			}

		} catch(err) {
			console.error(`Error: ${err}`);
			return;
		}
	}
</script>


<Modal
	bind:open
	modalHeading="Dados básicos"
	primaryButtonText="Confirmar"
	secondaryButtonText="Cancelar"
	on:click:button--secondary={() => (open = false)}
	on:open={() => {}}
	on:close={() => {
		// clearFormDelayed(800)
	}}
	on:submit={() => {
		submitForm()
		// clearFormDelayed(800);
	}}
>
	<FluidForm>
		<TextInput bind:value={formState.email} invalidText="Título inválido" invalid={emailIsInvalid} labelText="E-mail" required />
		<TextInput bind:value={formState.firstName} invalidText="Entrada inválida" invalid={firstNameIsInvalid} labelText="Primeiro nome" required />
		<TextInput bind:value={formState.lastName} invalidText="Entrada inválida" invalid={lastNameIsInvalid} labelText="Último nome" required />
		<TextInput bind:value={formState.initials} invalidText="Iniciais inválidas" invalid={initiaisAreInvalid} labelText="Iniciais" required />
		<h4 style="padding-top: 1em">Funções</h4>
		{#await splitRolesPromise}
		...
		{:then splitRoles}
			<Grid narrow padding>
				<Row>
					{#each splitRoles.reverse() as splitGroup}
						<Column>
							{#each splitGroup as role}
								{#if roleIsActive(role)}
									<Checkbox checked on:check={() => toggleRoleSelection(role)} labelText={role.name}/>
								{:else}
									<Checkbox on:check={() => toggleRoleSelection(role)} labelText={role.name}/>
								{/if}
							{/each}
						</Column>
					{/each}
				</Row>
		      </Grid>
		{/await}
	</FluidForm>
</Modal>