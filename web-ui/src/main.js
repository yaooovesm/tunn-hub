import {createApp} from 'vue'
import App from './App.vue'

//router
import router from "./router"

//axios
import axios from 'axios'
import VueAxios from 'vue-axios'

//element
import element from 'element-plus';
import 'element-plus/dist/index.css'

//icon
import "./assets/icon/iconfont.css"

//storage
import publicStorage from "@/public.storage";

//util
import utils from "@/utils";

//md5
import md5 from 'js-md5';


//配置请求拦截,添加token
axios.interceptors.request.use(function (config) {
    if (publicStorage.User !== undefined && publicStorage.User.token !== "") {
        config.headers["token"] = publicStorage.User.token
    }
    return config
})
axios.interceptors.response.use(res => {
    return res
}, err => {
    if (err.config.url === "/api/v1/token/check") {
        return Promise.reject(err)
    }
    let response = err.response
    if (response.status === 500) {
        if (response.data.error === undefined) {
            //删除登录状态,跳转到登录页,并弹出提示
            //删除状态
            localStorage.removeItem("tunnel_server_user")
            publicStorage.Load()
            //跳转
            router.push({path: "/login"})
            //提示
            utils.Warning("发生错误", "与服务器失去连接")
            return
        }
        if (response.data.error === "user not login") {
            //删除登录状态,跳转到登录页,并弹出提示
            //删除状态
            localStorage.removeItem("tunnel_server_user")
            publicStorage.Load()
            if (router.path !== "/login") {
                //跳转
                router.push({path: "/login"})
                //提示
                utils.Warning("登录超时", "请重新登录后继续")
            }
            return
        }
    }
    return Promise.reject(err)
})


const app = createApp(App);
//注册储存
app.config.globalProperties.$storage = publicStorage
//注册工具类
app.config.globalProperties.$utils = utils
//注册md5
app.config.globalProperties.$md5 = md5
//注册组件
app
    //router
    .use(router)
    //element
    .use(element)
    //axios
    .use(VueAxios, axios)
//挂载
app.mount('#app');
