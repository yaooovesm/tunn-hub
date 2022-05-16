<template>
  <div v-loading="loading">
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">服务器配置
        </div>
      </div>
      <div style="padding:17px 20px">
        <el-row>
          <el-col :span="24">
            <div style="text-align: left;">
              <el-descriptions
                  direction="vertical"
                  :column="3"
                  size="small"
                  border
              >
                <el-descriptions-item width="33.3%" label-class-name="overview-description-label"
                                      class-name="overview-description"
                                      label="服务器地址">
                  {{ config.address }}
                </el-descriptions-item>
                <el-descriptions-item width="33.3%" label-class-name="overview-description-label"
                                      class-name="overview-description"
                                      label="内网地址">
                  {{ config.cidr }}
                </el-descriptions-item>
                <el-descriptions-item width="33.3%" label-class-name="overview-description-label"
                                      class-name="overview-description"
                                      label="传输协议">
                  {{ config.protocol }}
                </el-descriptions-item>
                <el-descriptions-item width="33.3%" label-class-name="overview-description-label"
                                      class-name="overview-description"
                                      label="数据加密">
                  {{ config.encrypt }}
                </el-descriptions-item>
                <el-descriptions-item width="33.3%" label-class-name="overview-description-label"
                                      class-name="overview-description"
                                      label="MTU">
                  {{ config.mtu }}
                </el-descriptions-item>
                <el-descriptions-item width="33.3%" label-class-name="overview-description-label"
                                      class-name="overview-description"
                                      label="并行通道数">
                  {{ config.multi }}
                </el-descriptions-item>
                <el-descriptions-item width="33.3%" label-class-name="overview-description-label"
                                      class-name="overview-description"
                                      label="服务器版本">
                  {{ runtime.app }}
                </el-descriptions-item>
                <el-descriptions-item width="66.6%" label-class-name="overview-description-label"
                                      class-name="overview-description"
                                      label="运行平台">
                  <span>
                    <svg class="icon" aria-hidden="true" v-if="runtime.os==='windows'"
                         style="width: 1.1em;height: 1.0em;">
                      <use xlink:href="#icon-Windows"></use>
                    </svg>
                    <svg class="icon" aria-hidden="true" v-else-if="runtime.os==='linux'"
                         style="width: 1.2em;height: 1.1em;">
                      <use xlink:href="#icon-linux"></use>
                    </svg>
                    <svg class="icon" aria-hidden="true" v-else-if="runtime.os==='darwin'">
                      <use xlink:href="#icon-IOS"></use>
                    </svg>
                    {{ runtime.os === 'darwin' ? "OSX" : runtime.os }}_{{ runtime.arch }} {{
                      runtime.platform
                    }} {{ runtime.version }}
                  </span>
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
      runtime: {
        app: "",
        arch: "",
        os: "",
        platform: "",
        version: "",
      }
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
        this.runtime.app = response.data.runtime.app
        this.runtime.arch = response.data.runtime.arch
        this.runtime.os = response.data.runtime.os
        this.runtime.platform = response.data.runtime.platform
        this.runtime.version = response.data.runtime.version
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
  padding: 5px 10px !important;
  font-size: 13px !important;
}

.overview-description {
  padding: 5px 10px !important;
  font-size: 11px !important;
}
</style>