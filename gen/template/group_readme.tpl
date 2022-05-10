# {{.Title}} {{if .Service}} ({{.Service}}){{end}}

{{.Description}}

## api-groups
{{range $k,$v := .Groups}}
{{add $k 1}}. [{{$v.Title}}](./apis-{{$v.Group}}.md)
{{range $api := $v.Apis}}
    - [{{$api.Title}}](./apis-{{$api.Group}}.md#{{$api.Title}})
{{end}}{{end}}