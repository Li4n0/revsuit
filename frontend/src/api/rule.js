import request from './index'

export function getHttpRule(params) {
    return request({
        url: '/rule/http',
        params: params,
        method: 'get',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function upsertHttpRule(data) {
    return request({
        url: '/rule/http',
        data: data,
        method: 'post',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}


export function deleteHttpRule(data) {
    return request({
        url: '/rule/http',
        data: data,
        method: 'delete',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function getDnsRule(params) {
    return request({
        url: '/rule/dns',
        params: params,
        method: 'get',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function upsertDnsRule(data) {
    return request({
        url: '/rule/dns',
        data: data,
        method: 'post',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function deleteDnsRule(data) {
    return request({
        url: '/rule/dns',
        data: data,
        method: 'delete',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function getMysqlRule(params) {
    return request({
        url: '/rule/mysql',
        params: params,
        method: 'get',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function upsertMysqlRule(data) {
    return request({
        url: '/rule/mysql',
        data: data,
        method: 'post',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function deleteMysqlRule(data) {
    return request({
        url: '/rule/mysql',
        data: data,
        method: 'delete',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function getRmiRule(params) {
    return request({
        url: '/rule/rmi',
        params: params,
        method: 'get',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function upsertRmiRule(data) {
    return request({
        url: '/rule/rmi',
        data: data,
        method: 'post',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}

export function deleteRmiRule(data) {
    return request({
        url: '/rule/rmi',
        data: data,
        method: 'delete',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}