package templates

//Index is the index page template
const Index = `
{{template "header"}}

<!-- Main jumbotron for a primary marketing message or call to action -->
<div class="jumbotron">
	<div class="container">
		<h1>Adopt A Dog</h1>

		<p>Please have a look at all the wonderful dogs that are we have for adoption.<br/>Click on any of the pictures to
			get
			more detals.</p>
	</div>
</div>

<div class="container">
	<div class="row">
		{{range .dogs}}
		<div class="col-sm-6 col-md-4 dog">
			<a href="/dog/{{.ID}}">
			<div class="thumbnail">
				<img src="{{.PicURL}}">

				<div class="caption">
					<h3>{{.Name}}</h3>

					<p>{{.Breed.Name}}</p>
				</div>
			</div>
			</a>
		</div>
		{{end}}
	</div>

	<hr>

	<footer>
		<p>All dog breed names, descriptions and pictures are licenced under Creative Commons. Attribution can be found
		<a href="https://github.com/markmandel/recommendation-neo4j/blob/master/Attribution.md" target="_blank">here</a>.
		</p>
	</footer>
</div>
<!-- /container -->

{{template "footer"}}
`
