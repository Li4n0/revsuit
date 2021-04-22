<template>
  <a-table
      :columns="columns"
      :data-source="data"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      :rowClassName="(record, index) => index % 2 === 0 ? '' : 'gray-table-row'"
  >
    <div
        slot="filterDropdown"
        slot-scope="{ setSelectedKeys, selectedKeys, clearFilters, column }"
        style="padding: 8px"
    >
      <a-input
          :placeholder="`Search ${column.dataIndex}`"
          :value="selectedKeys[0]"
          style="width: 188px; margin-bottom: 8px; display: block;"
          @change="e => setSelectedKeys(e.target.value ? [e.target.value] : [])"
          @pressEnter="() => {
            filters[column.dataIndex] = selectedKeys[0];
            fetch()
          }"
      />
      <a-button
          type="primary"
          icon="search"
          size="small"
          style="width: 90px; margin-right: 8px"
          @click="() => {
            filters[column.dataIndex] = selectedKeys[0];
            fetch()
          }"
      >
        Search
      </a-button>
      <a-button size="small" style="width: 90px" @click="() =>{
        clearFilters();
        delete filters[column.dataIndex];
        fetch()
      }">
        Reset
      </a-button>
    </div>
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

import {getHttpRecord} from '@/api/record'
import {store} from '@/main'

const colors = {
  "GET": "green",
  "POST": "red",
  "HEAD": "pink",
  "PUT": "geekblue",
  "OPTIONS": "cyan",
  "DELETE": "purple",
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
    title: 'PATH',
    dataIndex: 'path',
    key: 'path',
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
      data: [],
      pagination: {current: 1},
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
      this.loading = true;
      let params = {
        ...this.filters,
        page: this.pagination.current,
        order: this.order
      }
      getHttpRecord(params).then(res => {
        let result = res.data.result
        this.data = result.data
        const pagination = {...this.pagination};
        // Read total count from server
        // pagination.total = data.totalCount;
        pagination.total = result.count;
        this.pagination = pagination;
        this.loading = false
      }).catch(e => {
        if (e.response.status === 403) {
          store.authed = false
          return []
        } else {
          console.error(e)
        }
      })
    }
  },
  mounted() {
    this.fetch({page: "1"});
  },
}
</script>