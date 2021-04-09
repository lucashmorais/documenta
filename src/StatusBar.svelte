<script>
	let event = new Date()
	const options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
	$: dateString = event.toLocaleDateString(undefined, options)

	let userName = "Aristotle da Silva"

	import "@fortawesome/fontawesome-free/js/all.js";
</script>
	

<style>
	.line {
		margin: 0 1em;
		display: flex;
		justify-content: space-between;
		font-size: 2.4vh;
		/* height: 5em; */
	}

	.container {
		padding: 1em 0;
		background-color: black;
		color: white;
	}

	.userInfoAndButton {
		display: flex
	}

	.userInfoAndButton * {
		margin-left: 1em;
	}
</style>

<div class="container">
	<div class="line">
		
		<i class="fas fa-bars"></i> <!-- uses solid style -->
		<div class="userInfoAndButton">
			<div>{userName}</div>
			<i class="fas fa-sign-in-alt"></i> <!-- uses solid style -->
		</div>
	</div>
</div>