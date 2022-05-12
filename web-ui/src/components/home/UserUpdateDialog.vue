<template>
  <el-dialog
      v-model="dialogVisible"
      width="40%"
      :close-on-click-modal="false"
      custom-class="default-dialog"
      :show-close="false"
      draggable
  >
    <template #title>
      <div class="title">
        <div class="title-text">修改用户</div>
      </div>
    </template>
    <el-row :gutter="10" v-loading="loading">
      <el-col :span="20" :offset="2">
        <el-form
            label-position="top"
            label-width="100px"
            :model="info"
        >
          <el-form-item label="账号">
            <el-input v-model="info.account" disabled/>
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="info.password" show-password/>
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="info.email"/>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
    <template #footer>
      <el-button @click="update" type="primary">确认</el-button>
      <el-button @click="dialogVisible = false">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script>
import {ElMessageBox} from "element-plus";
import axios from "axios";

export default {
  name: "UserUpdateDialog",
  data() {
    return {
      id: "",
      loading: false,
      dialogVisible: false,
      info: {
        account: "",
        password: "",
        email: "",
      }
    }
  },
  methods: {
    show: function () {
      this.dialogVisible = true
      let lo = localStorage.getItem("tunnel_server_user")
      if (lo !== "" && lo !== undefined && lo !== null) {
        let loUsr = JSON.parse(lo)
        if (loUsr.info !== undefined) {
          this.id = loUsr.info.id
          this.load()
          return
        }
      }
      this.dialogVisible = false
      ElMessageBox.alert('加载用户信息失败', '错误', {
        confirmButtonText: '确认',
        callback: () => {
          this.dialogVisible = false
        },
      })
    },
    load: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/user/info/" + this.id,
        data: {}
      }).then(res => {
        let response = res.data
        this.info = response.data
        this.loading = false
      }).catch(() => {
        ElMessageBox.alert('加载用户信息失败', '错误', {
          confirmButtonText: '确认',
          callback: () => {
            this.dialogVisible = false
          },
        })
        this.loading = false
      })
    },
    update: function () {
      this.loading = true
      ElMessageBox.prompt('请输入密码', '验证密码', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        inputPattern: /\S/,
        inputType: 'password',
        inputErrorMessage: '密码不能为空'
      }).then(({value}) => {
        let md5Value = this.$md5(value)
        this.$storage.Load()
        if (md5Value !== this.$storage.User.password) {
          this.$utils.Error("验证失败", "密码错误")
          this.loading = false
        } else {
          //更新
          let data = {
            id: this.info.id,
            email: this.info.email
          }
          if (this.info.password !== "") {
            data.password = this.info.password
          }
          axios({
            method: "post",
            url: "/api/v1/user/update/" + this.info.id,
            data: data
          }).then(() => {
            this.$utils.Success("修改成功", "改变将在下一次登录时生效")
            this.loading = false
            this.dialogVisible = false
          }).catch((err) => {
            this.$utils.HandleError(err)
            this.loading = false
          })
        }
      }).catch(() => {
        this.loading = false
      })
    }
  }
}
</script>

<style scoped>

</style>