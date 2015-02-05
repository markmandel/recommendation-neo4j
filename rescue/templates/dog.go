package templates

//Dog is the template for an individual dog view
const Dog = `
{{template "header" .}}

<div class="container">
	<div class="row">
		<div class="col-sm-12">
			<div class="page-header">
				<h1>{{.dog.Name}}
					<small>{{.dog.Breed.Name}}</small>
				</h1>
			</div>
		</div>
	</div>
	<div class="row">
		<div class="col-sm-3">
			<p>
				<a href="https://www.google.com/search?q=local+dog+adoption" target="_blank"
				   onclick="javascript: window.alert('Unfortunately this is a fake dog adoption site.\nIf you are interested in adopting a dog, please check your local shelters.')"
				   type="button" class="btn btn-success btn-lg btn-block">Adopt This Dog</a>
			</p>
			<p>{{.dog.Breed.Description}}</p>
		</div>
		<div class="col-sm-9">
			<img class="img-responsive" src="/resources/images/{{.dog.ID}}-{{.dog.Name}}.jpg">
		</div>
	</div>

	{{template "disclaimer"}}
</div>
{{define "extraJS"}}{{end}}
{{template "footer"}}
`
