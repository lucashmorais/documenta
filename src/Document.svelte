<script>
	import InfoLine from './InfoLine.svelte'
	import CommentArea from './CommentArea.svelte'
	import Minutes from './Minutes.svelte'
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

<Header company="CR" platformName="Documenta" bind:isSideNavOpen expandedByDefault=false persistentHamburgerMenu=true>
  <div slot="skip-to-content">
    <SkipToContent />
  </div>
  <HeaderUtilities>
	<!-- <HeaderGlobalAction aria-label="Settings" icon={Login20} /> -->
	<HeaderGlobalAction aria-label="Change user" icon={UserMultiple20} on:click={logout}/>
	<!-- <HeaderGlobalAction aria-label="Settings" /> -->
</HeaderUtilities>

  <!--
  <HeaderNav>
    <HeaderNavItem href="/" text="Link 1" />
    <HeaderNavItem href="/" text="Link 2" />
    <HeaderNavItem href="/" text="Link 3" />
    <HeaderNavMenu text="Menu">
      <HeaderNavItem href="/" text="Link 1" />
      <HeaderNavItem href="/" text="Link 2" />
      <HeaderNavItem href="/" text="Link 3" />
    </HeaderNavMenu>
    <HeaderNavItem href="/" text="Link 4" />
  </HeaderNav>
-->
</Header>

<SideNav bind:isOpen={isSideNavOpen}>
  <SideNavItems>
    <SideNavLink text="Link 1" />
    <SideNavLink text="Link 2" />
    <SideNavLink text="Link 3" />
	<!--
    <SideNavMenu text="Menu">
      <SideNavMenuItem href="/" text="Link 1" />
      <SideNavMenuItem href="/" text="Link 2" />
      <SideNavMenuItem href="/" text="Link 3" />
    </SideNavMenu>
    <SideNavDivider />
	-->
    <SideNavLink text="Link 4" />
  </SideNavItems>
</SideNav>

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
		<InfoLine />
		<h2>Resumo</h2>
		<Tile class="summary">
			Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce bibendum urna urna, sed tincidunt quam placerat id. In bibendum, velit ac tristique pretium, erat nulla facilisis elit, id tincidunt dolor eros nec tellus. Nullam pharetra lacus ligula, ac condimentum purus rutrum ac. Pellentesque fermentum felis nisi, ut aliquet erat porttitor egestas. Vivamus sapien mi, tincidunt ac aliquet eget, sodales nec mi. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Curabitur ullamcorper enim eu condimentum placerat. Quisque tempor mauris porttitor leo sagittis ultricies. Interdum et malesuada fames ac ante ipsum primis in faucibus. Nulla vestibulum, turpis et maximus rutrum, magna nisi finibus massa, at bibendum tortor mi cursus augue. Duis molestie leo ut nisl commodo accumsan. Ut sit amet eros ipsum. Cras quis malesuada mi. Nulla commodo porta lorem, sed imperdiet nunc feugiat quis. In consequat, massa a elementum varius, velit purus mattis diam, eu rhoncus mauris lorem eleifend purus. Sed lacinia dictum eros in lobortis. In tempus fermentum placerat. Morbi consectetur, ex et sagittis placerat, nisl massa gravida metus, et sollicitudin orci nisi sit amet erat. Nunc mattis luctus justo ac pretium. Praesent tincidunt nulla eget facilisis accumsan. Pellentesque turpis felis, blandit a diam vel, pellentesque rutrum libero. Nam non tristique diam. Sed lacinia, ipsum congue eleifend vestibulum, ligula risus egestas elit, vel efficitur leo tellus sit amet libero. Cras efficitur, urna quis tincidunt consequat, mauris lorem laoreet sem, at dignissim elit augue vel arcu. Nulla facilisi. Vestibulum id purus ultricies, sagittis nunc posuere, commodo urna. Morbi leo orci, facilisis ac erat id, vestibulum vestibulum risus. Nam nec arcu neque. Nam dapibus augue quis varius bibendum.
		</Tile>

		<h2>Anexos</h2>
		<AttachmentsArea/>

		<h2>Minutas</h2>
		<Minutes />

		<h2>Nova intervenção</h2>
		<InterventionForm on:commentWasPosted={refreshComments} />

		<h2>Intervenções</h2>
		<CommentArea bind:updateComments={coreRefreshComments}/>
	</div>
</Content>