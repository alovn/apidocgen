# {{.Title}} {{if .Service}} ({{.Service}}){{end}}

{{- if .Version}}

version: _@{{.Version}}_
{{- end}}

{{- if .Description}}

{{.Description}}
{{- end}}
{{range $k,$v := .Groups}}
{{add $k 1}}. [{{$v.Title}}](./apis-{{$v.Group}}.md)
{{range $k2,$api := $v.Apis}}
    {{add $k 1}}.{{add $k2 1}}. [{{$api.Title}}](./apis-{{$v.Group}}.md#{{add $k2 1}}-{{$api.Title}}) {{- if $api.Deprecated}}(Deprecated){{end}}
{{end}}{{end}}