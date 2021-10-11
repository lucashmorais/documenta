<script>
	import {
	  FluidForm,
	  TextInput,
	  PasswordInput,
	  Button,
	  ButtonSet,
	  ToastNotification,
	} from "carbon-components-svelte";
	import StatusBar from './StatusBar.svelte';
	import { writable } from 'svelte/store';
	import isEmail from 'validator/es/lib/isEmail';
	import isStrongPassword from 'validator/es/lib/isStrongPassword';
	import { getEndpointPrefix } from "./config-helper.js"
	
	let formState = {}
	
	formState.passwordValue = ""
	formState.passwordValue2 = ""
	formState.userValue = ""
	formState.firstName = ""
	formState.lastName = ""
	formState.initials = ""
	formState.title = ""
	
	let failedLastTime = false

	let passwordOptions = {
		minLength: 8,
		minLowercase: 0,
		minUppercase: 0,
		minNumbers: 0,
		minSymbols: 0,
		returnScore: false,
		pointsPerUnique: 1,
		pointsPerRepeat: 10,
		pointsForContainingLower: 100,
		pointsForContainingUpper: 1000,
		pointsForContainingNumber: 10000,
		pointsForContainingSymbol: 100000
	}
	
	let emailCheckIsEnabled = false;
	let passwordCheckIsEnabled = false;
	let password2CheckIsEnabled = false;
	let firstNameCheckIsEnabled = false;
	let lastNameCheckIsEnabled = false;
	let initialsCheckIsEnabled = false;
	
	function enableValidation() {
		emailCheckIsEnabled = true;
		passwordCheckIsEnabled = true;
		password2CheckIsEnabled = true;
		firstNameCheckIsEnabled = true;
		lastNameCheckIsEnabled = true;
		initialsCheckIsEnabled = true;
	}
	
	function disableValidation() {
		emailCheckIsEnabled = false;
		passwordCheckIsEnabled = false;
		password2CheckIsEnabled = false;
		firstNameCheckIsEnabled = false;
		lastNameCheckIsEnabled = false;
		initialsCheckIsEnabled = false;
	}
	
	$: coreUserInvalid = !isEmail(formState.userValue)
	$: corePasswordInvalid = !isStrongPassword(formState.passwordValue, passwordOptions)
	$: corePassword2Invalid = formState.passwordValue != formState.passwordValue2
	$: coreFirstNameInvalid = formState.firstName == ""
	$: coreLastNameInvalid = formState.lastName == ""
	$: coreInitialsInvalid = formState.initials == ""

	$: userInvalid = emailCheckIsEnabled && coreUserInvalid
	$: passwordInvalid = passwordCheckIsEnabled && corePasswordInvalid
	$: password2Invalid = password2CheckIsEnabled && corePassword2Invalid
	$: firstNameInvalid = firstNameCheckIsEnabled && coreFirstNameInvalid
	$: lastNameInvalid = lastNameCheckIsEnabled && coreLastNameInvalid
	$: initialsInvalid = initialsCheckIsEnabled && coreInitialsInvalid

	$: weakPasswordMessage = "A senha escolhida deve ter pelo menos 8 caracteres."

	$: someInputIsInvalid = coreUserInvalid || corePasswordInvalid || corePassword2Invalid || coreFirstNameInvalid || coreLastNameInvalid || coreInitialsInvalid
	
	const invalidEntryMessage = "Entrada inválida"
	let invalidPasswordMessage = "Senha inválida"
	
	$: invalidPasswordMessage = evaluatePasswords(formState.passwordValue, formState.passwordValue2)
	
	function evaluatePasswords(value1, value2) {
		if (value1 != value2) {
			// passwordInvalid = true;
			// password2Invalid = true;
			return "Senhas fornecidas não coincidem"
		} else {
			// passwordInvalid = false;
			// password2Invalid = false;
			return ""
		}
	}
	
	function clearForm() {
		const keys = Object.keys(formState)
		for (const key of keys) {
			formState[key] = ""
		}
	}

	async function buildErrorToastFromResponse(response) {
		response.json().then((json_response) => {
			fireToastNotification(json_response.cause);
		})
	}
	
	async function submitForm() {
		console.log("User: " + formState.userValue)
		console.log("Password: " + formState.passwordValue)
		
		if (someInputIsInvalid) {
			enableValidation()
			return
		}
		disableValidation()
		
		// if (formState.passwordValue != formState.passwordValue2) {
		// 	console.log("Passwords don't match")
		// 	return
		// }

		try {     
			const response = await fetch(getEndpointPrefix() + "/api/v1/user", {
					method: 'post',

					body: JSON.stringify({
						"phash": formState.passwordValue,
						"email": formState.userValue,
						"firstName": formState.firstName,
						"lastName": formState.lastName,
						"initials": formState.initials,
						"title": formState.title,
					}),

					headers: {
						'Content-type': 'application/json; charset=UTF-8'
					}
				}
			);
			
			if (response.status == 200) {
				console.log("[Register]: Successfully registered user");
				failedLastTime = false;
				fireToastNotification("success", {email: formState.userValue});
			} else {
				console.log("[Register]: Got valid response from server but user registering has failed.")
				failedLastTime = true;
				buildErrorToastFromResponse(response)
			}

		} catch(err) {
			console.error(`Error: ${err}`);
			return;
		}
	}
	
	document.addEventListener("keyup", function(event) {
		if (event.key === "Enter") {
			console.log("Enter is pressed!");
			submitForm();
		}
	})
	
	let notifications = writable({});
	$notifications = [];
	
	function fireToastNotification(kind, extras = null) {
		const s = new Date().toLocaleString()
		const l = $notifications.length;	// get our current items list count

		switch(kind) {
			case "success":
				$notifications[l] = {
					kind: "success",
					title: "Sucesso",
					subtitle: "O usuário com e-mail " + extras.email + " foi criado com sucesso.",
					caption: s,
					iconDescription: "Fechar notificação"
				}
				clearForm();
			break;
			case "email_was_taken":
				$notifications[l] = {
					kind: "error",
					title: "Erro",
					subtitle: "Já existe um usuário com o email indicado.",
					caption: s,
					iconDescription: "Fechar notificação"
				}
			break;
			default:
				$notifications[l] = {
					kind: "error",
					title: "Erro",
					subtitle: "O usuário não pôde ser criado. Verifique os dados fornecidos ou consulte o administrador do sistema.",
					caption: s,
					iconDescription: "Fechar notificação"
				}
		}
	}
	
	// setInterval(() => {
	// 	fireToastNotification("error")
	// }, 1000)
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
	margin-top: calc(5vh + 3rem);
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
	margin-bottom: 2em;
}

.absoluteToastWrapper {
	position: absolute;
	bottom: 1em;	
	right: 1em;
}

.stickyToast {
	position: sticky;
}
</style>

<StatusBar />
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
			<TextInput bind:value={formState.userValue} invalidText="Endereço inválido" invalid={userInvalid} labelText="Endereço de email" placeholder="abc@gmail.com" required />
			<PasswordInput
			  invalidText={weakPasswordMessage}
			  required
			  bind:value={formState.passwordValue}
			  invalid={passwordInvalid}
			  type="password"
			  labelText="Senha"
			  placeholder="Insira a senha"
			  showPasswordLabel="Exibir senha"
			  hidePasswordLabel="Ocultar senha"
			/>
			<PasswordInput
			invalidText={invalidPasswordMessage} 
			  required
			  bind:value={formState.passwordValue2}
			  invalid={password2Invalid}
			  type="password"
			  labelText="Confirmação da senha"
			  placeholder="Insira a senha outra vez"
			  showPasswordLabel="Exibir senha"
			  hidePasswordLabel="Ocultar senha"
			/>
			<TextInput bind:value={formState.firstName} invalid={firstNameInvalid} invalidText={invalidEntryMessage} labelText="Primeiro nome" placeholder="Rufus" required />
			<TextInput bind:value={formState.lastName} invalid={lastNameInvalid} invalidText={invalidEntryMessage} labelText="Último nome" placeholder="Silvius" required />
			<TextInput bind:value={formState.initials} invalid={initialsInvalid} invalidText={invalidEntryMessage} labelText="Iniciais" placeholder="RS" required />
			<!-- <TextInput bind:value={formState.title} invalid={titleInvalid} invalidText={invalidEntryMessage} labelText="Título" placeholder="Pe., Fr., etc" /> -->
		</FluidForm>

		<div class="button-holder">
			<ButtonSet>
				<!-- <a href="/home.html"><Button>Entrar</Button></a> -->
				<Button on:click={submitForm}>Registrar usuário</Button>
				<Button on:click={clearForm} kind="secondary">Limpar</Button>
			</ButtonSet>
		</div>
	</div>
</div>

<div class="absoluteToastWrapper">
	<div class="stickyToast">
	      {#each $notifications as toast}
	      		<svelte:component this={ToastNotification} {...toast}/>
	      {/each}
	      <!-- <ToastNotification
		kind="error"
		title="Erro"
		subtitle="Um erro interno impediu a criação do usuário."
		iconDescription="Fechar notificação"
		caption={new Date().toLocaleString()}
	      /> -->
	      <!-- <ToastNotification
		kind="info"
		title="New updates"
		subtitle="Restart to get the latest updates."
		caption={new Date().toLocaleString()}
	      />
	      <ToastNotification
		kind="info-square"
		title="New updates"
		subtitle="Restart to get the latest updates."
		caption={new Date().toLocaleString()}
	      /> -->
	      <!-- <ToastNotification
		kind="warning"
		title="Scheduled maintenance"
		subtitle="Maintenance will last 2-4 hours."
		caption={new Date().toLocaleString()}
	      />
	      <ToastNotification
		kind="warning-alt"
		title="Scheduled maintenance"
		subtitle="Maintenance will last 2-4 hours."
		caption={new Date().toLocaleString()}
	      /> -->
	</div>
</div>
