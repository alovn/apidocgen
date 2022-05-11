# {{.Title}} {{if .Service}} ({{.Service}}){{end}}

{{- if .Description}}

{{.Description}}
{{- end}}
{{range $k,$v := .Groups}}
{{add $k 1}}. [{{$v.Title}}](./apis-{{$v.Group}}.md)
{{range $k2,$api := $v.Apis}}
    - [{{$api.Title}}](./apis-{{$api.Group}}.md#{{$k2}}-{{$api.Title}}) {{- if $api.Deprecated}}(Deprecated){{end}}
{{end}}{{end}}