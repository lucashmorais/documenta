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
