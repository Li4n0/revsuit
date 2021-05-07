<template>
  <div>
    <a-card
        :tab-list="tabList"
        :active-tab-key="tabKey"
        @tabChange="key => this.tabKey = key"
        style="width: 100%;">
      <http v-if="tabKey==='HTTP'"></http>
      <dns v-else-if="tabKey==='DNS'"></dns>
      <rmi v-else-if="tabKey==='RMI'"></rmi>
      <mysql v-else-if="tabKey==='MYSQL'"></mysql>
      <ftp v-else-if="tabKey==='FTP'"></ftp>
      <notice v-else-if="tabKey==='NOTICE'"></notice>
      <rules v-else-if="tabKey==='RULES'"></rules>
    </a-card>
  </div>
</template>

<script>
import Http from "./Http"
import Dns from "./Dns"
import Rmi from "./Rmi"
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
        {key: "MYSQL", tab: "MYSQL"},
        {key: "FTP", tab: "FTP"},
        {key: "NOTICE", tab: "NOTICE"},
        {key: "RULES", tab: "RULES"}
      ],
      tabKey: "HTTP",
    }
  }
  ,
  methods: {},
  components: {
    Http,
    Dns,
    Rmi,
    Mysql,
    Ftp,
    Notice,
    Rules,
  }
}
</script>

<style scoped>

</style>