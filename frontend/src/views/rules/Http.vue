<template xmlns:a-col="http://www.w3.org/1999/html">
  <div>
    <a-button id="add-rule" type="primary" @click="addRule">
      <a-icon type="plus"/>
      New Rule
    </a-button>
    <!--    rule form-->
    <a-drawer
        :title="formAction+ ' HTTP rule'"
        :width="490"
        :visible="formVisible"
        :body-style="{ paddingBottom: '80px' }"
        @close="closeDrawer"
    >
      <a-form-model :model="form" ref="form" layout="vertical" @submit="handleSubmit">
        <BasicRule :form="form" :readOnly="formReadOnly" flagFormat="raw message"/>
        <a-row :gutter="24">
          <a-col :span="24">
            <a-form-model-item :rules="rules.response_status_code"
                               prop="response_status_code">
              <span slot="label">
                  Response Status Code
                  <a-tooltip
                      title="Number between 100-600, or template such as ${query.varname}/${body.varname}/${header.varname}/${flag}/${custom_varname}">
                    <a-icon type="question-circle-o"/>
                  </a-tooltip>
              </span>
              <a-input v-model="form.response_status_code"
                       style="width: 100%"
                       placeholder="200"
                       :disabled="formReadOnly"
              />
            </a-form-model-item>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="24">
            <a-form-model-item>
              <span slot="label">
                  Response Headers
                  <a-tooltip
                      title="Support template such as ${query.varname}/${body.varname}/${header.varname}">
                    <a-icon type="question-circle-o"/>
                  </a-tooltip>
              </span>
              <a-input-group compact v-for="headerKey in headerKeys" :key="headerKey">
                <a-auto-complete v-model="form['Header-'+headerKey]"
                                 style="width: 47%;margin-bottom: 5px"
                                 v-decorator="['Header-'+headerKey,]"
                                 :dataSource="headerSet"
                                 :filterOption="filterOption"
                                 :defaultOpen="false"
                                 placeholder="Header"
                                 :disabled="formReadOnly"
                ></a-auto-complete>

                <a-input v-model="form['Value-'+headerKey]"
                         @focus="()=>{ !formReadOnly&&(headerKey === headerKeys[headerKeys.length-1]) && form['Header-'+headerKey] ? addHeader(): null}"
                         style="width: 53%"
                         v-decorator="['Value-'+headerKey,]"
                         placeholder="Value"
                         :disabled="formReadOnly"
                >
                  <a-icon slot="addonAfter"
                          class="dynamic-delete-button"
                          type="minus-circle-o"
                          @click="() => !formReadOnly? removeHeader(headerKey):null"
                  />
                </a-input>
              </a-input-group>
            </a-form-model-item>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="24">
            <a-form-model-item>
              <span slot="label">
                  Response Body
                  <a-tooltip
                      title="Support template such as ${query.varname}/${body.varname}/${header.varname}">
                    <a-icon type="question-circle-o"/>
                  </a-tooltip>
              </span>
              <a-icon type="fullscreen" class="full-screen-icon" @click="fullScreenTextarea = true"/>
              <a-textarea v-model="form.response_body"
                          placeholder="Hello RevSuit!"
                          :readOnly="formReadOnly"
                          :auto-size="{ minRows: 10, maxRows: 30 }"
              />
              <a-modal
                  :dialog-style="{ top: '40px', bottom: '20px'}"
                  width="80rem"
                  :visible="fullScreenTextarea"
                  title='Response Body'
                  okText=''
                  @ok="fullScreenTextarea=false"
                  @cancel="fullScreenTextarea=false"
              >
                <a-textarea v-model="form.response_body"
                            placeholder="Hello RevSuit!"
                            :readOnly="formReadOnly"
                            style="height:100%"
                            :auto-size="{ minRows: fullScreenTextareaMinRows(), maxRows: fullScreenTextareaMinRows()*1.5 }"
                />
              </a-modal>
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
      style="overflow-x: auto;"
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
<script>

import {getHttpRule, upsertHttpRule, deleteHttpRule} from '@/api/rule'
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
  response_status_code: [{
    validator: (rule, code, callback) => {
      if (code === undefined || (!isNaN(code) && 100 < code && code < 600)) {
        return callback()
      } else if (/\${(query|body|header)\..+?}/.test(code)) {
        return callback()
      } else {
        return callback(new Error("please input legal response status code value"))
      }
    }, trigger: 'blur'
  }]
}

export default {
  name: 'HttpRules',
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
      rules: rules,
      formReadOnly: false,
      formAction: "", // View ,Create or Edit
      headerKeys: [1],
      headerSet: ["Accept-Patch",
        "Accept-Ranges",
        "Age",
        "Allow",
        "Cache-Control",
        "Connection",
        "Content-Disposition",
        "Content-Encoding",
        "Content-Language",
        "Content-Length",
        "Content-Location",
        "Content-Range",
        "Content-Type",
        "Date",
        "Delta-Base",
        "ETag",
        "Expires",
        "Last-Modified",
        "Link",
        "Location",
        "Pragma",
        "Proxy-Authenticate",
        "Public-Key-Pins",
        "Retry-After",
        "Server",
        "Set-Cookie",
        "Strict-Transport-Security",
        "Transfer-Encoding",
        "Upgrade",
        "Vary",
        "Via",
        "Warning",
        "WWW-Authenticate",
        "Content-Security-Policy",
        "Refresh",
        "X-Powered-By",
        "X-Request-ID",
        "X-UA-Compatible",
        "X-XSS-Protection",
        "Access-Control-Allow-Origin",
        "Access-Control-Allow-Credentials",
        "Access-Control-Expose-Headers",
        "Access-Control-Max-Age",
        "Access-Control-Allow-Methods",
        "Access-Control-Allow-Headers"],
      fullScreenTextarea: false,
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
      getHttpRule(params).then(res => {
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
    clickSwitch(record, prop) {
      record[prop] = !record[prop]
      upsertHttpRule(record).catch(e => {
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
      deleteHttpRule(record).then(() => {
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
      for (let k in this.form.response_headers) {
        this.form["Header-" + this.headerKeys.length] = k
        this.form["Value-" + this.headerKeys.length] = this.form.response_headers[k]
        this.addHeader()
      }
      if (this.formReadOnly) {
        this.removeHeader(this.headerKeys.length)
      }
    },
    closeDrawer() {
      this.formVisible = false;
      this.headerKeys = [1]
    },
    addHeader() {
      this.headerKeys.push(this.headerKeys[this.headerKeys.length - 1] + 1)
    },
    removeHeader(key) {
      if (this.headerKeys.length > 1) {
        this.headerKeys.splice(this.headerKeys.indexOf(key), 1)
      }
    },
    filterOption(input, option) {
      return (
          option.componentOptions.children[0].text.toUpperCase().indexOf(input.toUpperCase()) >= 0
      );
    },
    handleSubmit() {
      this.$refs.form.validate(valid => {
        if (valid) {
          let form = {}
          let headers = {}
          for (let k in this.form) {
            if (k.indexOf("Header-") === 0) {
              let i = k.substr("Header-".length)
              if (this.form["Value-" + i]) {
                headers[this.form[k]] = this.form["Value-" + i]
              }
            } else if (k.indexOf("Value-") === -1) {
              form[k] = this.form[k]
            }
          }
          form.response_headers = headers
          upsertHttpRule(form).then(() => {
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
        } else {
          return false;
        }
      });
    },
    handleCancel() {
      this.form = {}
      this.closeDrawer()
    },
    fullScreenTextareaMinRows() {
      return (document.body.clientHeight / window.getComputedStyle(document.body)["fontSize"].slice(0, -2)) * 0.5
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

<style scoped>
#add-rule {
  margin-bottom: 10px;
}

.full-screen-icon {
  position: absolute;
  z-index: 5;
  color: #9e9e9e;
  right: 1px;
  font-size: 1rem;
  cursor: pointer;
  transition: font-size 0.1s;
}

.full-screen-icon:hover {
  font-size: 1.1rem;
  transition: font-size 0.1s;
}
</style>