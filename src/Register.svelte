<script>
	import {
	  FluidForm,
	  TextInput,
	  PasswordInput,
	  Button,
	} from "carbon-components-svelte";
	
	let passwordValue = ""
	let passwordValue2 = ""
	let userValue = ""
	let firstName = ""
	let lastName = ""
	let initials = ""
	let title = ""
	
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

				response.text().then((text) => {
					console.log(text);
					window.location.href = "/home.html"
				});
			} else {
				console.log('[Login]: Got valid response from server but login has failed.')
			}

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
:global(html), :global(body) {
	height: 100%;
	width: 100%;
	margin: 0;
	max-width: 100%;
	overflow-x: hidden;
}

.centered-content {
	width: 100%;
	margin: 0;
	max-width: 100%;
	overflow-x: hidden;
}

h1 {
	margin-top: 5vh;
	margin-bottom: 5vh;
}

h1:focus-visible {
	outline-style: none;
}

h2 {
	margin-bottom: 1em;
}

.centered-content {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
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
	<h2>Novo usuário</h2>
<!-- 
type User struct {
	gorm.Model
	// ID       int
	Name        string `json: "name"`
	FirstName   string `json: "firstName"`
	LastName    string `json: "lastName"`
	Title       string `json: "title"`
	Initials    string `json: "initials"`
	Email       string `json: "email"`
	PHash       string `json: "phash"`
	Permissions []Permission
	Role        Role
} -->

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
			<PasswordInput
			  required
			  bind:value={passwordValue2}
			  type="password"
			  labelText="Confirmação da senha"
			  placeholder="Insira a senha outra vez"
			  showPasswordLabel="Exibir senha"
			  hidePasswordLabel="Ocultar senha"
			/>
			<TextInput bind:value={firstName} labelText="Primeiro nome" placeholder="Insira o usuário" required />
			<TextInput bind:value={lastName} labelText="Último nome" placeholder="Insira o usuário" required />
			<TextInput bind:value={initials} labelText="Iniciais" placeholder="ABC" required />
			<TextInput bind:value={title} labelText="Título" placeholder="Pe., Fr., etc" />
		</FluidForm>

		<div class="button-holder">
			<!-- <a href="/home.html"><Button>Entrar</Button></a> -->
			<Button on:click={submitForm}>Entrar</Button>
		</div>
	</div>
</div>