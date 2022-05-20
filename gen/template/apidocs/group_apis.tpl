{{- define "group_apis" -}}
---
title: {{.Title}}
navtitle: {{.Title}}
weight: 3
---


## {{.Title}}
{{- if .Description}}

{{.Description}}
{{- end}}
{{range $k,$api := $.Apis}}
{{add $k 1}}. [{{$api.Title}}](#{{add $k 1}}-{{$api.Title}}) {{- if $api.Deprecated}}(Deprecated){{end}}
{{- end}}

## apis
{{- range $k,$v := .Apis}}

### {{add $k 1}}. {{$v.Title}}

{{- if $v.Deprecated}}

___Deprecated___
{{- end}}

{{- if $v.Description}}

{{$v.Description}}
{{- end}}

{{- if $v.Author}}

author: _{{$v.Author}}_
{{- end}}

{{- if $v.Version}}

version: _{{$v.Version}}_
{{- end}}

```text
{{$v.HTTPMethod}} {{$v.FullURL}}
```
{{- if $v.Requests.Parameters}}

__Request__:

parameter|parameterType|dataType|required|validate|example|description
--|--|--|--|--|--|--
{{- range $p:= $v.Requests.Parameters}}
__{{$p.Name}}__|_{{$p.ParameterTypes}}_|{{$p.DataType}}|{{$p.Required}}|{{$p.Validate}}|{{$p.Example}}|{{$p.Description}}
{{- end}}
{{- end}}
{{- if $v.Requests.Body}}

_body_:

```{{if eq $v.Accept "json"}}javascript{{else}}{{$v.Accept}}{{end}}
{{$v.Requests.Body}}
```
{{- end}}
{{- if $v.Responses}}

__Response__:
{{- range $res := $v.Responses}}

```{{if eq $v.Format "json"}}javascript{{else}}{{$v.Format}}{{end}}
//StatusCode: {{$res.StatusCode}} {{$res.Description}}
{{$res.Body}}
```
{{- end}}

---
{{- end}}
{{- end}}
{{ template "footer" }}
{{ end }}