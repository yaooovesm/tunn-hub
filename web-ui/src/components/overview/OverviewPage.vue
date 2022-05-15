<template>
  <div style="padding-top: 40px">
    <el-row v-if="$storage.User.isLogin" :gutter="10">
      <el-col :span="22" :offset="1">
        <flow-overview/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="10" :lg="10" :xl="10" :offset="1" style="margin-top: 30px">
        <monitor-overview/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="6" :lg="6" :xl="6" style="margin-top: 30px">
        <i-p-pool-overview/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="6" :lg="6" :xl="6" style="margin-top: 30px">
        <users-overview/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="8" :lg="8" :xl="8" :offset="1" style="margin-top: 30px">
        <server-config-overview/>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import FlowOverview from "@/components/overview/FlowOverview";
import UsersOverview from "@/components/overview/UsersOverview";
import ServerConfigOverview from "@/components/overview/ServerConfigOverview";
import IPPoolOverview from "@/components/overview/IPPoolOverview";
import MonitorOverview from "@/components/overview/MonitorOverview";

export default {
  name: "OverviewPage",
  components: {
    MonitorOverview,
    IPPoolOverview,
    ServerConfigOverview,
    UsersOverview,
    FlowOverview
  }

}
</script>

<style scoped>

</style>