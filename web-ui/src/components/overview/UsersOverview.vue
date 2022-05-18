<template>
  <div v-loading="loading">
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">用户概况
        </div>
      </div>
      <div style="padding: 27px 20px">
        <el-row :gutter="30">
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" style="text-align: left">
            <div>
              <span style="text-align: left;display: inline-block;font-size: 13px">在线用户</span>
              <el-progress :percentage="Number((status.online_percentage.toFixed(1)))"
                           :color="customColors"
                           style="width: 100%">
                <template #default>
                  <span style="font-size: 10px;display: block;width: 140px;text-align: right">
                    <span style="display: inline-block">在线率
                      <span style="color: #007bbb">{{
                          ((status.online) / status.allow * 100).toFixed(1)
                        }}% </span>
                    </span>
                     当前 <span style="color: #007bbb">{{
                      status.online
                    }}</span> 在线
                  </span>
                </template>
              </el-progress>
            </div>
            <div style="margin-top: 20px">
              <span style="text-align: left;display: inline-block;font-size: 13px">已注册用户</span>
              <el-progress :percentage="Number(status.allow_percentage.toFixed(1))"
                           color="#1989fa"
                           style="width: 100%">
                <template #default>
                  <span style="font-size: 10px;display: block;width: 140px;text-align: right">
                    总量 <span style="color: #007bbb">{{ status.total }}</span> ，当前 <span
                      style="color: #007bbb">{{ status.allow }}</span> 可用</span>
                </template>
              </el-progress>
            </div>
            <el-divider/>
            <user-create-dialog ref="user_create"/>
            <div style="margin-bottom: 30px;text-align: left">
              <el-button type="primary" size="small" @click="$router.push({path:'/dashboard/users'})">用户管理</el-button>
              <el-button type="primary" size="small" @click="$refs.user_create.show()">快速创建</el-button>
            </div>
          </el-col>
        </el-row>
      </div>
      <div style="font-size: 12px;color: #808080;text-align: right;padding: 5px 10px">
        更新于
        {{ $utils.FormatDate("YYYY/mm/dd HH:MM:SS", updateTime) }}&nbsp;
        <el-button type="text" @click="update(false)" style="font-size: 12px;height: 12px;line-height: 13px">刷新
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import axios from "axios";
import UserCreateDialog from "@/components/control/UserCreateDialog";

export default {
  name: "UsersOverview",
  components: {UserCreateDialog},
  props: {
    passive: {
      type: Boolean,
      default: false
    }
  },
  mounted() {
    if (!this.passive) {
      this.update()
    }
  },
  data() {
    return {
      updateTime: new Date(),
      loading: false,
      status: {
        disabled: 0,
        offline: 0,
        online: 0,
        total: 0,
        allow: 0,
        allow_percentage: 0,
        online_percentage: 0,
      },
      customColors: [
        {color: '#5cb87a', percentage: 20},
        {color: '#5cb87a', percentage: 40},
        {color: '#1989fa', percentage: 60},
        {color: '#e6a23c', percentage: 80},
        {color: '#e6a23c', percentage: 100},
      ]
    }
  },
  methods: {
    update: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/user/general",
        data: {}
      }).then(res => {
        let response = res.data
        this.status = response.data
        this.status.offline = this.status.total - this.status.disabled - this.status.online
        this.status.allow = this.status.total - this.status.disabled
        this.status.online_percentage = (this.status.online / this.status.allow) * 100
        this.status.allow_percentage = (this.status.allow / this.status.total) * 100
        this.updateTime = new Date()
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
.subtitle {
  font-size: 12px;
  text-align: center;
  margin-bottom: 8px;
  color: #606266;
}

.subtitle-large {
  font-size: 20px;
  text-align: left;
  margin-bottom: 10px;
  color: #606266;
}

.dashboard-unit {
  color: #606266;
  background-color: #f2f2f2;
  border-radius: 4px;
  height: 150px;
  padding: 10px 0;
}

.dashboard-unit-text {
  color: #404040;
  font-size: 25px;
  text-align: center;
  margin-top: 50px;
}

.dashboard-unit-text-small {
  color: #404040;
  font-size: 12px;
  text-align: right;
  margin-top: 8px;
}
</style>