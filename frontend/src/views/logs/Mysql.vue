<template>
  <a-table
      :columns="columns"
      :data-source="data"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      :rowClassName="(record, index) => index % 2 === 0 ? '' : 'gray-table-row'"
  >
    <div v-if="record.files.length" slot="expandedRowRender" slot-scope="record" style="margin: 0">
      <b v-if="record.files.length" style="color: gray">FILES:</b><br>
      <a v-for="file in record.files" :key="file.name+record.id" :href="'/revsuit/api/file/mysql/'+file.id" target="_blank">{{ file.name }} </a>
    </div>
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
    <span slot="loadData" slot-scope="loadData">
       <a-tag v-if="loadData"
              color="green"
       >True</a-tag><a-tag v-else color="red">False</a-tag>
    </span>
    <span slot="fileNum" slot-scope="files">
       <a-tag v-if="files.length>=3"
              color="purple"
       >{{ files.length }}</a-tag>
      <a-tag v-else :color="colors[files.length]">
        {{ files.length }}
      </a-tag>
    </span>
  </a-table>
</template>
<style>
.gray-table-row {
  background-color: #f5f5f5;
}
</style>
<script>

import {getMysqlRecord} from '@/api/record'
import {store} from '@/main'

const colors = [
  "geekblue",
  "blue",
  "pink",
]

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
    title: 'USER',
    dataIndex: 'username',
    key: 'username',
    scopedSlots: {
      filterDropdown: 'filterDropdown',
      filterIcon: 'filterIcon',
    },
  },
  {
    title: 'LOAD DATA',
    dataIndex: 'load_local_data',
    key: 'load_local_data',
    scopedSlots: {
      customRender: "loadData",
    }
  },
  {
    title: 'FILE NUM',
    dataIndex: 'files',
    key: 'files',
    scopedSlots: {
      customRender: "fileNum",
    }
  },
  {
    title: 'CLIENT NAME',
    dataIndex: 'client_name',
    key: 'client_name',
    scopedSlots: {
      filterDropdown: 'filterDropdown',
      filterIcon: 'filterIcon',
    },
  },
  {
    title: 'CLIENT OS',
    dataIndex: 'client_os',
    key: 'client_os',
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
  name: 'MysqlLogs',
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
      getMysqlRecord(params).then(res => {
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