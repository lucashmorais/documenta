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
	import { decodeDate, decodeTime } from './utils.js'

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

	export let processID = 0;

	let minutesPromise;
	export function updateMinutes() {
		minutesPromise = new Promise((resolve, reject) => {
			console.log("[Minutes::updateMinutes]: Just entering minutesPromise's inner function")
			if (processID == "" || processID == "0") {
				resolve([])
			}
			fetch("http://localhost:3123/api/v1/minutes?processID=" + processID).
				then((response)=>response.json().
					then(function (minutes) {
						console.log("[Minutes::updateMinutes::minutes]: ", minutes)
						for (let i = 0; i < minutes.length; i++) {
							let a = minutes[i]
							console.log(a)
						}
						resolve(minutes)
					}
				)
			)
		})
	}
	

	let newVersionBlocks = [];
	let numTextHidingCallbacksAdded = 0;

	updateMinutes();
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
	{#await minutesPromise}
	...
	{:then minutes}
		{#each minutes as minute}
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
					<AccordionItem open title="{decodeDate(minute.UnixCreatedAt)}, {decodeTime(minute.UnixCreatedAt)}, {minute.User.FirstName} {minute.User.LastName}: {minute.Description}">
						<p>
						{minute.Content}
						</p>
					</AccordionItem>
				</Accordion>
			</Tile>
		{/each}
	{/await}
</div>