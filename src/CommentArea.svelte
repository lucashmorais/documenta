<script>
	import { Tile } from 'carbon-components-svelte'
	import { onMount } from 'svelte'
	import { decodeDate, decodeTime } from './utils.js'
	import { getEndpointPrefix } from "./config-helper.js"
	import { createEventDispatcher } from "svelte";
	
	export let currentUserPromise;

	var commentsPromise;
	export let processID;

	export function updateComments() {
		commentsPromise = new Promise((resolve, reject) => {
			fetch(getEndpointPrefix() + "/api/v1/comments?processID=" + processID).
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
	
	const dispatch = createEventDispatcher();
	
	function signalCommentModification() {
		dispatch("commentModification");
	}
	
	function buildDispatchFunction(referenceToComment) {
		return function() {
			referenceToComment.Content = document.getElementById(`comment${referenceToComment.ID}`).innerText;
			console.log("[_dynamicDispatchFunction::referenceToComment]: ", referenceToComment)
			dispatch("commentModification", {
				payload: referenceToComment
			})
		}
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
	{#await currentUserPromise then currentUser}
		{#each comments as comment}
			<div class="single-comment">
			{#if comment.UserID == currentUser.ID}
				<div id="comment{comment.ID}" on:input={buildDispatchFunction(comment)} class="text" contenteditable="true">
					{comment.Content}
				</div>
			{:else}
				<div class="text">
					{comment.Content}
				</div>
			{/if}
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
	{/await}
{:catch error}
{/await}