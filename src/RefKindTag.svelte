<script>
	import { Tag } from "carbon-components-svelte";

	export let isOutReference = false;
</script>

{#if isOutReference}
<Tag type="magenta">Saída</Tag>
{:else}
<Tag type="cyan">Entrada</Tag>
{/if}