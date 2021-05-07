import request from './index'

export function getHttpRule(params) {
    return request({
        url: '/rule/http',
        params: params,
        method: 'get'
    })
}

export function upsertHttpRule(data) {
    return request({
        url: '/rule/http',
        data: data,
        method: 'post'
    })
}


export function deleteHttpRule(data) {
    return request({
        url: '/rule/http',
        data: data,
        method: 'delete'
    })
}

export function getDnsRule(params) {
    return request({
        url: '/rule/dns',
        params: params,
        method: 'get'
    })
}

export function upsertDnsRule(data) {
    return request({
        url: '/rule/dns',
        data: data,
        method: 'post'
    })
}

export function deleteDnsRule(data) {
    return request({
        url: '/rule/dns',
        data: data,
        method: 'delete'
    })
}

export function getMysqlRule(params) {
    return request({
        url: '/rule/mysql',
        params: params,
        method: 'get'
    })
}

export function upsertMysqlRule(data) {
    return request({
        url: '/rule/mysql',
        data: data,
        method: 'post'
    })
}

export function deleteMysqlRule(data) {
    return request({
        url: '/rule/mysql',
        data: data,
        method: 'delete'
    })
}

export function getRmiRule(params) {
    return request({
        url: '/rule/rmi',
        params: params,
        method: 'get'
    })
}

export function upsertRmiRule(data) {
    return request({
        url: '/rule/rmi',
        data: data,
        method: 'post'
    })
}

export function deleteRmiRule(data) {
    return request({
        url: '/rule/rmi',
        data: data,
        method: 'delete'
    })
}

export function getFtpRule(params) {
    return request({
        url: '/rule/ftp',
        params: params,
        method: 'get'
    })
}

export function upsertFtpRule(data) {
    return request({
        url: '/rule/ftp',
        data: data,
        method: 'post'
    })
}

export function deleteFtpRule(data) {
    return request({
        url: '/rule/ftp',
        data: data,
        method: 'delete'
    })
}