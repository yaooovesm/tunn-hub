<template>
  <div>
    <el-select
        @click="load"
        v-model="selected"
        size="small"
        multiple
        loading-text="加载中"
        collapse-tags
        collapse-tags-tooltip
        placeholder="请选择网络"
        style="width: 300px"
        :loading="loading"
        :value="selected"
        no-data-text="没有可供导入的网络"
        no-match-text="没有可供导入的网络"
    >
      <el-option
          v-for="r in routes"
          :key="r.network"
          :label="r.name"
          :value="r.network"
          style="height: 50px;margin-bottom: 2px;padding: 0;"
          :disabled="r.provider === ''"
      >
        <el-popover
            placement="right-start"
            :width="200"
            :hide-after="0"
            trigger="hover"
        >
          <template #default>
            <div>
              <div class="detail-unit">
                <span>名称 </span>
                {{ r.name }}
              </div>
              <div class="detail-unit">
                <span>网络 </span>
                <div style="color: #2c3e50">{{ r.network }}</div>
              </div>
              <el-divider style="margin-top: 10px;margin-bottom: 10px"/>
              <div class="detail-unit">
                <span>提供者 </span>
                <div v-if="r.provider === ''" style="color: #990055;display: inline-block">未知用户</div>
                <div v-else style="color: #2c3e50;display: inline-block">{{ r.provider }}</div>
                <i style="color: #67c23a;font-size: 12px;display: inline-block;margin-left: 5px"
                   class="iconfont icon-safety-certificate"
                   v-if="r.certificated"></i>
              </div>
            </div>
          </template>
          <template #reference>
            <div style="line-height: 15px;padding: 5px 10px">
              <div style="font-size: 12px;margin-bottom: 8px;margin-top: 3px">
                <div style="color: #2c3e50;display: inline-block"
                     v-if="r.name.length>0">
                  {{ r.name }}
                </div>
                <div style="color: #909399;display: inline-block" v-else>未命名</div>
                <i style="color: #67c23a;font-size: 12px;display: inline-block;margin-left: 3px"
                   class="iconfont icon-safety-certificate"
                   v-if="r.certificated"></i>
              </div>
              <div>
                <div style="font-size: 12px;color: #007bbb">{{ r.network }}
                  <div style="display: inline-block;color: #909399;float: right;text-align: right;opacity: 0.7">
                    由
                    <span v-if="r.provider === ''" style="color: #990055">未知用户</span>
                    <span v-else-if="r.provider !== 'TunnHub'" style="color: #007bbb">用户</span>
                    <span v-else style="color: #67c23a">Hub</span>
                    <!--                <span v-if="r.provider === ''" style="color: #990055">未知用户</span>-->
                    <!--                <span v-else>{{ r.provider.length > 5 ? r.provider.substring(0, 5) + "..." : r.provider }}</span>-->
                    提供
                  </div>
                </div>
              </div>
            </div>
          </template>
        </el-popover>
      </el-option>
    </el-select>
    <el-button size="small" style="margin-left: 5px" @click="submit">添加</el-button>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "RouteImportSelector",
  props: {
    account: {
      type: String
    },
    imported: {
      type: Array
    },
    exported: {
      type: Array
    }
  },
  data() {
    return {
      loading: false,
      selected: [],
      routes: []
    }
  },
  mounted() {
    //this.load()
  },
  methods: {
    submit: function () {
      let selected = []
      for (let i = 0; i < this.routes.length; i++) {
        for (let j = 0; j < this.selected.length; j++) {
          if (this.routes[i].network === this.selected[j]) {
            selected.push(this.routes[i])
          }
        }
      }
      this.$emit("submit", selected)
    },
    load: function () {
      this.loading = true
      let accountParam = this.account === '' ? '' : '/' + this.account
      axios({
        method: "get",
        url: "/api/v1/cfg/route/available" + accountParam,
        data: {}
      }).then(res => {
        let response = res.data
        let routes = []
        for (let i = 0; i < response.data.length; i++) {
          if (!this.isDuplicate(response.data[i].network)) {
            routes.push(response.data[i])
          }
        }
        this.routes = routes
        this.loading = false
      }).catch((err) => {
        this.$utils.HandleError(err)
        this.loading = false
      })
    },
    isDuplicate: function (network) {
      for (let i = 0; i < this.imported.length; i++) {
        if (this.imported[i].network === network) {
          return true
        }
      }
      for (let i = 0; i < this.exported.length; i++) {
        if (this.exported[i].network === network) {
          return true
        }
      }
      return false
    }
  }
}
</script>

<style scoped>
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