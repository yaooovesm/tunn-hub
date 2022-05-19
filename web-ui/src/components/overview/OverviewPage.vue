<template>
  <div style="padding-top: 50px">
    <el-row v-if="$storage.User.isLogin" :gutter="20">
      <el-col :xs="22" :sm="22" :md="11" :lg="11" :xl="11" :offset="1">
        <flow-overview ref="flow" :passive="true"/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="11" :lg="11" :xl="11">
        <monitor-overview ref="monitor" :passive="true"/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="7" :lg="7" :xl="7" :offset="1" style="margin-top: 30px">
        <i-p-pool-overview ref="ippool" style="cursor: pointer" @click="$router.push({path:'/dashboard/control'})"
                           :passive="true"/>
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
import ReporterClient from "@/reporter.client";

export default {
  name: "OverviewPage",
  components: {
    MonitorOverview,
    IPPoolOverview,
    ServerConfigOverview,
    UsersOverview,
    FlowOverview
  },
  data() {
    return {
      loading: false
    }
  },
  mounted() {
    //this.connectToReporter()
  },
  unmounted() {
    //this.reporterClient.Close("ovw component")
  },
  methods: {
    connectToReporter: function () {
      this.loading = true
      let that = this
      this.reporterClient = new ReporterClient(
          {
            "flow": {
              name: "/api/v1/server/flow",
            },
            "ippool": {
              name: "/api/v1/server/ippool",
            },
            "monitor": {
              name: "/api/v1/server/monitor",
            }
          }, function (data) {
            let resp = JSON.parse(data)
            that.$refs.flow.set(resp.flow.Data)
            that.$refs.ippool.set(resp.ippool.Data)
            that.$refs.monitor.set(resp.monitor.Data)
          }, function () {
          }, function (err) {
            console.log(err)
          }, 2000
      )
      this.$nextTick(() => {
        this.reporterClient.Start("ovw component")
        this.loading = false
      });
    }
  }

}
</script>

<style scoped>

</style>