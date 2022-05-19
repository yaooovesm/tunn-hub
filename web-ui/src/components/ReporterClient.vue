<template>
  <div style="display: none"></div>
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
    }
  },
  mounted() {
  },
  unmounted() {
    this.Close()
  },
  methods: {
    Start: function () {
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