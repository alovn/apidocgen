# {{.Title}}

{{.Description}}

## Apis
{{range $k,$v := .Apis}}
### {{$v.Title}}
{{if $v.Author}}
author: {{$v.Author}}
{{end}}
```text
{{$v.HTTPMethod}} {{$v.Api}}
```

{{range $res := $v.Responses}}
Response:
```json
{{if $res.Description}}
// {{$res.Description}}{{end}}
// HTTP StatusCode: {{$res.StatusCode}}
{{$res.Examples}}
```
{{end}}
{{end}}