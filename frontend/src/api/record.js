import request from './index'

export function getHttpRecord(params) {
    return request({
        url: '/record/http',
        params: params,
        method: 'get',

    })
}

export function getDnsRecord(params) {
    return request({
        url: '/record/dns',
        params: params,
        method: 'get'
    })
}

export function getMysqlRecord(params) {
    return request({
        url: '/record/mysql',
        params: params,
        method: 'get'
    })
}

export function getRmiRecord(params) {
    return request({
        url: '/record/rmi',
        params: params,
        method: 'get'
    })
}

export function getFtpRecord(params) {
    return request({
        url: '/record/ftp',
        params: params,
        method: 'get'
    })
}