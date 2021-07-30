<script>
	import 'carbon-components-svelte/css/all.css';
	import StatusBar from './StatusBar.svelte'
	import DataTable from './DataTable/DataTable.svelte'
	import {
		Toolbar,
		Button,
		ButtonSet,
		ToolbarBatchActions,
		ToolbarContent,
		ToolbarSearch,
		FluidForm,
		
		PasswordInput,
		TextInput,
		TextArea,
		MultiSelect,
		Checkbox,
		Grid,
		Row,
		Column,
		Modal
	} from "carbon-components-svelte";
	// import Save16 from "carbon-icons-svelte/lib/Save16";
	import TrashCan16 from "carbon-icons-svelte/lib/TrashCan16";
	import Edit16 from "carbon-icons-svelte/lib/Edit16";
	import isEmail from 'validator/es/lib/isEmail';
	import isStrongPassword from 'validator/es/lib/isStrongPassword';

	let headers=[
		{ key: 'name', value: 'Nome' },
		{ key: 'description', value: 'Descrição' },
	];
									
	let open = false;
	let formState = {}
	
	formState.passwordValue = ""
	formState.passwordValue2 = ""
	formState.userValue = ""
	formState.firstName = ""
	formState.lastName = ""
	formState.initials = ""
	formState.title = ""
	
	let failedLastTime = false

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
	
	let emailCheckIsEnabled = false;
	let passwordCheckIsEnabled = false;
	let password2CheckIsEnabled = false;
	let firstNameCheckIsEnabled = false;
	let lastNameCheckIsEnabled = false;
	let initialsCheckIsEnabled = false;
	
	function enableValidation() {
		emailCheckIsEnabled = true;
		passwordCheckIsEnabled = true;
		password2CheckIsEnabled = true;
		firstNameCheckIsEnabled = true;
		lastNameCheckIsEnabled = true;
		initialsCheckIsEnabled = true;
	}
	
	function disableValidation() {
		emailCheckIsEnabled = false;
		passwordCheckIsEnabled = false;
		password2CheckIsEnabled = false;
		firstNameCheckIsEnabled = false;
		lastNameCheckIsEnabled = false;
		initialsCheckIsEnabled = false;
	}
	
	$: coreNameInvalid = formState.name == ""
	$: coreDescriptionInvalid = formState.description == ""

	$: nameInvalid = emailCheckIsEnabled && coreNameInvalid
	$: descriptionInvalid = firstNameCheckIsEnabled && coreDescriptionInvalid

	$: someInputIsInvalid = coreNameInvalid || coreDescriptionInvalid
	
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
	
	let rows=[];
	let selectedRowIds=[];
	
	async function submitBatchDeletion() {
		console.log("selectedRowIds", selectedRowIds)

		try {     
			let requestBody = JSON.stringify({
						"ids": selectedRowIds,
					});

			console.log("[submitForm:requestBody]: ", requestBody);

			const response = await fetch('http://localhost:3123/api/v1/roles', {
					method: 'delete',

					body: requestBody,
					headers: {
						'Content-type': 'application/json; charset=UTF-8'
					}
				}
			);
			
			if (response.status == 200) {
				console.log('[Add role]: Successfully registered role');
				failedLastTime = false;
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

	var usersPromise;
	export function updateUsers() {
		usersPromise = new Promise((resolve, reject) => {
			fetch("http://localhost:3123/api/v1/roles").
				then((response)=>response.json().
					then(function (users) {
						rows = []
						let userObj = {}
						for (const u of users) {
							userObj = {}
							// console.log(u)
							userObj.id = u.ID
							userObj.name = u.Name
							userObj.description = u.Description
							rows.push(userObj)
							// {"Name":"","Permissions":null}},{"ID":18,"CreatedAt":"2021-07-16T16:29:48.508153567-03:00","UpdatedAt":"2021-07-16T16:29:48.508153567-03:00","DeletedAt":null,"Name":"","FirstName":"","LastName":"","Title":"","Initials":"","Email":"bob5@gmail.com","PHash":"$s2$16384$8$1$Y6/11yOsr8lGANCNCgYjqgQt$j4cqxYraVArl+tIN0y7WZu7/YARYhkcQVbXpOIwNrFo=",
						}
						resolve(users)
					}
				)
			)
		})
	}
	
	let available_permissions = []
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

	updateUsers();
	let splitPermissionsPromise = getSplitPermissions(3);
	splitPermissionsPromise.then(value => console.log(value))
	
	function goToRegisterPage() {
		// window.open("/register.html", '_blank').focus();
		open = true;
	}
	
	// setInterval(() => console.log(selectedRowIds), 500)
</script>

<style>
	h1 {
		text-align: center;
		margin: 2em 0 1em 0;
	}

	h1:focus-visible {
		outline-style: none;
	}

	
	h2:not(.realh2) {
		font-size: 16px;
	}

	.content1 {
		margin: 0 5vw;
	}

	.content2 {
		/* margin: 1vh 5vw 4vh 5vw; */
		margin-top: 1em;
		margin-bottom: 2em;
	}

</style>

<StatusBar />
<h1>Documenta</h1>
<!-- <h2 class=realh2>Usuários</h2> -->

<div class="content1">
	<h2>Funções</h2>
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
			<TextInput bind:value={formState.name} invalidText="Título inválido" invalid={nameInvalid} labelText="Nome" required />
			<TextArea bind:value={formState.description} invalid={descriptionInvalid} invalidText="Descrição inválida" placeholder="Uma breve descrição da função" required />
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

		<!-- <div class="button-holder"> -->
			<!-- <ButtonSet> -->
				<!-- <a href="/home.html"><Button>Entrar</Button></a> -->
				<!-- <Button on:click={submitForm}>Registrar usuário</Button> -->
				<!-- <Button on:click={clearForm} kind="secondary">Limpar</Button> -->
			<!-- </ButtonSet> -->
		<!-- </div> -->
	      </Modal>
	      
		<div class="content2">
			<DataTable
				batchSelection
				bind:selectedRowIds
				headers={headers}
				rows={rows}
			>
			<Toolbar>
				<ToolbarBatchActions>
				  <!-- <Button icon={Save16}>Eliminar</Button> -->
				  <Button on:click={submitBatchDeletion} icon={TrashCan16}>Eliminar</Button>
				  <!-- <Button icon={Edit16}>Modificar</Button> -->
				</ToolbarBatchActions>
				<ToolbarContent>
				  <ToolbarSearch />
				  <!-- <ToolbarMenu>
				    <ToolbarMenuItem primaryFocus>Restart all</ToolbarMenuItem>
				    <ToolbarMenuItem href="https://cloud.ibm.com/docs/loadbalancer-service">
				      API documentation
				    </ToolbarMenuItem>
				    <ToolbarMenuItem danger>Stop all</ToolbarMenuItem>
				  </ToolbarMenu> -->
				  <Button on:click={goToRegisterPage}>Nova função</Button>
				</ToolbarContent>
		      </Toolbar>
			</DataTable>
		</div>
</div>
