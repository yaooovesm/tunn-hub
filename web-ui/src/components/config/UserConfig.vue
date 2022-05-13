<template>
  <div>
    <el-dialog
        v-model="dialogVisible"
        width="50%"
        top="18vh"
        :close-on-click-modal="false"
        custom-class="default-dialog"
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
            <el-tag
                v-for="route in importRoutes"
                :key="route"
                closable
                effect="dark"
                type="info"
                :disable-transitions="false"
                @close="handleDeleteImportRoute(route)"
                style="margin-right: 10px;margin-bottom: 5px"
            >
              {{ route.network }}
            </el-tag>
            <el-input
                v-if="addImportVisible"
                v-model="addImportValue"
                style="margin-top: 10px"
                size="small"
                placeholder="e.g. 192.168.1.0/24"
            >
              <template #append>
                <el-button @click="handleAddImport" size="small">添加</el-button>
              </template>
            </el-input>
            <el-button v-else style="margin-left: 5px" size="small" @click="addImportVisible = true">
              + 添加
            </el-button>
          </div>
        </el-card>
        <el-card shadow="always" body-style="padding:0" style="margin-top: 20px">
          <div class="title" style="margin-top: 20px;margin-bottom: 20px">
            <div class="title-text">网络暴露</div>
          </div>
          <div style="padding: 20px">
            <el-tag
                v-for="route in exportRoutes"
                :key="route"
                closable
                type="info"
                effect="dark"
                :disable-transitions="false"
                @close="handleDeleteExportRoute(route)"
                style="margin-right: 10px;margin-bottom: 5px"
            >
              {{ route.network }}
            </el-tag>
            <el-input
                v-if="addExportVisible"
                v-model="addExportValue"
                style="margin-top: 10px"
                placeholder="e.g. 192.168.1.0/24"
                size="small"
            >
              <template #append>
                <el-button @click="handleAddExport" size="small">添加</el-button>
              </template>
            </el-input>
            <el-button v-else style="margin-left: 5px" size="small" @click="addExportVisible = true">
              + 添加
            </el-button>
          </div>
        </el-card>
        <el-card shadow="always" body-style="padding:0" v-if="$storage.IsAdmin()" style="margin-top: 20px">
          <div class="title" style="margin-top: 20px;margin-bottom: 20px">
            <div class="title-text">高级设置</div>
          </div>
          <div style="padding: 20px">
            <el-checkbox v-model="enableStaticCIDR" label="地址静态分配" style="margin-bottom: 5px" size="small"/>
            <span
                style="margin-bottom: 5px;font-size: 13px;color: #404040;display: inline-block;transform: translateY(-2px);">
            <el-tooltip
                effect="dark"
                content='设置后客户端内网地址将会被静态分配以设置的值'
                placement="right">
              <i class="iconfont icon-exclamation-circle"
                 style="color: #909399;font-size: 10px;margin-left: 5px;line-height: 24px"></i>
            </el-tooltip>
          </span>
            <span v-if="enableStaticCIDR" style="display: block;margin-bottom: 5px;color: #909399;font-size: 12px">
              提示：若分配冲突则可能导致客户端无法接入网络，请确认后再修改。修改静态地址分配将在客户端重新连接后生效。
            </span>
            <el-input v-model="staticCIDR" v-if="enableStaticCIDR" placeholder="e.g. 192.168.1.1/24" size="small"/>
          </div>
        </el-card>
      </div>
      <template #footer>
        <el-button :loading="loading" type="primary" @click="save">保存</el-button>
        <el-button :loading="loading" @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import axios from "axios";
import {ElMessageBox} from "element-plus";

export default {
  name: "UserConfig",
  data() {
    return {
      loading: false,
      dialogVisible: false,
      addImportValue: "",
      addImportVisible: false,
      addExportValue: "",
      addExportVisible: false,
      configId: "",
      account: "",
      importRoutes: [],
      exportRoutes: [],
      staticCIDR: "",
      enableStaticCIDR: false,
    }
  },
  methods: {
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
    handleAddImport: function () {
      if (this.addImportValue !== '') {
        this.importRoutes.push({
          network: this.addImportValue,
          option: "import"
        })
      }
      this.addImportValue = ""
      this.addImportVisible = false
    },
    handleDeleteExportRoute: function (route) {
      for (let i in this.exportRoutes) {
        if (this.exportRoutes[i].network === route.network) {
          this.exportRoutes.splice(i, 1)
        }
      }
    },
    handleAddExport: function () {
      if (this.addExportValue !== '') {
        this.exportRoutes.push({
          network: this.addExportValue,
          option: "export"
        })
      }
      this.addExportValue = ""
      this.addExportVisible = false
    },
  }
}
</script>

<style scoped>

</style>