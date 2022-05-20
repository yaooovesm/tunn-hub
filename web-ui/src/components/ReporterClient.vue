<template>
  <div v-if="error!==''" style="margin-bottom: 10px">
    <el-alert v-if="error!==''" title="数据不可用" :description="error" type="error"
              style="box-shadow: 1px 1px 5px rgba(50, 50, 50, 0.2);">
      <template #default>
        <div style="text-align: left;color: #777777">{{ error }}</div>
      </template>
      <template #title>
        <i class="iconfont icon-times-circle" style="font-size: 10px"></i> <span
          style="font-size: 13px;color: #2c3e50;font-weight: bold">数据不可用</span>
      </template>
    </el-alert>
  </div>
</template>

<script>
import publicStorage from "@/public.storage";

export default {
  name: "ReporterClient",
  props: {
    resources: {
      type: Object,
    },
    interval: {
      type: Number,
    }
  },
  data() {
    return {
      ws: null,
      started: false,
      error: ""
    }
  },
  mounted() {
  },
  unmounted() {
    this.Close()
  },
  methods: {
    Start: function () {
      this.error = ""
      if ("WebSocket" in window) {
        let that = this
        publicStorage.Load()
        let ws = new WebSocket("ws://" + window.location.hostname + ":" + publicStorage.User.reporter + "/reporter")
        ws.onopen = function () {
          ws.send(JSON.stringify(
              {
                token: publicStorage.User.token,
                resources: that.resources,
                interval: that.interval <= 0 ? 5000 : that.interval
              }
          ))
        }
        ws.onmessage = function (e) {
          if (e.data instanceof Blob) {
            let blob = e.data;
            let reader = new FileReader();
            reader.readAsBinaryString(blob);
            reader.onload = function () {
              that.$emit("recv", reader.result)
            }
          }
        }
        ws.onclose = function () {
          that.$emit("closed")
        }
        this.ws = ws
        this.started = true
      } else {
        this.$emit("error", "unsupported")
        this.error = "浏览器不支持"
      }
    },
    Close: function () {
      if (!this.started) {
        return
      }
      try {
        this.ws.send("close")
        this.ws.close()
      } catch (err) {
        this.$emit("error", err)
      }
    }
  }
}
</script>

<style scoped>

</style>