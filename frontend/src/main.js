import Vue from 'vue';
import Antd from 'ant-design-vue';
import App from './App';
import router from "@/router";
import 'ant-design-vue/dist/antd.css';
import '@/utils/index'

Vue.config.productionTip = false;

Vue.use(Antd);
export const store = Vue.observable({
    authed: true,
})
/* eslint-disable no-new */
new Vue({
    router: router,
    render: h => h(App)
}).$mount('#app')

