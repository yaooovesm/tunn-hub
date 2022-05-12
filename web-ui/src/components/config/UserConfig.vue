<template>
  <div v-loading="loading">
    <el-dialog
        v-model="dialogVisible"
        width="50%"
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
      <el-card shadow="always" body-style="padding:0">
        <div class="title" style="margin-top: 20px;margin-bottom: 20px">
          <div class="title-text">路由导入</div>
        </div>
        <div style="padding: 20px">
          <el-tag
              v-for="route in importRoutes"
              :key="route"
              closable
              effect="dark"
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
              <el-button @click="handleAddImport">+</el-button>
            </template>
          </el-input>
          <el-button v-else style="margin-left: 5px" size="small" @click="addImportVisible = true">
            + 添加导入
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
              effect="dark"
              :disable-transitions="false"
              @close="handleDeleteImportRoute(route)"
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
              <el-button @click="handleAddExport">+</el-button>
            </template>
          </el-input>
          <el-button v-else style="margin-left: 5px" size="small" @click="addExportVisible = true">
            + 添加暴露
          </el-button>
        </div>
      </el-card>
      <el-card shadow="always" body-style="padding:0" v-if="$storage.IsAdmin()" style="margin-top: 20px">
        <div class="title" style="margin-top: 20px;margin-bottom: 20px">
          <div class="title-text">高级设置</div>
        </div>
        <div style="padding: 20px">
        </div>
      </el-card>
      <template #footer>
        <el-button type="primary">保存</el-button>
        <el-button @click="dialogVisible = false">关闭</el-button>
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
    load: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/cfg/" + this.configId,
        data: {}
      }).then(res => {
        let response = res.data
        console.log(response.data)
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
          option: "import"
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