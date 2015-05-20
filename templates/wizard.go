package templates

const Wizard =
`	<div class="content">
		<div class="row">
			<div class="columns communities-settings">
				<h1><i class="fi-widget"></i> Settings Wizard</h1>
				{{if .successMsg}}<div class="alert-box success">{{.successMsg}}</div>{{end}}
				{{if .errorMsg}}<div class="alert-box alert">{{.errorMsg}}</div>{{end}}
				<section>{{ template "wizardContent" .}}</section>
			</div>
		</div>
	</div>
	<script src="/static/js/settings.js"></script>
	<script src="/static/js/wizard.js"></script>
`
