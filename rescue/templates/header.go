package templates

//Header is the header template for all pages.
const Header = `
<!DOCTYPE html>
<html lang="en">
<head>
	<!-- Meta, title, CSS, favicons, etc. -->
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1">

	<title>Adopt A Dog :: {{.title}}</title>
	<link href="/resources/css/bootstrap.min.css" rel="stylesheet">

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
`
