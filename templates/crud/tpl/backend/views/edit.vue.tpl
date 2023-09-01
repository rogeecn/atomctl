<template>
  <div>
    <PageHeader title="{{ .Vars.moduleTitle }}" subtitle="{{ .Vars.title }}编辑" :back="true"/>

    <Container :loading="loading" :rows="3" class="pt-5">
      <a-form :model="form" @submit="handleSubmit" class="md:w-3/4 sm:w-full">
        <FormItems :form="form" />
        <a-form-item>
          <a-button type="primary" html-type="submit" size="large" :loading="submitting">提交</a-button>
        </a-form-item>
      </a-form>
    </Container>
  </div>
</template>

<script lang="ts" setup>
import { {{ .Model.Name }}Item, get{{ .Model.Name }}Item, update{{ .Model.Name }}Item } from '@/api/{{.Vars.module}}/{{ .Model.Filename }}';
import { Container, PageHeader } from "@/components/layout";
import useLoading from '@/hooks/loading';
import { Message } from '@arco-design/web-vue';
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import FormItems from './form-items.vue';

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