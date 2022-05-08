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

**Response**:
{{range $res := $v.Responses}}
```json
// StatusCode: {{$res.StatusCode}}
{{if $res.Description}}
// {{$res.Description}}{{end}}
{{$res.Examples}}
```
{{end}}{{end}}