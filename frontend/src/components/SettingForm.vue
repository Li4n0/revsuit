<template>
  <a-spin :spinning="spinning">
    <a-form-model :model="form" ref="form" layout="vertical">
      <a-row :gutter="24" v-for="(value,key) in form" :key="key">
        <a-col :span="24">
          <a-form-model-item :label="key+':'" :prop="key">
            <a-switch v-if="key==='Enable'" v-model="form[key]" size="large">
              <a-icon slot="checkedChildren" type="check"/>
              <a-icon slot="unCheckedChildren" type="close"/>
            </a-switch>
            <a-select v-else-if="key==='LogLevel'"
                      v-model="form[key]"
            >
              <a-select-option value='debug'>DEBUG</a-select-option>
              <a-select-option value='info'>INFO</a-select-option>
              <a-select-option value='warning'>WARNING</a-select-option>
              <a-select-option value='error'>ERROR</a-select-option>
              <a-select-option value='fatal'>FATAL</a-select-option>
            </a-select>
            <a-input :disabled="key ==='Addr' && parent==='Http'" v-else v-model="form[key]"/>
          </a-form-model-item>
        </a-col>
      </a-row>
      <a-space size="middle" v-if="Object.keys(form).length">
        <a-button @click="this.$parent.getConfig">Cancel</a-button>
        <a-popconfirm
            title="Are you sure update the config?
            It will stop running for at least two seconds, and may fail to restart."
            ok-text="Yes"
            cancel-text="No"
            @confirm="this.$parent.updateConfig"
        >
          <a-button type="danger">Update</a-button>
        </a-popconfirm>
      </a-space>
    </a-form-model>
  </a-spin>
</template>

<script>
export default {
  name: "SettingForm",
  data() {
    return {
      parent: this.$parent.$options.name
    }
  },
  props: ['form', 'spinning'],
}
</script>

<style scoped>

</style>