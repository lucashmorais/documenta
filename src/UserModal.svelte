<script>
	import {
		Modal,
		FluidForm,
		TextArea,
		TextInput,
		Grid,
		Column,
		Row,
		PasswordInput,
		Checkbox
	} from "carbon-components-svelte"
	import isEmail from "validator/es/lib/isEmail";
	import isStrongPassword from "validator/es/lib/isStrongPassword";
	import { getEndpointPrefix } from "./config-helper.js"

	import { createEventDispatcher } from "svelte";

	const dispatch = createEventDispatcher();

	function signalBackendModification() {
		dispatch("backendModification");
	}

	export let open = false;
	
	let formState = {
		"email": "",
		"first_name": "",
		"last_name": "",
		"initials": "",
	}

	let passwordOptions = {
		minLength: 8,
		minLowercase: 0,
		minUppercase: 0,
		minNumbers: 0,
		minSymbols: 0,
		returnScore: false,
		pointsPerUnique: 1,
		pointsPerRepeat: 10,
		pointsForContainingLower: 100,
		pointsForContainingUpper: 1000,
		pointsForContainingNumber: 10000,
		pointsForContainingSymbol: 100000
	}
	
	let validationIsEnabled = false;
	
	$: coreUserInvalid = formState != null && formState.email != null && !isEmail(formState.email)
	$: corePasswordInvalid = purpose == "registering" && formState != null && formState.passwordValue != null && !isStrongPassword(formState.passwordValue, passwordOptions)
	$: corePassword2Invalid = purpose == "registering" && formState != null && formState.passwordValue2 != null && formState.passwordValue != formState.passwordValue2
	$: coreFirstNameInvalid = formState != null && formState.firstName != null && formState.firstName == ""
	$: coreLastNameInvalid = formState != null && formState.lastName != null && formState.lastName == ""
	$: coreInitialsInvalid = formState != null && formState.initials != null && formState.initials == ""

	$: userIsInvalid = validationIsEnabled && coreUserInvalid
	$: passwordIsInvalid = validationIsEnabled && corePasswordInvalid
	$: password2IsInvalid = validationIsEnabled && corePassword2Invalid
	$: firstNameIsInvalid = validationIsEnabled && coreFirstNameInvalid
	$: lastNameIsInvalid = validationIsEnabled && coreLastNameInvalid
	$: initialsAreInvalid = validationIsEnabled && coreInitialsInvalid
	
	let failedLastTime = false

	$: weakPasswordMessage = "A senha escolhida deve ter pelo menos 8 caracteres."

	$: someInputIsInvalid = coreUserInvalid || corePasswordInvalid || corePassword2Invalid || coreFirstNameInvalid || coreLastNameInvalid || coreInitialsInvalid
	
	const invalidEntryMessage = "Entrada inválida"
	let invalidPasswordMessage = "Senha inválida"
	
	$: invalidPasswordMessage = evaluatePasswords(formState.passwordValue, formState.passwordValue2)
	
	function evaluatePasswords(value1, value2) {
		if (value1 != value2) {
			// passwordInvalid = true;
			// password2Invalid = true;
			return "Senhas fornecidas não coincidem"
		} else {
			// passwordInvalid = false;
			// password2Invalid = false;
			return ""
		}
	}
	
	export let purpose = "editing"
	$: if (purpose == "registering") { clearForm() }

	let purposePromise = null;
	function updatePurposePromise(ignored) {
		purposePromise = new Promise((resolve, reject) => {
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
	
	$: console.log("[UserModal::userInfo]: ", userInfo)
	$: updateFormState(userInfo)
	
	function enableValidation() {
		validationIsEnabled = true;
	}
	
	function disableValidation() {
		validationIsEnabled = false;
	}
	
	let available_roles = []
	$: splitRolesPromise = getSplitRoles(3, userInfo);
	
	export function updateRoles() {
		return new Promise((resolve, reject) => {
			fetch(getEndpointPrefix() + "/api/v1/roles").
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
		if (someInputIsInvalid) {
			console.log("[submitForm:someInputIsInvalid:formState]: ", formState)
			console.log("[submitForm]: [coreUserInvalid, coreFirstNameInvalid, coreLastNameInvalid, coreInitialsInvalid, corePasswordInvalid, corePassword2Invalid] = ", [coreUserInvalid, coreFirstNameInvalid, coreLastNameInvalid, coreInitialsInvalid, corePasswordInvalid, corePassword2Invalid])
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
							"phash": formState.passwordValue,
							"firstName": formState.firstName,
							"lastName": formState.lastName,
							"initials": formState.initials,
							"roles": getSelectedRoleIds(),
						});

				console.log("[submitForm:registering:requestBody]: ", requestBody);

				response = await fetch(getEndpointPrefix() + "/api/v1/user", {
						method: "post",

						body: requestBody,
						headers: {
							"Content-type": "application/json; charset=UTF-8"
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

				response = await fetch(getEndpointPrefix() + "/api/v1/user", {
						method: "put",

						body: requestBody,
						headers: {
							"Content-type": "application/json; charset=UTF-8"
						}
					}
				);
			} else {
				return;
			}
			
			if (response.status == 200) {
				console.log("[UserModal::submitForm]: Successfully performed action: ", purpose);
				open = false;
				clearForm();
				// updateRolesTable();
				signalBackendModification();
				// fireToastNotification("success", {email: formState.userValue});
			} else {
				console.log("[Add role]: Got valid response from server but user registration has failed.")
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
	on:open={() => {
		updateFormState(userInfo)
		
	}}
	on:close={() => {
		disableValidation()
		clearFormDelayed(800)
	}}
	on:submit={() => {
		submitForm()
	}}
>
	<FluidForm>
		{#await purposePromise then p}
			<TextInput bind:value={formState.email} invalidText="Endereço de e-mail inválido" invalid={userIsInvalid} labelText="E-mail" required />
			{#if p == "registering"}
				<PasswordInput bind:value={formState.passwordValue} invalid={passwordIsInvalid} invalidText={weakPasswordMessage} labelText="Senha" />
				<PasswordInput bind:value={formState.passwordValue2} invalid={password2IsInvalid} invalidText="Senha diferente da anterior" labelText="Confirmação da senha" />
			{/if}
			<TextInput bind:value={formState.firstName} invalidText="Entrada inválida" invalid={firstNameIsInvalid} labelText="Primeiro nome" required />
			<TextInput bind:value={formState.lastName} invalidText="Entrada inválida" invalid={lastNameIsInvalid} labelText="Último nome" required />
			<TextInput bind:value={formState.initials} invalidText="Iniciais inválidas" invalid={initialsAreInvalid} labelText="Iniciais" required />
		{/await}
		<h4 style="padding-top: 1em">Funções</h4>
		{#await splitRolesPromise then splitRoles}
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