<template>
  <div>
    <SettingForm :form="form" :spinning="spinning"></SettingForm>
    <a-result v-if="status"
              :status="status"
              :title="title"
              :sub-title="subTitle"
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
import {getNoticeConfig, updateNoticeConfig} from "@/api/settings"
import SettingForm from "@/components/SettingForm";

export default {
  name: "Notice",
  data() {
    return {
      form: {},
      spinning: false,
      status: "",
      title: "",
      subTitle: "",
      errors: []
    }
  },
  methods: {
    getConfig() {
      getNoticeConfig().then(res => {
        this.form = res.data
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
    updateConfig() {
      this.spinning = true
      let targetConfig = JSON.stringify(this.form)
      updateNoticeConfig(this.form).then(
          () => {
            setTimeout(() => {
              getNoticeConfig().then((res) => {
                    this.spinning = false
                    let nowConfig = JSON.stringify(res.data)
                    if (nowConfig !== targetConfig) {
                      this.status = "warning"
                      this.title = "Updating the configuration seems to have failed"
                      this.subTitle = "Please check your config options carefully."
                    }
                    this.form = res.data
                    this.status = ""
                    this.title = ""
                    this.subTitle = ""
                  }
              )
            }, 3000)
          }
      ).catch(e => {
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
    }
  },
  mounted() {
    this.getConfig()
  },
  components: {
    SettingForm
  }
}
</script>

<style scoped>

</style>