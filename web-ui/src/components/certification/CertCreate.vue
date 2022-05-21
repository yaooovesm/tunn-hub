<template>
  <div>
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">创建证书</div>
      </div>
      <div style="padding: 20px">
        <div style="text-align: left">
          <div class="current-val">当前证书 <span>{{ security.cert }}</span></div>
          <div class="current-val">当前秘钥 <span>{{ security.key }}</span></div>
          <div style="display: inline-block;margin-top: 20px">
            <el-button size="small">创建</el-button>
            <el-checkbox v-model="overwrite" :value="overwrite" size="small"
                         style="margin-left: 10px;transform: translateY(2px)">覆盖配置
            </el-checkbox>
          </div>

        </div>
        <div style="margin-top: 20px">
          <el-collapse>
            <el-collapse-item title="访问限制" name="1">
              <div style="text-align: left">
                <div class="current-val" style="font-size: 12px;color: black">允许以下IP地址连接</div>
                <el-tag
                    v-for="addr in addresses" :key="addr"
                    closable
                    effect="dark"
                    type="info"
                    :disable-transitions="false"
                    @close="handleAddressesDelete(addr)"
                    style="margin-right: 10px;margin-bottom: 5px"
                >
                  {{ addr }}
                </el-tag>
                <div style="margin-top: 10px">
                  <el-input v-model="addressAdd" size="small" style="width: calc(100% - 60px)"></el-input>
                  <el-button size="small" style="width: 50px;margin-left: 10px" @click="handleAddressesAdd">添加
                  </el-button>
                </div>
              </div>
              <div style="margin-top: 20px;text-align: left">
                <div class="current-val">允许以下域名连接</div>
                <el-tag
                    v-for="name in names" :key="name"
                    closable
                    effect="dark"
                    type="info"
                    :disable-transitions="false"
                    @close="handleNamesDelete(name)"
                    style="margin-right: 10px;margin-bottom: 5px"
                >
                  {{ name }}
                </el-tag>
                <div style="margin-top: 10px">
                  <el-input v-model="nameAdd" size="small" style="width: calc(100% - 60px)"></el-input>
                  <el-button size="small" style="width: 50px;margin-left: 10px" @click="handleNamesAdd">添加</el-button>
                </div>
              </div>
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "CertCreate",
  data() {
    return {
      overwrite: false,
      addresses: [],
      addressAdd: "",
      names: [],
      nameAdd: "",
      security: {
        cert: "",
        key: ""
      }
    }
  },
  mounted() {
    this.getCurrentConfig()
  },
  methods: {
    handleNamesAdd: function () {
      if (this.nameAdd !== "") {
        if (!this.isDomain(this.nameAdd)) {
          this.$utils.Warning("警告", "不是一个有效的域名")
          return
        }
        this.names.push(this.nameAdd)
        this.nameAdd = ""
      }
    },
    handleNamesDelete: function (name) {
      for (let i = 0; i < this.names.length; i++) {
        if (this.names[i] === name) {
          this.names.splice(i, 1)
          return
        }
      }
    },
    handleAddressesAdd: function () {
      if (this.addressAdd !== "") {
        if (!this.isIPv4(this.addressAdd)) {
          this.$utils.Warning("警告", "不是一个有效的IPv4地址")
          return
        }
        this.addresses.push(this.addressAdd)
        this.addressAdd = ""
      }
    },
    handleAddressesDelete: function (addr) {
      for (let i = 0; i < this.addresses.length; i++) {
        if (this.addresses[i] === addr) {
          this.addresses.splice(i, 1)
          return
        }
      }
    },
    isIPv4: function (str) {
      let reg = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/
      return reg.test(str);
    },
    isDomain: function (str) {
      let domain = /^([\w-]+\.)+((com)|(net)|(org)|(gov\.cn)|(info)|(cc)|(com\.cn)|(net\.cn)|(org\.cn)|(biz)|(tv)|(cn)|(mobi)|(name)|(sh)|(ac)|(io)|(tw)|(com\.tw)|(hk)|(com\.hk)|(ws)|(travel)|(us)|(tm)|(la)|(me\.uk)|(org\.uk)|(ltd\.uk)|(plc\.uk)|(in)|(eu)|(it)|(jp))$/;
      return domain.test(str);
    },
    getCurrentConfig: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/v1/server/config",
        data: {}
      }).then(res => {
        let response = res.data
        if (this.isDomain(response.data.global.address)) {
          this.names.push(response.data.global.address)
        } else {
          this.addresses.push(response.data.global.address)
        }
        this.security = response.data.security
        this.updateTime = new Date()
        this.loading = false
      }).catch((err) => {
        this.$utils.HandleError(err)
        this.loading = false
      })
    }
  }
}
</script>

<style scoped>
.current-val {
  font-size: 13px;
  line-height: 15px;
  margin-bottom: 10px;
  display: block;
  color: #404040;
}

.current-val span {
  color: #007bbb;
}
</style>