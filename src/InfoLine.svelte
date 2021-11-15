<script>
	import { DataTable, Dropdown } from "carbon-components-svelte";
	import { getEndpointPrefix } from "./config-helper.js"
	import { createEventDispatcher } from "svelte";

	let headers=[{ key: 'centro', value: 'Centro' }, { key: 'tipo', value: 'Tipo' }, { key: 'pend', value: 'Pendência Atual' }]
	let rows=[{ id: 'a', centro: 'sm', tipo: "Consulta", pend: 'Revisão do defensor' }]
	
	export let processPromise;
	
	const dispatch = createEventDispatcher();
	
	function signalChange() {
		dispatch("change");
	}

	let available_types;
	let available_centers;
	let current_type_dropdown_index;
	let backend_type_dropdown_index;
	let current_center_dropdown_index;
	let backend_center_dropdown_index;
	
	$: if (current_center_dropdown_index >= 0) {
		if (current_center_dropdown_index != backend_center_dropdown_index) {
			console.log(`[InfoLine::_reactive1::(current_center_dropdown_index, backend_center_dropdown_index)]: (${current_center_dropdown_index}, ${backend_center_dropdown_index})`)
			signalChange();				
		}
	}
	
	$: if (current_type_dropdown_index >= 0) {
		if (current_type_dropdown_index != backend_type_dropdown_index) {
			console.log(`[InfoLine::_reactive1::(current_type_dropdown_index, backend_type_dropdown_index)]: (${current_type_dropdown_index}, ${backend_type_dropdown_index})`)
			signalChange();				
		}
	}
	
	function updateProcessTypes() {
		return new Promise((resolve, reject) => {
			fetch(getEndpointPrefix() + "/api/v1/process_types").
				then((response)=>response.json().
					then(function (types) {
						available_types = []
						console.log("[updateProcessTypes::types]: ", types)
						for (const u of types) {
							let typeObj = {}
							typeObj.id = u.ID
							typeObj.text= u.Name
							typeObj.description = u.Description
							console.log(typeObj)
							available_types.push(typeObj)
						}
						resolve(types)
					})
				)
		})
	}

	function updateCenters() {
		return new Promise((resolve, reject) => {
			fetch(getEndpointPrefix() + "/api/v1/centers").
				then((response)=>response.json().
					then(function (centers) {
						available_centers = []
						console.log("[updateCenters::centers]: ", centers)
						for (const u of centers) {
							let centerObj = {}
							centerObj.id = u.ID
							centerObj.text= u.Name
							centerObj.shortName= u.ShortName
							centerObj.description = u.Description
							console.log(centerObj)
							available_centers.push(centerObj)
						}
						resolve(centers)
					})
				)
		})
	}
	
	export function refreshAndClear() {
		updateProcessTypes().then(() => {
			processPromise.then((process) => {
				backend_type_dropdown_index = available_types.findIndex((a) => a.id == process.ProcessTypeID);
				current_type_dropdown_index = backend_type_dropdown_index;
			})
		})

		updateCenters().then(() => {
			processPromise.then((process) => {
				backend_center_dropdown_index = available_centers.findIndex((a) => a.id == process.CenterID);
				current_center_dropdown_index = backend_center_dropdown_index;
			})
		})
	}
	refreshAndClear();
</script>

<style>
	.mainWrapper {
		margin-bottom: 5vh;
		width: 90%;
		margin: auto;
		justify-content: center;
	}
	
	/* :global(.dropdownWrapper > *) {
	    background-color: #f4f4f4!important;
	}
	
	:global(.dropdownWrapper:hover > *) {
	    background-color: #e5e5e5!important;
	}
	
	:global(*:hover > .dropdownWrapper:hover > *) {
	    background-color: #e5e5e5!important;
	}
	
	:global(*:hover > * > .dropdownWrapper:hover > *) {
	    background-color: #e5e5e5!important;
	} */
</style>

<div class="mainWrapper">
	<DataTable
		size="short"
		headers={headers}
		rows={rows}
	>
	<span slot="cell" let:cell>
		{#if cell.key === 'centro'}
			<div class=dropdownWrapper>
				<Dropdown
					titleText=""
					bind:selectedIndex={current_center_dropdown_index}
					items={available_centers}
				/>
			</div>
		{:else if cell.key === 'tipo'}
			<div class=dropdownWrapper>
				<Dropdown
					titleText=""
					bind:selectedIndex={current_type_dropdown_index}
					items={available_types}
				/>
			</div>
		{:else}
			{cell.value}
		{/if}
	</span>
	</DataTable>
</div>