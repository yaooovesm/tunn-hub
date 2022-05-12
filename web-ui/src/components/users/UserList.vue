<template>
  <div>
    <el-card shadow="always" body-style="padding:0" v-loading="loading">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px">
        <div class="title-text">用户列表</div>
      </div>
      <div style="padding: 20px">
        <el-table :data="users.slice((currentPage - 1) * pageSize, currentPage * pageSize)"
                  style="width: 100%"
                  stripe
                  :default-sort="{ prop: 'flow_count', order: 'descending' }"
                  scrollbar-always-on
        >
          <el-table-column prop="id" fixed label="ID" width="280"/>
          <el-table-column prop="account" fixed label="账号" width="120"/>
          <el-table-column prop="email" min-width="200" label="邮箱">
            <template #default="scope">
              {{ scope.row.email === "" ? "未设置" : scope.row.email }}
            </template>
          </el-table-column>
          <el-table-column label="角色" width="80" prop="auth">
            <template #default="scope">
              <el-tag size="small" type="warning" effect="dark" v-if="scope.row.auth==='admin'">管理员</el-tag>
              <el-tag size="small" type="success" effect="dark" v-else-if="scope.row.auth==='user'">用户</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="是否禁用" width="120" :align="'center'" prop="disabled" sortable>
            <template #default="scope">
              <el-tag size="small" type="success" effect="dark" v-if="scope.row.disabled===0">正常</el-tag>
              <el-tag size="small" type="danger" effect="dark" v-else-if="scope.row.disabled===1">禁用</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="流量统计" width="150" prop="flow_count" sortable>
            <template #default="scope">
              {{ $utils.FormatBytesSize(scope.row.flow_count) }}
            </template>
          </el-table-column>
          <el-table-column label="上次登录" width="160" prop="last_login" sortable>
            <template #default="scope">
              {{ scope.row.last_login === 0 ? "未曾登录" : $utils.UnixMilliToDate(scope.row.last_login, "") }}
            </template>
          </el-table-column>
          <el-table-column label="上次离线" width="160" prop="last_logout" sortable>
            <template #default="scope">
              {{ scope.row.last_logout === 0 ? "未曾离线" : $utils.UnixMilliToDate(scope.row.last_logout, "") }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="160" prop="last_logout" sortable>
            <template #default="scope">
              {{ $utils.SecondToDate(scope.row.created, "") }}
            </template>
          </el-table-column>
          <el-table-column width="200" fixed="right">
            <template #header>
              <div style="display: inline">
                <el-input v-model="search" size="small" placeholder="Type to search">
                  <template #append>
                    <el-button size="small" @click="searchUser">搜索</el-button>
                  </template>
                </el-input>
              </div>
            </template>
            <template #default="scope">
              <el-button size="small" @click="detailUser(scope.row)">详情</el-button>
              <el-dropdown type="primary" trigger="click" style="margin-left: 10px">
                <el-button size="small" type="primary">操作</el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="updateUser(scope.row)">修改</el-dropdown-item>
                    <el-dropdown-item divided :style="scope.row.disabled?'color:#67C23A':'color: #F56C6C'"
                                      @click="disableUser(scope.row)">{{ scope.row.disabled ? "解禁" : "禁用" }}
                    </el-dropdown-item>
                    <el-dropdown-item style="color: #F56C6C" @click="deleteUser(scope.row)">删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>
        <div style="margin-top: 30px; margin-bottom: 20px;display: inline-block;width: 100%">
          <div style="font-size: 12px;color: #808080;float: right;padding: 5px 10px;">
            更新于
            {{ $utils.FormatDate("YYYY/mm/dd HH:MM:SS", updateTime) }}&nbsp;
            <el-button type="text" @click="list" style="font-size: 12px;height: 12px;line-height: 13px">刷新
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
    <user-detail ref="detail"/>
    <user-update ref="update" @success="searchUser"/>
  </div>
</template>

<script>
import axios from "axios";
import UserDetail from "@/components/users/UserDetail";
import {ElMessageBox} from "element-plus";
import UserUpdate from "@/components/users/UserUpdate";

export default {
  name: "UserList",
  components: {UserUpdate, UserDetail},
  data() {
    return {
      search: "",
      loading: false,
      users: [],
      updateTime: new Date(),
      currentPage: 1,
      pageSize: 8,
    }
  },
  mounted() {
    this.list()
  },
  methods: {
    searchUser: function () {
      if (this.search === "") {
        this.list()
        return
      }
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/user/list/info",
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
          if (user.id !== "" && user.id.indexOf(this.search) !== -1 ||
              user.account !== "" && user.account.indexOf(this.search) !== -1 ||
              user.email !== "" && user.email.indexOf(this.search) !== -1) {
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
    disableUser: function (user) {
      let option = user.disabled === 1 ? "解禁" : "禁用"
      ElMessageBox.confirm(
          '是否' + option + '该用户',
          '警告',
          {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
          }
      ).then(() => {
        axios({
          method: "post",
          url: "/api/v1/user/disable" + user.id,
          data: {
            id: user.id,
            disabled: user.disabled === 1 ? 0 : 1
          }
        }).then(res => {
          let response = res.data
          this.$utils.Success("操作成功", response.msg)
        }).catch((err) => {
          this.$utils.HandleError(err)
        }).finally(() => {
          this.searchUser()
        })
      }).catch(() => {
      })
    },
    detailUser: function (user) {
      this.$refs.detail.show(user.id)
    },
    deleteUser: function (user) {
      ElMessageBox.confirm(
          '是否删除该用户',
          '警告',
          {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
          }
      ).then(() => {
        axios({
          method: "delete",
          url: "/api/v1/user/delete",
          data: {
            id: user.id,
          }
        }).then(res => {
          let response = res.data
          this.$utils.Success("操作成功", response.msg)
        }).catch((err) => {
          this.$utils.HandleError(err)
        }).finally(() => {
          this.searchUser()
        })
      }).catch(() => {
      })
    },
    updateUser: function (user) {
      this.$refs.update.show(user.id)
    },
    list: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/user/list/info",
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