<template>
  <div>
    <el-dialog
        v-model="dialogVisible"
        width="60%"
        custom-class="default-dialog"
        draggable
    >
      <template #title>
        <div class="title">
          <div class="title-text">用户配置文件
            <el-tooltip
                effect="dark"
                :content="account"
                placement="bottom-start"
                v-if="account.length>10"
            >
              <el-tag
                  type=""
                  effect="dark"
                  style="transform: translateY(-2px);margin-left: 10px;height: 25px"
              >
                {{ account.length > 10 ? account.substring(0, 10) + "..." : account }}
              </el-tag>
            </el-tooltip>
            <el-tag
                type=""
                effect="dark"
                style="transform: translateY(-2px);margin-left: 10px;height: 25px"
                v-else
            >
              {{ account }}
            </el-tag>
          </div>
        </div>
      </template>
      <el-scrollbar height="500px">
        <el-row v-loading="loading" style="padding-bottom: 30px;padding-top: 20px">
          <el-col :span="22" :offset="1">
            <el-descriptions
                border
                :column="5"
                size="small"
                direction="vertical"
                title="全局配置"
            >
              <el-descriptions-item width="20%" label="服务器地址">
                {{ config.global.address }}
              </el-descriptions-item>
              <el-descriptions-item width="20%" label="服务器端口">
                {{ config.global.port }}
              </el-descriptions-item>
              <el-descriptions-item width="20%" label="通信协议">
                {{ config.global.protocol }}
              </el-descriptions-item>
              <el-descriptions-item width="20%" label="MTU">
                {{ config.global.mtu }}
              </el-descriptions-item>
              <el-descriptions-item width="20%" label="连接数">
                {{ config.global.multi_connection }}
              </el-descriptions-item>
            </el-descriptions>
          </el-col>
          <el-col :span="22" :offset="1">
            <el-descriptions
                border
                :column="5"
                size="small"
                direction="vertical"
                title="虚拟网卡"
                style="margin-top: 20px"
            >
              <el-descriptions-item width="20%" label="CIDR">
                {{ config.device.cidr }}
              </el-descriptions-item>
              <el-descriptions-item width="20%" label="DNS">
                {{ config.device.dns }}
              </el-descriptions-item>
            </el-descriptions>
          </el-col>
          <el-col :span="22" :offset="1">
            <el-descriptions
                border
                :column="5"
                size="small"
                direction="vertical"
                title="限制"
                style="margin-top: 20px"
            >
              <el-descriptions-item width="20%" label="带宽限制">
                {{ config.limit.bandwidth === 0 ? "无限制" : config.limit.bandwidth+" Mbps" }}
              </el-descriptions-item>
            </el-descriptions>
          </el-col>
          <el-col :span="22" :offset="1">
            <el-descriptions
                border
                :column="5"
                size="small"
                direction="vertical"
                title="验证服务器"
                style="margin-top: 20px"
            >
              <el-descriptions-item width="20%" label="地址">
                {{ config.auth.Address }}
              </el-descriptions-item>
              <el-descriptions-item width="20%" label="端口">
                {{ config.auth.Port }}
              </el-descriptions-item>
            </el-descriptions>
          </el-col>
          <el-col :span="22" :offset="1">
            <el-descriptions
                border
                :column="5"
                size="small"
                direction="vertical"
                title="数据处理"
                style="margin-top: 20px"
            >
              <el-descriptions-item width="20%" label="加密方式">
                {{ config.data_process.encrypt }}
              </el-descriptions-item>
              <el-descriptions-item width="20%" label="秘钥">
                {{ config.data_process.key }}
              </el-descriptions-item>
            </el-descriptions>
          </el-col>
          <el-col :span="22" :offset="1">
            <el-descriptions
                border
                :column="1"
                size="small"
                direction="vertical"
                title="路由网络"
                style="margin-top: 20px"
            >
              <el-descriptions-item width="100%" label="网络导入">
                <span v-if="route.imports.length===0">无</span>
                <div v-else>
                  <span v-for="route in route.imports" :key="route">
                    <el-popover
                        placement="top-start"
                        :width="150"
                        trigger="hover"
                    >
                      <template #reference>
                        <el-tag
                            type="info"
                            effect="dark"
                            :disable-transitions="false"
                            style="margin-right: 10px;margin-bottom: 5px"
                        >
                          {{ route.network }}
                        </el-tag>
                      </template>
                      <template #default>
                        <div class="detail-unit">
                          <span>名称 </span> {{ route.name === '' ? '未命名' : route.name }}
                        </div>
                        <div class="detail-unit">
                          <span>网络 </span> {{ route.network }}
                        </div>
                      </template>
                    </el-popover>
                  </span>
                </div>
              </el-descriptions-item>
              <el-descriptions-item width="100%" label="网络暴露">
                <span v-if="route.exports.length===0">无</span>
                <div v-else>
                  <span v-for="route in route.exports" :key="route">
                    <el-popover
                        placement="top-start"
                        :width="150"
                        trigger="hover"
                    >
                      <template #reference>
                        <el-tag
                            type="info"
                            effect="dark"
                            :disable-transitions="false"
                            style="margin-right: 10px;margin-bottom: 5px"
                        >
                          {{ route.network }}
                        </el-tag>
                      </template>
                      <template #default>
                        <div class="detail-unit">
                          <span>名称 </span> {{ route.name === '' ? '未命名' : route.name }}
                        </div>
                        <div class="detail-unit">
                          <span>网络 </span> {{ route.network }}
                        </div>
                      </template>
                    </el-popover>
                  </span>
                </div>
              </el-descriptions-item>
            </el-descriptions>
          </el-col>
        </el-row>
      </el-scrollbar>


      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: "ShowConfig",
  data() {
    return {
      account: "",
      loading: false,
      dialogVisible: false,
      config: {
        global: {
          address: "",
          port: 0,
          protocol: "",
          mtu: 0,
          pprof: 0,
          default_route: false,
          multi_connection: 0,
          storage_path: ""
        },
        route: [],
        device: {
          cidr: "",
          dns: ""
        },
        auth: {
          Address: "",
          Port: 0
        },
        data_process: {
          encrypt: "",
          key: null
        }
      },
      route: {
        imports: [],
        exports: [],
      },
      limit: {
        bandwidth: 0
      }
    }
  },
  methods: {
    show: function (config, account) {
      this.account = account
      this.config = config
      let routes = config.route
      let imports = []
      let exports = []
      for (let i in routes) {
        if (routes[i].option === 'import') {
          imports.push(routes[i])
        } else if (routes[i].option === 'export') {
          exports.push(routes[i])
        }
      }
      this.route.imports = imports
      this.route.exports = exports
      this.dialogVisible = true
    }
  }
}
</script>

<style scoped>
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