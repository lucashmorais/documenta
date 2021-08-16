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

	export let open = false;
	
	let formState = {
		"email": "",
		"first_name": "",
		"last_name": "",
		"initiais": "",
	}
	
	let validationIsEnabled = false;
	let nameIsInvalid = false;
	let descriptionIsInvalid = false;

	export let userInfo = null;
	
	$: console.log("[UserModal::userInfo]: ", userInfo)
	
	function enableValidation() {
		validationIsEnabled = true;
	}
	
	function disableValidation() {
		validationIsEnabled = false;
	}
	
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
		const isActive = userInfo != null && userInfo.permissions.filter((p) => p == permission.summary).length > 0
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
</script>


<Modal
	bind:open
	modalHeading="Dados básicos"
	primaryButtonText="Confirmar"
	secondaryButtonText="Cancelar"
	on:click:button--secondary={() => (open = false)}
	on:open={() => {}}
	on:close={() => {}}
	on:submit={() => {}}
>
	<FluidForm>
		<TextInput bind:value={formState.email} invalidText="Título inválido" invalid={nameIsInvalid} labelText="E-mail" required />
		<TextInput bind:value={formState.first_name} invalidText="Entrada inválida" invalid={nameIsInvalid} labelText="Primeiro nome" required />
		<TextInput bind:value={formState.last_name} invalidText="Entrada inválida" invalid={nameIsInvalid} labelText="Último nome" required />
		<TextInput bind:value={formState.initials} invalidText="Iniciais inválidas" invalid={nameIsInvalid} labelText="Iniciais" required />
		<h4 style="padding-top: 1em">Funções</h4>
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