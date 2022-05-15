<template>
  <el-popover
      placement="bottom"
      title="磁盘占用"
      :width="230"
      trigger="hover"
  >
    <template #reference>
      <div>
        <el-progress type="dashboard" :percentage="Number(disk.usage.toFixed(1))" style="position: relative">
          <template #default="{ percentage }">
            <span class="percentage-value">{{ percentage }}%</span>
            <span class="percentage-label">磁盘</span>
          </template>
        </el-progress>
      </div>
    </template>
    <template #default>
      <div>
        <div class="detail-unit">
          <span>储存总量 </span> {{ $utils.FormatBytesSizeG(disk.total) }}
        </div>
        <div class="detail-unit">
          <span>消耗储存<span style="color: #007bbb;float: right">({{ disk.usage.toFixed(2) }}%)</span> </span>
          {{ $utils.FormatBytesSizeG(disk.used) }}
        </div>
      </div>
    </template>
  </el-popover>
</template>

<script>
export default {
  name: "DiskMonitor",
  data() {
    return {
      disk: {
        total: 0,
        used: 0,
        usage: 0,
        error: ""
      }
    }
  },
  methods: {
    set: function (data) {
      if (data === undefined) {
        this.disk = {
          total: 0,
          used: 0,
          usage: 0,
          error: "no data"
        }
        return
      }
      this.disk = data
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