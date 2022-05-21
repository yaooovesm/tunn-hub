<template>
  <div>
    <el-card shadow="always" body-style="padding:0" v-loading="loading">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">创建证书</div>
      </div>
      <div style="padding: 20px">
        <div style="text-align: left">
          <div class="current-val">当前证书 <span>{{ security.cert }}</span></div>
          <div class="current-val">当前秘钥 <span>{{ security.key }}</span></div>
          <div style="margin-top: 20px;font-size: 12px;color: #909399">
            <span style="color: #404040">备注：</span>
            <div>证书默认在7天后到期，如需要设置可在<span style="color: #007bbb"> 高级设置 </span>中设置。</div>
            <div>如果需要更改配置则需要勾选 <span style="color: #007bbb">覆盖配置</span>，未勾选时仅创建证书和私钥不修改配置。</div>
            <div>证书和秘钥更新后不会马上生效，需要 <span style="color: #007bbb">重启</span> 服务器后生效</div>
          </div>
          <div style="display: inline-block;margin-top: 20px">
            <el-button size="small" @click="handleCreate">创建</el-button>
            <el-checkbox v-model="overwrite" :value="overwrite" size="small"
                         style="margin-left: 10px;transform: translateY(2px)">覆盖配置
            </el-checkbox>
          </div>
        </div>
        <el-collapse style="margin-top: 20px">
          <el-collapse-item title="高级设置" name="1">
            <div style="padding: 10px">
              <div style="text-align: left">
                <div class="current-val" style="font-size: 12px;color: black">允许以下IP地址连接</div>
                <el-tag
                    v-if="addresses.length<=0"
                    effect="dark"
                    type="info"
                    :disable-transitions="false"
                    style="margin-right: 10px;margin-bottom: 5px"
                >
                  未设置
                </el-tag>
                <el-tag
                    v-for="addr in addresses" :key="addr"
                    closable
                    effect="dark"
                    type="success"
                    :disable-transitions="false"
                    @close="handleAddressesDelete(addr)"
                    style="margin-right: 10px;margin-bottom: 5px"
                >
                  {{ addr }}
                </el-tag>
                <div style="margin-top: 10px">
                  <el-input v-model="addressAdd" size="small" style="width: 220px"></el-input>
                  <el-button size="small" style="width: 50px;margin-left: 10px" @click="handleAddressesAdd">添加
                  </el-button>
                </div>
              </div>
              <div style="margin-top: 20px;text-align: left">
                <div class="current-val">允许以下域名连接</div>
                <el-tag
                    v-if="names.length<=0"
                    effect="dark"
                    type="info"
                    :disable-transitions="false"
                    style="margin-right: 10px;margin-bottom: 5px"
                >
                  未设置
                </el-tag>
                <el-tag
                    v-for="name in names" :key="name"
                    closable
                    effect="dark"
                    type="success"
                    :disable-transitions="false"
                    @close="handleNamesDelete(name)"
                    style="margin-right: 10px;margin-bottom: 5px"
                >
                  {{ name }}
                </el-tag>
                <div style="margin-top: 10px">
                  <el-input v-model="nameAdd" size="small" style="width: 220px"></el-input>
                  <el-button size="small" style="width: 50px;margin-left: 10px" @click="handleNamesAdd">添加</el-button>
                </div>
              </div>
              <div style="margin-top: 20px;text-align: left">
                <div class="current-val">证书过期时间</div>
                <div style="margin-top: 10px">
                  <el-date-picker
                      v-model="before"
                      type="datetime"
                      placeholder="选择证书过期时间"
                      :shortcuts="shortcuts"
                      size="small"
                  />
                </div>
              </div>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
    </el-card>
  </div>
</template>

<script>
import axios from "axios";
import {ElMessageBox} from "element-plus";

export default {
  name: "CertCreate",
  data() {
    return {
      loading: false,
      overwrite: false,
      addresses: [],
      addressAdd: "",
      names: [],
      nameAdd: "",
      before: new Date(),
      security: {
        cert: "",
        key: ""
      },
      shortcuts: [
        {
          text: '7天后',
          value: () => {
            const date = new Date()
            date.setTime(date.getTime() + 3600 * 1000 * 24 * 7)
            return date
          },
        },
        {
          text: '15天后',
          value: () => {
            const date = new Date()
            date.setTime(date.getTime() + 3600 * 1000 * 24 * 15)
            return date
          },
        },
        {
          text: '30天后',
          value: () => {
            const date = new Date()
            date.setTime(date.getTime() + 3600 * 1000 * 24 * 30)
            return date
          },
        },
      ]
    }
  },
  created() {
    this.before.setTime(this.before.getTime() + 3600 * 1000 * 24 * 7)
  },
  mounted() {
    this.getCurrentConfig()
  },
  methods: {
    handleCreate: function () {
      this.loading = true
      ElMessageBox.confirm(
          this.overwrite ? '创建证书，并<span style="color: #007bbb">修改</span>当前配置' : '仅创建证书，不修改配置',
          '是否创建证书',
          {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: this.overwrite ? 'warning' : 'info',
            dangerouslyUseHTMLString: true,
          }
      ).then(() => {
        let data = {
          overwrite: this.overwrite,
          addresses: this.addresses,
          names: this.names,
          before: this.before.getTime()
        }
        data = JSON.parse(JSON.stringify(data))
        //console.log(data)
        axios({
          method: "post",
          url: "/api/v1/server/cert/create",
          data: data
        }).then(res => {
          let response = res.data
          let name = response.data
          this.getCurrentConfig()
          this.$utils.Success("创建成功", "证书名称：" + name)
          this.loading = false
        }).catch((err) => {
          this.$utils.HandleError(err)
          this.loading = false
        })
      }).catch(() => {
        this.loading = false
      })
    },
    handleNamesAdd: function () {
      if (this.nameAdd !== "") {
        if (this.isIPv4(this.nameAdd)) {
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
    getCurrentConfig: function () {
      this.loading = true
      this.names = []
      this.addresses = []
      axios({
        method: "get",
        url: "/api/v1/server/config",
        data: {}
      }).then(res => {
        let response = res.data
        if (!this.isIPv4(response.data.global.address)) {
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