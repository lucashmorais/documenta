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
