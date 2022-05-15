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
            <div style="margin-bottom: 30px;text-align: center">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-popover
                      placement="bottom"
                      title="接收速率"
                      :width="230"
                      trigger="hover"
                  >
                    <template #reference>
                      <div>
                        <el-progress type="dashboard"
                                     :percentage="Number(status.rx.bandwidth_usage.toFixed(2))"
                                     style="position: relative">
                          <template #default>
                            <span class="percentage-value">{{ $utils.FormatBytesSpeed(status.rx.FlowSpeed) }}</span>
                            <span class="percentage-label">接收</span>
                          </template>
                        </el-progress>
                      </div>
                    </template>
                    <template #default>
                      <div>
                        <div class="detail-unit">
                          <span>消耗带宽
                            <span style="color: #007bbb;float: right">({{
                                status.rx.bandwidth_usage.toFixed(2)
                              }}%)</span>
                          </span> {{ status.rx.bandwidth.toFixed(1) }} Mbps
                        </div>
                        <div class="detail-unit">
                          <span>流量速率 </span> {{ $utils.FormatBytesSpeed(status.rx.FlowSpeed) }}
                        </div>
                        <div class="detail-unit">
                          <span>包速率</span>
                          {{ $utils.FormatPacketSpeed(status.rx.PacketSpeed) }}
                        </div>
                      </div>
                    </template>
                  </el-popover>
                </el-col>
                <el-col :span="12">
                  <el-popover
                      placement="bottom"
                      title="发送速率"
                      :width="230"
                      trigger="hover"
                  >
                    <template #reference>
                      <div>
                        <el-progress type="dashboard"
                                     :percentage="Number(status.tx.bandwidth_usage.toFixed(2))"
                                     style="position: relative">
                          <template #default>
                            <span class="percentage-value">{{ $utils.FormatBytesSpeed(status.tx.FlowSpeed) }}</span>
                            <span class="percentage-label">发送</span>
                          </template>
                        </el-progress>
                      </div>
                    </template>
                    <template #default>
                      <div>
                        <div class="detail-unit">
                          <span>消耗带宽
                            <span style="color: #007bbb;float: right">({{
                                status.tx.bandwidth_usage.toFixed(2)
                              }}%)</span>
                          </span>
                          {{ status.tx.bandwidth.toFixed(1) }} Mbps
                        </div>
                        <div class="detail-unit">
                          <span>流量速率 </span> {{ $utils.FormatBytesSpeed(status.tx.FlowSpeed) }}
                        </div>
                        <div class="detail-unit">
                          <span>包速率</span>
                          {{ $utils.FormatPacketSpeed(status.tx.PacketSpeed) }}
                        </div>
                      </div>
                    </template>
                  </el-popover>
                </el-col>
              </el-row>
            </div>
          </el-col>
          <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
            <div style="margin-bottom: 30px;text-align: center">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-popover
                      placement="bottom"
                      title="接收流量"
                      :width="230"
                      trigger="hover"
                  >
                    <template #reference>
                      <div>
                        <el-progress type="dashboard"
                                     :percentage="100"
                                     style="position: relative">
                          <template #default>
                            <span class="percentage-value">{{ $utils.FormatBytesSizeG(status.rx.Flow) }}</span>
                            <span class="percentage-label">接收流量</span>
                          </template>
                        </el-progress>
                      </div>
                    </template>
                    <template #default>
                      <div>
                        <div class="detail-unit">
                          <span>接收流量 </span> {{ $utils.FormatBytesSizeG(status.rx.Flow) }}
                        </div>
                        <div class="detail-unit">
                          <span>接收数据包 </span> {{ $utils.FormatPacketSize(status.rx.Packet) }}
                        </div>
                      </div>
                    </template>
                  </el-popover>
                </el-col>
                <el-col :span="12">
                  <el-popover
                      placement="bottom"
                      title="发送流量"
                      :width="230"
                      trigger="hover"
                  >
                    <template #reference>
                      <div>
                        <el-progress type="dashboard"
                                     :percentage="100"
                                     style="position: relative"
                                     :indeterminate="true"
                                     :duration="1">
                          <template #default>
                            <span class="percentage-value">{{ $utils.FormatBytesSizeG(status.tx.Flow) }}</span>
                            <span class="percentage-label">发送流量</span>
                          </template>
                        </el-progress>
                      </div>
                    </template>
                    <template #default>
                      <div>
                        <div class="detail-unit">
                          <span>发送流量 </span> {{ $utils.FormatBytesSizeG(status.tx.Flow) }}
                        </div>
                        <div class="detail-unit">
                          <span>发送数据包 </span> {{ $utils.FormatPacketSize(status.tx.Packet) }}
                        </div>
                      </div>
                    </template>
                  </el-popover>
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
      bandwidth: 1000,
      updateTime: new Date(),
      loading: false,
      timer: undefined,
      status: {
        rx: {
          Flow: 0,
          FlowSpeed: 0,
          Packet: 0,
          PacketSpeed: 0,
          bandwidth: 0,
          bandwidth_usage: 0,
        },
        tx: {
          Flow: 0,
          FlowSpeed: 0,
          Packet: 0,
          PacketSpeed: 0,
          bandwidth: 0,
          bandwidth_usage: 0,
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
        this.status.rx.bandwidth = this.status.rx.FlowSpeed / 1024 / 1024 * 8
        this.status.rx.bandwidth_usage = this.status.rx.bandwidth / this.bandwidth * 100
        this.status.tx.bandwidth = this.status.tx.FlowSpeed / 1024 / 1024 * 8
        this.status.tx.bandwidth_usage = this.status.tx.bandwidth / this.bandwidth * 100
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

.percentage-value {
  display: block;
  margin-top: 10px;
  font-size: 18px;
}

.percentage-label {
  display: block;
  margin-top: 10px;
  font-size: 12px;
}

.detail-unit {
  text-align: right;
  font-size: 12px;
  color: #007bbb;
}

.detail-unit span {
  color: #404040;
  float: left;
  display: inline-block;
}
</style>