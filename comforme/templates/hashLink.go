package templates

const HashLink = `	<div class="content">
		<div class="row">
			<div class="columns communities-settings">
				<h1><i class="fi-widget"></i> {{.pageTitle}}</h1>
				{{if .successMsg}}<div class="alert-box success">{{.successMsg}}</div>{{end}}
				{{if .errorMsg}}<div class="alert-box alert">{{.errorMsg}}</div>{{end}}
				<section>{{ template "wizardContent" .}}</section>
			</div>
		</div>
	</div>
`
