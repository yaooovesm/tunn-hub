<template>
  <div>
    <reporter-client :resources="res" :interval="5000" @recv="onRecv" ref="reporter_client" v-if="connect"/>
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px">
        <div class="title-text">在线列表</div>
      </div>
      <div style="padding: 20px">
        <el-table :data="users.slice((currentPage - 1) * pageSize, currentPage * pageSize)"
                  style="width: 100%"
                  stripe
                  :default-sort="showOffline?{ prop: 'flow_count', order: 'descending' }:{ prop: 'status.online', order: 'descending' }"
                  scrollbar-always-on
                  v-loading="loading"
        >
          <el-table-column label="平台" fixed width="105">
            <template #default="scope">
              <div>
                <svg class="icon" aria-hidden="true" v-if="scope.row.status.config.runtime.os==='windows'"
                     style="width: 1.1em;height: 1.0em;">
                  <use xlink:href="#icon-Windows"></use>
                </svg>
                <svg class="icon" aria-hidden="true" v-else-if="scope.row.status.config.runtime.os==='linux'"
                     style="width: 1.2em;height: 1.1em;">
                  <use xlink:href="#icon-linux"></use>
                </svg>
                <svg class="icon" aria-hidden="true" v-else-if="scope.row.status.config.runtime.os==='darwin'">
                  <use xlink:href="#icon-IOS"></use>
                </svg>
                {{
                  scope.row.status.config.runtime.os === 'darwin' ? "OSX" : scope.row.status.config.runtime.os
                }}
              </div>
            </template>
          </el-table-column>
          <el-table-column fixed label="客户端" width="90" prop="status.online" sortable>
            <template #default="scope">
              <div style="transform: translateY(2px);line-height: 23px">
                <div :class="scope.row.status.online?'circle online-circle':'circle offline-circle'"></div>
                <div :class="scope.row.status.online?'online-text':'offline-text'">
                  {{ scope.row.status.online ? "在线" : "离线" }}
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column fixed label="版本" width="80">
            <template #default="scope">
              <div v-if="scope.row.status.config.runtime.app !== ''">
                <span>
                {{
                    scope.row.status.config.runtime.app.substr(0, scope.row.status.config.runtime.app.lastIndexOf('.'))
                  }}
              </span>
                <el-tooltip
                    effect="dark"
                    placement="top"
                >
                  <template #content>{{ scope.row.status.config.runtime.app }}</template>
                  <i class="iconfont icon-info-circle"
                     style="font-size: 10px;margin-left: 5px;color: #007bbb;transform: translateY(-1px);display: inline-block"></i>
                </el-tooltip>
              </div>
              <div v-else style="font-size: 12px;color: #909399">
                unknown
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="account" fixed label="账号" show-overflow-tooltip>
            <template #default="scope">
              {{ scope.row.account }}
              <!--              <span style="color: #F56C6C;font-size: 12px">{{ scope.row.disabled === 1 ? "已禁用" : "" }}</span>-->
              <!--              <el-tag size="small" type="warning" effect="dark" v-if="scope.row.auth==='admin'">管理员</el-tag>-->
              <!--              <el-tag size="small" type="success" effect="dark" v-else-if="scope.row.auth==='user'">用户</el-tag>-->
              <!--              <el-tag size="small" type="danger" effect="dark" v-if="scope.row.disabled===1" style="margin-left: 5px">-->
              <!--                禁用-->
              <!--              </el-tag>-->
            </template>
          </el-table-column>
          <el-table-column label="带宽限制" width="120" prop="auth">
            <template #default="scope">
              <div>{{
                  scope.row.status.config.limit.bandwidth === 0 ? "无限制" : scope.row.status.config.limit.bandwidth + " Mbps"
                }}
              </div>
              <!--              <div>{{ $utils.FormatPacketSpeed(scope.row.status.rx.PacketSpeed) }}</div>-->
              <!--              <div>出方向：&nbsp;{{ $utils.FormatBytesSpeed(scope.row.status.tx.FlowSpeed) }}-->
              <!--              </div>-->
            </template>
          </el-table-column>
          <el-table-column label="传输速率" width="200">
            <template #default="scope">
              <el-popover
                  placement="bottom-end"
                  :title="'速率统计@'+scope.row.account"
                  :width="200"
                  trigger="hover"
              >
                <template #default>
                  <div style="margin-top: 20px" v-if="scope.row.status.online">
                    <div class="detail-unit-title">
                      <span>下行速率</span>
                    </div>
                    <div class="detail-unit">
                      <span>流量 </span>
                      {{ $utils.FormatBytesSpeed(scope.row.status.tx.FlowSpeed) }}
                    </div>
                    <div class="detail-unit">
                      <span>数据包 </span>
                      {{ $utils.FormatPacketSpeed(scope.row.status.tx.PacketSpeed) }}
                    </div>
                    <el-divider style="margin-top: 6px;margin-bottom: 6px"/>
                    <div class="detail-unit-title">
                      <span>上行速率</span>
                    </div>
                    <div class="detail-unit">
                      <span>流量 </span>
                      {{ $utils.FormatBytesSpeed(scope.row.status.rx.FlowSpeed) }}
                    </div>
                    <div class="detail-unit">
                      <span>数据包 </span>
                      {{ $utils.FormatPacketSpeed(scope.row.status.rx.PacketSpeed) }}
                    </div>
                    <el-divider style="margin-top: 6px;margin-bottom: 6px"/>
                    <div class="detail-unit-title">
                      <span>总速率</span>
                    </div>
                    <div class="detail-unit">
                      <span>流量 </span>
                      {{ $utils.FormatBytesSpeed(scope.row.status.rx.FlowSpeed + scope.row.status.tx.FlowSpeed) }}
                    </div>
                    <div class="detail-unit">
                      <span>数据包 </span>
                      {{ $utils.FormatPacketSpeed(scope.row.status.tx.PacketSpeed + scope.row.status.rx.PacketSpeed) }}
                    </div>
                  </div>
                  <div v-else>
                    <span style="color: #909399;font-size: 12px;text-align: center;margin-top: 20px">客户端已离线</span>
                  </div>
                </template>
                <template #reference>
                  <div>
                    <div style="display: inline-block;margin-right: 10px">
                      <span
                          style="display: inline-block;transform: translateY(-2px);margin-right: 3px;font-weight: bolder;color: #007bbb">↓</span>
                      <span>{{ $utils.FormatBytesSpeed(scope.row.status.tx.FlowSpeed) }}</span>
                    </div>
                    <div style="display: inline-block;">
                      <span
                          style="display: inline-block;transform: translateY(-2px);margin-right: 3px;font-weight: bolder;color: #007bbb">↑</span>
                      <span>{{ $utils.FormatBytesSpeed(scope.row.status.rx.FlowSpeed) }}</span>
                    </div>
                  </div>
                </template>
              </el-popover>
              <!--服务端RX对应客户端TX(上行)-->
              <!--              <div>{{ $utils.FormatBytesSpeed(scope.row.status.rx.FlowSpeed) }}-->
              <!--              </div>-->
              <!--              <div>{{ $utils.FormatPacketSpeed(scope.row.status.rx.PacketSpeed) }}</div>-->
              <!--              <div>出方向：&nbsp;{{ $utils.FormatBytesSpeed(scope.row.status.tx.FlowSpeed) }}-->
              <!--              </div>-->
            </template>
          </el-table-column>
          <el-table-column label="流量统计" width="200" prop="flow_count" sortable>
            <template #default="scope">
              <!--此处flow_count对应客户端RX-->
              <!--此处tx_count为客户端TX-->
              <el-popover
                  placement="bottom-end"
                  :title="'流量统计@'+scope.row.account"
                  :width="200"
                  trigger="hover"
              >
                <template #default>
                  <div>
                    <div class="detail-unit">
                      <span>下行流量 </span>
                      {{ $utils.FormatBytesSizeM(scope.row.flow_count) }}
                    </div>
                    <div class="detail-unit">
                      <span>上行流量 </span>
                      {{ $utils.FormatBytesSizeM(scope.row.tx_count) }}
                    </div>
                    <el-divider style="margin-top: 10px;margin-bottom: 10px"/>
                    <div class="detail-unit">
                      <span>流量总计 </span>
                      {{ $utils.FormatBytesSizeM(scope.row.tx_count + scope.row.flow_count) }}
                    </div>
                  </div>
                </template>
                <template #reference>
                  <div>
                    <div style="display: inline-block;margin-right: 10px">
                      <span
                          style="display: inline-block;transform: translateY(-2px);margin-right: 3px;font-weight: bolder;color: #007bbb">↓</span>
                      <span>{{ $utils.FormatBytesSizeM(scope.row.flow_count) }}</span>
                    </div>
                    <div style="display: inline-block">
                <span
                    style="display: inline-block;transform: translateY(-2px);margin-right: 3px;font-weight: bolder;color: #007bbb">↑</span>
                      <span>{{ $utils.FormatBytesSizeM(scope.row.tx_count) }}</span>
                    </div>
                  </div>
                </template>
              </el-popover>

            </template>
          </el-table-column>

          <!--          <el-table-column label="下行速率" width="120" align="center">-->
          <!--            <template #default="scope">-->
          <!--              &lt;!&ndash;服务端TX对应客户端RX(下行)&ndash;&gt;-->
          <!--              <div>{{ $utils.FormatBytesSpeed(scope.row.status.tx.FlowSpeed) }}-->
          <!--              </div>-->
          <!--              &lt;!&ndash;              <div>{{ $utils.FormatPacketSpeed(scope.row.status.tx.PacketSpeed) }}</div>&ndash;&gt;-->
          <!--              &lt;!&ndash;              <div>出方向：&nbsp;{{ $utils.FormatBytesSpeed(scope.row.status.tx.FlowSpeed) }}&ndash;&gt;-->
          <!--              &lt;!&ndash;              </div>&ndash;&gt;-->
          <!--            </template>-->
          <!--          </el-table-column>-->
          <!--          <el-table-column label="流量监控" width="170" prop="auth">-->
          <!--            <template #default="scope">-->
          <!--              <div>入方向：&nbsp;{{ $utils.FormatBytesSpeed(scope.row.status.rx.FlowSpeed) }}-->
          <!--              </div>-->
          <!--              <div>出方向：&nbsp;{{ $utils.FormatBytesSpeed(scope.row.status.tx.FlowSpeed) }}-->
          <!--              </div>-->
          <!--            </template>-->
          <!--          </el-table-column>-->
          <!--          <el-table-column label="数据包监控" width="170" prop="auth">-->
          <!--            <template #default="scope">-->
          <!--              <div>入方向：&nbsp;{{ $utils.FormatPacketSpeed(scope.row.status.rx.PacketSpeed) }}-->
          <!--              </div>-->
          <!--              <div>出方向：&nbsp;{{ $utils.FormatPacketSpeed(scope.row.status.tx.PacketSpeed) }}-->
          <!--              </div>-->
          <!--            </template>-->
          <!--          </el-table-column>-->
          <el-table-column prop="status.config.device.cidr" label="内网地址" width="150"/>
          <el-table-column prop="status.address" label="客户端地址" width="180"/>
          <el-table-column label="状态" width="100">
            <template #default="scope">
              <el-tag size="small" type="danger" effect="dark" v-if="scope.row.disabled===1">禁用</el-tag>
              <el-tag size="small" type="success" effect="dark" v-else-if="scope.row.disabled===0">正常</el-tag>
              <!--              <span style="margin-left: 5px">-->
              <!--                <el-tag size="small" type="success" effect="dark" v-if="scope.row.status.online">在线</el-tag>-->
              <!--                <el-tag size="small" type="info" effect="dark" v-else>离线</el-tag>-->
              <!--              </span>-->
            </template>
          </el-table-column>
          <!--          <el-table-column label="流量发送" width="150" prop="tx_count" show-overflow-tooltip>-->
          <!--            <template #default="scope">-->
          <!--              -->
          <!--            </template>-->
          <!--          </el-table-column>-->
          <!--          <el-table-column label="上次登录" width="160" prop="last_login" sortable>-->
          <!--            <template #default="scope">-->
          <!--              {{ scope.row.last_login === 0 ? "未曾登录" : $utils.UnixMilliToDate(scope.row.last_login, "") }}-->
          <!--            </template>-->
          <!--          </el-table-column>-->
          <!--          <el-table-column label="上次离线" width="160" prop="last_logout" sortable>-->
          <!--            <template #default="scope">-->
          <!--              {{ scope.row.last_logout === 0 ? "未曾离线" : $utils.UnixMilliToDate(scope.row.last_logout, "") }}-->
          <!--            </template>-->
          <!--          </el-table-column>-->
          <el-table-column width="230" fixed="right">
            <template #header>
              <div style="display: inline">
                <el-input v-model="search" size="small" placeholder="Type to search">
                  <template #append>
                    <el-button size="small" @click="searchUser(false)">搜索</el-button>
                  </template>
                </el-input>
              </div>
            </template>
            <template #default="scope">
              <el-button size="small" @click="detailUser(scope.row)">
                详情
              </el-button>
              <el-button size="small" @click="showConfig(scope.row.status.config,scope.row.account)">
                属性
              </el-button>
              <el-dropdown size="small" trigger="click" style="margin-left: 10px">
                <el-button size="small">操作</el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="$refs.config_set.show(scope.row.config_id,scope.row.account)">设置
                    </el-dropdown-item>
                    <el-dropdown-item @click="disconnect(scope.row.id,scope.row.account)">断开</el-dropdown-item>
                    <el-dropdown-item @click="reconnect(scope.row.id,scope.row.account)">重置</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <!--              <el-button size="small"-->
              <!--                         @click="$refs.config_set.show(scope.row.config_id,scope.row.account)">-->
              <!--                设置-->
              <!--              </el-button>-->
              <!--              <el-button size="small" @click="disconnect(scope.row.id,scope.row.account)"-->
              <!--                         v-if="scope.row.status.online">-->
              <!--                断开-->
              <!--              </el-button>-->
            </template>
          </el-table-column>
        </el-table>
        <div style="margin-top: 30px; margin-bottom: 20px;display: inline-block;width: 100%">
          <div style="font-size: 12px;color: #808080;float: right;padding: 5px 10px;">
            更新于
            {{ $utils.FormatDate("YYYY/mm/dd HH:MM:SS", updateTime) }}&nbsp;
            <el-switch
                v-model="showOffline"
                inline-prompt
                active-text="是"
                inactive-text="否"
                style="margin-left: 10px"
                size="small"
                @change="searchUser(false)"
            />
            显示离线
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
    <user-detail ref="detail"/>
    <user-update ref="update" @success="searchUser(false)"/>
    <show-config ref="config"/>
    <user-config ref="config_set" style="text-align: left"/>
  </div>
</template>

<script>
import axios from "axios";
import UserDetail from "@/components/users/UserDetail";
import UserUpdate from "@/components/users/UserUpdate";
import ShowConfig from "@/components/control/ShowConfig";
import {ElMessageBox} from "element-plus";
import UserConfig from "@/components/config/UserConfig";
import ReporterClient from "@/components/ReporterClient";

export default {
  name: "UserControlList",
  components: {ShowConfig, UserUpdate, UserDetail, UserConfig, ReporterClient},
  data() {
    return {
      showOffline: false,
      search: "",
      loading: false,
      users: [],
      updateTime: new Date(),
      currentPage: 1,
      pageSize: 8,
      connect: false,
      res: {
        "list": {
          name: "/api/v1/user/list",
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
      let response = JSON.parse(data).list
      let search = []
      for (let i = 0; i < response.Data.length; i++) {
        if (response.Data[i].status.online === false) {
          //过滤离线
          if (!this.showOffline) {
            continue
          }
        }
        let user = response.Data[i].info
        user.status = response.Data[i].status
        if (user.id !== "" && user.id.indexOf(this.search) !== -1 ||
            user.account !== "" && user.account.indexOf(this.search) !== -1 ||
            user.email !== "" && user.email.indexOf(this.search) !== -1) {
          search.push(user)
        }
        this.users = search
        this.updateTime = new Date()
      }
    },
    showConfig: function (config, account) {
      this.$refs.config.show(config, account)
    },
    searchUser: function (silence) {
      if (this.search === "") {
        this.list(silence)
        return
      }
      if (!silence) {
        this.loading = true
      }
      axios({
        method: "get",
        url: "/api/v1/user/list",
        data: {}
      }).then(res => {
        let response = res.data
        if (response === undefined) {
          this.$utils.Error("搜索失败", "未找到数据")
          return
        }
        let search = []
        for (let i = 0; i < response.data.length; i++) {
          if (response.data[i].status.online === false) {
            //过滤离线
            if (!this.showOffline) {
              continue
            }
          }
          let user = response.data[i].info
          user.status = response.data[i].status
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
    detailUser: function (user) {
      this.$refs.detail.show(user.id)
    },
    list: function (silence) {
      if (!silence) {
        this.loading = true
      }
      axios({
        method: "get",
        url: "/api/v1/user/list",
        data: {}
      }).then(res => {
        let response = res.data
        let users = []
        for (let i = 0; i < response.data.length; i++) {
          if (response.data[i].status.online === false) {
            //过滤离线
            if (!this.showOffline) {
              continue
            }
          }
          let user = response.data[i].info
          user.status = response.data[i].status
          users.push(user)
        }
        this.users = users
        this.updateTime = new Date()
        this.loading = false
      }).catch((err) => {
        this.$utils.HandleError(err)
        this.loading = false
      })
    },
    disconnect: function (id, account) {
      if (id === "") {
        this.$utils.Warning("操作失败", "无法获取用户ID")
        return
      }
      ElMessageBox.confirm(
          '是否断开与用户' + account + "的连接",
          '警告',
          {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
          }
      ).then(() => {
        axios({
          method: "get",
          url: "/api/v1/server/disconnect/id/" + id,
          data: {}
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
    reconnect: function (id, account) {
      if (id === "") {
        this.$utils.Warning("操作失败", "无法获取用户ID")
        return
      }
      ElMessageBox.confirm(
          '是否重置与用户' + account + "的连接",
          '警告',
          {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
          }
      ).then(() => {
        axios({
          method: "get",
          url: "/api/v1/server/reconnect/id/" + id,
          data: {}
        }).then(res => {
          let response = res.data
          this.$utils.Success("操作成功,请等待生效", response.msg)
        }).catch((err) => {
          this.$utils.HandleError(err)
        }).finally(() => {
          this.searchUser()
        })
      }).catch(() => {
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
.detail-unit-title {
  font-size: 12px;
  font-weight: 600;
  color: #404040;
  margin-bottom: 7px;
}

.detail-unit {
  text-align: right;
  font-size: 12px;
  color: #007bbb;
}

.detail-unit span {
  color: #404040;
  float: left;
  display: inline-block;
}
</style>