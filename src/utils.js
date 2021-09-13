export function decodeDate(s) {
	const dtFormat = new Intl.DateTimeFormat('pt-BR', {
		// timeStyle: 'long',
		dateStyle: 'long',
		timeZone: 'America/Sao_Paulo',
	});

	return dtFormat.format(new Date(s * 1000));
}

export function decodeTime(s) {
	const dtFormat = new Intl.DateTimeFormat('pt-BR', {
		timeStyle: 'long',
		// dateStyle: 'full',
		timeZone: 'America/Sao_Paulo',
	});

	return dtFormat.format(new Date(s * 1000));
}

export function getAvailableCenters() {
	return new Promise((resolve, reject) => {
		fetch('http://localhost:3123/api/v1/centers').then((response) =>
			response.json().then((centers) => resolve(centers))
		);
	});
}

export function getAvailableUsers(preProcessingFunc = null) {
	return new Promise((resolve, reject) => {
		fetch('http://localhost:3123/api/v1/users').then((response) =>
			response.json().then((users) => {
				if (preProcessingFunc) preProcessingFunc(users);
				resolve(users);
			})
		);
	});
}

// Function that calls the /api/v1/user_sequence endpoint to create a new UserSequence
export function createUserSequence(userSequence) {
	return new Promise((resolve, reject) => {
		fetch('http://localhost:3123/api/v1/user_sequence', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(userSequence),
		}).then((response) => response.json().then((userSequence) => resolve(userSequence)));
	});
}

// Function that calls the /api/v1/user_sequences endpoint to fetch UserSequences
// If processID is provided, it will use it as a URL query parameter to only fetch the UserSequences with such processID
export function getUserSequences(processID) {
	return new Promise((resolve, reject) => {
		let url = 'http://localhost:3123/api/v1/user_sequences';
		if (processID) url += `?processID=${processID}`;
		fetch(url).then((response) =>
			response.json().then((userSequences) => resolve(userSequences))
		);
	});
}

export function getNameFromUser(user) {
	return `${user.FirstName} ${user.LastName}`;
}
