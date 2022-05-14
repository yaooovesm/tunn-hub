<template>
  <div v-loading="loading">
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">服务器配置
        </div>
      </div>
      <div style="padding: 20px">
        <el-row>
          <el-col :span="24">
            <div style="margin-bottom: 30px;text-align: left;padding-bottom: 9px">
              <el-descriptions
                  direction="vertical"
                  :column="3"
                  size="small"
                  border
              >
                <el-descriptions-item label-class-name="overview-description-label" label="服务器地址" width="33.3%">
                  {{ config.address }}
                </el-descriptions-item>
                <el-descriptions-item label-class-name="overview-description-label" label="内网地址" width="33.3%">
                  {{ config.cidr }}
                </el-descriptions-item>
                <el-descriptions-item label-class-name="overview-description-label" label="传输协议" width="33.3%">
                  {{ config.protocol }}
                </el-descriptions-item>
                <el-descriptions-item label-class-name="overview-description-label" label="数据加密" width="33.3%">
                  {{ config.encrypt }}
                </el-descriptions-item>
                <el-descriptions-item label-class-name="overview-description-label" label="MTU" width="33.3%">
                  {{ config.mtu }}
                </el-descriptions-item>
                <el-descriptions-item label-class-name="overview-description-label" label="并行通道数" width="33.3%">
                  {{ config.multi }}
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

export default {
  name: "ServerConfigOverview",
  mounted() {
    this.update()
  },
  data() {
    return {
      updateTime: new Date(),
      loading: false,
      config: {
        protocol: "",
        address: "",
        cidr: "",
        encrypt: "",
        mtu: 0,
        multi: 0,
      },
    }
  },
  methods: {
    update: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/server/config",
        data: {}
      }).then(res => {
        let response = res.data
        this.config.protocol = response.data.global.protocol
        this.config.address = response.data.global.address
        this.config.cidr = response.data.device.cidr
        this.config.encrypt = response.data.data_process.encrypt
        this.config.mtu = response.data.global.mtu
        this.config.multi = response.data.global.multi_connection
        this.raw = response.data
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

<style>
.overview-description-label {
  background-color: #f2f2f2 !important;
}
</style>