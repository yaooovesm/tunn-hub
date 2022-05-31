<template>
  <div id="app" v-loading="loading">
    <router-view/>
  </div>
</template>

<script>

import axios from "axios";

export default {
  name: 'app',
  data() {
    return {
      loading: false,
      checkDone: false,
      isLogin: false
    }
  },
  mounted() {
    window.onload = function () {
      if (!window.sessionStorage["login"]) {
        localStorage.removeItem("tunnel_server_user")
        location.reload()
      } else {
        window.sessionStorage.removeItem("login");
      }
    };
    window.onunload = function () {
      window.sessionStorage["login"] = true;
    };
    window.onbeforeunload = function () {
      window.sessionStorage["login"] = true;
    };
    this.checkLogin()
    window.onbeforeunload = (e) => {
      this.$storage.Load()
      e.stopPropagation()
    };
  },
  components: {},
  methods: {
    checkLogin: function () {
      this.loading = true
      let usr = this.$storage.User
      if (usr.isLogin && usr.token !== "" && usr.info !== undefined) {
        this.$router.push({path: "/dashboard/home"})
      } else {
        //查看本地是否有登录信息
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
          }).catch(() => {
            //验证失败删除后跳转到登录页
            //this.$utils.Warning("登录超时", err.response.data.msg)
            this.$router.push({path: "/login"})
            localStorage.removeItem("tunnel_server_user")
          })
        } else {
          //没有登录信息跳转到登录页
          this.$router.push({path: "/login"})
        }
      }
      this.loading = false
    }
  }
}
</script>

<style>
.icon {
  width: 1em;
  height: 1em;
  vertical-align: -0.15em;
  fill: currentColor;
  overflow: hidden;
}

.default-dialog {
  border-radius: 4px !important;
}

.default-dialog .el-dialog__header {
  padding-top: 20px;
  padding-left: 0;
  height: 50px;

}

.title {
  text-align: left;
  border-left: solid 8px rgba(0, 123, 187, 0.8);
  padding-left: 20px;
  height: 25px;
  line-height: 25px;
  margin-bottom: 50px;
}

.title-text {
  font-size: 20px;
}

.circle {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
  margin-bottom: 1px;
}

.online-text {
  color: #67C23A;
  height: 25px;
  margin-left: 3px;
  display: inline-block;
  font-size: 12px;
  transform: translateY(-1px);
}

.offline-text {
  color: #909399;
  height: 25px;
  margin-left: 3px;
  display: inline-block;
  font-size: 12px;
  transform: translateY(-1px);
}

.online-circle {
  background-color: #67C23A;
  box-shadow: 0 0 3px #67C23A;
}

.offline-circle {
  background-color: #909399;
  box-shadow: 0 0 3px #909399;
}

body {
  -webkit-tap-highlight-color: transparent;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  padding: 0;
  margin: 0;
  height: 100%;
  width: 100%;
  position: absolute;
  user-select: none;
}

#app {
  height: 100%;
  width: 100%;
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
