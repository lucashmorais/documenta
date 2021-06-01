<script>
  import {
    Header,
    HeaderNav,
    HeaderNavItem,
    HeaderNavMenu,
    SideNav,
    SideNavItems,
    SideNavMenu,
    SideNavMenuItem,
    SideNavLink,
    SideNavDivider,
    SkipToContent,
    Content,
    Grid,
    Row,
    Column,
	Tile,
	TextArea,
	ButtonSet,
	Button,
	Accordion,
	AccordionItem,
	Dropdown
  } from "carbon-components-svelte";
	import ActionsBlock from './ActionsBlock.svelte';
	import InterventionForm from './InterventionForm.svelte';
	import { onMount } from 'svelte';

	function pushNodeReferenceToArray(node) {
		newVersionBlocks.push(node)
		console.log(newVersionBlocks.length)
	}

	function hideNodeAndPushReference(node) {
		node.style.display = "none";
		pushNodeReferenceToArray(node);
	}

	function removeOuterPadding(node) {
		node.parentNode.style["padding-right"] = "0!important";
		node.style["background-color"] = "red!important";
		node.parentNode.style["background-color"] = "blue!important";
	}

	function setupHidingCallback(node) {
		console.log("We are now going to add an event listener to the following node:");
		console.log(node);
		console.log("We would like this node to be able to toggle the display of the following text input node:");
		node.callbackID = numTextHidingCallbacksAdded;
		console.log(newVersionBlocks[node.callbackID]);
		numTextHidingCallbacksAdded++;
		node.addEventListener("click", () => newVersionBlocks[node.callbackID].style.display = "");
	}

	onMount(() => {
		document.querySelectorAll(".removeOuterPadding").forEach((node) => node.parentNode.style["padding-right"] = "0");
	})

	let newVersionBlocks = [];
	let numTextHidingCallbacksAdded = 0;
</script>

<style>
	.interventions {
		display: grid;
		grid-template-columns: minmax(150px, 1fr) 4fr;
		grid-template-rows: 1fr;
		gap: 0px 0px;
		grid-template-areas:
			". .";
		 gap: 0.5em 0.5em;
	}
</style>

<div class="interventions">
	<Tile><ActionsBlock city={"Campinas"} editAction={setupHidingCallback}/></Tile>
	<Tile>
		<Accordion>
			<div use:hideNodeAndPushReference>
				<AccordionItem open title="Nova versão">
					<div class="removeOuterPadding">
						<InterventionForm />
					</div>
				</AccordionItem>
			</div>
			<AccordionItem open title="10 de março de 2021, 18:55, Defensor">
				<p>
				Natural Language Classifier uses advanced natural language processing and
				machine learning techniques to create custom classification models. Users
				train their data and the service predicts the appropriate category for the
				inputted text.
				</p>
			</AccordionItem>
			<AccordionItem title="10 de março de 2021, 15:40, Vogal de sr">
				<p>
				Natural Language Classifier uses advanced natural language processing and
				machine learning techniques to create custom classification models. Users
				train their data and the service predicts the appropriate category for the
				inputted text.
				</p>
			</AccordionItem>
			<AccordionItem title="10 de março de 2021, 12:30, Vigário Regional: Adicionando sugestão de pessoas por cuidar dos novos eventos com os pais">
				<p>
				Analyze text to extract meta-data from content such as concepts, entities,
				emotion, relations, sentiment and more.
				</p>
			</AccordionItem>
			<AccordionItem title="10 de março de 2021, 11:55, Vogal de sr: Texto inicial">
				<p>
				Translate text, documents, and websites from one language to another.
				Create industry or region-specific translations via the service's
				customization capability.
				</p>
			</AccordionItem>
		</Accordion>
	</Tile>
	<Tile><ActionsBlock city={"Centro de Estudos"} editAction={setupHidingCallback}/></Tile>
	<Tile>
		<Accordion>
			<div use:hideNodeAndPushReference>
				<AccordionItem open title="Nova versão">
					<div class="removeOuterPadding">
						<InterventionForm />
					</div>
				</AccordionItem>
			</div>
			<AccordionItem open title="10 de março de 2021, 18:55, Defensor: atualizando texto de acordo com a sugestão recebida sobre tipos de atividades por realizar no novo Centro">
				<p>
				Natural Language Classifier uses advanced natural language processing and
				machine learning techniques to create custom classification models. Users
				train their data and the service predicts the appropriate category for the
				inputted text.
				</p>
			</AccordionItem>
			<AccordionItem title="10 de março de 2021, 11:55, Vogal de sr: Texto inicial">
				<p>
				Translate text, documents, and websites from one language to another.
				Create industry or region-specific translations via the service's
				customization capability.
				</p>
			</AccordionItem>
		</Accordion>
	</Tile>
	<Tile><ActionsBlock city={"Niterói"} editAction={setupHidingCallback}/></Tile>
	<Tile>
		<Accordion>
			<div use:hideNodeAndPushReference>
				<AccordionItem open title="Nova versão">
					<div class="removeOuterPadding">
						<InterventionForm />
					</div>
				</AccordionItem>
			</div>
			<AccordionItem open title="10 de março de 2021, 18:55, Defensor: Texto inicial">
				<p>
				Natural Language Classifier uses advanced natural language processing and
				machine learning techniques to create custom classification models. Users
				train their data and the service predicts the appropriate category for the
				inputted text.
				</p>
			</AccordionItem>
		</Accordion>
	</Tile>
</div>