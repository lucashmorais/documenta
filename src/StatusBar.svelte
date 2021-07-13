<script>
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