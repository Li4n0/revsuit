import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home
    },
    {
        path: '/logs/http',
        name: 'HttpLogs',
        component: () => import(/* webpackChunkName: "about" */ '../views/logs/Http')
    },
    {
        path: '/logs/dns',
        name: 'DnsLogs',
        component: () => import(/* webpackChunkName: "about" */ '../views/logs/Dns')
    },
    {
        path: '/logs/mysql',
        name: 'MysqlLogs',
        component: () => import(/* webpackChunkName: "about" */ '../views/logs/Mysql')
    },
    {
        path: '/logs/rmi',
        name: 'RmiLogs',
        component: () => import(/* webpackChunkName: "about" */ '../views/logs/Rmi')
    },
    {
        path: '/rules/http',
        name: 'HttpRules',
        component: () => import(/* webpackChunkName: "about" */ '../views/rules/Http')
    },
    {
        path: '/rules/dns',
        name: 'DnsRules',
        component: () => import(/* webpackChunkName: "about" */ '../views/rules/Dns')
    },
    {
        path: '/rules/mysql',
        name: 'MysqlRules',
        component: () => import(/* webpackChunkName: "about" */ '../views/rules/Mysql')
    },
    {
        path: '/rules/rmi',
        name: 'RmiRules',
        component: () => import(/* webpackChunkName: "about" */ '../views/rules/Rmi')
    }
]

const router = new VueRouter({
    routes
})

export default router
