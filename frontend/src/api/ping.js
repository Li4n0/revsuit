import request from './index'

export function ping() {
    return request({
        url: '/ping',
        method: 'get',
        validateStatus: function (status) {
            return status >= 200 && status < 300 // 默认的
        }
    })
}