<template>
  <div style="padding-top: 50px">
    <el-row v-if="$storage.User.isLogin" :gutter="20">
      <el-col :span="22" :offset="1">
        <cert-create/>
      </el-col>
      <el-col :span="22" :offset="1" style="margin-top: 30px">
        <cert-download/>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import CertCreate from "@/components/certification/CertCreate";
import CertDownload from "@/components/certification/CertDownload";

export default {
  name: "CertPanel",
  components: {CertDownload, CertCreate}
}
</script>

<style scoped>

</style>