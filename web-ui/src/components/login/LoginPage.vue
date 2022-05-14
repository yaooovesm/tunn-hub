<template>
  <div>
    <el-container style="height: 100vh;background-color: #f2f2f2">
      <el-main>
        <el-row>
          <el-col :span="24">
            <login-box style="margin-top: 21vh;margin-left: 50%;transform: translateX(-50%);background-color: #ffffff"
                       :auto="true"/>
          </el-col>
        </el-row>
      </el-main>
    </el-container>
  </div>
</template>

<script>

import LoginBox from "@/components/login/LoginBox";

export default {
  name: "LoginPage",
  components: {
    LoginBox
  },
  mounted() {
  },
  data() {
    return {}
  },
  methods: {}
}
</script>

<style scoped>

</style>