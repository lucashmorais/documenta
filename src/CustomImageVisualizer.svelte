<script>
	import { onMount } from 'svelte'
	export let imageSrc = "https://upload.wikimedia.org/wikipedia/commons/thumb/9/97/The_Earth_seen_from_Apollo_17.jpg/599px-The_Earth_seen_from_Apollo_17.jpg"
	
	function draw() {
		let canvasElement = document.getElementById('canvas')
		var ctx = canvasElement.getContext('2d');
		var img = new Image();

		img.onload = function() {
			let naturalWidth = img.width
			let naturalHeight = img.height

			let width = naturalWidth
			let height = naturalHeight
			
			canvasElement.width = width
			canvasElement.height = height

			ctx.drawImage(img, 0, 0, width, height);
		};
		img.src = imageSrc;
	}

	onMount(() => {
		draw();
	})
</script>

<style>
	#test {
		display:flex!important;
		justify-content: center!important;
		align-items: center!important;
	}
</style>

<div class="canvasWrapper" id="test">
	<!-- <img src={imageSrc}/> -->
	<canvas id="canvas"></canvas>
</div>
