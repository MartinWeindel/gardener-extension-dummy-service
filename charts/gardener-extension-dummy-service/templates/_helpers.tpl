{{- define "name" -}}
gardener-extension-dummy-service
{{- end -}}

{{- define "extensionconfig" -}}
---
apiVersion: dummy-service.extensions.config.gardener.cloud/v1alpha1
kind: Configuration
bar: {{ required ".Values.serviceConfig.bar is required" .Values.serviceConfig.bar }}
{{- end }}

{{-  define "image" -}}
  {{- if .Values.skaffoldImage }}
  {{- .Values.skaffoldImage }}
  {{- else }}
    {{- if hasPrefix "sha256:" .Values.image.tag }}
    {{- printf "%s@%s" .Values.image.repository .Values.image.tag }}
    {{- else }}
    {{- printf "%s:%s" .Values.image.repository .Values.image.tag }}
    {{- end }}
  {{- end }}
{{- end }}

{{- define "leaderelectionid" -}}
extension-dummy-service-leader-election
{{- end -}}