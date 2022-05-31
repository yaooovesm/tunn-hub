<template>
  <div>
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px">
        <div class="title-text">创建用户</div>
      </div>
      <div style="padding: 20px" v-loading="loading">
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
        <div style="margin-top: 20px;text-align: right">
          <el-button type="primary" @click="create">创建</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "UserCreate",
  data() {
    return {
      loading: false,
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