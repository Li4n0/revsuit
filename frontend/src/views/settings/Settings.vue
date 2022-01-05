<template>
  <div>
    <a-card
        :tab-list="tabList"
        :active-tab-key="tabKey"
        @tabChange="key => this.tabKey = key"
        style="width: 100%;">
      <http ref="settings" v-if="tabKey==='HTTP'"></http>
      <dns ref="settings" v-else-if="tabKey==='DNS'"></dns>
      <rmi ref="settings" v-else-if="tabKey==='RMI'"></rmi>
      <ldap ref="settings" v-else-if="tabKey==='LDAP'"></ldap>
      <mysql ref="settings" v-else-if="tabKey==='MYSQL'"></mysql>
      <ftp ref="settings" v-else-if="tabKey==='FTP'"></ftp>
      <notice ref="settings" v-else-if="tabKey==='NOTICE'"></notice>
      <rules ref="settings" v-else-if="tabKey==='RULES'"></rules>
    </a-card>
  </div>
</template>

<script>
import Http from "./Http"
import Dns from "./Dns"
import Rmi from "./Rmi"
import Ldap from "./Ldap"
import Mysql from "./Mysql"
import Ftp from "./Ftp"
import Notice from "./Notice"
import Rules from "./Rules"

export default {
  name: "Settings",
  data() {
    return {
      tabList: [
        {key: "HTTP", tab: "PLATFORM|HTTP"},
        {key: "DNS", tab: "DNS"},
        {key: "RMI", tab: "RMI"},
        {key: "LDAP", tab: "LDAP"},
        {key: "MYSQL", tab: "MYSQL"},
        {key: "FTP", tab: "FTP"},
        {key: "NOTICE", tab: "NOTICE"},
        {key: "RULES", tab: "RULES"}
      ],
      tabKey: "HTTP",
    }
  }
  ,
  methods: {
    fetch() {
      this.$refs.settings.getConfig();
    }
  },
  components: {
    Http,
    Dns,
    Rmi,
    Ldap,
    Mysql,
    Ftp,
    Notice,
    Rules,
  }
}
</script>

<style scoped>

</style>