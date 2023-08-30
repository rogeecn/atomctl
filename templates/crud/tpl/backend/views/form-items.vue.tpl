<template>
      {{- range .Model.Fields }}
        {{- if or (eq .Name "ID") (eq .Name "CreatedAt") (eq .Name "UpdatedAt") (eq .Name "DeletedAt") }}
        {{- else }}
          {{- if or (eq .Type "int")  (eq .Type "int8")  (eq .Type "int16")  (eq .Type "int32")  (eq .Type "int64")  (eq .Type "uint")  (eq .Type "uint8")  (eq .Type "uint16")  (eq .Type "uint32")  (eq .Type "uint64") }}
  <a-form-item field="{{ .Tag }}" label="{{ .Comment }}">
    <a-input-number v-model="form.{{ .Tag }}" placeholder="请输入{{ .Comment }}" />
  </a-form-item>
          {{ else if or (eq .Type "bool") }}
  <a-form-item field="{{ .Tag }}" label="{{ .Comment }}">
    <a-switch type="round" v-model="form.{{ .Tag }}" checked-color="#14C9C9" unchecked-color="#F53F3F" />
  </a-form-item>
          {{ else }}
  <a-form-item field="{{ .Tag }}" label="{{ .Comment }}">
    <a-input v-model="form.{{ .Tag }}" placeholder="请输入{{ .Comment }}" />
  </a-form-item>
          {{- end }}
        {{- end }}
      {{- end }}
</template>

<script lang="ts" setup>
import { {{ .Model.Name }}Item } from '@/api/{{ .Model.Filename }}';

const props = defineProps<{
  form: {{ .Model.Name }}Item;
}>();
</script>