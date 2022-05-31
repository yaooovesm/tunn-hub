<template>
  <div>
    <el-container style="height: 100vh">
      <el-aside width="auto" style="height: 100%;">
        <div style="height: 100%;text-align: left;">
          <el-menu :default-active="$route.path"
                   ref="menu"
                   :collapse="isCollapse"
                   style="box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);height: 100vh"
                   :router="true"
          >
            <el-menu-item class="main-title" :index="$route.path"
                          style="height: 60px;text-align: center;font-size: 16px">
              <div v-if="isCollapse" style="width: 20px;text-align: center">T</div>
              {{ isCollapse ? "" : "TunnHub WebUI" }}
              <template v-slot:title>
                <div style="text-align: center">{{ isCollapse ? "TunnHub WebUI" : "" }}</div>
              </template>
            </el-menu-item>
            <!--控制台-->
            <el-menu-item index="/dashboard/overview" style="height: 50px" v-if="$storage.IsAdmin()">
              <i class="iconfont icon-yibiaopan" style="font-size: 20px;margin-right: 5px"></i>
              <template v-slot:title>
                <span>概况</span>
              </template>
            </el-menu-item>
            <!--主页-->
            <el-menu-item index="/dashboard/home" style="height: 50px">
              <i class="iconfont icon-guanli" style="font-size: 20px;margin-right: 5px"></i>
              <template v-slot:title>
                <span>主页</span>
              </template>
            </el-menu-item>
            <!--控制台-->
            <el-menu-item index="/dashboard/control" style="height: 50px" v-if="$storage.IsAdmin()">
              <i class="iconfont icon-jiancexitong" style="font-size: 20px;margin-right: 5px"></i>
              <template v-slot:title>
                <span>控制台</span>
              </template>
            </el-menu-item>
            <!--用户管理-->
            <el-menu-item index="/dashboard/users" style="height: 50px" v-if="$storage.IsAdmin()">
              <i class="iconfont icon-renyuan" style="font-size: 20px;margin-right: 5px"></i>
              <template v-slot:title>
                <span>用户管理</span>
              </template>
            </el-menu-item>
            <!--认证管理-->
            <el-menu-item index="/dashboard/cert" style="height: 50px;" v-if="$storage.IsAdmin()">
              <i class="iconfont icon-renyuanmiyue" style="font-size: 20px;margin-right: 5px"></i>
              <template v-slot:title>
                <span>认证管理</span>
              </template>
            </el-menu-item>
            <el-menu-item :index="$route.path"
                          :style="$storage.IsAdmin()?'height: 50px;margin-top: calc(100vh - 363px)':'height: 50px;margin-top: calc(100vh - 163px)'"
                          @click="collapse">
              <i :class="'iconfont icon-angle-'+icon" style="font-size: 11px;margin-right: 10px;margin-left: 4px"></i>
              <template v-slot:title>
                <span>{{ isCollapse ? "展开" : "折叠" }}</span>
              </template>
            </el-menu-item>
          </el-menu>
        </div>
      </el-aside>
      <el-container>
        <el-header style="width: auto;padding: 0">
          <page-header/>
        </el-header>
        <el-main>
          <el-scrollbar height="100%" style="width: 100%">
            <router-view style="padding-bottom: 30px;margin-left: 1px;width: calc(100% - 11px)"/>
          </el-scrollbar>
        </el-main>
        <el-footer v-if="$route.path ==='/dashboard/home' || $route.path ==='/dashboard/overview'">
          <link-overview/>
        </el-footer>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import PageHeader from "@/components/dashboard/PageHeader";
import LinkOverview from "@/components/overview/LinkOverview";

export default {
  name: "DashboardPage",
  components: {LinkOverview, PageHeader},
  data() {
    return {
      icon: "right",
      isCollapse: true
    }
  },
  methods: {
    collapse: function () {
      if (this.isCollapse) {
        this.icon = "left"
        this.isCollapse = false
      } else {
        this.icon = "right"
        this.isCollapse = true
      }
    }
  }
}
</script>

<style scoped>
.el-menu-item.is-active.main-title {
  color: #222222 !important;
  cursor: auto;
}

.el-menu-item.is-active.main-title:hover {
  color: #222222 !important;
  background-color: transparent;
}

</style>