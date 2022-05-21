<template>
  <div>
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">下载证书</div>
      </div>
      <div style="padding: 20px;text-align: left">
        <div style="font-size: 12px;color: #909399;margin-bottom: 20px">
          下载证书有泄漏风险，请谨慎操作
        </div>
        <el-button size="small" :disabled="!allow">下载</el-button>
        <el-checkbox v-model="allow" :value="allow" size="small"
                     style="margin-left: 10px;">我已知晓风险
        </el-checkbox>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: "CertDownload",
  data() {
    return {
      allow: false
    }
  }
}
</script>

<style scoped>

</style>