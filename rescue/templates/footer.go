package templates

//Footer is the footer template for all pages.
const Footer = `
<script src="/resources/js/jquery.min.js"></script>
<script src="/resources/js/bootstrap.min.js"></script>
<script src="/resources/js/jquery.matchHeight-min.js"></script>
{{template "extraJS"}}
</body>
</html>
`
