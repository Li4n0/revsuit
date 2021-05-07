import request from './index'

export function getHttpConfig() {
    return request({
        url: '/setting/getHttpConfig',
    })
}

export function updateHttpConfig(data) {
    return request({
        url: '/setting/updateHttpConfig',
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