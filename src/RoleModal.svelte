<script>
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
	
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	function signalBackendModification() {
		dispatch('backendModification');
	}
	
	export let open = false;
	let formState = {}
	
	let failedLastTime = false

	$: coreNameInvalid = formState.name == ""
	$: coreDescriptionInvalid = formState.description == ""
	$: someInputIsInvalid = coreNameInvalid || coreDescriptionInvalid

	let available_permissions = []
	let splitPermissionsPromise = getSplitPermissions(3);
	splitPermissionsPromise.then(value => console.log(value))
	
	export function updatePermissions() {
		return new Promise((resolve, reject) => {
			fetch("http://localhost:3123/api/v1/permissions").
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
	}
	
	function disableValidation() {
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
	
	function clearForm() {
		const keys = Object.keys(formState)
		for (const key of keys) {
			formState[key] = ""
		}
		disableValidation()
		unselectAllPermissions()
		setTimeout(() => splitPermissionsPromise = getSplitPermissions(3), 700)
		// splitPermissionsPromise = null
	}
	
	async function submitForm() {
		if (someInputIsInvalid) {
			enableValidation()
			return
		}
		disableValidation()

		try {     
			let requestBody = JSON.stringify({
						"name": formState.name,
						"description": formState.description,
						"permissions": getSelectedPermissionIds(),
					});

			console.log("[submitForm:requestBody]: ", requestBody);

			const response = await fetch('http://localhost:3123/api/v1/role', {
					method: 'post',

					body: requestBody,
					headers: {
						'Content-type': 'application/json; charset=UTF-8'
					}
				}
			);
			
			if (response.status == 200) {
				console.log('[Add role]: Successfully registered role');
				failedLastTime = false;
				open = false;
				clearForm();
				// updateRolesTable();
				signalBackendModification();
				// fireToastNotification("success", {email: formState.userValue});
			} else {
				console.log('[Add role]: Got valid response from server but role registration has failed.')
				failedLastTime = true;
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
	on:open
	on:close={() => clearForm()}
	on:submit={() => submitForm()}
>
	<FluidForm>
		<TextInput bind:value={formState.name} invalidText="Título inválido" invalid={coreNameInvalid} labelText="Nome" required />
		<TextArea bind:value={formState.description} invalid={coreDescriptionInvalid} invalidText="Descrição inválida" placeholder="Uma breve descrição da função" required />
		<h4 style="padding-top: 1em">Permissões</h4>
		{#await splitPermissionsPromise}
		derp
		{:then splitPermissions}
			<Grid narrow padding>
				<Row>
					{#each splitPermissions.reverse() as splitGroup}
						<Column>
							{#each splitGroup as permission}
								<Checkbox on:check={() => togglePermissionSelection(permission)} labelText={permission.summary}/>
							{/each}
						</Column>
					{/each}
				</Row>
		      </Grid>
		{/await}
	</FluidForm>
</Modal>