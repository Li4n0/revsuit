<template xmlns:a-col="http://www.w3.org/1999/html">
  <div>
    <a-button id="add-rule" type="primary" @click="addRule">
      <a-icon type="plus"/>
      New Rule
    </a-button>
    <!--    rule form-->
    <a-drawer
        :title="formAction+ ' MySQL rule'"
        :width="490"
        :visible="formVisible"
        :body-style="{ paddingBottom: '80px' }"
        @close="closeDrawer"
    >
      <a-form-model :model="form" ref="form" layout="vertical" @submit="handleSubmit">
        <BasicRule :form="form" :readOnly="formReadOnly" :flagField="'user or schema'"/>
        <a-row :gutter="24">
          <a-col :span="24">
            <a-form-model-item label="Files" :rules="rules.files">
              <a-input
                  v-model="form.files"
                  style="width: 100%"
                  placeholder="please enter file name,use ';' to split multiple file names"
                  :readOnly="formReadOnly"
                  :disabled="form.exploit_jdbc_client"
              />
            </a-form-model-item>
          </a-col>
        </a-row>
        <a-row>
          <a-form-model-item>
            <div class="ant-form-item-label">
              <label for="exploit-jdbc-client">Exploit Jdbc Client
                <a-tooltip placement="topLeft"
                           title="Whether test to exploit jdbc client.">
                  <a-icon type="question-circle"/>
                </a-tooltip>
              </label>
            </div>
            <a-switch v-model="form.exploit_jdbc_client" id="exploit-jdbc-client" :disabled="formReadOnly"/>
          </a-form-model-item>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="24">
            <a-form-model-item>
              <label for="payload">Payload
                <a-tooltip placement="topLeft" title="Need to set ${payload} variable by flag format first. Payload which key is same as ${payload} will be used.">
                  <a-icon type="question-circle"/>
                </a-tooltip>
              </label>
              <a-input-group id="payload" compact v-for="payloadKey in payloadKeys" :key="payloadKey">
                <a-input v-model="form['Key-'+payloadKey]"
                         style="width: 47%;margin-bottom: 5px"
                         v-decorator="['Key-'+payloadKey,]"
                         :defaultOpen="false"
                         placeholder="Key"
                         :disabled="!form.exploit_jdbc_client"
                         :readOnly="formReadOnly"
                ></a-input>

                <a-input v-model="form['Value-'+payloadKey]"
                         @focus="()=>{ !formReadOnly&&(payloadKey === payloadKeys[payloadKeys.length-1]) && form['Key-'+payloadKey] ? addPayload(): null}"
                         style="width: 53%"
                         v-decorator="['Value-'+payloadKey,]"
                         placeholder="Base64 encoded payload Value"
                         :disabled="!form.exploit_jdbc_client"
                         :readOnly="formReadOnly"
                >
                  <a-icon slot="addonAfter"
                          class="dynamic-delete-button"
                          type="minus-circle-o"
                          @click="() => form.exploit_jdbc_client? removePayload(payloadKey):null"
                  />
                </a-input>
              </a-input-group>
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
<style scoped>
#add-rule {
  margin-bottom: 10px;
}
</style>
<script>

import {getMysqlRule, upsertMysqlRule, deleteMysqlRule} from '@/api/rule'
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
      customRender: "rank"
    }
  },
  {
    title: 'FILES',
    dataIndex: 'files',
    key: 'files',
    ellipsis: true
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
const rules = {
  files: [{required: false, message: "please enter file name"}],
}
export default {
  name: 'MysqlRules',
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
      }, filters: {},
      loading: false,
      columns,
      form: {},
      rules: rules,
      payloadKeys: [1],
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
      getMysqlRule(params).then(res => {
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
          this.$notification.error({
            message: 'Unknown error: ' + e.response.status,
            style: {
              width: '100px',
              marginLeft: `${335 - 600}px`,
            },
            duration: 4
          });
        }
      })
    },
    clickSwitch(record, prop) {
      record[prop] = !record[prop]
      upsertMysqlRule(record).then().catch(e => {
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
      deleteMysqlRule(record).then(() => {
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

      for (let k in this.form.payloads) {
        this.form["Key-" + this.payloadKeys.length] = k
        this.form["Value-" + this.payloadKeys.length] = this.form.payloads[k]
        this.addPayload()
      }
      if (this.formReadOnly) {
        this.removePayload(this.payloadKeys.length)
      }
    },
    closeDrawer() {
      this.formVisible = false;
      this.payloadKeys = [1]
    },
    addPayload() {
      this.payloadKeys.push(this.payloadKeys[this.payloadKeys.length - 1] + 1)
    },
    removePayload(key) {
      if (this.payloadKeys.length > 1) {
        this.payloadKeys.splice(this.payloadKeys.indexOf(key), 1)
      }
    },
    handleSubmit() {
      this.$refs.form.validate(valid => {
        if (valid) {
          let payloads = {}
          let form = {}
          for (let k in this.form) {
            if (k.indexOf("Key-") === 0) {
              let i = k.substr("Key-".length)
              if (this.form["Value-" + i]) {
                payloads[this.form[k]] = this.form["Value-" + i]
              }
            } else if (k.indexOf("Value-") === -1) {
              form[k] = this.form[k]
            }
          }
          form.payloads = payloads
          upsertMysqlRule(form).then(() => {
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
    }
  },
  mounted() {
    this.fetch({page: "1"});
  },
  components: {
    BasicRule,
  }
}
</script>