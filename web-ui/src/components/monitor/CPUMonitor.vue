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
          <el-progress type="dashboard" :percentage="Number(cpu.usage.toFixed(1))" style="position: relative">
            <template #default="{ percentage }">
              <span class="percentage-value">{{ percentage }}%</span>
              <span class="percentage-label">处理器</span>
            </template>
          </el-progress>
          <el-progress color="black" type="dashboard" :percentage="Number(cpu.app_used.toFixed(1))"
                       style="position: absolute;transform:translateX(-126px);opacity: 0.45;">
            <template #default>
              <span></span>
            </template>
          </el-progress>
        </div>
      </template>
      <template #default>
        <div>
          <div class="detail-unit">
            <span>总利用率 </span> {{ cpu.usage.toFixed(2) }}%
          </div>
          <el-divider style="margin-top: 10px;margin-bottom: 10px"/>
          <div class="detail-unit">
            <span>TunnHub占用</span>
            {{ cpu.app_used.toFixed(2) }}%
          </div>
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
        usage: 0,
        app_used: 0,
        error: ""
      },
    }
  },
  methods: {
    set: function (data) {
      if (data === undefined) {
        this.cpu = {
          usage: 0,
          app_used: 0,
          error: "no data"
        }
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
  font-size: 25px;
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