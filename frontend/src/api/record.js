import request from './index'

export function getHttpRecord(params) {
    return request({
        url: '/record/http',
        params: params,
        method: 'get',

    })
}

export function deleteHttpRecord(params) {
    return request({
        url: '/record/http',
        params: params,
        method: 'delete',

    })
}

export function getDnsRecord(params) {
    return request({
        url: '/record/dns',
        params: params,
        method: 'get'
    })
}

export function deleteDnsRecord(params) {
    return request({
        url: '/record/dns',
        params: params,
        method: 'delete'
    })
}

export function getMysqlRecord(params) {
    return request({
        url: '/record/mysql',
        params: params,
        method: 'get'
    })
}

export function deleteMysqlRecord(params) {
    return request({
        url: '/record/mysql',
        params: params,
        method: 'delete'
    })
}

export function getRmiRecord(params) {
    return request({
        url: '/record/rmi',
        params: params,
        method: 'get'
    })
}

export function deleteRmiRecord(params) {
    return request({
        url: '/record/rmi',
        params: params,
        method: 'delete'
    })
}

export function getFtpRecord(params) {
    return request({
        url: '/record/ftp',
        params: params,
        method: 'get'
    })
}

export function deleteFtpRecord(params) {
    return request({
        url: '/record/ftp',
        params: params,
        method: 'delete'
    })
}