<template>
  <div>
    <el-dialog
        v-model="dialogVisible"
        width="500px"
        top="30vh"
        :close-on-click-modal="false"
        custom-class="default-dialog"
        destroy-on-close
        draggable
    >
      <template #title>
        <div class="title">
          <div class="title-text">添加网络暴露</div>
        </div>
      </template>
      <el-row :gutter="10" v-loading="loading">
        <el-col :span="20" :offset="2">
          <el-form
              label-position="top"
              label-width="100px"
              :model="route"
          >
            <el-form-item label="名称">
              <el-input v-model="route.name"/>
            </el-form-item>
            <el-form-item label="网络">
              <el-input v-model="route.addr" style="width: calc(100% - 100px);"/>
              <span
                  style="display: inline-block;margin-right: 10px;margin-left: 10px;font-size: 15px;width: 10px">/</span>
              <el-select v-model="route.maskNum" style="width: 70px" placeholder="24" value="">
                <el-option
                    v-for="m in masks"
                    :key="m"
                    :label="m"
                    :value="m"
                />
              </el-select>
            </el-form-item>
          </el-form>
        </el-col>
      </el-row>
      <template #footer>
        <el-button @click="confirm" type="primary">确认</el-button>
        <el-button @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: "RouteExportDialog",
  data() {
    return {
      loading: false,
      dialogVisible: false,
      masks: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
        17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32],
      route: {
        maskNum: 24,
        addr: "",
        name: "",
      }
    }
  },
  methods: {
    confirm: function () {
      let route = {
        name: this.route.name,
        network: this.route.addr + "/" + this.route.maskNum,
        option: "export"
      }
      this.$emit("submit", route)
    },
    show: function () {
      this.dialogVisible = true
    },
    close: function () {
      this.dialogVisible = false
    }
  }
}
</script>

<style scoped>

</style>