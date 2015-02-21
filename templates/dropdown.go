package templates

const Dropdown = `<select name="{{.name}}">
{{$selected := .selected}}
{{range $id, $option := $.options}}
    <option value="{{$id}}"{{if eq $id string($selected)}} selected=""{{end}}>{{$option}}</option>
{{end}}
</select>`
