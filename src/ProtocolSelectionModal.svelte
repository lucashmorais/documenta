<script>
	import { Modal, ContentSwitcher, Switch, NumberInput } from "carbon-components-svelte";	
	import { getEndpointPrefix } from "./config-helper.js";
	
	export let open = false;
	
	function getClassTitle(pClass) {
		return (pClass.prefix != '' ? pClass.prefix + ' ' : '') + `${pClass.base}`
	}
	
	let protocolClasses = [
		{prefix: '', base: 100, maxNumElements: 100},
		{prefix: '', base: 200, maxNumElements: 100},
		{prefix: '', base: 300, maxNumElements: 200},
		{prefix: '', base: 500, maxNumElements: 100000},
		{prefix: 'B', base: 1, maxNumElements: 100000},
	];
	
	let selectedIndex = 0;
	let coreNextProtocolNumber = 0;
	let selectedProtocolNumber = 0;
	let lastSelectedProtocolNumber = -1;
	let mutex = false;

	$: nextProtocolNumber = coreNextProtocolNumber > 0 ? coreNextProtocolNumber : protocolClasses[selectedIndex].base
	$: if (nextProtocolNumber) {
		selectedProtocolNumber = nextProtocolNumber
		console.log("[__onNextProtocolNumberChange::nextProtocolNumber]: ", nextProtocolNumber)
	}
	
	function getNextOutboundProtocolNumber(startNumber = 0, inverseSearch = false) {
		let c = protocolClasses[selectedIndex]

		let prefix = c.prefix
		let baseNumber = c.base
		let maxElements = c.maxNumElements

		let url = getEndpointPrefix() + `/api/v1/next_protocol_number?baseNumber=${baseNumber}&maxValues=${maxElements}` 

		if (prefix != '') {
			url = url + `&prefix=${prefix}`
		}

		if (startNumber != 0) {
			url = url + `&startNumber=${startNumber}`
		}

		if (inverseSearch) {
			url = url + '&inverseSearch=true'
		}
		
		console.log("[getNextOutboundProtocolNumber::url]: " + url)

		return new Promise((resolve) => {
			fetch(url).then(
				(response) => response.json().then(
					(nextProtocolNumber) => resolve(nextProtocolNumber)
				)
			);
		});
	}
	
	async function handleClassChange() {
		console.log(`[handleClassChange::selectedIndex]: ${selectedIndex}`)
		coreNextProtocolNumber = await getNextOutboundProtocolNumber();
		console.log(`[handleClassChange::nextProtocolNumber]: ${coreNextProtocolNumber}`)
	}
	
	async function handleNumberChange() {
		if (mutex) {
			return
		}
		mutex = true;

		let inverseSearch = lastSelectedProtocolNumber > selectedProtocolNumber
		let nextAvailableProtocolNumber = await getNextOutboundProtocolNumber(selectedProtocolNumber, inverseSearch);

		console.log(`[handleClassChange::nextProtocolNumber]: ${nextAvailableProtocolNumber}`)
		if (selectedProtocolNumber != nextAvailableProtocolNumber) {
			selectedProtocolNumber = nextAvailableProtocolNumber			
		}

		lastSelectedProtocolNumber = selectedProtocolNumber

		mutex = false;
	}
</script>

<style>

</style>

<Modal
  preventCloseOnClickOutside
  bind:open
  modalHeading="Selecionar protocolo de saída"
  primaryButtonText="Confirmar"
  secondaryButtonText="Cancelar"
  on:click:button--secondary
  on:open
  on:close
  on:submit
>
	<ContentSwitcher on:change={handleClassChange} bind:selectedIndex size="sm">
		{#each protocolClasses as pClass}
			<Switch text={getClassTitle(pClass)} />
		{/each}
	</ContentSwitcher>
	<NumberInput
		min={nextProtocolNumber}
		max={protocolClasses[selectedIndex].base + protocolClasses[selectedIndex].maxNumElements}
		bind:value={selectedProtocolNumber}
		on:change={handleNumberChange}
		invalidText="Valor fora dos limites"
		helperText="Parte numérica do protocolo de saída da minuta"
	/>
</Modal>