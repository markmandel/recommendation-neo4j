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
			<p>{{.dog.Breed.Description}}</p>
		</div>
		<div class="col-sm-9">
			<img src="{{.dog.PicURL}}"/>
		</div>
	</div>

	{{template "disclaimer"}}
</div>

{{template "footer"}}
`
