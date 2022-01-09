import request from "@/api/index";

export function getVersion() {
    return request({
        url: '/version',
        method: 'get',
    })
}

export function getUpgrade() {
    return request({
        url: '/getUpgrade',
        method: 'get',
    })
}