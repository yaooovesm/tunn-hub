<template>
  <div>
    <el-dialog
        v-model="dialogVisible"
        width="50%"
        top="18vh"
        :close-on-click-modal="false"
        custom-class="default-dialog"
        destroy-on-close
        :append-to-body="true"
        draggable
    >
      <template #title>
        <div class="title">
          <div class="title-text">用户配置
            <el-tag
                type=""
                effect="dark"
                style="transform: translateY(-2px);margin-left: 10px;height: 25px"
            >
              用户&nbsp;{{ account }}
            </el-tag>
          </div>
        </div>
      </template>
      <div v-loading="loading">
        <el-card shadow="always" body-style="padding:0">
          <div class="title" style="margin-top: 20px;margin-bottom: 20px">
            <div class="title-text">网络导入</div>
          </div>
          <div style="padding: 20px">
            <div v-if="importRoutes.length>0">
              <span v-for="route in importRoutes" :key="route">
                <el-popover
                    placement="top-start"
                    :width="150"
                    trigger="hover"
                >
                  <template #reference>
                    <el-tag
                        closable
                        effect="dark"
                        type="info"
                        :disable-transitions="false"
                        @close="handleDeleteImportRoute(route)"
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
            <div v-else>
              <span style="font-size: 12px;color: #909399">还没有导入网络</span>
            </div>
            <route-import-selector
                style="margin-top: 8px"
                :imported="importRoutes"
                :exported="exportRoutes"
                @submit="handleAddImport"
            />
          </div>
        </el-card>
        <el-card shadow="always" body-style="padding:0" style="margin-top: 20px">
          <div class="title" style="margin-top: 20px;margin-bottom: 20px">
            <div class="title-text">网络暴露</div>
          </div>
          <div style="padding: 20px">
            <div v-if="exportRoutes.length>0">
              <span v-for="route in exportRoutes"
                    :key="route">
                <el-popover
                    placement="top-start"
                    :width="150"
                    trigger="hover"
                >
                  <template #reference>
                    <el-tag
                        closable
                        type="info"
                        effect="dark"
                        :disable-transitions="false"
                        @close="handleDeleteExportRoute(route)"
                        @click="$refs.edit_export.show(route)"
                        style="margin-right: 10px;margin-bottom: 5px;cursor: pointer"
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
            <div v-else>
              <span style="font-size: 12px;color: #909399">还没有暴露网络</span>
            </div>
            <div>
              <route-export-dialog ref="export_dialog" @submit="handleAddExport"/>
              <el-button @click="$refs.export_dialog.show()" size="small" style="margin-top: 8px">添加
              </el-button>
            </div>
          </div>
        </el-card>
        <el-card shadow="always" body-style="padding:0" v-if="$storage.IsAdmin()" style="margin-top: 20px">
          <div class="title" style="margin-top: 20px;margin-bottom: 20px">
            <div class="title-text">高级设置</div>
          </div>
          <div style="padding: 20px">
            <div>
              <span style="font-size: 12px;display: inline-block;transform: translateY(-2px);">地址静态</span>
              <el-checkbox v-model="enableStaticCIDR" label="启用" style="margin-bottom: 5px;margin-left: 10px"
                           size="small"/>
              <span
                  style="margin-bottom: 5px;font-size: 13px;color: #404040;display: inline-block;transform: translateY(-2px);">
                <el-tooltip
                    effect="dark"
                    content='设置后客户端内网地址将会被静态分配以设置的值'
                    placement="right">
                  <i class="iconfont icon-exclamation-circle"
                     style="color: #007bbb;font-size: 10px;margin-left: 5px;line-height: 24px"></i>
                </el-tooltip>
              </span>
            </div>

            <span v-if="enableStaticCIDR" style="display: block;margin-bottom: 5px;color: #909399;font-size: 12px">
              提示：若分配冲突则可能导致客户端无法接入网络，请确认后再修改。修改静态地址分配将在客户端重新连接后生效。
            </span>
            <el-input v-model="staticCIDR" v-if="enableStaticCIDR" placeholder="e.g. 192.168.1.1/24" size="small"/>
            <div style="margin-top: 20px;font-size: 12px">
              带宽限制
              <el-input size="small" :min="0" :max="1000" v-model="limit.bandwidth"
                        style="width: 60px;margin: 0 5px;"
                        placeholder="带宽">
              </el-input>
              Mbps
              <el-tooltip
                  effect="dark"
                  content='此处的设置将会同时影响上行和下行速率，带宽=上行带宽+下行带宽。设置 "0" 则不限制。'
                  placement="right">
                <i class="iconfont icon-exclamation-circle"
                   style="color: #007bbb;font-size: 10px;margin-left: 5px;line-height: 24px"></i>
              </el-tooltip>
              <span style="display: block;margin-top: 8px;color: #909399;font-size: 12px">
              提示：此处的设置将会同时影响上行和下行速率，带宽=上行带宽+下行带宽。设置 "0" 则不限制。设置将在客户端重新连接后生效。
              </span>
            </div>
            <div style="margin-top: 20px;font-size: 12px">
              流量限制
              <el-input size="small" :min="0" v-model="limit.flow"
                        style="width: 60px;margin: 0 5px;"
                        placeholder="流量">
              </el-input>
              M
              <el-tooltip
                  effect="dark"
                  content='流量限制单位为M，限制对象为服务器出方向(客户端入方向)流量。设置 "0" 则不限制。'
                  placement="right">
                <i class="iconfont icon-exclamation-circle"
                   style="color: #007bbb;font-size: 10px;margin-left: 5px;line-height: 24px"></i>
              </el-tooltip>
              <!--              <span style="display: block;margin-top: 8px;color: #909399;font-size: 12px">-->
              <!--              提示：流量限制指的是限制Hub端的发送即客户端接收，单位为M。设置 "0" 则不限制。设置将在客户端重新连接后生效。-->
              <!--              </span>-->
            </div>
          </div>
        </el-card>
      </div>
      <template #footer>
        <el-button :loading="loading" style="margin-left: 10px" @click="reset" type="warning">重置</el-button>
        <el-button :loading="loading" type="primary" @click="save">保存</el-button>
        <el-button :loading="loading" @click="close">关闭</el-button>
      </template>
      <route-export-edit-dialog ref="edit_export" @submit="handleExportModify"/>
    </el-dialog>
  </div>
</template>

<script>
import axios from "axios";
import {ElMessageBox} from "element-plus";
import RouteImportSelector from "@/components/config/RouteImportSelector";
import RouteExportDialog from "@/components/config/RouteExportDialog";
import RouteExportEditDialog from "@/components/config/RouteExportEditDialog";

export default {
  name: "UserConfig",
  components: {RouteExportEditDialog, RouteExportDialog, RouteImportSelector},
  data() {
    return {
      loading: false,
      dialogVisible: false,
      limit: {
        bandwidth: 0,
        flow: 0,
      },
      addImportValue: "",
      addExportValue: "",
      configId: "",
      account: "",
      importRoutes: [],
      exportRoutes: [],
      staticCIDR: "",
      enableStaticCIDR: false,
    }
  },
  methods: {
    handleExportModify: function (r) {
      for (let i = 0; i < this.exportRoutes.length; i++) {
        let origin = this.exportRoutes[i]
        if (origin.name === r.name) {
          this.exportRoutes[i] = r
          break
        }
      }
      this.$refs.edit_export.close()
    },
    reset: function () {
      this.loading = true
      ElMessageBox.confirm(
          '是否重置该配置文件',
          '警告',
          {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
          }
      ).then(() => {
        axios({
          method: "get",
          url: "/api/v1/cfg/reset/" + this.configId,
          data: {}
        }).then(res => {
          let response = res.data
          this.$utils.Success("操作成功", response.msg)
        }).catch((err) => {
          this.$utils.HandleError(err)
        }).finally(() => {
          this.loading = false
          this.load()
        })
      }).catch(() => {
      }).finally(() => {
        this.loading = false
      })
    },
    show: function (configId, account) {
      if (configId === "" || account === "") {
        this.$utils.Error("配置丢失", "请联系管理员")
        return
      }
      this.configId = configId
      this.account = account
      this.dialogVisible = true
      this.load()
    },
    close: function () {
      this.dialogVisible = false
      this.addExportValue = ""
      this.addImportValue = ""
    },
    save: function () {
      this.loading = true
      let data = {
        id: this.configId,
        device: {
          cidr: "",
          dns: "",
        }
      }
      if (this.enableStaticCIDR && this.staticCIDR !== "") {
        data.device.cidr = this.staticCIDR
      }
      data.routes = [...this.importRoutes, ...this.exportRoutes]
      if (this.$storage.IsAdmin()) {
        let bandwidth = Number(this.limit.bandwidth)
        if (bandwidth < 0 || bandwidth > 1000) {
          this.$utils.Error("设置错误", "带宽值范围在0-1000")
          return
        }
        data.limit = {
          bandwidth: bandwidth,
          flow: Number(this.limit.flow)
        }
      }
      axios({
        method: "post",
        url: "/api/v1/cfg/update",
        data: data
      }).then(() => {
        this.$utils.Success("提示", "更新配置成功")
        this.load()
        this.loading = false
      }).catch((err) => {
        this.$utils.HandleError(err)
        this.loading = false
      })
    },
    load: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/cfg/" + this.configId,
        data: {}
      }).then(res => {
        let response = res.data
        let routes = response.data.routes
        let imports = []
        let exports = []
        for (let i in routes) {
          if (routes[i].option === 'import') {
            imports.push(routes[i])
          } else if (routes[i].option === 'export') {
            exports.push(routes[i])
          }
        }
        this.importRoutes = imports
        this.exportRoutes = exports
        if (response.data.device.cidr !== "") {
          this.staticCIDR = response.data.device.cidr
          this.enableStaticCIDR = true
        } else {
          this.staticCIDR = ""
          this.enableStaticCIDR = false
        }
        this.limit = response.data.limit
        this.loading = false
      }).catch(() => {
        ElMessageBox.alert('加载用户配置失败', '错误', {
          confirmButtonText: '确认',
          callback: () => {
            this.dialogVisible = false
          },
        })
        this.loading = false
      })
    },
    handleDeleteImportRoute: function (route) {
      for (let i in this.importRoutes) {
        if (this.importRoutes[i].network === route.network) {
          this.importRoutes.splice(i, 1)
        }
      }
    },
    handleAddImport: function (routes) {
      for (let i = 0; i < routes.length; i++) {
        this.importRoutes.push({
          network: routes[i].network,
          option: "import"
        })
      }
    },
    handleDeleteExportRoute: function (route) {
      for (let i in this.exportRoutes) {
        if (this.exportRoutes[i].network === route.network) {
          this.exportRoutes.splice(i, 1)
        }
      }
    },
    handleAddExport: function (route) {
      this.exportRoutes.push(route)
      this.$refs.export_dialog.close()
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