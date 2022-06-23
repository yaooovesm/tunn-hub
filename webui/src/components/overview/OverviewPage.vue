<template>
  <div style="padding-top: 30px">
    <el-row v-if="$storage.User.isLogin" :gutter="20">
      <el-col :span="22" :offset="1" style="margin-bottom: 10px" v-if="connect">
        <reporter-client :resources="res" :interval="2000" @recv="onRecv" ref="reporter_client" v-if="connect"/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="22" :lg="22" :xl="22" :offset="1">
        <traffic-recorder-overview/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="11" :lg="11" :xl="11" :offset="1">
        <flow-overview ref="flow" :passive="true"/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="11" :lg="11" :xl="11">
        <monitor-overview ref="monitor" :passive="true"/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="7" :lg="7" :xl="7" :offset="1" style="margin-top: 30px">
        <i-p-pool-overview ref="ippool" :passive="true"/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="7" :lg="7" :xl="7" style="margin-top: 30px">
        <users-overview/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="8" :lg="8" :xl="8" style="margin-top: 30px">
        <server-config-overview/>
      </el-col>
      <el-col :xs="22" :sm="22" :md="22" :lg="22" :xl="22" :offset="1" style="margin-top: 80px">
        <link-overview style="left: 50%"/>
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
import ReporterClient from "@/components/ReporterClient";
import TrafficRecorderOverview from "@/components/overview/TrafficRecorderOverview";
import LinkOverview from "@/components/overview/LinkOverview";

export default {
  name: "OverviewPage",
  components: {
    LinkOverview,
    TrafficRecorderOverview,
    ReporterClient,
    MonitorOverview,
    IPPoolOverview,
    ServerConfigOverview,
    UsersOverview,
    FlowOverview
  },
  data() {
    return {
      connect: false,
      loading: false,
      res: {
        "flow": {
          name: "/api/v1/server/flow",
        },
        "ippool": {
          name: "/api/v1/server/ippool",
        },
        "monitor": {
          name: "/api/v1/server/monitor",
        }
      }
    }
  },
  mounted() {
    this.loading = true
    this.connect = true
    this.$nextTick(() => {
      this.$refs.reporter_client.Start()
      this.loading = false
    })
  },
  unmounted() {
    this.connect = false
  },
  methods: {
    onRecv: function (data) {
      let resp = JSON.parse(data)
      this.$refs.flow.set(resp.flow.Data)
      this.$refs.ippool.set(resp.ippool.Data)
      this.$refs.monitor.set(resp.monitor.Data)
    },
  }

}
</script>

<style scoped>

</style>