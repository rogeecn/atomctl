<template>
  <div>
    <PageHeader title="{{ .ModuleTitle }}" subtitle="{{ .Title }}列表">
      <ActionCreate :to="{ name: '{{ .Model.Name }}Create' }" />
      <ActionImport />
      <ActionExport />
      <ActionRefresh @click="fetchData" />
      <ActionDensity v-model:size="size" />
      <ActionColumn
        :columns="allColumns"
        :hidden="hiddenColumns"
        v-model:clone="showColumns"
      />
    </PageHeader>

    <QueryForm
      class="m-5 pt-5"
      ref="queryForm"
      @search="fetchData"
      :filters="table{{ .Model.Name }}Filters()"
    />

    <Container class="m-5">
      <a-table
        row-key="id"
        :hoverable="true"
        :stripe="true"
        :row-selection="rowSelection"
        v-model:selectedKeys="selectedKeys"
        :loading="loading"
        :pagination="pagination"
        :columns="showColumns"
        :data="renderData"
        :bordered="false"
        :size="size"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #operations="{ record }">
          <RowOperations
            :record="record"
            :reload="fetchData"
            edit="{{ .Model.Name }}Edit"
            view="{{ .Model.Name }}View"
            :params="{ id: record.id }"
            :deleteAction="delete{{ .Model.Name }}Item"
          />
        </template>
      </a-table>
    </Container>
  </div>
</template>

<script lang="ts" setup>
import {
{{ .Model.Name }}Item,
{{ .Model.Name }}ListQuery,
delete{{ .Model.Name }}Item,
query{{ .Model.Name }}List,
table{{ .Model.Name }}Columns,
table{{ .Model.Name }}Filters,
} from "@/api/{{ .Model.Filename }}";
import { Container, PageHeader } from "@/components/layout";
import {
ActionColumn,
ActionCreate,
ActionDensity,
ActionExport,
ActionImport,
ActionRefresh,
QueryForm,
RowOperations,
} from "@/components/table";
import { SizeProps } from "@/components/table/types";
import useDatetime from "@/hooks/datetime";
import useLoading from "@/hooks/loading";
import { Pagination } from "@/types/global";
import {
PaginationProps,
TableColumnData,
TableRowSelection,
} from "@arco-design/web-vue";
import { onMounted, reactive, ref } from "vue";

const { date } = useDatetime();

const size = ref<SizeProps>("large");
const queryForm = ref();

const { columns, hidden } = table{{ .Model.Name }}Columns();
const allColumns = ref<TableColumnData[]>(columns ?? []);
const hiddenColumns = ref<string[]>(hidden ?? []);
const showColumns = ref<TableColumnData[]>([]);

// load columns
onMounted(() => {
  fetchData(basePagination);
});

// fetch table data
const { loading, setLoading } = useLoading(true);
const renderData = ref<{{ .Model.Name }}Item[]>([]);

const basePagination: Pagination = { current: 1, pageSize: 20 };
const pagination = reactive<PaginationProps>({
  ...basePagination,
  showPageSize: true,
  showJumper: true,
  showTotal: true,
  simple: false,
  hideOnSinglePage: true,
  showMore: true,
});

const fetchData = async (pg: Pagination = basePagination) => {
  setLoading(true);
  try {
    const pages = { page: pg.current, limit: pg.pageSize };
    const params = ({
      ...pages,
      ...queryForm.value?.filterItems(),
    } as unknown) as {{ .Model.Name }}ListQuery;
    console.log(params);

    const { data } = await query{{ .Model.Name }}List(params);
    renderData.value = data.items;
    pagination.total = data.total;
    // pagination.current = data.page;
    // pagination.pageSize = data.limit;
  } catch (err) {
    // you can report use errorHandler or other
  } finally {
    setLoading(false);
  }
};

const onPageChange = (current: number) => {
  pagination.current = current;
  fetchData({ current: current, pageSize: pagination.pageSize });
};

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  fetchData({ current: pagination.current, pageSize: pageSize });
};

// row selection
const selectedKeys = ref([]);

const rowSelection = reactive<TableRowSelection>({
  type: "checkbox",
  showCheckedAll: true,
  onlyCurrent: false,
});
</script>
@/components/table/table@/components/table
