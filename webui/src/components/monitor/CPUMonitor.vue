<template>
  <div>
    <el-popover
        placement="bottom"
        title="CPU占用"
        :width="230"
        trigger="hover"
    >
      <template #reference>
        <div>
          <el-progress :color="customColors" type="dashboard" :percentage="Number(cpu.usage.toFixed(1))"
                       style="position: relative">
            <template #default="{ percentage }">
              <span class="percentage-value" v-if="cpu.error===''">{{ percentage }}%</span>
              <span class="percentage-value" v-else>
                <i class="iconfont icon-exclamation-circle" style="color: #f56c6c;"></i>
              </span>
              <span class="percentage-label">处理器</span>
            </template>
          </el-progress>
          <el-progress :color="customColors" type="dashboard" :percentage="Number(cpu.app_used.toFixed(1))"
                       style="position: absolute;transform:translateX(-126px);opacity: 0.45;">
            <template #default>
              <span></span>
            </template>
          </el-progress>
        </div>
      </template>
      <template #default>
        <div v-if="cpu.error===''">
          <div class="detail-unit">
            <span>总利用率 </span> {{ cpu.usage.toFixed(2) }}%
          </div>
          <el-divider style="margin-top: 10px;margin-bottom: 10px"/>
          <div class="detail-unit">
            <span>TunnHub占用</span>
            {{ cpu.app_used.toFixed(2) }}%
          </div>
        </div>
        <div v-else style="font-size: 12px;color: #909399">
          {{ cpu.error }}
        </div>
      </template>
    </el-popover>
  </div>
</template>

<script>
export default {
  name: "CPUMonitor",
  data() {
    return {
      cpu: {
        usage: 0.0,
        app_used: 0.0,
        error: "",
      },
      customColors: [
        {color: '#5cb87a', percentage: 20},
        {color: '#5cb87a', percentage: 40},
        {color: '#1989fa', percentage: 60},
        {color: '#e6a23c', percentage: 90},
        {color: '#f56c6c', percentage: 100},
      ]
    }
  },
  methods: {
    set: function (data) {
      if (data === undefined) {
        this.cpu = {
          usage: 0.0,
          app_used: 0.0,
          error: "no data"
        }
        return
      }
      this.cpu = data
    }
  }
}
</script>

<style scoped>
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