# {{.Title}}

{{.Description}}

## apis
{{range $k,$v := .Apis}}
### {{$v.Title}}
{{if $v.Author}}
author: _{{$v.Author}}_
{{end}}
```text
{{$v.HTTPMethod}} {{$v.FullURL}}
```

{{if $v.Requests.Parameters}}**Request**:

parameters|type|required|validate|example|description
--|--|--|--|--|--{{range $p:= $v.Requests.Parameters}}
**{{$p.Name}}**|_{{$p.Types}}_|{{$p.Required}}|{{$p.Validate}}|{{$p.Example}}|{{$p.Description}}{{end}}{{end}}

{{if $v.Responses}}**Response**:
{{range $res := $v.Responses}}
```json
// StatusCode: {{$res.StatusCode}} {{$res.Description}}
{{$res.Examples}}
```
{{end}}{{end}}{{end}}