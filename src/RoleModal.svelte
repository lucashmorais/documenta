<script>
	import { purpose } from "./stores.js"
	import {
		FluidForm,
		TextInput,
		TextArea,
		Checkbox,
		Grid,
		Row,
		Column,
		Modal
	} from "carbon-components-svelte";
	import { getEndpointPrefix } from "./config-helper.js"
	
	import { createEventDispatcher } from "svelte";

	const dispatch = createEventDispatcher();

	function signalBackendModification() {
		dispatch("backendModification");
	}
	
	export let open = false;
	export let selectedRole = null;

	let formState = {name: "", description: ""}
	
	let validationIsEnabled = false

	$: nameIsInvalid = nameInvalid(formState.name, validationIsEnabled)
	$: descriptionIsInvalid = descriptionInvalid(formState.description, validationIsEnabled)
	$: console.log("[RoleModal:purpose]: ", $purpose)
	
	$: console.log("[RoleModal::selectedRole]: ", selectedRole)
	
	function coreNameInvalid() {
		return formState.name == ""
	}
	
	function nameInvalid() {
		nameIsInvalid = coreNameInvalid() && validationIsEnabled
		console.log("[nameInvalid]: ", nameIsInvalid)
		return nameIsInvalid
	}
	
	function coreDescriptionInvalid() {
		return formState.description == ""
	}
	
	function descriptionInvalid() {
		descriptionIsInvalid = coreDescriptionInvalid() && validationIsEnabled
		console.log("[descriptionInvalid]: ", descriptionIsInvalid)
		return descriptionIsInvalid
	}

	function someInputIsInvalid() {
		return coreNameInvalid() || coreDescriptionInvalid()
	}
	
	function updateFormBasedOnPurpose() {
		console.log("[updateFormBasedOnPurpose] selectedRole: ", selectedRole)
		if (selectedRole != null) {
			console.log("[RoleModal::updateFormBasedOnPurpose]: Updating formState")
			formState.name = selectedRole.name
			formState.description = selectedRole.description
			formState = formState
			console.log("[RoleModal::updateFormBasedOnPurpose]: Updated formState: ", formState)
		}
		splitPermissionsPromise = getSplitPermissions(3);
	}

	let available_permissions = []
	let splitPermissionsPromise = getSplitPermissions(3);
	splitPermissionsPromise.then(value => console.log(value))
	
	export function updatePermissions() {
		return new Promise((resolve, reject) => {
			fetch(getEndpointPrefix() + "/api/v1/permissions").
				then((response)=>response.json().
					then(function (permissions) {
						available_permissions = []
						let permissionObj = {}
						for (const u of permissions) {
							permissionObj = {}
							// console.log(u)
							permissionObj.id = u.ID
							permissionObj.summary= u.Summary
							permissionObj.selected = false
							console.log(permissionObj)
							available_permissions.push(permissionObj)
							// {"Name":"","Permissions":null}},{"ID":18,"CreatedAt":"2021-07-16T16:29:48.508153567-03:00","UpdatedAt":"2021-07-16T16:29:48.508153567-03:00","DeletedAt":null,"Name":"","FirstName":"","LastName":"","Title":"","Initials":"","Email":"bob5@gmail.com","PHash":"$s2$16384$8$1$Y6/11yOsr8lGANCNCgYjqgQt$j4cqxYraVArl+tIN0y7WZu7/YARYhkcQVbXpOIwNrFo=",
						}
						resolve(permissions)
					}
				)
			)
		})
	}
	
	function enableValidation() {
		validationIsEnabled = true;
	}
	
	function disableValidation() {
		validationIsEnabled = false;
	}
	
	// TODO: ENSURE THIS WORKS FOR ALL POSSIBLE REMAINDER VALUES
	function getSplitPermissions(numBlocks) {
		return new Promise(async function(resolve) {
			await updatePermissions()
			let numPermissions = available_permissions.length
			let numPerBlock = numPermissions / numBlocks
			
			let splitPermissions = []
			for (let i = 0; i < numPermissions; i += numPerBlock) {
				splitPermissions.push(available_permissions.slice(i, i + numPerBlock))
			}
			resolve(splitPermissions)
		});
	}
	
	function getSelectedPermissions() {
		return available_permissions.filter(p => p.selected)
	}
	
	function getSelectedPermissionIds() {
		return getSelectedPermissions().map(p => p.id)
	}
	
	function unselectAllPermissions() {
		available_permissions.map(p => p.selected = false)
	}

	function togglePermissionSelection(permission) {
		permission.selected = !permission.selected
		console.log(getSelectedPermissions())
		console.log(getSelectedPermissionIds())
	}
	
	function permissionIsActive(permission) {
		console.log("[RoleModal::permissionIsActive::selectedRole]: ", selectedRole)
		const isActive = selectedRole != null && selectedRole.permissions.filter((p) => p == permission.summary).length > 0
		if (isActive) {
			togglePermissionSelection(permission)
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
		unselectAllPermissions()
		setTimeout(() => splitPermissionsPromise = getSplitPermissions(3), 700)
		// splitPermissionsPromise = null
	}
	
	async function submitForm() {
		if (someInputIsInvalid()) {
			console.log("[submitForm:someInputIsInvalid:formState]: ", formState)
			enableValidation()
			return
		}
		disableValidation()
		console.log("[submitForm]: Building request for answering purpose: ", $purpose)

		try {     
			let requestBody;
			let response;
			if ($purpose == "registering") {
				requestBody = JSON.stringify({
							"name": formState.name,
							"description": formState.description,
							"permissions": getSelectedPermissionIds(),
						});

				console.log("[submitForm:registering:requestBody]: ", requestBody);

				response = await fetch(getEndpointPrefix() + "/api/v1/role", {
						method: "post",

						body: requestBody,
						headers: {
							"Content-type": "application/json; charset=UTF-8"
						}
					}
				);
			} else if ($purpose == "editing") {
				requestBody = JSON.stringify({
							"id": selectedRole.id,
							"name": formState.name,
							"description": formState.description,
							"permissions": getSelectedPermissionIds(),
						});

				console.log("[submitForm:editing:requestBody]: ", requestBody);

				response = await fetch(getEndpointPrefix() + "/api/v1/role", {
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
				console.log("[Add role]: Successfully registered role");
				open = false;
				clearForm();
				// updateRolesTable();
				signalBackendModification();
				// fireToastNotification("success", {email: formState.userValue});
			} else {
				console.log("[Add role]: Got valid response from server but role registration has failed.")
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
		if ($purpose == "registering") {
			clearForm()
		} else {
			updateFormBasedOnPurpose()
		}
	}}
	on:close={() => clearForm()}
	on:submit={() => submitForm()}
>
	<FluidForm>
		<TextInput bind:value={formState.name} invalidText="Título inválido" invalid={nameIsInvalid} labelText="Nome" required />
		<TextArea bind:value={formState.description} invalid={descriptionIsInvalid} invalidText="Descrição inválida" placeholder="Uma breve descrição da função" required />
		<h4 style="padding-top: 1em">Permissões</h4>
		{#await splitPermissionsPromise}
		...
		{:then splitPermissions}
			<Grid narrow padding>
				<Row>
					{#each splitPermissions.reverse() as splitGroup}
						<Column>
							{#each splitGroup as permission}
								{#if permissionIsActive(permission)}
									<Checkbox checked on:check={() => togglePermissionSelection(permission)} labelText={permission.summary}/>
								{:else}
									<Checkbox on:check={() => togglePermissionSelection(permission)} labelText={permission.summary}/>
								{/if}
							{/each}
						</Column>
					{/each}
				</Row>
		      </Grid>
		{/await}
	</FluidForm>
</Modal>