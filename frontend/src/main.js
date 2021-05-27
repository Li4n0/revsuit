import Vue from 'vue';
import Antd from 'ant-design-vue';
import App from './App';
import router from "@/router";
import 'ant-design-vue/dist/antd.css';
import '@/utils/index'
import visibility from 'vue-visibility-change';


Vue.config.productionTip = false;

Vue.use(Antd);
Vue.use(visibility);

export const store = Vue.observable({
    authed: true,
    pageSize: localStorage.getItem("pageSize") ? parseInt(localStorage.getItem("pageSize")) : 10
})
/* eslint-disable no-new */
new Vue({
    router: router,
    render: h => h(App)
}).$mount('#app')

