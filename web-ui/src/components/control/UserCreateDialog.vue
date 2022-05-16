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
        <div class="title-text">创建用户</div>
      </div>
    </template>
    <el-row :gutter="10" v-loading="loading">
      <el-col :span="20" :offset="2">
        <el-form :model="user" label-position="top">
          <el-row :gutter="30">
            <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
              <el-form-item label="账号">
                <el-input v-model="user.account" placeholder="account"/>
              </el-form-item>
              <el-form-item label="密码">
                <el-input v-model="user.password" placeholder="password" show-password/>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
              <el-form-item label="邮箱">
                <el-input v-model="user.email" placeholder="email (option)"/>
              </el-form-item>
              <el-form-item label="是否启用">
                <el-select v-model="user.disabled" placeholder="默认禁用账号" value="">
                  <el-option label="启用" :value="0"/>
                  <el-option label="禁用" :value="1"/>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-col>
    </el-row>
    <template #footer>
      <el-button @click="create" type="primary">确认</el-button>
      <el-button @click="dialogVisible = false">关闭</el-button>
    </template>
  </el-dialog>

</template>

<script>
import axios from "axios";

export default {
  name: "UserCreateDialog",
  data() {
    return {
      loading: false,
      dialogVisible: false,
      user: {
        account: "",
        password: "",
        email: "",
        disabled: 1,
        config_id: ""
      }
    }
  },
  methods: {
    show: function () {
      this.dialogVisible = true
    },
    create: function () {
      if (this.user.account === "") {
        this.$utils.Warning("提示", "用户不能为空")
        return
      }
      if (this.user.password === "") {
        this.$utils.Warning("提示", "密码不能为空")
        return
      }
      this.loading = true
      axios({
        method: "put",
        url: "/api/v1/user/create",
        data: this.user
      }).then(res => {
        let response = res.data
        this.$utils.Success(response.msg, "ID: " + response.data.id)
        this.user = {
          account: "",
          password: "",
          email: "",
          disabled: 1,
          config_id: ""
        }
        this.loading = false
        this.$emit("success")
      }).catch((err) => {
        this.$utils.HandleError(err)
        this.loading = false
      })
    }
  }
}
</script>

<style scoped>

</style>