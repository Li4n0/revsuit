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
          <a-menu-item key="/logs/mysql">
            <router-link to="/logs/mysql">MySQL Logs</router-link>
          </a-menu-item>
        </a-sub-menu>
        <a-sub-menu key="rules">
          <span slot="title"><a-icon type="radar-chart"/><span>Rules</span></span>
          <a-menu-item key="/rules/http">
            <router-link to="/rules/http">HTTP Rules</router-link>
          </a-menu-item>
          <a-menu-item key="/rules/dns">
            <router-link to="/rules/dns">DNS Rules</router-link>
          </a-menu-item>
          <a-menu-item key="/rules/mysql">
            <router-link to="/rules/mysql">MySQL Rules</router-link>
          </a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header style="background: #fff; padding: 0">
        <a-icon
            class="trigger"
            :type="collapsed ? 'menu-unfold' : 'menu-fold'"
            @click="() => (collapsed = !collapsed)"
        />
      </a-layout-header>
      <a-layout-content
          :style="{ margin: '24px 16px', padding: '24px', borderRadius: '20px',background: '#fff', minHeight: 'initial' }"
      >
        <transition name="fade-transform">
          <router-view></router-view>
        </transition>
      </a-layout-content>
    </a-layout>
    <Auth></Auth>
  </a-layout>
</template>
<script>
import Auth from '@/components/Auth'

export default {
  data() {
    return {
      collapsed: false,
      openKeys: ['logs', "rules"],
    };
  },
  components: {
    Auth
  }
};
</script>
<style scoped>
html, body {
  height: 100%;
  margin: 0;
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
