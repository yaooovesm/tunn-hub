<template>
  <div v-loading="loading">
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">流量监控
        </div>
      </div>
      <div style="padding: 20px">
        <el-row :gutter="30">
          <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
            <div style="margin-bottom: 30px;text-align: left">
              <div class="subtitle-large">速率</div>
              <el-row :gutter="20">
                <el-col :span="12">
                  <div class="dashboard-unit">
                    <div class="subtitle" style="text-align: left;line-height: 10px;margin-left: 10px">
                      服务器接收(RX)
                    </div>
                    <div class="dashboard-unit-text">
                      <span style="color: rgba(0,123,187,0.8)">{{ $utils.FormatBytesSpeed(status.rx.FlowSpeed) }}</span>
                    </div>
                    <div class="dashboard-unit-text-small"
                         style="color: #909399;margin-top: 8px;text-align: center;">
                      包速率 {{ $utils.FormatPacketSpeed(status.rx.PacketSpeed) }}
                    </div>
                  </div>
                </el-col>
                <el-col :span="12">
                  <div class="dashboard-unit">
                    <div class="subtitle" style="text-align: left;line-height: 10px;margin-left: 10px">
                      服务器发送(TX)
                    </div>
                    <div class="dashboard-unit-text">
                      <span style="color: rgba(0,123,187,0.8)">
                      {{ $utils.FormatBytesSpeed(status.tx.FlowSpeed) }}
                      </span>
                    </div>
                    <div class="dashboard-unit-text-small"
                         style="color: #909399;margin-top: 8px;text-align: center;">
                      包速率 {{ $utils.FormatPacketSpeed(status.tx.PacketSpeed) }}
                    </div>
                  </div>
                </el-col>
              </el-row>
            </div>
          </el-col>
          <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
            <div style="margin-bottom: 30px;text-align: left">
              <div class="subtitle-large">流量统计</div>
              <el-row :gutter="20">
                <el-col :span="12">
                  <div class="dashboard-unit">
                    <div class="subtitle" style="text-align: left;line-height: 10px;margin-left: 10px">
                      服务器接收(RX)
                    </div>
                    <div class="dashboard-unit-text">
                      <span style="color: rgba(0,123,187,0.8)">
                        {{ $utils.FormatBytesSize(status.rx.Flow) }}
                      </span>
                    </div>
                    <div class="dashboard-unit-text-small"
                         style="color: #909399;margin-top: 8px;text-align: center;">
                      数据包 {{ $utils.FormatPacketSize(status.rx.Packet) }}
                    </div>
                  </div>
                </el-col>
                <el-col :span="12">
                  <div class="dashboard-unit">
                    <div class="subtitle" style="text-align: left;line-height: 10px;margin-left: 10px">
                      服务器发送(TX)
                    </div>
                    <div class="dashboard-unit-text">
                      <span style="color: rgba(0,123,187,0.8)">
                        {{ $utils.FormatBytesSize(status.tx.Flow) }}
                      </span>
                    </div>
                    <div class="dashboard-unit-text-small"
                         style="color: #909399;margin-top: 8px;text-align: center;">
                      数据包 {{ $utils.FormatPacketSize(status.tx.Packet) }}
                    </div>
                  </div>
                </el-col>
              </el-row>
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
  name: "FlowOverview",
  data() {
    return {
      updateTime: new Date(),
      loading: false,
      timer: undefined,
      status: {
        rx: {
          Flow: 0,
          FlowSpeed: 0,
          Packet: 0,
          PacketSpeed: 0
        },
        tx: {
          Flow: 0,
          FlowSpeed: 0,
          Packet: 0,
          PacketSpeed: 0
        }
      },
    }
  },
  mounted() {
    this.update(false)
    this.timer = setInterval(() => {
      this.update(true)
    }, 5000)
  },
  unmounted() {
    clearInterval(this.timer)
  },
  methods: {
    update: function (silence) {
      if (!silence) {
        this.loading = true
      }
      axios({
        method: "get",
        url: "/api/v1/server/flow",
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