<template>
  <div style="padding-top: 50px">
    <el-row>
      <el-col :span="22" :offset="1">
        <user-list ref="list"/>
      </el-col>
      <el-col :span="22" :offset="1" style="margin-top: 30px">
        <user-create @success="createSuccess"/>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import UserList from "@/components/users/UserList";
import UserCreate from "@/components/users/UserCreate";

export default {
  name: "UserPanel",
  components: {UserCreate, UserList},
  data() {
    return {
      activeName: "list",
    }
  },
  methods: {
    createSuccess: function () {
      this.$refs.list.searchUser()
    }
  }
}
</script>

<style scoped>

</style>