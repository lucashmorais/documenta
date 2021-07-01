<script>
  import { getContext } from 'svelte';
  import ImageVisualizerModal from './ImageVisualizerModal.svelte';
  import { fly } from 'svelte/transition';
  import { Button } from "carbon-components-svelte";

  const { open } = getContext('simple-modal');

  const showSurprise = () => {
    open(ImageVisualizerModal, { message: "It's a modal!" }, 		  {
				transitionWindow: fly,
				transitionWindowProps: {
					y: 100,
					duration: 1000
				},
			},);
  };
</script>

<Button on:click={showSurprise}>Primary button</Button>