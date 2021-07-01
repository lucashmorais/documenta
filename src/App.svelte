<script>
	import { Router, Link, Route } from "svelte-navigator";
	import Home from "./Home.svelte";
	import Login from "./Login.svelte";
	import Document from "./Document.svelte";
	// import About from "./routes/About.svelte";
	// import Blog from "./routes/Blog.svelte";
	// import Search from "./routes/Search.svelte";
</script>

<Router>
	<!-- <nav> -->
		<!-- <Link to="/">Home</Link> -->
		<!-- <Link to="/document/">Document</Link> -->
		<!-- <Link to="about">About</Link>
		<Link to="blog">Blog</Link> -->
	<!-- </nav> -->
	<!-- <div> -->
		<Route path="/">
			<Login />
		</Route>
		<Route path="/home/">
			<Home />
		</Route>
		<Route path="/document/">
			<Document />
		</Route>
		<!-- <Route path="about" component={About} />
		<Route path="blog/*">
			<Route path="/">
				<Blog />
			</Route>
			<Route path=":id" component={BlogPost} />
		</Route>
		<Route path="search/:query" let:params>
			<Search query={params.query} />
		</Route> -->
	<!-- </div> -->
</Router>