<template>
  <a-table
      style="overflow-x: auto;"
      :columns="columns"
      :data-source="data"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      :rowClassName="(record, index) => index % 2 === 0 ? '' : 'gray-table-row'"
  >
    <filter-dropdown
        slot="filterDropdown"
        slot-scope="{ setSelectedKeys, selectedKeys, clearFilters, column }"
        :set-selected-keys="setSelectedKeys"
        :selected-keys="selectedKeys"
        :clear-filters="clearFilters"
        :column="column"
        :filters="filters"
        :fetch="fetch"
    />
    <a-icon
        slot="filterIcon"
        slot-scope="filtered"
        type="search"
        :style="{ color: filtered ? '#108ee9' : undefined }"
    />

    <span slot="time" slot-scope="time">
        {{ new Date(time).format("yyyy-MM-dd hh:mm:ss") }}
    </span>
    <span slot="method" slot-scope="method">
      <a-tag
          :color="colors[method]"
      >
        {{ method.toUpperCase() }}
      </a-tag>
    </span>
    <code slot="expandedRowRender" slot-scope="record" style="margin: 0">
      <b style="color: gray">RAW REQUEST:</b><br>
      <hr/>
      <span style="white-space: pre-line">{{ record.raw_request }}</span>
    </code>
  </a-table>
</template>
<style>
.gray-table-row {
  background-color: #f5f5f5;
}
</style>
<script>

import {getHttpRecord, deleteHttpRecord} from '@/api/record'
import {store} from '@/main'
import FilterDropdown from '@/components/FilterDropdown'

const colors = {
  "GET": "#52c41a",
  "POST": "#f5222d",
  "PUT": "#eb2f96",
  "HEAD": "#02a7ff",
  "OPTIONS": "#13c2c2",
  "DELETE": "#722ed1",
}

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    sorter: true,
    sortDirections: ['descend', 'ascend'],
  },
  {
    title: 'REQUEST TIME',
    dataIndex: 'request_time',
    key: 'request_time',
    scopedSlots: {customRender: 'time'},
  },
  {
    title: 'RULE',
    dataIndex: 'rule_name',
    key: 'rule_name',
    scopedSlots: {
      filterDropdown: 'filterDropdown',
      filterIcon: 'filterIcon',
    },
  },
  {
    title: 'FLAG',
    dataIndex: 'flag',
    key: 'flag',
    scopedSlots: {
      filterDropdown: 'filterDropdown',
      filterIcon: 'filterIcon',
    },
  },
  {
    title: 'METHOD',
    dataIndex: 'method',
    key: 'method',
    scopedSlots: {
      customRender: 'method',
      filterDropdown: 'filterDropdown',
      filterIcon: 'filterIcon',
    },
  },
  {
    title: 'URI',
    dataIndex: 'uri',
    key: 'uri',
    ellipsis: true,
    scopedSlots: {
      filterDropdown: 'filterDropdown',
      filterIcon: 'filterIcon',
    },
  },
  {
    title: 'REMOTE IP',
    key: 'remote_ip',
    dataIndex: 'remote_ip',
    scopedSlots: {
      filterDropdown: 'filterDropdown',
      filterIcon: 'filterIcon',
    },
  },
  {
    title: 'IP AREA',
    key: 'ip_area',
    dataIndex: 'ip_area'
  }
];

export default {
  name: 'HttpLogs',
  data() {
    return {
      store,
      data: [],
      pagination: {
        current: 1, showSizeChanger: true, pageSize: store.pageSize,
        onShowSizeChange: (current, size) => {
          store.pageSize = size
        }
      },
      filters: {},
      order: "desc",
      loading: false,
      columns,
      colors
    };
  },
  methods: {
    handleTableChange(pagination, filters, sorter) {
      const pager = {...this.pagination};
      pager.current = pagination.current;
      this.pagination = pager;
      this.order = sorter.order === "ascend" ? "asc" : "desc"
      this.fetch();
    },
    fetch: function () {
      let params = {
        ...this.filters,
        page: this.pagination.current,
        pageSize: this.pagination.pageSize,
        order: this.order
      }
      this.loading = true;
      getHttpRecord(params).then(res => {
        let result = res.data.result
        this.data = result.data
        const pagination = {...this.pagination};

        pagination.total = result.count;
        this.pagination = pagination;
        this.loading = false
      }).catch(e => {
        if (e.response.status !== 403) {
          this.$message.error('Unknown error with status code: ' + e.response.status)
        }
      })
    },
    delete: function () {
      let params = {
        ...this.filters,
      }
      this.loading = true;
      deleteHttpRecord(params).then(() => {
        this.$message.success('Deleted successfully')
        this.filters = {}
        this.fetch()
      }).catch(e => {
        if (e.response.status !== 403) {
          this.$message.error('Failed to delete: ' + e.response.data.error)
        }
      })
    },
  },
  mounted() {
    this.fetch();
  },
  components: {
    FilterDropdown
  }
}
</script>