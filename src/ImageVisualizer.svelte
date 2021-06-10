<script>
	import { Carousel } from "svelte-images";
	const { Modal, open, close } = Carousel;
	
	export let images = [
	{
		src: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/97/The_Earth_seen_from_Apollo_17.jpg/599px-The_Earth_seen_from_Apollo_17.jpg"
	},
	// {
	// 	src: "/borat.png"
	// },
	];
	
	const popModal = idx =>
	setTimeout(() => {
		open(images, idx);
	}, 0);
</script>

<style>
	.gallery {
		display: flex;
		flex-flow: row wrap;
	}
	
	img {
		width: 30%;
		height: auto !important;
		cursor: pointer;
	}
	img:hover {
		opacity: 0.5;
		filter: grayscale(0.5) blur(1px);
	}
</style>

<Modal />

<div class="gallery">
	{#each images as image, i}
	<img
	{...image}
	src={image.thumbnail || image.src}
	alt={image.alt || ''}
	on:click={() => popModal(i)} />
	{/each}
</div>