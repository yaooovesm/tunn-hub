<template>
  <div>
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px">
        <div class="title-text">地址池</div>
      </div>
      <div style="padding: 20px">
        <el-table :data="users.slice((currentPage - 1) * pageSize, currentPage * pageSize)"
                  style="width: 100%"
                  stripe
                  scrollbar-always-on
                  v-loading="loading"
        >
          <el-table-column label="分配地址" fixed width="180">
            <template #default="scope">
              {{ scope.row.info.Address }}
            </template>
          </el-table-column>
          <el-table-column label="所在网络" fixed width="180">
            <template #default="scope">
              {{ scope.row.info.Network }}
            </template>
          </el-table-column>
          <el-table-column label="分配类型" width="100" fixed>
            <template #default="scope">
              <span v-if="scope.row.info.IsDynamic" style="color: #0077aa">动态分配</span>
              <span v-else>静态分配</span>
            </template>
          </el-table-column>
          <el-table-column label="连接ID" width="100" fixed>
            <template #default="scope">
              {{ scope.row.info.UUID.substr(0, 8) }}
            </template>
          </el-table-column>
          <el-table-column label="账号" fixed>
            <template #default="scope">
              {{ scope.row.account }}
            </template>
          </el-table-column>
          <el-table-column label="分配时间" width="200" prop="info.Date">
            <template #default="scope">
              {{ $utils.UnixMilliToDate(scope.row.info.Date, "") }}
            </template>
          </el-table-column>
          <el-table-column label="过期时间" width="200">
            <template #default="scope">
              {{ scope.row.info.Expire === 0 ? "不过期" : $utils.UnixMilliToDate(scope.row.info.Expire, "") }}
            </template>
          </el-table-column>
          <el-table-column width="230" fixed="right">
            <template #header>
              <div style="display: inline">
                <el-input v-model="searchText" size="small" placeholder="Type to search">
                  <template #append>
                    <el-button size="small" @click="search">搜索</el-button>
                  </template>
                </el-input>
              </div>
            </template>
          </el-table-column>
        </el-table>
        <div style="margin-top: 30px; margin-bottom: 20px;display: inline-block;width: 100%">
          <div style="font-size: 12px;color: #808080;float: right;padding: 5px 10px;">
            更新于
            {{ $utils.FormatDate("YYYY/mm/dd HH:MM:SS", updateTime) }}&nbsp;
            <el-button type="text" @click="list(false)"
                       style="font-size: 12px;height: 12px;line-height: 13px;transform: translateY(-1px)">刷新
            </el-button>

          </div>
          <el-pagination
              layout="prev,pager,next,jumper"
              background
              :total="users.length"
              :page-size="pageSize"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              :current-page="currentPage"
              small
              style="text-align:center;float: right"
          />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "PoolControlList",
  data() {
    return {
      searchText: "",
      loading: false,
      users: [],
      updateTime: new Date(),
      currentPage: 1,
      pageSize: 8,
    }
  },
  mounted() {
    this.search()
  },
  methods: {
    search: function () {
      if (this.searchText === "") {
        this.list()
        return
      }
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/server/ippool/list",
        data: {}
      }).then(res => {
        let response = res.data
        if (response === undefined) {
          this.$utils.Error("搜索失败", "未找到数据")
          return
        }
        let search = []
        for (let i = 0; i < response.data.length; i++) {
          let user = response.data[i]
          if (user.account === this.searchText ||
              user.info.Address === this.searchText ||
              user.info.UUID === this.searchText) {
            search.push(user)
          }
        }
        this.users = search
        this.updateTime = new Date()
        this.loading = false
      }).catch((err) => {
        this.$utils.HandleError(err)
        this.loading = false
      })
    },
    list: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/server/ippool/list",
        data: {}
      }).then(res => {
        let response = res.data
        this.users = response.data
        this.updateTime = new Date()
        this.loading = false
      }).catch((err) => {
        this.$utils.HandleError(err)
        this.loading = false
      })
    },
    handleSizeChange: function (size) {
      this.pageSize = size;
    },
    handleCurrentChange: function (currentPage) {
      this.currentPage = currentPage;
    }
  }
}
</script>

<style scoped>

</style>