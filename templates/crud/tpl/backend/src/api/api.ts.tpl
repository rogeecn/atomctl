import { Columns, Filter, FilterType, Pagination, PaginationResp } from '@/types/global';
import axios from 'axios';


export interface {{ .Model.Name }}ListQuery extends Pagination {
	{{- range .Model.Fields }}
    {{- if or (eq .Name "CreatedAt") (eq .Name "UpdatedAt") (eq .Name "DeletedAt") }}
      {{ .Tag }}?: Date; // {{ .Comment }}
    {{- else }}
      {{- if or (eq .Type "int")  (eq .Type "int8")  (eq .Type "int16")  (eq .Type "int32")  (eq .Type "int64")  (eq .Type "uint")  (eq .Type "uint8")  (eq .Type "uint16")  (eq .Type "uint32")  (eq .Type "uint64") }}
        {{ .Tag }}?: number; // {{ .Comment }}
      {{- else if or (eq .Type "bool") }}
        {{ .Tag }}?: boolean; // {{ .Comment }}
      {{- else }}
        {{ .Tag }}?: string; // {{ .Comment }}
      {{- end }}
    {{- end }}
	{{- end }}
}

export interface {{ .Model.Name }}Item {
	{{- range .Model.Fields }}
    {{- if or (eq .Name "CreatedAt") (eq .Name "UpdatedAt") (eq .Name "DeletedAt") }}
      {{- if eq .Name "DeletedAt" }}
      {{- else }}
      {{ .Tag }}?: Date; // {{ .Comment }}
      {{- end }}
    {{- else }}
      {{- if or (eq .Type "int")  (eq .Type "int8")  (eq .Type "int16")  (eq .Type "int32")  (eq .Type "int64")  (eq .Type "uint")  (eq .Type "uint8")  (eq .Type "uint16")  (eq .Type "uint32")  (eq .Type "uint64") }}
        {{ .Tag }}?: number; // {{ .Comment }}
      {{- else if or (eq .Type "bool") }}
        {{ .Tag }}?: boolean; // {{ .Comment }}
      {{- else }}
        {{ .Tag }}?: string; // {{ .Comment }}
      {{- end }}
    {{- end }}
	{{- end }}
}

export const table{{ .Model.Name }}Filters = ():Filter[] => {
  return [
    { type: FilterType.List, name: "ids", label: "ID" },
   	{{- range .Model.Fields }}
    {{- if or (eq .Name "CreatedAt") (eq .Name "UpdatedAt") (eq .Name "DeletedAt") }}
      { type: FilterType.Date, name: '{{ .Tag }}', label: '{{ .Comment }}' },
    {{- else if or (eq .Name "ID") }}
    {{- else }}
      {{- if or (eq .Type "int")  (eq .Type "int8")  (eq .Type "int16")  (eq .Type "int32")  (eq .Type "int64")  (eq .Type "uint")  (eq .Type "uint8")  (eq .Type "uint16")  (eq .Type "uint32")  (eq .Type "uint64") }}
      { type: FilterType.Number, name: '{{ .Tag }}', label: '{{ .Comment }}' },
      {{- else if or (eq .Type "bool") }}
        { type: FilterType.Bool, name: '{{ .Tag }}', label: '{{ .Comment }}' },
      {{- else }}
        { type: FilterType.String, name: '{{ .Tag }}', label: '{{ .Comment }}' },
      {{- end }}
    {{- end }}
	{{- end }} 
  ]
}

export const table{{ .Model.Name }}Columns = ():Columns => {
  return {
    columns: [
      {{- range .Model.Fields }}
        { title: '{{ .Comment }}', dataIndex: '{{ .Tag }}', slotName: '{{ .Tag }}' },
      {{- end }} 
      { title: '操作', dataIndex: 'operations' ,slotName: 'operations', align:'right' },
    ],
    hidden: [
      'uuid', 'created_at', 'updated_at', 'deleted_at'
    ],
  }
}

export const table{{ .Model.Name }}Labels = (): Record<string,string> => {
  return {
  {{- range .Model.Fields }}
    '{{ .Tag }}': '{{ .Comment }}',
	{{- end }} 
  }
}

export function query{{ .Model.Name }}List(params: {{ .Model.Name }}ListQuery) {
  return axios.get<PaginationResp<{{ .Model.Name }}Item>>('/{{ .Model.RouteName }}', { params });
}

export function create{{ .Model.Name }}Item(data: {{ .Model.Name }}Item) {
  return axios.post(`/{{ .Model.RouteName }}`, data);
}

export function update{{ .Model.Name }}Item(id: number, data: {{ .Model.Name }}Item) {
  return axios.put(`/{{ .Model.RouteName }}/${id}`, data);
}

export function get{{ .Model.Name }}Item(id: number) {
  return axios.get(`/{{ .Model.RouteName }}/${id}`);
}

export function delete{{ .Model.Name }}Item(id: number) {
  return axios.delete(`/{{ .Model.RouteName }}/${id}`);
}