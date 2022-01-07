import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'Home',
        redirect: '/logs/http'
    },
    // logs router
    {
        path: '/logs/http',
        name: 'HttpLogs',
        component: () => import( '../views/logs/Http')
    },
    {
        path: '/logs/dns',
        name: 'DnsLogs',
        component: () => import( '../views/logs/Dns')
    },
    {
        path: '/logs/mysql',
        name: 'MysqlLogs',
        component: () => import( '../views/logs/Mysql')
    },
    {
        path: '/logs/rmi',
        name: 'RmiLogs',
        component: () => import( '../views/logs/Rmi')
    },
    {
        path: '/logs/ldap',
        name: 'LdapLogs',
        component: () => import( '../views/logs/Ldap')
    },
    {
        path: '/logs/ftp',
        name: 'FtpLogs',
        component: () => import( '../views/logs/Ftp')
    },
    // rules router
    {
        path: '/rules/http',
        name: 'HttpRules',
        component: () => import( '../views/rules/Http')
    },
    {
        path: '/rules/dns',
        name: 'DnsRules',
        component: () => import( '../views/rules/Dns')
    },
    {
        path: '/rules/mysql',
        name: 'MysqlRules',
        component: () => import( '../views/rules/Mysql')
    },
    {
        path: '/rules/rmi',
        name: 'RmiRules',
        component: () => import( '../views/rules/Rmi')
    },
    {
        path: '/rules/ldap',
        name: 'LdapRules',
        component: () => import( '../views/rules/Ldap')
    },
    {
        path: '/rules/ftp',
        name: 'FtpRules',
        component: () => import( '../views/rules/Ftp')
    },
    //Settings router
    {
        path: '/settings',
        name: 'Settings',
        component: () => import( '../views/settings/Settings')
    }
]

const router = new VueRouter({
    routes
})

export default router
