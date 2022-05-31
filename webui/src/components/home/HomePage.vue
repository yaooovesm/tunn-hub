<template>
  <el-row style="padding-top: 50px" v-if="$storage.User.isLogin">
    <el-col :span="22" :offset="1">
      <user-card :id="id"/>
    </el-col>
    <el-col :span="22" :offset="1" style="margin-top: 30px">
      <client-card :id="id"/>
    </el-col>
  </el-row>
</template>

<script>
import UserCard from "./UserCard";
import ClientCard from "./ClientCard";

export default {
  name: "HomePage",
  components: {ClientCard, UserCard},
  created() {
    this.$storage.Load()
    this.id = this.$storage.User.info.id
  },
  data() {
    return {
      id: ""
    }
  },
}
</script>

<style scoped>

</style>