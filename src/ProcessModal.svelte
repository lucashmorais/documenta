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
		Checkbox
	} from "carbon-components-svelte"
	import isEmail from 'validator/es/lib/isEmail';
	import isStrongPassword from 'validator/es/lib/isStrongPassword';
	
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
	
	let titleIsInvalid = false;
	let summaryIsInvalid = false;
	
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
	
	export let purpose = "registering"
	$: if (purpose == 'registering') { clearForm() }
	
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
	
	let available_types;
	
	function updateProcessTypes() {
		return new Promise((resolve, reject) => {
			fetch("http://localhost:3123/api/v1/process_types").
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
	
	function clearForm() {
		const keys = Object.keys(formState)
		for (const key of keys) {
			formState[key] = ""
		}
		formState = formState
		disableValidation()
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
				
				response = await fetch('http://localhost:3123/api/v1/user', {
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
				console.log('[Add role]: Got valid response from server but user registration has failed.')
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
		{#await purposePromise}
			Loading...
		{:then p}
			<TextInput invalidText="Endereço de e-mail inválido" invalid={titleIsInvalid} labelText="Título" required />
			<!-- <Dropdown
				titleText="Tipo de processo"
				selectedIndex={0}
				items={[{ id: '0', text: 'Slack' }, { id: '1', text: 'Email' }, { id: '2', text: 'Fax' }]}
			/> -->
			<!-- Ensure race conditions involving `available_types` do not cause any trouble -->
			<Dropdown
				titleText="Tipo de processo"
				selectedIndex={0}
				items={available_types}
			/>
			<TextArea invalidText="Endereço de e-mail inválido" invalid={summaryIsInvalid} labelText="Resumo" required />
		{/await}
</Modal>