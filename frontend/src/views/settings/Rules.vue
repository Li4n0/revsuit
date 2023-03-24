<template>
  <div>
    <a-space size="middle">
      <a-upload
          name="rules"
          action="../api/setting/importRules"
          @change="handleChange"
          :showUploadList="false"
      >
        <a-button type="primary" icon="upload">Import</a-button>
      </a-upload>
      <a-button type="primary" icon="download" onclick="window.open('../api/setting/exportRules')">Export
      </a-button>
    </a-space>
    <a-result v-if="status"
              :status="status"
              :title="title"
    >
      <template #extra>
        <a-row>
          <a-col :span="12" :offset="7">
            <div style="text-align: left">
              <li v-for="err in errors" :key="err">{{ err }}</li>
            </div>
          </a-col>
        </a-row>
      </template>
    </a-result>
  </div>
</template>

<script>
export default {
  name: "Rules",
  data() {
    return {
      status: "",
      title: "",
      errors: []
    }
  },
  methods: {
    handleChange(info) {
      let resp = info.file.response
      if (info.file.status === 'done' || info.file.status === 'error') {
        if (resp.error) {
          this.status = "warning"
        } else {
          this.status = "success"
        }
        this.title = resp.result
        this.errors = resp.error
      }
    },
  }
}
</script>

<style scoped>

</style>