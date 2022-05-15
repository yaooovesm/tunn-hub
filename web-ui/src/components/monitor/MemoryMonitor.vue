<template>
  <div>
    <el-popover
        placement="bottom"
        title="内存占用"
        :width="230"
        trigger="hover"
    >
      <template #reference>
        <div>
          <el-progress type="dashboard" :percentage="total_percentage" style="position: relative">
            <template #default="{ percentage }">
              <span class="percentage-value">{{ percentage }}%</span>
              <span class="percentage-label">内存</span>
            </template>
          </el-progress>
          <el-progress color="black" type="dashboard" :percentage="app_percentage"
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
            <span>内存总量 </span> {{ $utils.FormatBytesSizeM(memory.total) }}
          </div>
          <div class="detail-unit">
            <span>消耗内存<span style="color: #007bbb;float: right">({{ memory.usage.toFixed(2) }}%)</span> </span>
            {{ $utils.FormatBytesSizeM(memory.used) }}
          </div>
          <el-divider style="margin-top: 10px;margin-bottom: 10px"/>
          <div class="detail-unit">
            <span>TunnHub消耗 <span style="color: #007bbb;float: right">({{
                ((this.memory.app_used / this.memory.total) * 100).toFixed(2)
              }}%)</span> </span>
            {{ $utils.FormatBytesSizeM(memory.app_used) }}
          </div>
        </div>
      </template>
    </el-popover>
  </div>
</template>

<script>
export default {
  name: "MemoryMonitor",
  data() {
    return {
      memory: {
        total: 0,
        used: 0,
        usage: 0,
        app_used: 0,
        swap_total: 0,
        swap_used: 0,
        swap_usage: 0,
        error: ""
      },
      total_percentage: 0,
      app_percentage: 0
    }
  },
  methods: {
    set: function (data) {
      if (data === undefined) {
        this.memory = {
          total: 0,
          used: 0,
          usage: 0,
          app_used: 0,
          swap_total: 0,
          swap_used: 0,
          swap_usage: 0,
          error: "no data"
        }
        return
      }
      this.memory = data
      this.total_percentage = Number(this.memory.usage.toFixed(1))
      this.app_percentage = Number((this.memory.app_used / this.memory.total).toFixed(1)) * 100
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