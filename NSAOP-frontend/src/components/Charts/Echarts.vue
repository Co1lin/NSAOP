<template>
  <div class="container">
    <div id="chart" ref="chart" :class="inHome? 'homeChart' : 'chart'">
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts'

export default {
  name: "Echarts",
  props: {
    inHome: {
      type: Boolean,
      default() {
        return false
      }
    },
    traffics: {
      type: Array,
      default() {
        return []
      }
    },
    currentDate: {
      type: Date,
      default() {
        return new Date()
      }
    },
    wrapperStyle: {
      type: Object,
      default() {
        return {
          width: 600,
          height: 500,
        }
      }
    }
  },
  data() {
    return {
      chart: {},
    }
  },
  computed: {
    chartStyle() {
      return {
        width: this.wrapperStyle.width + 'px',
        height: this.wrapperStyle.height + 'px',
      }
    },
    trafficsFixed() {
      let ret = this.traffics.map(item => (item / 1024).toFixed(2))
      if (ret.length === 1){
        ret[1] = ret[0]
        ret[0] = 0
        // console.log(ret)
      }
      return ret
    },
    totalTraffic() {
      return (this.traffics.reduce((prev, cur) => {
        return prev + cur
      }, 0) / 1024).toFixed(2)
    },
    option() {
      return {
        title: {
          text: '近30天内流量统计',
          subtext: '总流量' + this.totalTraffic + 'GB'
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            label: {
              backgroundColor: '#6a7985'
            }
          }
        },
        xAxis: {
          type: 'category',
          show: false,
          boundaryGap: false,
          data: this.dateSpan
        },
        yAxis: {
          type: 'value',
          axisLabel: {
            formatter: '{value}GB',
            showMinLabel: false,
            showMaxLabel: false,
          },
          min: (value) => {
            return value.min * 0.9
          },
          max: (value) => {
            return value.max * 1.1
          }
        },
        series: [{
          data: this.trafficsFixed,
          name: '流量',
          type: 'line',
          showSymbol: false,
          lineStyle: {
            width: 0
          },
          areaStyle: {
            opacity: 0.8,
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{
              offset: 0,
              color: 'rgba(128, 255, 165)'
            }, {
              offset: 1,
              color: 'rgba(1, 191, 236)'
            }])
          },
          smooth: true,
        }],
        grid: {
          left: '10%',
          top: '15%',
        },
      };
    },
    dateSpan() {
      if (this.traffics !== null && this.traffics !== undefined && this.traffics.length !== 0) {
        let ret = [],  date = this.currentDate
        let length = (this.traffics.length == 1) ? this.traffics.length + 1 : this.traffics.length
        while (length > 0) {
          ret.push(this.transformDate(date))
          date.setDate(date.getDate() - 1)
          length = length - 1
        }
        return ret.reverse()
      }
      return []
    },
  },
  watch: {
    option() {
      this.chart.setOption(this.option)
    }
  },
  mounted() {
    this.chart = echarts.init(this.$refs['chart'])
    setTimeout(() => {
      this.chart.setOption(this.option)
      window.onresize = () => this.chart.resize()
    }, 300)
  },
  methods: {
    transformDate(date) {
      return (date.getMonth() + 1) + '月' + date.getDate() + '日'
    }
  }
}
</script>

<style lang="scss" scoped>
.container {
  display: flex;
  justify-content: center;
  text-align: center;
}

.chart {
  width: $base-app-column-width;
  height: $base-app-row1-height;
  min-height: 500px;
}

.homeChart {
  width: $base-app-column-width;
  height: $base-app-row2-height;
  min-height: 400px;
}
</style>
