{{- define "group_index" -}}
---
title: {{.Title}}
weight: 1
---

## {{.Title}} ({{.Service}})

{{- if .Version}}

version: _@{{.Version}}_
{{- end}}

{{- if .Description}}

{{.Description}}
{{- end}}
{{- range $k,$v := .Groups}}

{{add $k 1}}. [{{$v.Title}}](./apis-{{$v.Group}})
{{- range $k2,$api := $v.Apis}}

    {{add $k 1}}.{{add $k2 1}}. [{{$api.Title}}](./apis-{{$api.Group}}#{{add $k2 1}}-{{$api.Title}}) {{- if $api.Deprecated}}(Deprecated){{end}}
{{- end }}
{{- end }}
{{ template "footer" }}
{{ end }}