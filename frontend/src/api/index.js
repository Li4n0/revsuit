import axios from 'axios'

const service = axios.create({
    baseURL: '/revsuit/api/', // api的base_url
    timeout: 5000, // request timeout
    validateStatus: function (status) {
        return status >= 200 && status < 300 // 默认的
    }
})

export default service