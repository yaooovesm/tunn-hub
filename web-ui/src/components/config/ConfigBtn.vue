<template>
  <div>
    <div v-if="configId!==''">
      <el-button size="small" type="primary" @click="$refs.config.show(configId,account)">用户配置</el-button>
      <user-config ref="config"/>
    </div>
    <div v-else>
      <el-button size="small" type="primary" disabled>用户配置</el-button>
      <el-tooltip
          effect="dark"
          content='配置丢失请联系管理员'
          placement="right-start"
      >
        <i class="iconfont icon-exclamation-circle"
           style="color: #F56C6C;font-size: 12px;margin-left: 5px;line-height: 24px"></i>
      </el-tooltip>
    </div>
  </div>
</template>

<script>
import UserConfig from "@/components/config/UserConfig";

export default {
  name: "ConfigBtn",
  components: {UserConfig},
  props: {
    configId: {
      type: String,
    },
    account: {
      type: String,
    }
  },
}
</script>

<style scoped>

</style>