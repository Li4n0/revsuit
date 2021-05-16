import request from './index'

export function auth(token) {
    return request({
        url: '/auth',
        method: 'get',
        headers: {"Token": token},
        validateStatus: function (status) {
            return status >= 200 && status < 300
        }
    })
}

export function getVersion() {
    return request({
        url: '/version',
        method: 'get',
        validateStatus: function (status) {
            return status >= 200 && status < 300
        }
    })
}