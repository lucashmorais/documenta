<script>
  import { Button } from "carbon-components-svelte";
	import { onMount } from 'svelte';

	onMount(async () => {
		draw()
	});
	
	function draw() {
		var ctx = document.getElementById('canvas').getContext('2d');
		var img = new Image();
		img.onload = function() {
			ctx.drawImage(img, 0, 0);
		};
		img.src = 'https://upload.wikimedia.org/wikipedia/commons/thumb/3/39/I-69_exit_70_MI.jpg/800px-I-69_exit_70_MI.jpg';
	}
</script>

<style>
	.canvasContainer {
		display: flex
	}
</style>

<div class="canvasContainer">
	<canvas id="canvas" width="800" height="600"></canvas>
</div>
<Button>Derp</Button>