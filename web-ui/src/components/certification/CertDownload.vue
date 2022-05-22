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
        <el-button @click="download" v-loading="loading" size="small" :disabled="!allow">下载</el-button>
        <el-checkbox v-model="allow" :value="allow" size="small"
                     style="margin-left: 10px;">我已知晓风险
        </el-checkbox>
      </div>
    </el-card>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "CertDownload",
  data() {
    return {
      allow: false,
      loading: false
    }
  },
  methods: {
    download: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/server/cert/download",
        data: {},
      }).then(res => {
        let blob = new Blob([res.data], {type: "application/octet-stream;"});
        let url = window.URL.createObjectURL(blob); // 创建一个临时的url指向blob对象
        let a = document.createElement("a");
        a.href = url;
        let filename
        let disposition = res.headers["content-disposition"]
        filename = disposition.replace("attachment;filename=", "")
        a.download = filename.replaceAll("\"", "");
        a.click();
        // 释放这个临时的对象url
        window.URL.revokeObjectURL(url);
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