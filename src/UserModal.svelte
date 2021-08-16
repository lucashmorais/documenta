<script>
	import {
		Modal,
		FluidForm,
		TextArea,
		TextInput
	} from "carbon-components-svelte"

	export let open = true
	
	let formState = {"name": "Derp", "description": "Aaaaa"}
	
	let nameIsInvalid = false;
	let descriptionIsInvalid = false;
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
		<TextInput bind:value={formState.name} invalidText="Título inválido" invalid={nameIsInvalid} labelText="Nome" required />
		<TextArea bind:value={formState.description} invalid={descriptionIsInvalid} invalidText="Descrição inválida" placeholder="Uma breve descrição da função" required />
		<!-- <h4 style="padding-top: 1em">Permissões</h4>
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
		{/await} -->
	</FluidForm>
</Modal>