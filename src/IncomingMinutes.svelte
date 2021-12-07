<script>
	import "carbon-components-svelte/css/all.css";
	import StatusBar from "./StatusBar.svelte"
	import { getNameFromUser } from "./utils.js"
	import {
		Tile
	} from "carbon-components-svelte";
	import { getEndpointPrefix } from "./config-helper.js"

	export function getMinutes() {
		return new Promise((resolve, reject) => {
			fetch(getEndpointPrefix() + "/api/v1/minutes").
				then((response)=>response.json().
					then(function (minutes) {
						// console.log("[Minutes::updateMinutes::minutes]: ", minutes)
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
	let minutesPromise = getMinutes()
</script>

<style>
	h1 {
		text-align: center;
		margin: 2em 0 1em 0;
	}

	h1:focus-visible {
		outline-style: none;
	}

	
	h2:not(.realh2) {
		font-size: 20px;
	}
	
	h3 {
		font-size: 16px;
		font-weight: bold;
		margin-top: 1em;
	}
	
	h3:first-child {
		margin-top: 0;
	}

	.content1 {
		margin: 0 5vw;
	}
	
	.content2 {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1fr;
		grid-template-rows: 1fr 1fr 1fr;
		gap: 2em 2em;
		grid-auto-flow: row;
		grid-template-areas: ". . . ." ". . . ." ". . . .";
		margin-top: 2em;
	}
	
	.element {
		/* max-width: 200px; */
	}

</style>

<StatusBar />
<h1>Documenta</h1>

<div class="content1">
	<h2>Minutas não alocadas</h2>
		<div class="content2">
			{#await minutesPromise then minutes}
				<!-- {#each [...Array(30).keys()] as i} -->
				{#each minutes as minute}
					<div class="element">
						<Tile>
							<h3>
								Autor
							</h3>
							<p>
								{getNameFromUser(minute.User)}
							</p>
							<h3>
								Centro
							</h3>
							<p>
								{minute.Center.Name}
							</p>
							<h3>
								Resumo breve
							</h3>
							<p>
								{minute.Content}
							</p>
						</Tile>
					</div>
				{/each}
			{/await}
		</div>
</div>
