import request from './index'

export function auth(token) {
    return request({
        url: '/auth',
        method: 'get',
        headers: {"Token": token},
    })
}

export function getVersion() {
    return request({
        url: '/version',
        method: 'get',
    })
}