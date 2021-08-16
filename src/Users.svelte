<script>
	import 'carbon-components-svelte/css/all.css';
	import StatusBar from './StatusBar.svelte'
	import DataTable from './DataTable/DataTable.svelte'
	import {
		Toolbar,
		Button,
		ToolbarBatchActions,
		ToolbarContent,
		ToolbarSearch
	} from "carbon-components-svelte";
	// import Save16 from "carbon-icons-svelte/lib/Save16";
	import TrashCan16 from "carbon-icons-svelte/lib/TrashCan16";
	import Edit16 from "carbon-icons-svelte/lib/Edit16";

	let headers=[
		{ key: 'email', value: 'E-mail' },
		{ key: 'firstName', value: 'Primeiro nome' },
		{ key: 'lastName', value: 'Último nome' },
		{ key: 'initials', value: 'Iniciais' },
		{ key: 'passwordHash', value: 'Hash da senha' },
	];
	
	let rows=[];
	let selectedRowIds=[];
	$: console.log("selectedRowIds: ", selectedRowIds)

	// let rows=[
	// 	{ id: 'a', assunto: 'Assunto 0', centro: 'cs', tipo: 3000, pend: 'Revisão do defensor' },
	// 	{ id: 'b', assunto: 'Assunto 1', centro: 'brs', tipo: 443, pend: 'Revisão inicial do Secretário' },
	// ];

	var usersPromise;
	export function updateUsers() {
		usersPromise = new Promise((resolve, reject) => {
			fetch("http://localhost:3123/api/v1/users").
				then((response)=>response.json().
					then(function (users) {
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
							rows.push(userObj)
							// {"Name":"","Permissions":null}},{"ID":18,"CreatedAt":"2021-07-16T16:29:48.508153567-03:00","UpdatedAt":"2021-07-16T16:29:48.508153567-03:00","DeletedAt":null,"Name":"","FirstName":"","LastName":"","Title":"","Initials":"","Email":"bob5@gmail.com","PHash":"$s2$16384$8$1$Y6/11yOsr8lGANCNCgYjqgQt$j4cqxYraVArl+tIN0y7WZu7/YARYhkcQVbXpOIwNrFo=",
						}
						resolve(users)
					}
				)
			)
		})
	}

	updateUsers();
	
	function goToRegisterPage() {
		window.open("/register.html", '_blank').focus();
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
				  <Button icon={TrashCan16}>Eliminar</Button>
				  {#if selectedRowIds.length < 2}
					  <Button on:click={() => console.log("Showing edit button")} icon={Edit16}>Editar</Button>
				  {/if}
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
				  <Button on:click={goToRegisterPage}>Novo usuário</Button>
				</ToolbarContent>
		      </Toolbar>
			</DataTable>
		</div>
</div>
