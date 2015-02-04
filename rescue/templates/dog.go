package templates

//Dog is the template for an individual dog view
const Dog = `
{{template "header"}}

<div class="container">
	<div class="row">
		 <div class="col-sm-12">
		 	<div class="page-header">
  				<h1>{{.dog.Name}} <small>{{.dog.Breed.Name}}</small></h1>
			</div>
		 </div>
	</div>

	<hr>

	<footer>
		<p>All dog breed names, descriptions and pictures are licenced under Creative Commons. Attribution can be found
		<a href="https://github.com/markmandel/recommendation-neo4j/blob/master/Attribution.md" target="_blank">here</a>.
		</p>
	</footer>
</div>

{{template "footer"}}
`
