<script>
  export let key;
  export let task;
  import { ListItem, Avatar } from 'svelte-materialify';
</script>

<ListItem>
	<div slot="prepend">
	  <Avatar><img src="//picsum.photos/100?random=1" alt="Avatar" /></Avatar>
	</div>
	{ task._id }
	<span slot="subtitle"> <b>John Doe</b> - {task.text} </span>
</ListItem>