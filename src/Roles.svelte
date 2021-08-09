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
		Tag,
		Link,
		Modal
	} from "carbon-components-svelte";
	// import Save16 from "carbon-icons-svelte/lib/Save16";
	import TrashCan16 from "carbon-icons-svelte/lib/TrashCan16";
	import Edit16 from "carbon-icons-svelte/lib/Edit16";
	import Launch16 from "carbon-icons-svelte/lib/Launch16";
	import isEmail from 'validator/es/lib/isEmail';
	import isStrongPassword from 'validator/es/lib/isStrongPassword';
	import RoleModal from './RoleModal.svelte'

	let headers=[
		{ key: 'name', value: 'Nome' },
		{ key: 'description', value: 'Descrição' },
		{ key: 'permissions', value: 'Permissões' },
	];
									
	let open = false;
	let deleteConfirmationOpen = false;
	let formState = {}
	
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
		deleteConfirmationOpen = false;
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
				console.log('[Add role]: Successfully deleted roles');
				selectedRowIds = [];
				updateRolesTable();
				// fireToastNotification("success", {email: formState.userValue});
			} else {
				console.log('[Add role]: Got valid response from server but role deletion has failed.')
				console.log(response)
				// buildErrorToastFromResponse(response)
			}

		} catch(err) {
			console.error(`Error: ${err}`);
			return;
		}
	}

	var rolesPromise;
	export function updateRolesTable() {
		rolesPromise = new Promise((resolve, reject) => {
			fetch("http://localhost:3123/api/v1/roles").
				then((response)=>response.json().
					then(function (roles) {
						rows = []
						let roleObj = {}
						for (const r of roles) {
							console.log("[updateRolesTable]: Printing a role object obtained from the backend")
							console.log(r)
							roleObj = {}
							roleObj.id = r.ID
							roleObj.name = r.Name
							roleObj.description = r.Description
							roleObj.permissions = r.Permissions.map((p) => p.Summary)
							rows.push(roleObj)
						}
						resolve(roles)
					}
				)
			)
		})
	}

	updateRolesTable();
	
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
	
		<RoleModal
			bind:open={open}
			on:backendModification={updateRolesTable}
		/>

	      <Modal
		bind:open={deleteConfirmationOpen}
		modalHeading="Atenção!"
		primaryButtonText="Confirmar"
		secondaryButtonText="Cancelar"
		on:click:button--secondary={() => (deleteConfirmationOpen = false)}
		on:open
		on:close
		on:submit={submitBatchDeletion}
	      >
	      	{#if selectedRowIds.length > 1}
			Você tem certeza que gostaria de remover {selectedRowIds.length} funções?
		{:else}
			Você tem certeza que gostaria de deletar a função selecionada?
		{/if}
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
				  <Button on:click={() => deleteConfirmationOpen = true} icon={TrashCan16}>Remover</Button>
				  {#if selectedRowIds.length < 2}
					  <Button on:click={() => console.log("derp")} icon={Edit16}>Editar</Button>
				  {/if}
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
				<span slot="cell" let:row let:cell>
					{#if cell.key === 'permissions'}
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
