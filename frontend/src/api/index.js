import axios from 'axios'

const service = axios.create({
    baseURL: '/revsuit/api/', // api的base_url
    timeout: 5000 // request timeout
})

service.interceptors.request.use(function (config) {
    const token = localStorage.getItem("token")
    if (token !== null) {
        config.headers['Token'] = token
    }
    // 在发送请求之前做些什么
    return config
})

export default service