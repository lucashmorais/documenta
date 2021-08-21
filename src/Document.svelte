<script>
	import InfoLine from './InfoLine.svelte'
	import CommentArea from './CommentArea.svelte'
	import Minutes from './Minutes.svelte'
	import StatusBar from './StatusBar.svelte'
	import AttachmentsArea from './AttachmentsArea.svelte'
	import InterventionForm from './InterventionForm.svelte'
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

	let processPromise = updateProcess();
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

<StatusBar />

<Content>
	<div class="contents">
		<h1>
			<div>Definição do Nome do Clube Infantil do Sumaré</div>
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

		<h2>Anexos</h2>
		<AttachmentsArea processID={processID}/>

		<h2>Minutas</h2>
		<Minutes />

		<h2>Nova intervenção</h2>
		<InterventionForm on:commentWasPosted={refreshComments} />

		<h2>Intervenções</h2>
		<CommentArea bind:updateComments={coreRefreshComments}/>
	</div>
</Content>