<script>
	import InfoLine from './InfoLine.svelte'
	import CommentArea from './CommentArea.svelte'
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
  } from "carbon-components-svelte";

  let isSideNavOpen = false;
</script>

<style>
	h1 {
		margin-top: 1em;
		margin-bottom: 1em;
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

	.interventions {
		display: grid;
		grid-template-columns: 4fr 0.5fr;
		grid-template-rows: 1fr;
		gap: 0px 0px;
		grid-template-areas:
			". .";
		 gap: 0.5em 0.5em;
	}

	.interventionButtonSet {
		margin-top: 0.5em;
	}

	:global(.summary) {
		font-size: 20px;
		line-height: 150%;
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
</Header>Go to prayer

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
	<h1>Definição do Nome do Clube Infantil do Sumaré</h1>
	<div class="contents">
		<InfoLine />
		<h2>Resumo</h2>
		<Tile class="summary">
			Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce bibendum urna urna, sed tincidunt quam placerat id. In bibendum, velit ac tristique pretium, erat nulla facilisis elit, id tincidunt dolor eros nec tellus. Nullam pharetra lacus ligula, ac condimentum purus rutrum ac. Pellentesque fermentum felis nisi, ut aliquet erat porttitor egestas. Vivamus sapien mi, tincidunt ac aliquet eget, sodales nec mi. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Curabitur ullamcorper enim eu condimentum placerat. Quisque tempor mauris porttitor leo sagittis ultricies. Interdum et malesuada fames ac ante ipsum primis in faucibus. Nulla vestibulum, turpis et maximus rutrum, magna nisi finibus massa, at bibendum tortor mi cursus augue. Duis molestie leo ut nisl commodo accumsan. Ut sit amet eros ipsum. Cras quis malesuada mi. Nulla commodo porta lorem, sed imperdiet nunc feugiat quis. In consequat, massa a elementum varius, velit purus mattis diam, eu rhoncus mauris lorem eleifend purus. Sed lacinia dictum eros in lobortis. In tempus fermentum placerat. Morbi consectetur, ex et sagittis placerat, nisl massa gravida metus, et sollicitudin orci nisi sit amet erat. Nunc mattis luctus justo ac pretium. Praesent tincidunt nulla eget facilisis accumsan. Pellentesque turpis felis, blandit a diam vel, pellentesque rutrum libero. Nam non tristique diam. Sed lacinia, ipsum congue eleifend vestibulum, ligula risus egestas elit, vel efficitur leo tellus sit amet libero. Cras efficitur, urna quis tincidunt consequat, mauris lorem laoreet sem, at dignissim elit augue vel arcu. Nulla facilisi. Vestibulum id purus ultricies, sagittis nunc posuere, commodo urna. Morbi leo orci, facilisis ac erat id, vestibulum vestibulum risus. Nam nec arcu neque. Nam dapibus augue quis varius bibendum.
		</Tile>

		<h2>Minutas</h2>
		<!-- <h2>Intervenções</h2> -->
		<div class="interventions">
			<Tile>
				<Accordion>
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
			<Tile>Campinas</Tile>
			<Tile>
				<Accordion>
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
			<Tile>Centro de Estudos</Tile>
			<Tile>
				<Accordion>
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
			<Tile>Niterói</Tile>
		</div>

		<h2>Intervenções</h2>
		<CommentArea/>

		<h2>Nova intervenção</h2>
		<TextArea
		rows={10}
		placeholder="Comentário por adicionar..."
		/>

		<div class="interventionButtonSet">
			<ButtonSet>
			<Button kind="secondary">Limpar</Button>
			<Button>Enviar</Button>
			</ButtonSet>
		</div>
	</div>
</Content>