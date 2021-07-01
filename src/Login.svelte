<script>
	import {
	  FluidForm,
	  TextInput,
	  PasswordInput,
	  Button,
	} from "carbon-components-svelte";
</script>
      
<style>
.centered-content, :global(html), :global(body) {
	height: 100%;
	width: 100%;
	margin: 0;
	max-width: 100%;
	overflow-x: hidden;
}

h1 {
	margin-bottom: 10vh;
}

h1:focus-visible {
	outline-style: none;
}

.centered-content {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	padding-bottom: 12vh;
}

.form {
	width: 90%;	
}

.button-holder {
	margin-top: 1em;
}
</style>

<div class="centered-content">
	<h1>Documenta</h1>

	<div class="form">
		<FluidForm>
			<TextInput labelText="Usuário" placeholder="Insira o usuário" required />
			<PasswordInput
			  required
			  type="password"
			  labelText="Senha"
			  placeholder="Insira a senha"
			  showPasswordLabel="Exibir senha"
			  hidePasswordLabel="Ocultar senha"
			/>
		</FluidForm>

		<div class="button-holder">
			<Button>Entrar</Button>
		</div>
	</div>
</div>