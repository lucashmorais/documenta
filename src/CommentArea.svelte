<script>
	import { Tile } from 'carbon-components-svelte'
	import { onMount } from 'svelte'
	import { decodeDate, decodeTime } from './utils.js'

	var commentsPromise;
	export let processID;

	export function updateComments() {
		commentsPromise = new Promise((resolve, reject) => {
			fetch("http://localhost:3123/api/v1/comments?processID=" + processID).
				then((response)=>response.json().
					then(function (comments) {
						for (const c of comments) {
							console.log(c)
							console.log(c.Content) 
							console.log(c.UpdatedAt) 
						}
						resolve(comments)
					}
				)
			)
		})
	}

	updateComments();
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

	.core-date {
		grid-area: core-date;
		line-height: 1.5em;
	}


	.text {
		grid-area: text;
		line-height: 1.5em;
	}

</style>

{#await commentsPromise then comments}
	{#each comments as comment}
		<div class="single-comment">
		<div class="text">
			<!-- Aenean eros nulla, feugiat vulputate velit ac, mollis scelerisque tortor. Aliquam elementum sollicitudin mauris. Pellentesque mollis consectetur orci, sit amet vulputate arcu finibus pulvinar. Praesent euismod mi sit amet est malesuada, vel consequat est blandit. Duis tincidunt eu ipsum quis suscipit. Donec faucibus massa id augue dictum, at tempor magna auctor. Curabitur suscipit erat velit, vel cursus mi scelerisque a. Duis vestibulum tristique quam, at tempus nisl fringilla et. -->
			{comment.Content}
		</div>
		<div class="initials">
			<div class="grid-container">
				<div class="core-initials">
					{comment.User.Initials}
				</div>
				<div class="core-date">
					<!-- {unixDecodeTime(comment.UnixCreatedAt)} -->
					{decodeDate(comment.UnixCreatedAt)}
					<br>
					{decodeTime(comment.UnixCreatedAt)}
				</div>
			</div>
		</div>
		</div>
	{/each}
{:catch error}
{/await}