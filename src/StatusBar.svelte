<script>
	import Login20 from "carbon-icons-svelte/lib/Login20";
	import UserMultiple20 from "carbon-icons-svelte/lib/UserMultiple20";
	import Cookie from "js-cookie";
	import { getNameFromUser } from "./utils";
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
	
	// Function that fetches a JSON object describing the currently logged user from the `localhost:3123/api/v1/current_user` endpoint
	function getCurrentUser() {
		return fetch("/api/v1/current_user", {
			method: "GET",
			headers: {
				"Content-Type": "application/json",
				"Authorization": "Bearer " + Cookie.get("documentaLoginToken")
			}
		})
		.then(response => response.json())
		.then(json => {
			console.log("[getCurrentUser.then.then::json]: ", json)
			return json;
		});
	}
	
	export let currentUserPromise = getCurrentUser();
</script>

<style>
	.userNameDisplay {
		color: white;
		display: flex;
		align-items: center;
	}
</style>

<Header company="CR" platformName="Documenta" bind:isSideNavOpen expandedByDefault=false persistentHamburgerMenu=true>
  <div slot="skip-to-content">
    <SkipToContent />
  </div>
  <HeaderUtilities>
	<!-- <HeaderGlobalAction aria-label="Settings" icon={Login20} /> -->
	{#await currentUserPromise}-
	{:then user}
		<div class="userNameDisplay">{getNameFromUser(user)}</div>
	{/await}
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