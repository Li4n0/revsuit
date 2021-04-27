<template>
  <a-table
      :columns="columns"
      :data-source="data"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      :rowClassName="(record, index) => index % 2 === 0 ? '' : 'gray-table-row'"
  >

    <div slot="selectDropdown"
         slot-scope="{ setSelectedKeys, selectedKeys, clearFilters, column }"
         style="padding: 8px">
      <a-checkbox
          :checked="filters[column.dataIndex] === 'CRASHED'"
          @change="(e)=>{e.target.checked?filters[column.dataIndex] = 'CRASHED':filters[column.dataIndex] = '';fetch()}">
        CRASHED
      </a-checkbox>
      <br/>
      <a-checkbox
          :checked="filters[column.dataIndex] === 'FINISHED'"
          @change="(e)=>{e.target.checked?filters[column.dataIndex] = 'FINISHED':filters[column.dataIndex] = '';fetch()}">
        FINISHED
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
    <span slot="method" slot-scope="method">
     <a-tag v-if="method" :color="method==='UPLOAD'?'#eb2f96':'#02a7ff'">{{ method }}</a-tag>
    </span>
    <span slot="file" slot-scope="file">
     <a-tag
         v-if="file.id"
         color="#f5222d"
     ><a target="_blank" :href="'/revsuit/api/file/ftp/'+file.id">TRUE</a> </a-tag>
     <a-tag v-else color="#722ed1">
      FALSE
      </a-tag>
    </span>

    <span slot="time" slot-scope="time">
        {{ new Date(time).format("yyyy-MM-dd hh:mm:ss") }}
    </span>
    <span slot="status" slot-scope="status">
      <a-tag
          :color="colors[status]"
      >
        {{ status }}
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

import {getFtpRecord} from '@/api/record'
import {store} from '@/main'
import FilterDropdown from '@/components/FilterDropdown'

const colors = {
  "CRASHED": "#f50",
  "FINISHED": "#52c41a"
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
    title: 'USER',
    dataIndex: 'user',
    key: 'user',
    scopedSlots: {
      filterDropdown: 'filterDropdown',
      filterIcon: 'filterIcon',
    },
  },
  {
    title: 'PASSWORD',
    dataIndex: 'password',
    key: 'password',
    scopedSlots: {
      filterDropdown: 'filterDropdown',
      filterIcon: 'filterIcon',
    },
  },
  {
    title: 'PATH',
    dataIndex: 'path',
    key: 'path',
    ellipsis: true,
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
      filterDropdown: 'selectDropdown',
      filterIcon: 'filterIcon',
    },
  },
  {
    title: 'FILE',
    dataIndex: 'file',
    key: 'file',
    scopedSlots: {
      customRender: 'file'
    },
  },
  {
    title: 'STATUS',
    dataIndex: 'status',
    key: 'status',
    scopedSlots: {
      customRender: 'status',
      filterDropdown: 'selectDropdown',
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
  name: 'FtpLogs',
  data() {
    return {
      store,
      data: [],
      pagination: {current: 1},
      filters: {},
      order: "desc",
      loading: false,
      columns,
      colors: colors
    };
  },
  methods: {
    aa(file) {
      console.log(file)
    },
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
        order: this.order
      }
      this.loading = true;
      getFtpRecord(params).then(res => {
        let result = res.data.result
        this.data = result.data
        const pagination = {...this.pagination};

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
    },
  },
  mounted() {
    this.fetch();
  },
  watch: {
    'store.authed'() {
      this.fetch()
    }
  },
  components: {
    FilterDropdown
  }
}
</script>