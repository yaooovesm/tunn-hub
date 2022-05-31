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
                <!--接收-->
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
                <!--发送-->
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
                <!--处理速度-->
                <el-col :span="12">
                  <el-popover
                      placement="bottom"
                      title="网络状态"
                      :width="230"
                      trigger="hover"
                  >
                    <template #reference>
                      <div>
                        <el-progress type="dashboard"
                                     :percentage="Number(status.total.bandwidth_usage.toFixed(2))"
                                     style="position: relative">
                          <template #default>
                            <span class="percentage-value">{{
                                $utils.FormatPacketSpeed(status.total.PacketSpeed)
                              }}</span>
                            <span class="percentage-label">处理速度</span>
                          </template>
                        </el-progress>
                      </div>
                    </template>
                    <template #default>
                      <div>
                        <div class="detail-unit">
                          <div class="detail-unit">
                          <span>消耗带宽
                            <span style="color: #007bbb;float: right">({{
                                status.total.bandwidth_usage.toFixed(2)
                              }}%)</span>
                          </span>
                            {{ status.total.bandwidth.toFixed(1) }} Mbps
                          </div>
                          <span>数据包总量 </span> {{ $utils.FormatPacketSize(status.total.Packet) }}
                        </div>
                        <div class="detail-unit">
                          <span>数据包速度 </span> {{ $utils.FormatPacketSpeed(status.total.PacketSpeed) }}
                        </div>
                        <el-divider style="margin-top: 10px;margin-bottom: 10px"/>
                      </div>
                    </template>
                  </el-popover>
                </el-col>
                <!--流量统计-->
                <el-col :span="12">
                  <el-popover
                      placement="bottom"
                      title="流量统计"
                      :width="230"
                      trigger="hover"
                  >
                    <template #reference>
                      <div>
                        <el-progress type="dashboard"
                                     :percentage="Number(status.total.bandwidth_usage.toFixed(2))"
                                     style="position: relative">
                          <template #default>
                            <span class="percentage-value" style="font-size: 14px"><span style="color: #67c23a">↓</span> {{
                                $utils.FormatBytesSizeG(status.rx.Flow)
                              }}</span>
                            <span class="percentage-value" style="font-size: 14px;margin-top: 4px"><span
                                style="color: #67c23a">↑</span> {{
                                $utils.FormatBytesSizeG(status.tx.Flow)
                              }}</span>
                            <span class="percentage-label" style="margin-top: 5px">流量</span>
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
                        <el-divider style="margin-top: 10px;margin-bottom: 10px"/>
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
        },
        total: {
          bandwidth: 0,
          bandwidth_usage: 0,
          Packet: 0,
          PacketSpeed: 0,
        }
      },
    }
  },
  props: {
    passive: {
      type: Boolean,
      default: false
    }
  },
  mounted() {
    if (!this.passive) {
      this.update(false)
      this.timer = setInterval(() => {
        this.update(true)
      }, 5000)
    }
  },
  unmounted() {
    if (!this.passive) {
      clearInterval(this.timer)
    }
  },
  methods: {
    set: function (data) {
      this.status = data
      this.status.rx.bandwidth = this.status.rx.FlowSpeed / 1024 / 1024 * 8
      this.status.rx.bandwidth_usage = this.status.rx.bandwidth / this.bandwidth * 100
      this.status.tx.bandwidth = this.status.tx.FlowSpeed / 1024 / 1024 * 8
      this.status.tx.bandwidth_usage = this.status.tx.bandwidth / this.bandwidth * 100
      let total_bandwidth = this.status.tx.bandwidth + this.status.rx.bandwidth
      this.status.total = {
        bandwidth: total_bandwidth,
        bandwidth_usage: total_bandwidth / this.bandwidth * 100,
        Packet: this.status.rx.Packet + this.status.tx.Packet,
        PacketSpeed: this.status.rx.PacketSpeed + this.status.tx.PacketSpeed,
      }
      this.updateTime = new Date()
    },
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
        this.set(response.data)
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