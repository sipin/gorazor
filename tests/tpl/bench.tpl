<html>
	<head><title>test</title></head>
	<body>
		<ul>
		{{range .}}
			{{if .Print}}
				<li>ID={{.ID}}, Message={{.Message}}</li>
			{{end}}
		{{end}}
		</ul>
	</body>
</html>
