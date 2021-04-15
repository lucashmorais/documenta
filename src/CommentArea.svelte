<script>
	import { Tile } from 'carbon-components-svelte'
</script>
<style>
	.grid {
		margin-top: 2em;
		display: grid;
		grid-template-columns: 4fr 0.5fr;
		grid-template-rows: 1fr;
		gap: 0px 0px;
		grid-template-areas:
			". .";
		 gap: 0.5em 0.5em;
	}

	.single-comment {
		display: grid;
		grid-template-columns: 0.6fr 1.4fr;
		grid-template-rows: 1fr;
		gap: 0px 0px;
		grid-template-areas:
			"initials text";
		margin-bottom: 2em;
	}

	.initials {
		grid-area: initials;
		display: grid;
		grid-template-columns: 1fr;
		grid-template-rows: 1.2fr 0.8fr;
		gap: 0px 0px;
		grid-template-areas:
			"core-initials"
			"core-date";
	}

	.core-initials {
		grid-area: core-initials;
		font-size: 40px;
		margin-bottom: 0.3em;
	}

	.core-date { grid-area: core-date; }


	.text {
		grid-area: text;
		line-height: 1.5em;
	}

</style>

{#each [1,2,3] as iteration}
	<div class="single-comment">
	<div class="text">
		Aenean eros nulla, feugiat vulputate velit ac, mollis scelerisque tortor. Aliquam elementum sollicitudin mauris. Pellentesque mollis consectetur orci, sit amet vulputate arcu finibus pulvinar. Praesent euismod mi sit amet est malesuada, vel consequat est blandit. Duis tincidunt eu ipsum quis suscipit. Donec faucibus massa id augue dictum, at tempor magna auctor. Curabitur suscipit erat velit, vel cursus mi scelerisque a. Duis vestibulum tristique quam, at tempus nisl fringilla et.
	</div>
	<div class="initials">
		<div class="grid-container">
			<div class="core-initials">
				AA
			</div>
			<div class="core-date">
				07/10/2020
				15:43:10
			</div>
		</div>
	</div>
	</div>
{/each}
	<div class="single-comment">
	<div class="text">
		Aenean eros nulla, feugiat vulputate velit ac, mollis scelerisque tortor. Aliquam elementum sollicitudin mauris. Pellentesque mollis consectetur orci, sit amet vulputate arcu finibus pulvinar. Praesent euismod mi sit amet est malesuada, vel consequat est blandit. Duis tincidunt eu ipsum quis suscipit. Donec faucibus massa id augue dictum, at tempor magna auctor. Curabitur suscipit erat velit, vel cursus mi scelerisque a. Duis vestibulum tristique quam, at tempus nisl fringilla et.
		Aenean eros nulla, feugiat vulputate velit ac, mollis scelerisque tortor. Aliquam elementum sollicitudin mauris. Pellentesque mollis consectetur orci, sit amet vulputate arcu finibus pulvinar. Praesent euismod mi sit amet est malesuada, vel consequat est blandit. Duis tincidunt eu ipsum quis suscipit. Donec faucibus massa id augue dictum, at tempor magna auctor. Curabitur suscipit erat velit, vel cursus mi scelerisque a. Duis vestibulum tristique quam, at tempus nisl fringilla et.
		Aenean eros nulla, feugiat vulputate velit ac, mollis scelerisque tortor. Aliquam elementum sollicitudin mauris. Pellentesque mollis consectetur orci, sit amet vulputate arcu finibus pulvinar. Praesent euismod mi sit amet est malesuada, vel consequat est blandit. Duis tincidunt eu ipsum quis suscipit. Donec faucibus massa id augue dictum, at tempor magna auctor. Curabitur suscipit erat velit, vel cursus mi scelerisque a. Duis vestibulum tristique quam, at tempus nisl fringilla et.
		Aenean eros nulla, feugiat vulputate velit ac, mollis scelerisque tortor. Aliquam elementum sollicitudin mauris. Pellentesque mollis consectetur orci, sit amet vulputate arcu finibus pulvinar. Praesent euismod mi sit amet est malesuada, vel consequat est blandit. Duis tincidunt eu ipsum quis suscipit. Donec faucibus massa id augue dictum, at tempor magna auctor. Curabitur suscipit erat velit, vel cursus mi scelerisque a. Duis vestibulum tristique quam, at tempus nisl fringilla et.
	</div>
	<div class="initials">
		<div class="grid-container">
			<div class="core-initials">
				AA
			</div>
			<div class="core-date">
				07/10/2020
				15:43:10
			</div>
		</div>
	</div>
	</div>

<div class="grid">
</div>