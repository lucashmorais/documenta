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
	
	let available_roles = []
	let splitRolesPromise = getSplitRoles(3);
	splitRolesPromise.then(value => console.log(value))
	
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
							roleObj.summary= u.Name
							roleObj.selected = false
							console.log(roleObj)
							available_roles.push(roleObj)
							// {"Name":"","Roles":null}},{"ID":18,"CreatedAt":"2021-07-16T16:29:48.508153567-03:00","UpdatedAt":"2021-07-16T16:29:48.508153567-03:00","DeletedAt":null,"Name":"","FirstName":"","LastName":"","Title":"","Initials":"","Email":"bob5@gmail.com","PHash":"$s2$16384$8$1$Y6/11yOsr8lGANCNCgYjqgQt$j4cqxYraVArl+tIN0y7WZu7/YARYhkcQVbXpOIwNrFo=",
						}
						resolve(roles)
					}
				)
			)
		})
	}
	
	// TODO: ENSURE THIS WORKS FOR ALL POSSIBLE REMAINDER VALUES
	function getSplitRoles(numBlocks) {
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
		const isActive = userInfo != null && userInfo.roles.filter((p) => p == role.summary).length > 0
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
		{#await splitRolesPromise}
		...
		{:then splitRoles}
			<Grid narrow padding>
				<Row>
					{#each splitRoles.reverse() as splitGroup}
						<Column>
							{#each splitGroup as role}
								{#if roleIsActive(role)}
									<Checkbox checked on:check={() => toggleRoleSelection(role)} labelText={role.summary}/>
								{:else}
									<Checkbox on:check={() => toggleRoleSelection(role)} labelText={role.summary}/>
								{/if}
							{/each}
						</Column>
					{/each}
				</Row>
		      </Grid>
		{/await}
	</FluidForm>
</Modal>