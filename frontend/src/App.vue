<template>
  <a-layout id="nav">
    <a-layout-sider v-model="collapsed" :trigger="null" collapsible>
      <div class="logo"><b>R</b><span v-if="!collapsed"><b>ev</b>Suit</span></div>
      <a-menu theme="dark" mode="inline" :selectedKeys="[this.$route.path]" :open-keys.sync='openKeys'>
        <!--        <a-menu-item key="/">-->
        <!--          <router-link to="/">-->
        <!--            <a-icon type="dashboard"/>-->
        <!--            <span>Dashboard</span>-->
        <!--          </router-link>-->
        <!--        </a-menu-item>-->
        <a-sub-menu key="logs">
          <span slot="title"><a-icon type="bar-chart"/><span>Logs</span></span>
          <a-menu-item key="/logs/http">
            <router-link to="/logs/http">HTTP Logs</router-link>
          </a-menu-item>
          <a-menu-item key="/logs/dns">
            <router-link to="/logs/dns">DNS Logs</router-link>
          </a-menu-item>
          <a-menu-item key="/logs/rmi">
            <router-link to="/logs/rmi">RMI Logs</router-link>
          </a-menu-item>
          <a-menu-item key="/logs/mysql">
            <router-link to="/logs/mysql">MySQL Logs</router-link>
          </a-menu-item>
          <a-menu-item key="/logs/ftp">
            <router-link to="/logs/ftp">FTP Logs</router-link>
          </a-menu-item>
        </a-sub-menu>
        <a-sub-menu key="rules">
          <span slot="title"><rule-icon/><span>Rules</span></span>
          <a-menu-item key="/rules/http">
            <router-link to="/rules/http">HTTP Rules</router-link>
          </a-menu-item>
          <a-menu-item key="/rules/dns">
            <router-link to="/rules/dns">DNS Rules</router-link>
          </a-menu-item>
          <a-menu-item key="/rules/rmi">
            <router-link to="/rules/rmi">RMI Rules</router-link>
          </a-menu-item>
          <a-menu-item key="/rules/mysql">
            <router-link to="/rules/mysql">MySQL Rules</router-link>
          </a-menu-item>
          <a-menu-item key="/rules/ftp">
            <router-link to="/rules/ftp">FTP Rules</router-link>
          </a-menu-item>
        </a-sub-menu>
        <a-menu-item  key="/settings">
          <router-link to="/settings">
            <a-icon type="setting"/>
            Settings
          </router-link>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header style="background: #fff; padding: 0">
        <a-icon
            class="trigger"
            :type="collapsed ? 'menu-unfold' : 'menu-fold'"
            @click="() => (collapsed = !collapsed)"
        />
        <!--        <transition name="list1">-->
        <div v-if="isLogMode" style="float: right; width:45% ;padding: 12px 0;line-height: 24px;">
          <a-row :gutter="24" type="flex">
            <a-col :span="21">
              <a-form-model v-show="showSettings" ref="settings" layout="inline">
                <a-row :gutter="24" type="flex">
                  <a-col :span="9" :offset="2">
                    <a-form-model-item label="Auto Refresh">
                      <a-switch id="auto-refresh" v-model="autoRefresh"></a-switch>
                    </a-form-model-item>
                  </a-col>
                  <a-col :span="13">
                    <a-form-model-item label="Refresh Interval">
                      <a-input-number style="margin-right: -3rem;" id="refresh-interval"
                                      v-model="refreshInterval"
                                      :disabled="!autoRefresh"></a-input-number>
                    </a-form-model-item>
                  </a-col>
                </a-row>
              </a-form-model>
            </a-col>
            <a-col :span="3">
              <a-icon
                  :style="'font-size: 18px;padding: 12px 0;'+(showSettings?'color: #1b90ff;':'')"
                  type="setting"
                  @click="showSettings = !showSettings"
              />
            </a-col>
          </a-row>
        </div>
        <!--        </transition>-->
      </a-layout-header>
      <a-layout-content
          :style="{ margin: '24px 16px', padding: '24px', borderRadius: '20px',background: '#fff', minHeight: 'initial' }"
      >
        <transition name="fade-transform">
          <router-view ref='content'/>
        </transition>
      </a-layout-content>
    </a-layout>
    <Auth></Auth>
  </a-layout>
</template>
<script>
import Auth from '@/components/Auth'
import RuleIcon from "@/components/Icon";

export default {
  data() {
    return {
      autoRefresh: localStorage.getItem("autoRefresh") === "true",
      refreshInterval: localStorage.getItem("refreshInterval"),
      collapsed: false,
      showSettings: false,
      openKeys: ['logs', "rules"],
    };
  },
  computed: {
    isLogMode() {
      return this.$route.path.includes('logs')
    }
  },
  methods: {
    timing() {
      if (this.timer !== null) {
        clearInterval(this.timer)
      }
      this.timer = setInterval(() => {
        this.$refs.content.fetch()
      }, this.refreshInterval * 1000)
    }
  },
  mounted() {
    if (this.autoRefresh && this.isLogMode) {
      this.timing()
    }
  },
  destroyed() {
    clearInterval(this.timer)
  },
  watch: {
    autoRefresh(val) {
      if (!val) {
        clearInterval(this.timer)
      } else {
        this.timing()
      }
      localStorage.setItem('autoRefresh', val)
    },
    isLogMode(val) {
      if (!val) {
        clearInterval(this.timer)
      } else {
        this.timing()
      }
      localStorage.setItem('autoRefresh', val)
    },
    refreshInterval(val) {
      clearInterval(this.timer)
      this.timing()
      localStorage.setItem('refreshInterval', val)
    }
  },
  components: {
    RuleIcon,
    Auth
  }
};
</script>
<style scoped>
html, body {
  height: 100%;
  margin: 0;
}

.fade-enter.fade-enter-active,
.fade-appear.fade-appear-active {
  -webkit-animation-name: none;
  animation-name: none;
}

.fade-leave.fade-leave-active {
  -webkit-animation-name: none;
  animation-name: none;
}

/* fade-transform */
.fade-transform-leave-active,
.fade-transform-enter-active {
  transition: all .3s;
  opacity: 0;
}

.fade-transform-enter {
  opacity: 0;
}

.fade-transform-leave {
  opacity: 0;
}

.fade-transform-leave-to {
}

.fade-transform-enter-to {
  opacity: 0;
}

.ant-menu-item > span > a {
  color: rgba(255, 255, 255, 0.65);
}

.ant-menu-item-selected > span > a {
  color: white;
}

#nav {
  height: 100%;
}

#nav .trigger {
  font-size: 18px;
  line-height: 64px;
  padding: 0 24px;
  cursor: pointer;
  transition: color 0.3s;
}

#nav .trigger:hover {
  color: #1890ff;
}

#nav .logo {
  height: 32px;
  background: #0a1d2d;
  margin: 16px;
  text-align: center;
  font-size: 1.2rem;
  color: white;
  padding-bottom: 5px;
  border-bottom: 2px solid #b6befa;
}
</style>
