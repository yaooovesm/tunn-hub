<template>
  <div style="margin-bottom: 30px">
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">流量监控
        </div>
      </div>
      <div style="padding: 20px">
        <div v-if="!noData" :id="id" style="width: 100%;height: 200px"/>
        <div v-else
             style="width: 100%;padding-top: 80px;padding-bottom:80px;font-size: 15px;color: #aaaaaa;text-align: center;">
          <div style="height: 40px">
            无数据或数据量不足
            <el-tooltip
                effect="dark"
                placement="right"
                content="记录时间过短(<=2分钟)或记录功能已关闭"
            >
              <i class="iconfont icon-info-circle"
                 style="font-size: 10px;margin-left: 5px;color: #007bbb;transform: translateY(-1px);display: inline-block"></i>
            </el-tooltip>
          </div>
        </div>
      </div>
      <div style="font-size: 12px;color: #808080;text-align: right;padding: 5px 10px">
        <el-button type="text" @click="changeRange"
                   style="font-size: 12px;height: 12px;line-height: 13px">{{ range === '1h' ? "近24小时数据" : "近1小时数据" }}
        </el-button>
        <el-button type="text" @click="update()" style="font-size: 12px;height: 12px;line-height: 13px">刷新</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import * as echarts from "echarts";
import axios from "axios";

export default {
  name: "TrafficRecorderOverview",
  created() {
    this.id = "traffic-recorder-chart-" + this.guid()
  },
  mounted() {
    this.chart = echarts.init(document.getElementById(this.id));
    // this.chart.showLoading({
    //   text: "加载中",
    //   x: "center",
    //   y: "center",
    //   textStyle: {
    //     color: "#409EFF",
    //     fontSize: 14
    //   },
    //   effect: "spin"
    // })
    // 绘制图表
    this.chart.setOption(this.option)
    this.update(false)
    this.timer = setInterval(() => {
      this.update(true)
    }, 60000)
  },
  unmounted() {
    clearInterval(this.timer)
  },
  data() {
    this.chart = null
    return {
      range: "1h",
      id: "",
      noData: false,
      data: [],
      timer: undefined,
      option: {
        tooltip: {
          trigger: 'axis',
          formatter: function (params) {
            //console.log(params)
            //let timestamp = echarts.format.formatTime('yyyy-MM-dd hh:mm:ss', params[0].data[0]);
            return "<div style='font-size: 12px;width: 200px'>" +
                "<div style='text-align: left;'>接收流量 <div style='float: right;'><span style='color: #007bbb'>" + params[0].value + "</span> M</div></div>" +
                "<div style='text-align: left'>发送流量  <div style='float: right;'><span style='color: #007bbb'>" + params[1].value + "</span> M</div>" +
                "<div style='text-align: left'>接收速率(AVG)  <div style='float: right;'><span style='color: #007bbb'>" + (params[0].value / 60).toFixed(3) + "</span> M/s</div></div>" +
                "<div style='text-align: left'>发送速率(AVG)  <div style='float: right;'><span style='color: #007bbb'>" + (params[1].value / 60).toFixed(3) + "</span> M/s</div></div>" +
                "<div style='text-align: left'>记录时间  <div style='float: right;'><span style='color: #007bbb'>" +
                echarts.format.formatTime('yyyy/MM/dd hh:mm', new Date(parseInt(params[0].axisValue))) +
                "-" + echarts.format.formatTime('hh:mm', new Date(parseInt(params[0].axisValue) + 60000)) +
                "</span></div></div>" +
                "</div>"
          }
        },
        grid: {
          left: "50px",
          top: "12px",
          right: "20px",
          bottom: "30px"
        },
        xAxis: {
          data: [],
          type: 'category',
          nameLocation: 'center',
          boundaryGap: false,
          nameGap: '40',
          axisLine: {
            show: true
          },
          axisTick: {
            show: false
          },
          axisLabel: {
            show: true,
            formatter: function (value) {
              return echarts.format.formatTime('hh:mm', new Date(parseInt(value)))
            },
            //interval: 0
          },
          splitLine: {
            show: true,    // 是否显示分隔线。默认数值轴显示，类目轴不显示
            //interval: '0',    // 坐标轴刻度标签的显示间隔，在类目轴中有效.0显示所有
            // color分隔线颜色，可设置单个颜色，也可设置颜色数组，分隔线会按数组中颜色顺序依次循环设置颜色
            color: ['#ccc'],
            width: 1, // 分隔线线宽
            type: 'solid', // 坐标轴线线的类型（solid实线类型；dashed虚线类型；dotted点状类型）
          },
          // splitArea: {
          //   show: true, // 是否显示分隔区域
          //   interval: '0', // 坐标轴刻度标签的显示间隔，在类目轴中有效.0显示所有
          //   areaStyle: {
          //     // color分隔区域颜色。分隔区会按数组中颜色顺序依次循环设置颜色。默认是一个深浅的间隔色
          //     color: ['rgba(250,250,250,0.3)', 'rgba(200,200,200,0.5)'],
          //     opacity: 1, // 图形透明度。支持从0到1的数字，为0时不绘制该图形
          //   },
          // },
        },
        yAxis: [
          {
            type: 'value',
            // max: function (value) {
            //   return parseInt(value.max) + 1
            // },
            min: 0,
            nameLocation: 'center',
            nameRotate: '90',
            nameGap: '20',
            nameTextStyle: {
              padding: [0, 0, 20, 100]    // 四个数字分别为上右下左与原位置距离
            },
            // axisLine: {
            //   show: false,
            // },
            // splitLine: {
            //   show: false
            // },
            axisLabel: {
              formatter: '{value} M',
            },
          },
        ],
        series: [
          {
            name: "接收流量",
            data: [],
            type: 'line',
            showSymbol: false,
            itemStyle: {
              normal: {
                color: "#249F68",
                lineStyle: {
                  width: 1
                },
              }
            },
          },
          {
            name: "发送流量",
            data: [],
            type: 'line',
            showSymbol: false,
            itemStyle: {
              normal: {
                color: "#2A7DDE",
                lineStyle: {
                  width: 1
                },
              }
            },
          },
        ]
      },
    }
  },
  methods: {
    changeRange: function () {
      this.range === '1h' ? this.range = '24h' : this.range = '1h'
      this.update(false)
    },
    update: function (silence) {
      if (!silence) {
        this.chart.showLoading({
          text: "加载中",
          x: "center",
          y: "center",
          textStyle: {
            color: "#409EFF",
            fontSize: 14
          },
          effect: "spin"
        })
        axios({
          method: "get",
          url: "/api/v1/server/monitor/traffic/" + this.range,
          data: {}
        }).then(res => {
          let response = res.data
          if (response.data == null || response.data.length <= 2) {
            this.noData = true
          } else {
            this.noData = false
            let ts = []
            //let traffic = []
            let rxFlow = []
            let txFlow = []
            for (let i in response.data) {
              ts.push(response.data[i].timestamp)
              rxFlow.push(this.formatMb(response.data[i].rx_flow))
              txFlow.push(this.formatMb(response.data[i].tx_flow))
              // traffic.push({
              //   RxFlowSpeed: response.data[i].RxFlowSpeed,
              //   TxFlowSpeed: response.data[i].TxFlowSpeed,
              //   RxPacketSpeed: response.data[i].RxPacketSpeed,
              //   TxPacketSpeed: response.data[i].TxPacketSpeed,
              //   Timestamp: response.data[i].Timestamp
              // })
            }
            this.option.xAxis.data = ts
            this.option.series[0].data = rxFlow
            this.option.series[1].data = txFlow
            this.chart.setOption(this.option)
          }

        }).catch((err) => {
          this.$utils.HandleError(err)
        }).finally(() => {
          if (!silence) {
            this.chart.hideLoading();
          }
        })
      }
    },
    formatMb: function (data) {
      return (data / 1024 / 1024).toFixed(2)
    },
    formatDate: function (time, format) {
      let t = new Date(time);
      let tf = function (i) {
        return (i < 10 ? "0" : "") + i;
      };
      return format.replace(/yyyy|MM|dd|HH|mm|ss/g, function (a) {
        switch (a) {
          case "yyyy":
            return tf(t.getFullYear());
          case "MM":
            return tf(t.getMonth() + 1);
          case "mm":
            return tf(t.getMinutes());
          case "dd":
            return tf(t.getDate());
          case "HH":
            return tf(t.getHours());
          case "ss":
            return tf(t.getSeconds());
        }
      });
    },
    guid: function () {
      return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        let r = Math.random() * 16 | 0,
            v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
      });
    }
  }
}
</script>

<style scoped>
.test {
}
</style>