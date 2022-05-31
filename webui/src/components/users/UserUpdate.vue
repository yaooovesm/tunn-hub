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
        <div class="title-text">修改用户
          <el-tooltip
              effect="dark"
              :content="id"
              placement="bottom-start"
              v-if="id.length>10"
          >
            <el-tag
                type=""
                effect="dark"
                style="transform: translateY(-2px);margin-left: 10px;height: 25px"
            >
              ID&nbsp;{{ id.length > 10 ? id.substring(0, 10) + "..." : id }}
            </el-tag>
          </el-tooltip>
          <el-tag
              type=""
              effect="dark"
              style="transform: translateY(-2px);margin-left: 10px;height: 25px"
              v-else
          >
            ID&nbsp;{{ id }}
          </el-tag>
        </div>
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
  name: "UserUpdate",
  emits: ['success'],
  props: {},
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
    show: function (id) {
      this.id = id
      this.dialogVisible = true
      this.load()
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
        url: "/api/v1/user/update/" + this.id,
        data: data
      }).then(() => {
        this.$utils.Success("修改成功", "ID " + this.id)
        this.$emit("success")
        this.id = ""
        this.info = {
          account: "",
          password: "",
          email: "",
        }
        this.dialogVisible = false
        this.loading = false
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