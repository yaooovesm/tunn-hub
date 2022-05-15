<template>
  <div class="header-box">
    <div style="float: left;padding-left: 20px;padding-top: 20px;height: 20px">
      <el-page-header content="" style="height: 20px">
        <template #title>
          <el-button v-if="$route.path!=='/dashboard/overview'"
                     @click="$router.push({path: '/dashboard/overview'})"
                     type="text"
                     class="btn"
                     style="padding: 0;height: 20px;line-height: 20px;transform: translateY(-2px);">
            <span>
              <i class="iconfont icon-angle-left"></i>&nbsp;概况
            </span>
          </el-button>
          <span v-else style="cursor: default">TunnHub</span>
        </template>
        <template #icon>
          <span></span>
        </template>
        <template #content>
          <span style="font-size: 16px;color: #404040">{{ getPageName($route.path) }}</span>
        </template>
      </el-page-header>
    </div>
    <div style="float: right;padding-left: 20px;padding-right: 40px">
      <el-dropdown>
        <div class="header-block">
          <i class="iconfont icon-account"></i> {{ account }}
        </div>
        <template v-slot:dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="$refs.UpdateDialog.show()">修改信息</el-dropdown-item>
            <el-dropdown-item @click="$router.push({path: '/dashboard/home'})">我的用户</el-dropdown-item>
            <el-dropdown-item divided style="color: #F56C6C" @click="confirmLogout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    <UserUpdateDialog ref="UpdateDialog"/>
  </div>
</template>

<script>
import axios from "axios";
import {ElMessageBox} from "element-plus";
import UserUpdateDialog from "@/components/home/UserUpdateDialog";

export default {
  name: "PageHeader",
  components: {UserUpdateDialog},
  mounted() {
    this.update()
  },
  data() {
    return {
      account: "",
      auth: "",
    }
  },
  methods: {
    getPageName: function (route) {
      switch (route) {
        case "/dashboard/overview":
          return "Overview"
        case "/dashboard/home":
          return "Home"
        case "/dashboard/control":
          return "ControlPanel"
        case "/dashboard/users":
          return "Accounts"
        default:
          return ""
      }
    },
    update: function () {
      let lo = localStorage.getItem("tunnel_server_user")
      if (lo !== "" && lo !== undefined && lo !== null) {
        let loUsr = JSON.parse(lo)
        this.$storage.User.isLogin = true
        this.$storage.User.info = loUsr.info
        this.$storage.User.token = loUsr.token
      }
      let usr = this.$storage.User
      this.account = usr.info.account
      this.auth = usr.info.auth
    },
    confirmLogout: function () {
      ElMessageBox.confirm(
          '<div>' +
          '<div style="font-size: 14px;line-height: 20px;margin-bottom: 15px;margin-top: 5px;color: #303133;font-weight: 600">' +
          '<i class="iconfont icon-exclamation-circle" style="color: #E6A23C"></i> <span style="font-weight: 500">是否退出登录</span> </div>' +
          '</div>',
          "确认操作",
          {
            dangerouslyUseHTMLString: true,
            draggable: true,
            closeOnClickModal: false,
            showClose: false,
            confirmButtonText: '确认',
            cancelButtonText: "取消",
          }
      ).then((action) => {
        if (action !== 'cancel') {
          this.logout()
        }
      })
    },
    logout: function () {
      let usr = this.$storage.User
      axios({
        method: "get",
        url: "/api/v1/user/logout/" + usr.info.account,
        data: {}
      }).then(() => {
        this.$storage.User = {
          isLogin: false,
          token: "",
          info: undefined,
        }
        localStorage.removeItem("tunnel_server_user")
        this.$router.push({path: "/login"})
      }).catch(err => {
        console.log("error : " + err)
        this.$router.push({path: "/login"})
      })
    }
  }
}
</script>

<style scoped>
.btn {
  color: #007bbb;
  transition-duration: 0.5s;
}

.btn:hover {
  color: #0ebbba;
}

.header-box {
  border-bottom: solid 1px rgba(60, 60, 60, 0.1);
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  height: 60px
}

.header-block {
  vertical-align: center;
  text-align: center;
  height: 60px;
  line-height: 60px;
  font-size: 14px;
  transition-duration: 0.5s;
}

.header-block:hover {
  color: rgba(0, 123, 187, 0.8);
}
</style>