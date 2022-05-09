# {{.Title}}

{{.Description}}

## Apis
{{range $k,$v := .Apis}}
### @api {{$v.Title}}
{{if $v.Author}}
author: {{$v.Author}}
{{end}}
```text
{{$v.HTTPMethod}} {{$v.Api}}
```
{{if $v.Requests.Parameters}}
**Request**:
parameters|type|required|validate|example|description
--|--|--|--|--|--{{range $p:= $v.Requests.Parameters}}
**{{$p.Name}}**|*{{$p.Types}}*|{{$p.Required}}|{{$p.Validate}}|{{$p.Example}}|{{$p.Description}}{{end}}{{end}}
**Response**:
{{range $res := $v.Responses}}
```json
// StatusCode: {{$res.StatusCode}}
{{if $res.Description}}
// {{$res.Description}}{{end}}
{{$res.Examples}}
```
{{end}}{{end}}