import request from './index'

export function getPlatformConfig() {
    return request({
        url: '/setting/getPlatformConfig',
    })
}

export function updatePlatformConfig(data) {
    return request({
        url: '/setting/updatePlatformConfig',
        method: "post",
        data: data
    })
}

export function getDnsConfig() {
    return request({
        url: '/setting/getDnsConfig',
    })
}

export function updateDnsConfig(data) {
    return request({
        url: '/setting/updateDnsConfig',
        method: "post",
        data: data,
    })
}

export function getRmiConfig() {
    return request({
        url: '/setting/getRmiConfig',
    })
}

export function updateRmiConfig(data) {
    return request({
        url: '/setting/updateRmiConfig',
        method: "post",
        data: data,
    })
}

export function getMySQLConfig() {
    return request({
        url: '/setting/getMySQLConfig',
    })
}

export function updateMySQLConfig(data) {
    return request({
        url: '/setting/updateMySQLConfig',
        method: "post",
        data: data,
    })
}

export function getFtpConfig() {
    return request({
        url: '/setting/getFtpConfig',
    })
}

export function updateFtpConfig(data) {
    return request({
        url: '/setting/updateFtpConfig',
        method: "post",
        data: data,
    })
}

export function getNoticeConfig() {
    return request({
        url: '/setting/getNoticeConfig',
    })
}

export function updateNoticeConfig(data) {
    return request({
        url: '/setting/updateNoticeConfig',
        method: "post",
        data: data,
    })
}