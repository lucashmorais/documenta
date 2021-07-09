<script>
	import {
	  FluidForm,
	  TextInput,
	  PasswordInput,
	  Button,
	} from "carbon-components-svelte";
	
	let passwordValue = ""
	let userValue = ""
	
	async function submitForm() {
		console.log("User: " + userValue)
		console.log("Password: " + passwordValue)

		try {     
			const response = await fetch('http://localhost:3123/api/v1/login', {
					method: 'post',

					body: JSON.stringify({
						Email: userValue,
						Password: passwordValue
					}),

					headers: {
						'Content-type': 'application/json; charset=UTF-8'
					}
				}
			);
			
			if (response.status == 200) {
				console.log('[Login]: Successfully logged in with valid credentials.');
			} else {
				console.log('[Login]: Got valid response from server but login has failed.')
			}

			response.text().then((text) => {
				console.log(text);
			});

		} catch(err) {
			console.error(`Error: ${err}`);
			return;
		}
	}
	
	document.addEventListener("keyup", function(event) {
		if (event.key === 'Enter') {
			console.log('Enter is pressed!');
			submitForm();
		}
	});
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
			<TextInput bind:value={userValue} labelText="Usuário" placeholder="Insira o usuário" required />
			<PasswordInput
			  required
			  bind:value={passwordValue}
			  type="password"
			  labelText="Senha"
			  placeholder="Insira a senha"
			  showPasswordLabel="Exibir senha"
			  hidePasswordLabel="Ocultar senha"
			/>
		</FluidForm>

		<div class="button-holder">
			<!-- <a href="/home.html"><Button>Entrar</Button></a> -->
			<Button on:click={submitForm}>Entrar</Button>
		</div>
	</div>
</div>