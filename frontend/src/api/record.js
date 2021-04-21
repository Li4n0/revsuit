import request from './index'

export function getHttpRecord(params) {
    return request({
        url: '/record/http',
        params: params,
        method: 'get',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function getDnsRecord(params) {
    return request({
        url: '/record/dns',
        params: params,
        method: 'get',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function getMysqlRecord(params) {
    return request({
        url: '/record/mysql',
        params: params,
        method: 'get',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}