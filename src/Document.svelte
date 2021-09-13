<script>
	import { getAvailableCenters } from './utils.js'
	import InfoLine from './InfoLine.svelte'
	import CommentArea from './CommentArea.svelte'
	import Minutes from './Minutes.svelte'
	import StatusBar from './StatusBar.svelte'
	import AttachmentsArea from './AttachmentsArea.svelte'
	import InterventionForm from './InterventionForm.svelte'
	import SequenceTable from './SequenceTable.svelte'
	import Login20 from "carbon-icons-svelte/lib/Login20";
	import UserMultiple20 from "carbon-icons-svelte/lib/UserMultiple20";
	import Cookie from "js-cookie";
	import {
	Header,
	HeaderNav,
	HeaderNavItem,
	HeaderNavMenu,
	HeaderUtilities,
	HeaderGlobalAction,
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

	let isSideNavOpen = false;

	var coreRefreshComments;
	var coreRefreshMinutes;
	
	function refreshMinutes() {
		coreRefreshMinutes();
	}
	
	function refreshComments() {
		coreRefreshComments();
	}
	
	function logout() {
		Cookie.remove("documentaLoginToken")
		window.location.href = "/";
	}

	let urlParams = new URLSearchParams(window.location.search);
	let processID = urlParams.get('id')
	
	console.log("processID: ", processID)

	export function updateProcess() {
		let processPromise = new Promise((resolve, reject) => {
			if (processID != null) {
				fetch("http://localhost:3123/api/v1/processes?processID=" + processID).
					then((response)=>response.json().
						then(function (wrappedProcess) {
							console.log("[updateProcess]: ", wrappedProcess[0])
							resolve(wrappedProcess[0])
						}
					)
				)
			} else {
				resolve (null);
			}
		})
		return processPromise
	}
	
	let sequencePromise;
	let currentUserPromise;
	
	function printCurrentUser(p) {
		if (p) {
			p.then((user) => {
				console.log("[Document::printCurrentUser:then]: ", user)			
			})
		}
	}
	
	$: printCurrentUser(currentUserPromise)

	let processPromise = updateProcess();

	let availableCenters = [];
	getAvailableCenters().then((centers) => {
		console.log("[getAvailableCentersCallback::centers]: ", centers)
		availableCenters = centers
	})
	
	function currentUserHasModificationWrights(seqP, userP) {
		if (!seqP || !userP){
			return null;
		}
		return new Promise((resolve, reject) => {
			seqP.then((sequence) => {
				userP.then((user) => {
					console.log("[currentUserHasModificationWrights::then::then::sequence]: ", sequence)
					console.log("[currentUserHasModificationWrights::then::then::user]: ", user)
					let coreSeq = sequence.sequence
					let numCompletions = coreSeq.NumCompletions
					console.log("[currentUserHasModificationWrights::then::then::numCompletions]: ", numCompletions)
					
					let result;
					if (numCompletions >= coreSeq.Users.length) {
						result = false
					} else {
						result = coreSeq.Users[numCompletions].ID == user.ID
					}
					console.log("[currentUserHasModificationWrights::then::then::result]: ", result)
					resolve(result)
				})
			})
		})				
	}
	
	$: currentUserHasModificationWrightsPromise = currentUserHasModificationWrights(sequencePromise, currentUserPromise)
</script>

<style>
	h1 {
		margin-top: 1em;
		margin-bottom: 1em;
		text-align: left;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	h2 {
		font-size: 22px;
		margin: 3em 0 1em 0;
	}

	.contents {
		margin: 0 15vw;
	}
	
	.records {

	}

	:global(.summary) {
		font-size: 15px;
		line-height: 150%;
	}
</style>

<StatusBar bind:currentUserPromise/>

<Content>
	<div class="contents">
		<h1>
		{#await processPromise}
		...
		{:then process}
			<div>{process.Title}</div>
		{/await}
		<div style="margin-left: 0.5em">
			<Dropdown
				selectedIndex={0}
				items={[{ id: '0', text: 'Rascunho' }, { id: '1', text: 'Ativo' }, { id: '2', text: 'Concluído' }]}
			/>
		</div>
		</h1>
		<InfoLine processPromise={processPromise}/>
		<h2>Resumo</h2>
		<Tile class="summary">
			{#await processPromise}
			<div></div>
			{:then process}
				{#if process != null}
					{process.Summary}
				{:else}
					Lorem ipsum dolor sit amet
				{/if}
			{/await}
		</Tile>

		<h2>Sequência de análise</h2>
		<SequenceTable bind:sequencePromise processID={processID}/>

		<h2>Anexos</h2>
		<AttachmentsArea processID={processID}/>

		<h2>Minutas</h2>
		<Minutes processID={processID} bind:updateMinutes={coreRefreshMinutes} on:minuteWasPosted={refreshMinutes} availableCenters={availableCenters}/>

		<h2>Nova intervenção</h2>
		<InterventionForm processID={processID} on:commentWasPosted={refreshComments} on:minuteWasPosted={refreshMinutes} availableCenters={availableCenters}/>

		<h2>Intervenções</h2>
		<CommentArea processID={processID} bind:updateComments={coreRefreshComments}/>
	</div>
</Content>