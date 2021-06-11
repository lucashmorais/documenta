<script>
	// import PdfViewer from 'svelte-pdf';
	import PdfViewer from './custom_libs/svelte-pdf';
	let urlParams = new URLSearchParams(window.location.search);
	let id = urlParams.get('id')
</script>

<!-- <PdfViewer url='./build/test.pdf' /> -->

<PdfViewer url={'./api/v1/file/' + String(id)} />