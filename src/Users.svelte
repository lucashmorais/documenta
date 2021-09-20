<script>
	import 'carbon-components-svelte/css/all.css';
	import StatusBar from './StatusBar.svelte'
	import DataTable from './DataTable/DataTable.svelte'
	import UserModal from './UserModal.svelte'
	import SimpleConfirmationModal from './SimpleConfirmationModal.svelte'
	import {
		Toolbar,
		Button,
		ToolbarBatchActions,
		ToolbarContent,
		ToolbarSearch,
		Tag
	} from "carbon-components-svelte";
	// import Save16 from "carbon-icons-svelte/lib/Save16";
	import TrashCan16 from "carbon-icons-svelte/lib/TrashCan16";
	import Edit16 from "carbon-icons-svelte/lib/Edit16";

	let headers=[
		{ key: 'email', value: 'E-mail' },
		{ key: 'firstName', value: 'Primeiro nome' },
		{ key: 'lastName', value: 'Último nome' },
		{ key: 'initials', value: 'Iniciais' },
		// { key: 'passwordHash', value: 'Hash da senha' },
		{ key: 'roleNames', value: 'Funções' },
	];
	
	let rows=[];
	let selectedRowIds=[];
	let selectedRow = null;
	
	function updateSelectedRow(ignored_param) {
		if (selectedRowIds.length == 1) {
			selectedRow = rows.filter((row) => row.id == selectedRowIds[0])[0]
		} else {
			selectedRow = null;
		}
		
		console.log("[updateSelectedRow]: selectedRow: ", selectedRow)
	}

	$: updateSelectedRow(selectedRowIds)
	
	let editModalIsOpen = false;
	let deleteModalIsOpen = false;
	
	let editModalPurpose = null
	
	function updateEditModalPurposeBasedOnSelection(ignored) {
		if (selectedRowIds.length == 1) {
			editModalPurpose = "editing"
		} else {
			editModalPurpose = "registering"
		}
	}
	$: updateEditModalPurposeBasedOnSelection(selectedRowIds)
	$: console.log("[Users::editModalPurpose]: ", editModalPurpose)

	var usersPromise;
	export function updateUsers() {
		usersPromise = new Promise((resolve, reject) => {
			fetch("http://localhost:3123/api/v1/users").
				then((response)=>response.json().
					then(function (users) {
						console.log(users)
						rows = []
						let userObj = {}
						for (const u of users) {
							userObj = {}
							console.log(u)
							userObj.id = u.ID
							userObj.email = u.Email
							userObj.firstName = u.FirstName
							userObj.lastName = u.LastName
							userObj.lastName = u.LastName
							userObj.initials = u.Initials
							userObj.passwordHash = u.PHash
							userObj.roles = u.Roles
							userObj.roleNames = u.Roles.map((r) => r.Name)
							userObj.rolePermissions = u.Roles.map((r) => r.Permissions)
							console.log("[updateUsers]: User just built: ", userObj)
							rows.push(userObj)
						}
						resolve(users)
					}
				)
			)
		})
		return usersPromise
	}

	updateUsers();
	
	function goToRegisterPage() {
		window.open("/register.html", '_blank').focus();
	}
	
	async function submitBatchDeletion() {
		deleteModalIsOpen = false;
		console.log("[submitBatchDeletion::selectedRowIds]: ", selectedRowIds)

		try {     
			let requestBody = JSON.stringify({
						"ids": selectedRowIds,
					});

			console.log("[submitForm:requestBody]: ", requestBody);

			const response = await fetch('http://localhost:3123/api/v1/users', {
					method: 'delete',

					body: requestBody,
					headers: {
						'Content-type': 'application/json; charset=UTF-8'
					}
				}
			);
			
			if (response.status == 200) {
				console.log('[Add role]: Successfully deleted users');
				selectedRowIds = [];
				updateUsers();
			} else {
				console.log('[Add role]: Got valid response from server but user deletion has failed.')
				console.log(response)
				// buildErrorToastFromResponse(response)
			}

		} catch(err) {
			console.error(`Error: ${err}`);
			return;
		}
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

<SimpleConfirmationModal
	bind:open={deleteModalIsOpen}
	bind:selectedItems={selectedRowIds}
	on:actionConfirmed={submitBatchDeletion}
	singularString="usuário"
	pluralString="usuários"
/>

<UserModal 
	bind:open={editModalIsOpen}
	bind:userInfo={selectedRow}
	bind:purpose={editModalPurpose}
	on:backendModification={() => {
		console.log("[Users::backendModificationHandler]: Just got a new modification signal")
		updateUsers().then(() => {
			updateSelectedRow()
		})
	}}
/>

<StatusBar />
<h1>Documenta</h1>
<!-- <h2 class=realh2>Usuários</h2> -->

<div class="content1">
	<h2>Usuários</h2>
		<div class="content2">
			<DataTable
				batchSelection
				bind:selectedRowIds
				headers={headers}
				rows={rows}
			>
			<Toolbar>
				<ToolbarBatchActions>
				  <Button on:click={() => deleteModalIsOpen = true} icon={TrashCan16}>Eliminar</Button>
				  {#if selectedRowIds.length < 2}
					  <Button on:click={() => editModalIsOpen = true} icon={Edit16}>Editar</Button>
				  {/if}
				</ToolbarBatchActions>
				<ToolbarContent>
				  <ToolbarSearch />
				  <!-- <Button on:click={goToRegisterPage}>Novo usuário</Button> -->
				  <Button on:click={() => editModalIsOpen = true}>Novo usuário</Button>
				</ToolbarContent>
		      </Toolbar>
				<span slot="cell" let:row let:cell>
					{#if cell.key === 'roleNames'}
						{#each cell.value as v}
							  <Tag>
							    {v}
							  </Tag>
						{/each}
					{:else}
						{cell.value}
					{/if}
			      </span>
			</DataTable>
		</div>
</div>
