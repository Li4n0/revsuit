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
    <div v-if="record.files.length" slot="expandedRowRender" slot-scope="record" style="margin: 0">
      <b v-if="record.files.length" style="color: gray">FILES:</b><br>
      <a v-for="file in record.files" :key="file.name+record.id" :href="'../api/file/mysql/'+file.id"
         target="_blank">{{ file.name }} </a>
    </div>
    <div slot="selectDropdown"
         slot-scope="{ setSelectedKeys, selectedKeys, clearFilters, column }"
         style="padding: 8px">
      <a-checkbox
          :checked="filters[column.dataIndex] === 'true'"
          @change="(e)=>{e.target.checked?filters[column.dataIndex] = 'true':filters[column.dataIndex] = '';fetch()}">
        True
      </a-checkbox>
      <br/>
      <a-checkbox
          :checked="filters[column.dataIndex] === 'false'"
          @change="(e)=>{e.target.checked?filters[column.dataIndex] = 'false':filters[column.dataIndex] = '';fetch()}">
        False
      </a-checkbox>
    </div>
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
    <span slot="loadData" slot-scope="loadData">
       <a-tag v-if="loadData"
              color="#eb2f96"
       >TRUE</a-tag><a-tag v-else color="#f5222d">FALSE</a-tag>
    </span>
    <span slot="fileNum" slot-scope="files">
       <a-tag v-if="files.length>=3"
              color="#722ed1"
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

import {deleteMysqlRecord, getMysqlRecord} from '@/api/record'
import {store} from '@/main'
import FilterDropdown from '@/components/FilterDropdown'

const colors = [
  "#13c2c2",
  "#52c41a",
  "#02a7ff",
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
    title: 'SCHEMA',
    dataIndex: 'schema',
    key: 'schema',
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
      filterDropdown: 'selectDropdown',
      filterIcon: 'filterIcon',
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
      getMysqlRecord(params).then(res => {
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
      deleteMysqlRecord(params).then(() => {
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