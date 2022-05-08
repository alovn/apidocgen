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
{{if $res.Description}}
// {{$res.Description}}{{end}}
// StatusCode: {{$res.StatusCode}}
{{$res.Examples}}
```
{{end}}{{end}}