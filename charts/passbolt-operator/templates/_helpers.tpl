{{/*
Expand the name of the chart.
*/}}
{{- define "passbolt-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "passbolt-operator.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "passbolt-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "passbolt-operator.labels" -}}
helm.sh/chart: {{ include "passbolt-operator.chart" . }}
{{ include "passbolt-operator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "passbolt-operator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "passbolt-operator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: control-plane
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "passbolt-operator.serviceAccountName" -}}
{{- include "passbolt-operator.fullname" . }}
{{- end }}

{{- define "passbolt-operator.role.manager" -}}
{{- printf "%s-manager-role" (include "passbolt-operator.fullname" .) | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "passbolt-operator.role.leaderelection" -}}
{{- printf "%s-leader-election" (include "passbolt-operator.fullname" .) | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "passbolt-operator.rolebinding.manager" -}}
{{- printf "%s-manager" (include "passbolt-operator.fullname" .) | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "passbolt-operator.rolebinding.leaderelection" -}}
{{- printf "%s-leader-election" (include "passbolt-operator.fullname" .) | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "passbolt-operator.secret.name" -}}
{{- if .Values.secret.name }}
{{- .Values.secret.name }}
{{- else }}
{{- printf "%s-" (include "passbolt-operator.fullname" .) | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}