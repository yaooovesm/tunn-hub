<template>
  <div v-loading="loading">
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">地址池
        </div>
      </div>
      <div style="padding: 20px">
        <el-row :gutter="30">
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
            <div style="margin-bottom: 30px;text-align: left">
              <div class="dashboard-unit">
                <div class="subtitle" style="text-align: left;line-height: 12px;margin-left: 10px">
                  地址池 ({{ status.network === "" ? "未启用" : status.network }})
                  <el-tooltip
                      effect="dark"
                      placement="top"
                      v-if="status.network !== ''"
                  >
                    <template #content> 起始地址：{{ status.start }}<br/>结束地址：{{ status.end }}</template>
                    <i class="iconfont icon-info-circle"
                       style="font-size: 10px;margin-left: 5px;transform: translateY(2px)"></i>
                  </el-tooltip>
                </div>
                <div class="dashboard-unit-text">
                  <span style="color: #67C23A">{{ status.size - status.used }}</span> / {{ status.size }}
                </div>
                <div class="dashboard-unit-text-small"
                     style="color: #909399;margin-top: 8px;text-align: center;">
                  已分配 {{ status.used }}
                </div>
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
        used: 0
      },
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