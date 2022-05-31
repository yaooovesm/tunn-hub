<template>
  <div v-loading="loading">
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">性能监控
        </div>
      </div>
      <div style="padding: 35px 20px">
        <el-row :gutter="10">
          <el-col :span="8">
            <c-p-u-monitor ref="cpu"/>
          </el-col>
          <el-col :span="8">
            <memory-monitor ref="memory"/>
          </el-col>
          <el-col :span="8">
            <disk-monitor ref="disk"/>
          </el-col>
          <!--          <el-col :span="9">-->
          <!--            <span></span>-->
          <!--          </el-col>-->
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
import MemoryMonitor from "@/components/monitor/MemoryMonitor";
import axios from "axios";
import CPUMonitor from "@/components/monitor/CPUMonitor";
import DiskMonitor from "@/components/monitor/DiskMonitor";

export default {
  name: "MonitorOverview",
  components: {DiskMonitor, CPUMonitor, MemoryMonitor},
  data() {
    return {
      updateTime: new Date(),
      loading: false,
      timer: undefined,
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
      if (this.$refs.memory !== null) {
        this.$refs.memory.set(data.memory)
      }
      if (this.$refs.cpu != null) {
        this.$refs.cpu.set(data.cpu)
      }
      if (this.$refs.disk != null) {
        this.$refs.disk.set(data.disk)
      }
      this.updateTime = new Date()
    },
    update: function (silence) {
      if (!silence) {
        this.loading = true
      }
      axios({
        method: "get",
        url: "/api/v1/server/monitor",
        data: {}
      }).then(res => {
        let response = res.data
        this.set(response.data)
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
</style>