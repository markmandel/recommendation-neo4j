package templates

//Index is the index page template
const Index = `
<!DOCTYPE html>
<html lang="en">
<head>
	<!-- Meta, title, CSS, favicons, etc. -->
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1">

	<title>Adopt A Dog :: {{.title}}</title>
	<link href="//maxcdn.bootstrapcdn.com/bootswatch/3.3.2/cerulean/bootstrap.min.css" rel="stylesheet">

	<style type="text/css">
		/* Move down content because we have a fixed navbar that is 50px tall */
		body {
			padding-top: 50px;
			padding-bottom: 20px;
		}
	</style>
</head>
<body>

<nav class="navbar navbar-inverse navbar-fixed-top">
	<div class="container">
		<div class="navbar-header">
			<a class="navbar-brand" href="/">Home</a>
		</div>

	</div>
</nav>

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
			<div class="thumbnail">
				<img src="{{.PicURL}}">

				<div class="caption">
					<h3>{{.Name}}</h3>

					<p>{{.Breed.Name}}</p>
				</div>
			</div>
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

<script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
<script src="//maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js"></script>
<script src="http://brm.io/js/libs/matchHeight/jquery.matchHeight-min.js"></script>
<script type="text/javascript">
	$(function() {
		$('.dog').matchHeight();
	});
</script>
</body>
</html>
`
