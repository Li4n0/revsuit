import axios from 'axios'
import {store} from "@/main";

const service = axios.create({
    baseURL: location.pathname.slice(0, -'/admin/'.length) + "/api", // api的base_url
    timeout: 5000, // request timeout
    validateStatus: function (status) {
        if (status === 403) {
            store.authed = false
        }
        return status >= 200 && status < 300 // 默认的
    }
})

export default service