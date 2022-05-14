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
          style="height: 45px;margin-bottom: 2px"
      >
        <div style="line-height: 15px;padding: 5px 0">
          <div style="font-size: 13px;" v-if="r.name.length>0">{{ r.name }}</div>
          <div style="font-size: 13px;color: #909399" v-else>未命名</div>
          <div style="font-size: 12px;margin-top: 5px;color: #007bbb">{{ r.network }}</div>
        </div>
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
      axios({
        method: "get",
        url: "/api/v1/cfg/route/available",
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

</style>