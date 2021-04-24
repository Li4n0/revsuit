<template>
  <div>
    <a-modal
        title="Auth with token"
        :visible="visible"
        :confirm-loading="confirmLoading"
        @cancel="cancel"
        @ok="auth"
    >
      <a-input v-model.lazy="token" @pressEnter="auth" placeholder="Your token"/>
    </a-modal>
  </div>
</template>
<script>
import {store} from "@/main";
import {ping} from "@/api/ping"

export default {
  data() {
    return {
      confirmLoading: false,
      token: localStorage.getItem("token")
    };
  },
  computed: {
    visible: function () {
      return !store.authed
    },
  },
  watch: {
    token: function (val) {
      localStorage.setItem("token", val)
    }
  },
  methods: {
    auth() {
      this.confirmLoading = true;
      ping().then(() => {
        store.authed = true;
        this.confirmLoading = false;
      }).catch(e => {
        if (e.response.status === 403) {
          this.$notification.error({
            message: 'Wrong token',
            description:
                'Your token is wrong, please check your server config file.',
            style: {
              width: '600px',
              marginLeft: `${335 - 600}px`,
            },
            duration: 2.5
          });
        }
        this.confirmLoading = false;
      })

    },
    cancel() {
      store.authed = true;
    }
  },
};
</script>