package templates

//Index is the index page template
const Index = `
{{template "header" .}}

<!-- Main jumbotron for a primary marketing message or call to action -->
<div class="jumbotron">
	<div class="container">
		<h1>Adopt A Dog</h1>

		<p>Please have a look at all the wonderful dogs that are we have for adoption.<br/>Click on any of the pictures to
			get
			more details.</p>

		{{if .recommendations}}
		<div class="row" id="lookedat">
			<div class="col-sm-12">
				<h4>Some dogs we thought you might like...</h4>
			</div>
			{{range .recommendations}}
			<div class="col-xs-4 col-md-2 recommendation">
				<a href="/dog/{{.ID}}" class="thumbnail" title="{{.Name}} - {{.Breed.Name}}">
					<img src="/resources/images/{{.ID}}-{{.Name}}.jpg" alt="{{.Name}} - {{.Breed.Name}}">
				</a>
			</div>
			{{end}}
		</div>
		{{end}}

	</div>
</div>

<div class="container">
	<div class="row">
		{{range .dogs}}
		<div class="col-sm-6 col-md-4 dog" id="{{.Name}}">
			<a href="/dog/{{.ID}}">
				<div class="thumbnail">
					<img src="/resources/images/{{.ID}}-{{.Name}}.jpg">

					<div class="caption">
						<h3>{{.Name}}</h3>

						<p>{{.Breed.Name}}</p>
					</div>
				</div>
			</a>
		</div>
		{{end}}
	</div>

	{{template "disclaimer"}}
</div>
<!-- /container -->

{{define "extraJS"}}
<script type="text/javascript">
	$(function() {
		$('.dog').matchHeight();
		$('.recommendation').matchHeight();
	});
</script>
{{end}}

{{template "footer"}}
`
