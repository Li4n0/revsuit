<template xmlns:a-col="http://www.w3.org/1999/html">
  <div>
    <a-button id="add-rule" type="primary" @click="addRule">
      <a-icon type="plus"/>
      New Rule
    </a-button>
    <!--    rule form-->
    <a-drawer
        :title="formAction+ ' FTP rule'"
        :width="490"
        :visible="formVisible"
        :body-style="{ paddingBottom: '80px' }"
        @close="closeDrawer"
    >
      <a-form-model :model="form" ref="form" layout="vertical" @submit="handleSubmit">
        <BasicRule :form="form" :readOnly="formReadOnly"/>
        <a-row :gutter="24">
          <a-col :span="24">
            <a-form-model-item>
              <span slot="label">
                  Pasv Address
                  <a-tooltip
                      title="1. Support template such as ${user}/${password}/${varname}.
2.For setting the rebind mode, please use ',' to separate addresses.">
                    <a-icon type="question-circle-o"/>
                  </a-tooltip>
              </span>
              <a-input
                  v-model="form.pasv_address"
                  style="width: 100%"
                  placeholder="Use external IP by default"
                  :readOnly="formReadOnly"
              />
            </a-form-model-item>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="24">
            <a-form-model-item>
              <span slot="label">
                  Data
                  <a-tooltip
                      title="1. The data to be returned when the client executes the download request.
                      2. Please use base64 encoding.
                      3. Support template such as ${user}/${password}/${varname}">
                    <a-icon type="question-circle-o"/>
                  </a-tooltip>
              </span>
              <a-textarea v-model="form.data"
                          placeholder="The data to be returned when the client executes the download request. Please use base64 encoding."
                          :readOnly="formReadOnly"
                          :auto-size="{ minRows: 10, maxRows: 30 }"
              />
            </a-form-model-item>
          </a-col>
        </a-row>
      </a-form-model>
      <div
          :style="{
          position: 'absolute',
          right: 0,
          bottom: 0,
          width: '100%',
          borderTop: '1px solid #e9e9e9',
          padding: '10px 16px',
          background: '#fff',
          textAlign: 'right',
          zIndex: 1,
        }"
      >
        <a-button :style="{ marginRight: '8px' }" @click="handleCancel">
          Cancel
        </a-button>
        <a-button type="primary" :disabled="formReadOnly" @click="handleSubmit">
          Submit
        </a-button>
      </div>
    </a-drawer>
    <!--    rule table -->
    <a-table
        :columns="columns"
        :data-source="data"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
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
            @pressEnter="() => {filters[column.dataIndex] = selectedKeys[0];fetch()}"
        />
        <a-button
            type="primary"
            icon="search"
            size="small"
            style="width: 90px; margin-right: 8px"
            @click="() => {filters[column.dataIndex] = selectedKeys[0];fetch()}"
        >
          Search
        </a-button>
      </div>
      <a-icon
          slot="filterIcon"
          slot-scope="filtered"
          type="search"
          :style="{ color: filtered ? '#108ee9' : undefined }"
      />
      <span slot="rank" slot-scope="rank">
      <a-tag
          :color="'#'+(0x2db7f5+rank*80).toString(16)"
      >
        {{ rank }}
      </a-tag>
      </span>

      <span slot="switchRender" slot-scope="checked,record,index,dataIndex">
        <a-switch :checked="checked" @click="clickSwitch(record,dataIndex.dataIndex)"></a-switch>
      </span>
      <span slot="valueRender" slot-scope="values">
        <span v-for="value in values.split(',')" :key="value">{{ value }}<br/></span>
      </span>
      <span slot="action" slot-scope="text,record,index">
        <a-button @click="viewRule(record)" style="
                color: #67C23A;
            background-color: transparent;
            border-color: #67C23A;
            text-shadow: none;
            margin:0 10px 3px 0;
        " size="small" ghost>View</a-button>
        <a-button @click="editRule(record,index)" style="
        color: #909399;
    background-color: transparent;
    border-color: #909399;
    text-shadow: none;
    margin:0 10px 3px 0;
" size="small" ghost>Edit</a-button>
        <a-popconfirm
            title="Are you sure delete this task?"
            ok-text="Yes"
            cancel-text="No"
            @confirm="deleteRule(record,index)"
        >
        <a-button type="danger" size="small" ghost>Delete</a-button>
        </a-popconfirm>
      </span>
    </a-table>
  </div>
</template>
<script>

import {getFtpRule, upsertFtpRule, deleteFtpRule} from '@/api/rule'
import {store} from '@/main'
import BasicRule from "@/components/BasicRule";

const VIEW = "View"
const EDIT = "Edit"
const CREATE = "Create"

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    sorter: true,
    sortDirections: ['descend', 'ascend'],
  },
  {
    title: 'NAME',
    dataIndex: 'name',
    key: 'name',
    scopedSlots: {
      filterDropdown: 'filterDropdown',
      filterIcon: 'filterIcon',
    },
  },
  {
    title: 'FLAG FORMAT',
    dataIndex: 'flag_format',
    key: 'flag_format',
    ellipsis: true,
  },
  {
    title: 'RANK',
    dataIndex: 'rank',
    key: 'rank',
    scopedSlots: {
      customRender: 'rank',
    },
  },
  {
    title: 'PASV ADDRESS',
    dataIndex: 'pasv_address',
    key: 'pasv_address',
    ellipsis: true,
  },
  {
    title: 'PUSH TO CLIENT',
    dataIndex: 'push_to_client',
    key: 'push_to_client',
    scopedSlots: {
      customRender: 'switchRender',
    }
  },
  {
    title: 'NOTICE',
    dataIndex: 'notice',
    key: 'notice',
    scopedSlots: {
      customRender: 'switchRender',
    }
  },
  {
    title: 'Action',
    key: 'action',
    scopedSlots: {customRender: 'action'},
  },
];
export default {
  name: 'FtpRules',
  data() {
    return {
      store,
      data: [],
      formVisible: false,
      pagination: {
        current: 1, showSizeChanger: true, pageSize: store.pageSize,
        onShowSizeChange: (current, size) => {
          store.pageSize = size
        }
      },      filters: {},
      loading: false,
      columns,
      form: {},
      formReadOnly: false,
      formAction: "", // View ,Create or Edit
    }
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
      getFtpRule(params).then(res => {
        let result = res.data.result
        this.data = result.data
        const pagination = {...this.pagination};

        pagination.total = result.count;
        this.pagination = pagination;
        this.loading = false
      }).catch(e => {
        this.$notification.error({
            message: 'Unknown error: ' + e.response.status,
            style: {
              width: '100px',
              marginLeft: `${335 - 600}px`,
            },
            duration: 4
          });
      })
    },
    clickSwitch(record, prop) {
      record[prop] = !record[prop]
      upsertFtpRule(record).then().catch(e => {
        this.$notification.error({
          message: 'Edit failed',
          description:
          e.response.data.error,
          style: {
            width: '600px',
            marginLeft: `${335 - 600}px`,
          },
          duration: 4
        });
      })
    },
    addRule() {
      this.form = {}
      this.showForm(CREATE)
    },
    viewRule(record) {
      this.form = record
      this.showForm(VIEW)
    },
    editRule(record) {
      this.form = JSON.parse(JSON.stringify(record))
      this.showForm(EDIT)
    },
    deleteRule(record, index) {
      deleteFtpRule(record).then(() => {
        this.data.splice(index, 1)
      }).catch(e => {
        this.$notification.error({
          message: 'Error',
          description:
          e.response.data.error,
          style: {
            width: '600px',
            marginLeft: `${335 - 600}px`,
          },
          duration: 4
        });
      })
    },
    showForm(action) {
      this.formAction = action
      this.formReadOnly = action === VIEW;
      this.formVisible = true;
    },
    closeDrawer() {
      this.formVisible = false;
    },
    handleSubmit() {
      this.$refs.form.validate(valid => {
        if (valid) {
          upsertFtpRule(this.form).then(() => {
            this.closeDrawer()
            this.fetch({page: this.pagination.current});
            this.$notification.info({
              message: 'Success',
              style: {
                width: '600px',
                marginLeft: `${335 - 600}px`,
              },
              duration: 2.5
            });
          }).catch(e => {
            this.$notification.error({
              message: this.formAction + ' failed',
              description:
              e.response.data.error,
              style: {
                width: '600px',
                marginLeft: `${335 - 600}px`,
              },
              duration: 4
            });
          })
        }
      })
    },
    handleCancel() {
      this.form = {}
      this.closeDrawer()
    },
  },
  mounted() {
    this.fetch({page: "1"});
  },
  components: {
    BasicRule,
  }
}
</script>
<style scoped>
#add-rule {
  margin-bottom: 10px;
}
</style>