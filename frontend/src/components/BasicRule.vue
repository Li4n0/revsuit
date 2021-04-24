<template>
  <div>
    <a-row :gutter="24">
      <a-col :span="24">
        <a-form-model-item label="Name" :rules="rules.name" prop="name">
          <a-input v-model="form.name"
                   placeholder="Please enter rule name"
                   :readOnly="readOnly"
          />
        </a-form-model-item>
      </a-col>
    </a-row>
    <a-row :gutter="24">
      <a-col :span="24">
        <a-form-model-item :rules="rules.flagFormat" prop="flag_format">
              <span slot="label">
        Flag Format&nbsp;
        <a-tooltip title="Basic usage:
        1. Only when the request contains content that satisfies the flag format, the request will be captured.
        2. Use regular expression syntax.
        3. The character '*' means to capture all requests.
        Advanced usage:
        1. When the regex uses grouping without group name, the platform will only notify the user or push to the client when the first group appears for the first time.
        2. When the regex uses grouping with group name, you can get these submatches through template variables and use them in other fields of the rule.">
          <a-icon type="question-circle-o"/>
        </a-tooltip>
      </span>
          <a-input
              v-model="form.flag_format"
              style="width: 100%"
              placeholder="please enter flag format"
              :readOnly="readOnly"
          />
        </a-form-model-item>
      </a-col>
    </a-row>
    <a-row :gutter="24">
      <a-col :span="24">
        <a-form-model-item prop="rank">
              <span slot="label">
        Rank
        <a-tooltip title="When request match multiple rules, high-rank rules will be matched first">
          <a-icon type="question-circle-o"/>
        </a-tooltip>
      </span>
          <a-input-number style="width: 100%"
                          v-model="form.rank"
                          v-decorator="['rank']"
                          :disabled="readOnly"
                          placeholder="0"
          >
          </a-input-number>
        </a-form-model-item>
      </a-col>
    </a-row>
    <a-row :gutter="16">
      <a-col :span="12">
        <a-form-model-item>
          <div class="ant-form-item-label">
            <label for="push-to-client">Push to Client
              <a-tooltip placement="topLeft" title="Whether push to client when capture flag with this rule.">
                <a-icon type="question-circle"/>
              </a-tooltip>
            </label>
          </div>
          <a-switch v-model="form.push_to_client" id="push-to-client" :disabled="readOnly"/>
        </a-form-model-item>
      </a-col>
      <a-col :span="12">
        <a-form-model-item>
          <div class="ant-form-item-label">
            <label for="notice">Notice
              <a-tooltip placement="topLeft" title="Whether notice with bot when capture flag with this rule.">
                <a-icon type="question-circle"/>
              </a-tooltip>
            </label>
          </div>
          <a-switch v-model="form.notice" id="notice" :disabled="readOnly"/>
        </a-form-model-item>
      </a-col>
    </a-row>
  </div>
</template>

<script>
export default {
  name: "BasicRule",
  data() {
    return {
      rules: {
        name: [
          {required: true, message: 'Please input rule name', trigger: 'blur'},
        ],
        flagFormat: [
          {required: true, message: 'Please input flag format', trigger: 'blur'},
        ]
      },
    }
  },
  props: ['form', 'readOnly']
}
</script>

<style scoped>

</style>