<template>
  <div style="padding-top: 50px">
    <el-row v-if="$storage.User.isLogin" :gutter="20">
      <el-col :xs="22" :sm="22" :md="11" :lg="11" :xl="11" :offset="1">
        <flow-overview/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="11" :lg="11" :xl="11">
        <monitor-overview/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="7" :lg="7" :xl="7" :offset="1" style="margin-top: 30px">
        <i-p-pool-overview style="cursor: pointer" @click="$router.push({path:'/dashboard/control'})"/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="7" :lg="7" :xl="7" style="margin-top: 30px">
        <users-overview style="cursor: pointer" @click="$router.push({path:'/dashboard/users'})"/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="8" :lg="8" :xl="8" style="margin-top: 30px">
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