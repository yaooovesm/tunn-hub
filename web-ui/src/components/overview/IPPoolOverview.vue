<template>
  <div v-loading="loading">
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">地址池
        </div>
      </div>
      <div style="padding: 23px 20px">
        <el-row :gutter="30">
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
            <div style="text-align: left">
              <div>
                <span style="text-align: left;display: inline-block;font-size: 13px">地址池容量：<span style="color: #007bbb">{{
                    status.size
                  }}</span></span>
                <el-progress :percentage="Number(status.pool_percentage.toFixed(1))"
                             :color="customColors"
                             style="width: 100%">
                  <template #default="{percentage}">
                  <span
                      style="font-size: 10px;display: block;width: 140px;text-align: right">
                    <span style="color: #007bbb">{{ status.used }}</span> 已用,剩余 {{
                      percentage
                    }}% 可用</span>
                  </template>
                </el-progress>
              </div>
              <div style="margin-top: 20px">
                <span style="text-align: left;display: inline-block;font-size: 13px">网络可用地址：<span
                    style="color: #007bbb">{{ status.networkTotal }}</span></span>
                <el-progress
                    :percentage="Number(status.network_percentage.toFixed(1))"
                    :color="customColors"
                    style="width: 100%">
                  <template #default="{percentage}">
                  <span
                      style="font-size: 10px;display: block;width: 140px;text-align: right">
                    <span style="color: #007bbb">{{ status.used + status.static }}</span> 已用,剩余 {{
                      percentage
                    }}% 可用</span>
                  </template>
                </el-progress>
              </div>
              <el-divider/>
              <div style="margin-top: 20px;text-align: left">
                <el-row :gutter="30">
                  <el-col :span="12">
                    <div style="text-align: left;font-size: 13px">网络：<span
                        style="color: #007bbb;display: block;float: right">{{
                        status.network
                      }}</span></div>
                    <div style="text-align: left;font-size: 13px;margin-top: 6px">起始地址：<span
                        style="color: #007bbb;display: block;float: right">{{
                        status.start
                      }}</span></div>
                    <div style="text-align: left;font-size: 13px;margin-top: 6px">结束地址：<span
                        style="color: #007bbb;display: block;float: right">{{
                        status.end
                      }}</span></div>
                  </el-col>
                  <el-col :span="12">
                    <div style="text-align: left;font-size: 13px;">地址池利用比：<span
                        style="color: #007bbb;display: block;float: right">
                      {{ status.used }}
                      /
                      {{
                        status.size
                      }}</span></div>
                    <div style="text-align: left;font-size: 13px;margin-top: 6px">网络利用比：
                      <span
                          style="color: #007bbb;display: block;float: right">
                        {{ status.static + status.used }}
                        /
                        {{
                          status.networkTotal
                        }}</span></div>
                  </el-col>
                </el-row>
              </div>
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

export default {
  name: "IPPoolOverview",
  mounted() {
    this.update()
  },
  data() {
    return {
      updateTime: new Date(),
      loading: false,
      status: {
        end: "",
        network: "",
        size: 0,
        start: "",
        used: 0,
        networkTotal: 0,
        static: 0,
        network_percentage: 0,
        pool_percentage: 0,
      },
      customColors: [
        {color: '#f56c6c', percentage: 20},
        {color: '#e6a23c', percentage: 40},
        {color: '#1989fa', percentage: 60},
        {color: '#5cb87a', percentage: 80},
        {color: '#5cb87a', percentage: 100},
      ]
    }
  },
  methods: {
    update: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/server/ippool",
        data: {}
      }).then(res => {
        let response = res.data
        this.status = response.data
        let bit = this.status.network.substring(this.status.network.indexOf("/") + 1)
        this.status.networkTotal = Math.pow(2, 32 - bit) - 2
        this.status.network_percentage = (this.status.networkTotal - this.status.static - this.status.used) / this.status.networkTotal * 100
        this.status.pool_percentage = ((this.status.size - this.status.used) / this.status.size) * 100
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