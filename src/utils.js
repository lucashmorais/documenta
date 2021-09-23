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

// Function that calls the /api/v1/user_sequence/count_completion POST endpoint,
// passing `processID`  as a query parameter
export function countCompletion(processID) {
	return new Promise((resolve, reject) => {
		fetch(
			`http://localhost:3123/api/v1/user_sequence/count_completion?processID=${processID}`,
			{
				method: 'POST',
			}
		).then((response) => response.json().then((updatedSequence) => resolve(updatedSequence)));
	});
}

// Function that calls the /api/v1/process/:process_id/set_status/:status PUT endpoint to set the status of a `Process`
export function setProcessStatus(processID, status) {
	return new Promise((resolve, reject) => {
		fetch(`http://localhost:3123/api/v1/process/${processID}/set_status/${status}`, {
			method: 'PUT',
		}).then((response) => response.json().then((updatedProcess) => resolve(updatedProcess)));
	});
}
