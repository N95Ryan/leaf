{{- /*
Custom gotestsum template with emojis for Leaf project
Usage: gotestsum --format testdox --format-hide-empty-pkg
*/ -}}

{{- range .}}
{{- if eq .Action "pass"}}
✅ {{.Package}} {{.Test}} ({{.Elapsed}})
{{- else if eq .Action "fail"}}
❌ {{.Package}} {{.Test}} ({{.Elapsed}})
{{if .Output}}{{.Output}}{{end}}
{{- else if eq .Action "skip"}}
⏭️  {{.Package}} {{.Test}}
{{- end}}
{{- end}}
