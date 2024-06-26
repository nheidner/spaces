{{- define "server.name" -}}
{{ default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }} 
{{- end }}

{{- define "server.dbHost" -}}
{{- if .Values.dbHost }}
{{- .Values.dbHost }}
{{- else }}
{{- .Release.Name }}-postgresql.{{ .Release.Namespace }}.svc
{{- end }}
{{- end }}

{{- define "server.host" -}}
{{- if eq .Release.Namespace "prod" }}
{{- .Values.hosts.prod }}
{{- else if eq .Release.Namespace "staging" }}
{{- .Values.hosts.staging }}
{{- else }}
{{- .Values.hosts.dev }}
{{- end }}
{{- end }}

{{- define "server.chart" -}}
{{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "server.selectorLabels" -}}
app.kubernetes.io/name: {{ include "server.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "server.labels" -}}
helm.sh/chart: {{ include "server.chart" . }}
{{ include "server.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

