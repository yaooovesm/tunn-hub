<template>
  <div style="padding-top: 50px">
    <el-row>
      <el-col :span="22" :offset="1">
        <user-control-list/>
      </el-col>
      <el-col :span="22" :offset="1" style="margin-top: 30px">
        <pool-control-list/>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import UserControlList from "@/components/control/UserControlList";
import PoolControlList from "@/components/control/PoolControlList";

export default {
  name: "ControlPanel",
  components: {PoolControlList, UserControlList}
}
</script>

<style scoped>

</style>