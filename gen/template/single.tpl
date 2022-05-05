# {{.Title}}

{{.Description}}

## apis
{{range $k,$v := .Apis}}
### {{$v.Title}}
{{if $v.Author}}
author: {{$v.Author}}
{{end}}
```text
{{$v.HTTPMethod}} {{$v.FullURL}}
```
{{end}}