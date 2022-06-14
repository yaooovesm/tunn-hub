<template>
  <div v-loading="loading">
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px">
        <div class="title-text">用户信息</div>
      </div>
      <div style="padding: 20px">
        <el-descriptions
            border
            :column="4"
            size="small"
            direction="vertical"
        >
          <el-descriptions-item width="20%" label="用户名">{{ User.account }}</el-descriptions-item>
          <el-descriptions-item width="20%" label="角色">
            <el-tag size="small" type="warning" effect="dark" v-if="User.auth==='admin'">管理员</el-tag>
            <el-tag size="small" type="success" effect="dark" v-else-if="User.auth==='user'">用户</el-tag>
          </el-descriptions-item>
          <el-descriptions-item width="20%" label="是否禁用">
            <el-tag size="small" type="success" effect="dark" v-if="User.disabled===0">正常</el-tag>
            <el-tag size="small" type="danger" effect="dark" v-else-if="User.disabled===1">禁用</el-tag>
          </el-descriptions-item>
          <el-descriptions-item width="20%" label="">
            <el-popover
                placement="bottom-end"
                :width="200"
                trigger="hover"
            >
              <template #default>
                <div>
                  <div class="detail-unit">
                    <span>下行流量 </span>
                    {{ $utils.FormatBytesSizeM(User.flow_count) }}
                  </div>
                  <div class="detail-unit">
                    <span>上行流量 </span>
                    {{ $utils.FormatBytesSizeM(User.tx_count) }}
                  </div>
                  <el-divider style="margin-top: 10px;margin-bottom: 10px"/>
                  <div class="detail-unit">
                    <span>流量总计 </span>
                    {{ $utils.FormatBytesSizeM(User.tx_count + User.flow_count) }}
                  </div>
                </div>
              </template>
              <template #reference>
                <div>
                  <div style="display: inline-block;margin-right: 10px">
                      <span
                          style="display: inline-block;transform: translateY(-2px);margin-right: 3px;font-weight: bolder;color: #007bbb">↓</span>
                    <span>{{ $utils.FormatBytesSizeM(User.flow_count) }}</span>
                  </div>
                  <div style="display: inline-block">
                <span
                    style="display: inline-block;transform: translateY(-2px);margin-right: 3px;font-weight: bolder;color: #007bbb">↑</span>
                    <span>{{ $utils.FormatBytesSizeM(User.tx_count) }}</span>
                  </div>
                </div>
              </template>
            </el-popover>
            <template v-slot:label>
              流量统计
              <!--              <el-tooltip-->
              <!--                  effect="dark"-->
              <!--                  content="此处统计的流量为服务器出方向流量，服务器入方向(客户端出方向)不计算入内。"-->
              <!--                  placement="top-end"-->
              <!--              >-->
              <!--                <i style="font-size: 12px;font-weight: 500; color: rgba(0,123,187,0.8)"-->
              <!--                   class="iconfont icon-question-circle"></i>-->
              <!--              </el-tooltip>-->
            </template>
          </el-descriptions-item>
          <!--          <el-descriptions-item width="20%" label="邮箱">{{-->
          <!--              User.email === "" ? "未设置" : User.email-->
          <!--            }}-->
          <!--          </el-descriptions-item>-->
          <el-descriptions-item width="25%" label="创建时间">{{
              $utils.SecondToDate(User.created, "")
            }}
          </el-descriptions-item>
          <el-descriptions-item width="25%" label="上次登录">{{
              User.last_login === 0 ? "未曾登录" :
                  $utils.UnixMilliToDate(User.last_login, "")
            }}
          </el-descriptions-item>
          <el-descriptions-item width="25%" label="上次离线">{{
              User.last_logout === 0 ? "未曾离线" :
                  $utils.UnixMilliToDate(User.last_logout, "")
            }}
          </el-descriptions-item>
          <el-descriptions-item width="25%" label="用户配置">
            <config-btn :config-id="User.config_id" :account="User.account"/>
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <div style="font-size: 12px;color: #808080;text-align: right;padding: 5px 10px">
        截止至
        {{ $utils.UnixMilliToDate(User.updated, "") }}&nbsp;&nbsp;
        更新于
        {{ $utils.FormatDate("YYYY/mm/dd HH:MM:SS", updateTime) }}&nbsp;
        <el-button type="text" @click="update" style="font-size: 12px;height: 12px;line-height: 13px">刷新</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import axios from "axios";
import ConfigBtn from "@/components/config/ConfigBtn";

export default {
  name: "UserCard",
  components: {
    ConfigBtn
  },

  props: {
    id: {
      type: String
    },
  },
  data() {
    return {
      updateTime: new Date(),
      loading: false,
      User: {
        account: "",
        auth: "",
        config_id: "",
        created: 0,
        email: "",
        flow_count: 0,
        tx_count: 0,
        last_login: 0,
        last_logout: 0,
        updated: 0,
      }
    }
  },
  mounted() {
    this.update()
  },
  methods: {
    update: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/user/info/" + this.id,
        data: {}
      }).then(res => {
        let response = res.data
        this.User = response.data
        this.updateTime = new Date()
        this.loading = false
      }).catch(() => {
        this.loading = false
      })
    }
  }
}
</script>

<style scoped>
.tool-box {
  text-align: right;
  margin-bottom: 5px;
}
</style>