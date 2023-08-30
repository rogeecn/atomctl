<template>
  <div>
    <PageHeader title="INPUT_TITLE" subtitle="编辑" :back="true"/>

    <Container :loading="loading" :rows="3" class="pt-5">
      <a-form :model="form" @submit="handleSubmit" class="md:w-3/4 sm:w-full">
        <!-- form start -->

      {{- range .Model.Fields }}
        {{- if or (eq .Name "ID") (eq .Name "CreatedAt") (eq .Name "UpdatedAt") (eq .Name "DeletedAt") }}
        {{- else }}
          {{- if or (eq .Type "int")  (eq .Type "int8")  (eq .Type "int16")  (eq .Type "int32")  (eq .Type "int64")  (eq .Type "uint")  (eq .Type "uint8")  (eq .Type "uint16")  (eq .Type "uint32")  (eq .Type "uint64") }}
          <a-form-item field="{{ .Tag }}" label="{{ .Comment }}">
            <a-number v-model="form.{{ .Tag }}" placeholder="请输入{{ .Comment }}" />
          </a-form-item>
          {{- else if or (eq .Type "bool") }}
          <a-form-item field="{{ .Tag }}" label="{{ .Comment }}">
            <a-switch type="round" v-model="form.{{ .Tag }}" checked-color="#14C9C9" unchecked-color="#F53F3F" />
          </a-form-item>
          {{- else }}
          <a-form-item field="{{ .Tag }}" label="{{ .Comment }}">
            <a-input v-model="form.{{ .Tag }}" placeholder="请输入{{ .Comment }}" />
          </a-form-item>
          {{- end }}
        {{- end }}
      {{- end }}

        <a-form-item>
          <a-button type="primary" html-type="submit" size="large" :loading="submitting">提交</a-button>
        </a-form-item>
        <!-- form end -->
      </a-form>
    </Container>
  </div>
</template>

<script lang="ts" setup>
import { {{ .Model.Name }}Item, get{{ .Model.Name }}Item, update{{ .Model.Name }}Item } from '@/api/{{ .Model.Filename }}';
import { Container, PageHeader } from "@/components/layout";
import useLoading from '@/hooks/loading';
import { Message } from '@arco-design/web-vue';
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
const { loading, setLoading } = useLoading()

onMounted(() => {
  fetchData()
})
const form = ref<{{ .Model.Name }}Item>({});

const fetchData = async () => {
  try {
    setLoading(true);
    const { data } = await get{{ .Model.Name }}Item(Number(route.params.id))
    form.value = data
  } catch (e: any) {
    Message.error(e.message)
  } finally {
    setLoading(false)
  }
}

const { loading: submitting, setLoading: setSubmitting } = useLoading()
// form
const handleSubmit = async ({ values, errors }: any) => {
  console.log(values, errors);
  try {
    setSubmitting(true);
    await update{{ .Model.Name }}Item(Number(route.params.id), values)

    Message.success("更新成功")
    router.back()
  } catch (e: any) {
    // Message.error(e.message)
  } finally {
    setSubmitting(false)
  }
};
</script>