<script>
	import "carbon-components-svelte/css/all.css";
	import StatusBar from "./StatusBar.svelte"
	import DataTable from "./DataTable/DataTable.svelte"
	import UserModal from "./UserModal.svelte"
	import SimpleConfirmationModal from "./SimpleConfirmationModal.svelte"
	import {
		Toolbar,
		Button,
		ToolbarBatchActions,
		ToolbarContent,
		ToolbarSearch,
		Tag,
		Tile
	} from "carbon-components-svelte";
	// import Save16 from "carbon-icons-svelte/lib/Save16";
	import TrashCan16 from "carbon-icons-svelte/lib/TrashCan16";
	import Edit16 from "carbon-icons-svelte/lib/Edit16";
	import { getEndpointPrefix } from "./config-helper.js"

	let headers=[
		{ key: "email", value: "E-mail" },
		{ key: "firstName", value: "Primeiro nome" },
		{ key: "lastName", value: "Último nome" },
		{ key: "initials", value: "Iniciais" },
		// { key: "passwordHash", value: "Hash da senha" },
		{ key: "roleNames", value: "Funções" },
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
			fetch(getEndpointPrefix() + "/api/v1/users").
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
		window.open("/register.html", "_blank").focus();
	}
	
	async function submitBatchDeletion() {
		deleteModalIsOpen = false;
		console.log("[submitBatchDeletion::selectedRowIds]: ", selectedRowIds)

		try {     
			let requestBody = JSON.stringify({
						"ids": selectedRowIds,
					});

			console.log("[submitForm:requestBody]: ", requestBody);

			const response = await fetch(getEndpointPrefix() + "/api/v1/users", {
					method: "delete",

					body: requestBody,
					headers: {
						"Content-type": "application/json; charset=UTF-8"
					}
				}
			);
			
			if (response.status == 200) {
				console.log("[Add role]: Successfully deleted users");
				selectedRowIds = [];
				updateUsers();
			} else {
				console.log("[Add role]: Got valid response from server but user deletion has failed.")
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
		font-size: 20px;
	}
	
	h3 {
		font-size: 16px;
		font-weight: bold;
		margin-top: 1em;
	}
	
	h3:first-child {
		margin-top: 0;
	}

	.content1 {
		margin: 0 5vw;
	}
	
	.content2 {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1fr;
		grid-template-rows: 1fr 1fr 1fr;
		gap: 2em 2em;
		grid-auto-flow: row;
		grid-template-areas: ". . . ." ". . . ." ". . . .";
		margin-top: 2em;
		/* margin-bottom: 2em; */
	}
	
	.element {
		/* max-width: 200px; */
	}

</style>

<SimpleConfirmationModal
	bind:open={deleteModalIsOpen}
	bind:selectedItems={selectedRowIds}
	on:actionConfirmed={submitBatchDeletion}
	singularString="usuário"
	pluralString="usuários"
/>

<StatusBar />
<h1>Documenta</h1>
<!-- <h2 class=realh2>Usuários</h2> -->

<div class="content1">
	<h2>Minutas não alocadas</h2>
		<div class="content2">
			{#each [...Array(30).keys()] as i}
				<div class="element">
					<Tile>
						<h3>
							Autor
						</h3>
						<p>
							Fulano de Tal
						</p>
						<h3>
							Centro
						</h3>
						<p>
							Bauru
						</p>
						<h3>
							Resumo breve
						</h3>
						<p>
							Beli Orao (Serbo-Croatian for 'White Eagle') was a royal yacht built in 1938–39 for the Yugoslav Royal Navy, which intended her to serve as a patrol boat, escort, or guard ship in wartime. Upon completion, she was pressed into service as the admiralty yacht – used by senior admirals for transport and to review fleet exercises. She was captured in April 1941 by the Italians during the World War II Axis invasion of Yugoslavia. The Regia Marina (Italian Royal Navy) replaced her guns and used her as a gunboat for harbour protection and coastal escort duties, briefly as Alba then Zagabria. She was then used to train anti-submarine warfare specialists from the naval base at La Spezia. 
						</p>
					</Tile>
				</div>
			{/each}
		</div>
</div>
