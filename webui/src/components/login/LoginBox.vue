<template>
  <div v-loading="loading"
       class="box-outer"
       style="width: 420px;height:400px;padding-top:30px;padding-bottom: 8px;transition-duration: 0.5s">
    <div class="title">
      <div class="title-text">登录到控制台</div>
    </div>
    <el-row>
      <el-col :span="20" :offset="2">
        <el-form
            label-position="top"
            label-width="100px"
            :model="loginData"
            style="max-width: 460px"
        >
          <el-form-item label="用户">
            <el-input v-model="loginData.account" @keyup.enter="login"/>
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="loginData.password" @keyup.enter="login" show-password/>
          </el-form-item>
        </el-form>
        <el-button type="primary" @click="login" style="margin-top: 40px;width: 100%">Login</el-button>
      </el-col>
    </el-row>
    <el-divider style="margin-top: 50px;margin-bottom: 30px"/>
    <div class="footer">
      Tunnhub v{{ version }}
      <span v-if="develop">[开发版本]
        <el-tooltip
            effect="dark"
            content="当前版本为开发版，可能存在缺陷，请勿使用"
            placement="bottom-end"
        >
        <i style="font-size: 8px;color: rgba(0,123,187,0.8)" class="iconfont icon-question-circle"></i>
        </el-tooltip>
      </span>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "LoginBox",
  components: {},
  mounted() {
    this.getVersion();
    if (this.auto) {
      this.autoLogin()
    }
  },
  props: {
    auto: {
      type: Boolean,
      default: true
    }
  },
  data() {
    return {
      loading: false,
      version: "",
      develop: false,
      loginData: {
        account: "",
        password: "",
      }
    }
  },
  methods: {
    getVersion: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/server/version",
        data: {}
      }).then(res => {
        let response = res.data
        this.version = response.data.version
        this.develop = response.data.develop
        this.loading = false
      }).catch(() => {
        this.version = "获取版本失败"
        this.loading = false
        //this.$utils.Error("获取版本失败", "请检查网络连接")
      })
    },
    autoLogin: function () {
      // const loading = ElLoading.service({
      //   lock: true,
      //   text: 'Loading',
      //   background: 'rgba(255, 255, 255, 0.5)',
      // })
      this.loading = true
      let lo = localStorage.getItem("tunnel_server_user")
      if (lo !== "" && lo !== undefined && lo !== null) {
        //本地有登录信息，验证
        let loUsr = JSON.parse(lo)
        axios({
          method: "post",
          url: "/api/v1/token/check",
          headers: {
            token: loUsr.token
          },
          data: {
            account: loUsr.info.account
          }
        }).then(() => {
          this.$storage.User.isLogin = true
          this.$storage.User.info = loUsr.info
          this.$storage.User.token = loUsr.token
          this.$router.push({path: "/dashboard/home"})
          //验证成功跳转
          this.loading = false
        }).finally(() => {
          this.loading = false
        });
      } else {
        this.loading = false
      }
    },
    login: function () {
      if (this.auto) {
        this.autoLogin()
      }
      if (this.loginData.account === "") {
        this.$utils.Warning("提示", "用户名不能为空")
        return
      }
      if (this.loginData.password === "") {
        this.$utils.Warning("提示", "密码不能为空")
        return
      }
      this.loading = true
      axios({
        method: "post",
        url: "/api/v1/user/login",
        data: {
          account: this.loginData.account,
          password: this.loginData.password
        }
      }).then(res => {
        let response = res.data
        if (response.data !== undefined) {
          if (response.data.info !== undefined && response.data.token !== "") {
            this.$storage.User.info = response.data.info
            this.$storage.User.token = response.data.token
            this.$storage.User.isLogin = true
            this.$storage.User.reporter = response.data.reporter
            this.$storage.User.password = this.$md5(this.loginData.password)
            if (response.data.info.auth === "admin") {
              this.$router.push({path: "/dashboard/overview"})
            } else {
              this.$router.push({path: "/dashboard/home"})
            }
            localStorage.setItem("tunnel_server_user", JSON.stringify(this.$storage.User))
          } else {
            this.$utils.Error("登录失败", response.msg)
          }
        } else {
          this.$utils.Error("登录失败", "未知错误")
        }
        this.loading = false
      }).catch(err => {
        let response = err.response.data
        this.$utils.Error(response.msg, response.error)
        this.loading = false
      })
    }
  }

}
</script>

<style scoped>
.box-outer {
  box-shadow: 1px 1px 5px rgba(50, 50, 50, 0.5);
  border-radius: 10px;
}

.footer {
  font-size: 11px;
  color: #808080;
  text-align: right;
  padding-right: 10px;
  opacity: 0.5;
}
</style>