<template>
  <div v-loading="loading">
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">客户端
          <div class="online-status" style="margin-left: 10px">
            <div :class="Status.online?'circle online-circle':'circle offline-circle'"></div>
            <div :class="Status.online?'online-text':'offline-text'">
              {{ Status.online ? "在线 (" + Status.address + ")" : "离线" }}
            </div>
          </div>
        </div>
      </div>
      <div style="padding: 20px">
        <el-row :gutter="30">
          <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
            <!--RX数据显示-->
            <div style="margin-bottom: 30px">
              <div class="subtitle">服务器接收(RX)</div>
              <el-descriptions
                  border
                  :column="3"
                  size="small"
                  direction="vertical"
              >
                <el-descriptions-item width="33.3%" label="流量">{{
                    $utils.FormatBytesSize(Status.rx.Flow)
                  }}
                </el-descriptions-item>
                <el-descriptions-item width="33.3%" label="数据包">{{
                    $utils.FormatPacketSize(Status.rx.Packet)
                  }}
                </el-descriptions-item>
                <el-descriptions-item width="33.3%" label="速度">{{ $utils.FormatBytesSpeed(Status.rx.FlowSpeed) }} ({{
                    $utils.FormatPacketSpeed(Status.rx.PacketSpeed)
                  }})
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </el-col>
          <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
            <!--TX数据显示-->
            <div style="margin-bottom: 30px">
              <div class="subtitle">服务器发送(TX)</div>
              <el-descriptions
                  border
                  :column="3"
                  size="small"
                  direction="vertical"
              >
                <el-descriptions-item width="33.3%" label="流量">{{
                    $utils.FormatBytesSize(Status.tx.Flow)
                  }}
                </el-descriptions-item>
                <el-descriptions-item width="33.3%" label="数据包">{{
                    $utils.FormatPacketSize(Status.tx.Packet)
                  }}
                </el-descriptions-item>
                <el-descriptions-item width="33.3%" label="速度">{{ $utils.FormatBytesSpeed(Status.tx.FlowSpeed) }} ({{
                    $utils.FormatPacketSpeed(Status.tx.PacketSpeed)
                  }})
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" v-if="Status.online">
            <!--连接信息显示-->
            <div style="margin-bottom: 30px">
              <div class="subtitle">连接信息
                <el-tooltip
                    effect="dark"
                    :content="Status.uuid"
                    placement="bottom-start"
                >
                  <i v-if="Status.uuid!==''" class="iconfont icon-info-circle"
                     style="color: #007bbb;font-size: 12px"></i>
                </el-tooltip>
              </div>
              <el-descriptions
                  border
                  :column="4"
                  size="small"
                  direction="vertical"
              >
                <el-descriptions-item width="25%" label="内网地址">{{ Status.config.device.cidr }}</el-descriptions-item>
                <el-descriptions-item width="25%" label="通信协议">{{
                    Status.config.global.protocol
                  }}
                </el-descriptions-item>
                <el-descriptions-item width="25%" label="数据处理">{{
                    Status.config.data_process.encrypt
                  }}
                </el-descriptions-item>
                <el-descriptions-item width="25%" label="客户端版本">{{
                    Status.config.runtime.app
                  }}
                </el-descriptions-item>
                <el-descriptions-item width="25%" label="平台">
                  <svg class="icon" aria-hidden="true" v-if="Status.config.runtime.os==='windows'">
                    <use xlink:href="#icon-Windows"></use>
                  </svg>
                  <svg class="icon" aria-hidden="true" v-else-if="Status.config.runtime.os==='linux'">
                    <use xlink:href="#icon-linux"></use>
                  </svg>
                  <svg class="icon" aria-hidden="true" v-else-if="Status.config.runtime.os==='darwin'">
                    <use xlink:href="#icon-IOS"></use>
                  </svg>
                  {{
                    Status.config.runtime.os
                  }}
                </el-descriptions-item>
                <el-descriptions-item width="25%" label="架构">
                  {{ Status.config.runtime.arch }}
                </el-descriptions-item>
                <el-descriptions-item width="50%" label="系统">
                  {{ Status.config.runtime.platform }}
                </el-descriptions-item>
              </el-descriptions>
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
import "../../assets/icon/iconfont"


export default {
  name: "ClientCard",
  props: {
    id: {
      type: String
    },
  },
  data() {
    return {
      timer: undefined,
      updateTime: new Date(),
      loading: false,
      Status: {
        address: "",
        online: false,
        rx: {
          Flow: 0,
          FlowSpeed: 0,
          Packet: 0,
          PacketSpeed: 0,
        },
        tx: {
          Flow: 0,
          FlowSpeed: 0,
          Packet: 0,
          PacketSpeed: 0,
        },
        config: {
          global: {
            protocol: "",
            mtu: 0,
            restart: false,
            multi_connection: 0,
          },
          device: {
            cidr: "",
          },
          data_process: {
            encrypt: "",
          },
        }
      }

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
        url: "/api/v1/user/status/" + this.id,
        data: {}
      }).then(res => {
        let response = res.data
        this.Status = response.data
        this.updateTime = new Date()
        this.loading = false
      }).catch((err) => {
        this.$utils.HandleError(err)
        clearInterval(this.timer)
        this.loading = false
      })
    }
  }
}
</script>

<style scoped>
.subtitle {
  font-size: 12px;
  text-align: left;
  margin-bottom: 8px;
  color: #606266;
}

.online-status {
  display: inline-block;
  height: 25px;
  vertical-align: center;
  transform: translateY(-1px);
}
</style>