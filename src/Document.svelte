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
  } from "carbon-components-svelte";

  let isSideNavOpen = false;
</script>

<style>
	
	h1 {
		padding-bottom: 1em;
	}
	.textTile {
		padding-top: 10vh;
		margin-top: 2em;
	}
	
</style>

<Header company="CR" platformName="Documenta" bind:isSideNavOpen expandedByDefault=false persistentHamburgerMenu=true>
  <div slot="skip-to-content">
    <SkipToContent />
  </div>

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
	<h1>Documenta</h1>
	<div class="textTile">
		<Tile class="childClass">Content</Tile>
	</div>
</Content>