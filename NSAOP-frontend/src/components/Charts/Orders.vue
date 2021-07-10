<template>
  <div class="Echarts">
    <div id="main" class="chart" />
  </div>
</template>

<script>
import * as echarts from 'echarts'

export default {
  name: 'Echarts',
  props: {
    pieData:{
      type: Array,
      default() {
        return[
          {value: 948, name: '运行中'},
          {value: 484, name: '暂停中'},
          {value: 300, name: '已取消'},
          {value: 735, name: '待审批'},
          {value: 580, name: '待部署'},
        ]
      }
    }
  },
  mounted() {
    this.myEcharts()
  },
  watch: {
    pieData: {
      handler: function(){
        if(this.pieData.length === 5){
          this.myEcharts()
        }
      },
      deep: true
    }
  },
  methods:{
    myEcharts(){
      let chartDom = document.getElementById('main');
      let myChart = echarts.init(chartDom);
      window.onresize = function () {
        myChart.resize();
      };
      let option;
      let pieData = this.pieData
      let colorList = ['#9EC97F', '#587EC0', '#DE6E7B', '#84DFDB','#F3C9'];
      option = {
        legend: [
          {
            icon: "none",
            bottom: 'center',
            left: '55%',
            orient: 'horizontal',
            textStyle:{
              fontSize: 16,//字体大小
              rich: {
                a:{
                  color:colorList[0],
                  padding:[0,10],
                  fontSize:16
                },
                b:{
                  color:colorList[1],
                  padding:[0,10],
                  fontSize:16
                },
                c:{
                  color:colorList[2],
                  padding:[0,10],
                  fontSize:16
                },
                d:{
                  color:colorList[3],
                  padding:[0,10],
                  fontSize:16
                },
                e:{
                  color:colorList[4],
                  padding:[0,10],
                  fontSize:16
                },
              }
            },
            itemGap:50,
            formatter: function(name) {
              let data = pieData;
              let tarValue = 0;
              let tar = 0;
              for (let i = 0, l = data.length; i < l; i++) {
                if (data[i].name == name) {
                  tarValue = data[i].value;
                  tar = i;
                }
              }
              let prefix = '';
              if (tar == 0) {
                prefix = 'a|'
              }
              if (tar == 1) {
                prefix = 'b|'
              }
              if (tar == 2) {
                prefix = 'c|'
              }
              if (tar == 3) {
                prefix = 'd|'
              }
              if (tar == 4) {
                prefix = 'e|'
              }
              let arr = [
                name+'         ',
                '{'+prefix+tarValue+'}',
              ]
              return arr.join('')
            },
            data: pieData
          },
        ],
        series: [
          {
            name: '状态分布',
            legendHoverLink: false,
            center: ["25%", "50%"],
            type: 'pie',
            radius: ['35%', '60%'],
            avoidLabelOverlap: false,
            itemStyle: {
              normal:{
                borderRadius: 10,
                borderColor: '#fff',
                borderWidth: 4,
                color:function(params) {
                  return colorList[params.dataIndex]
                }
              },
            },
            label: {
              show: false,
              position: 'center'
            },
            emphasis: {
              label: {
                show: true,
                fontSize: '25',
                fontWeight: 'bold',
                formatter: [
                  '{b}',
                  '{c}'
                ].join('\n'),
              }
            },
            labelLine: {
              show: false
            },
            data: pieData
          }
        ]
      };
      option && myChart.setOption(option);
      window.addEventListener("resize", () => {
        myChart.resize();
      });
      //自定义事件，当鼠标移动上是，则设置 title 不显示
      myChart.on("mouseover", params => {
        myChart.setOption({
          title: {
            show: false //当鼠标移动上是，则设置 title 不显示
          },
          series: {
            label: {
              emphasis: {
                rich: {
                  c: {
                    //获取 该区域的颜色值，为对应的企业数字体设置颜色
                    color: params.color
                  }
                }
              }
            }
          }
        });
      });
      myChart.dispatchAction({
        type: 'highlight',
        seriesIndex: 0,
        dataIndex: 0
      });
      myChart.on('mouseover', function(e) {
        //当检测到鼠标悬停事件，取消默认选中高亮
        if (e.dataIndex === 0) {
          return
        } else {
          myChart.dispatchAction({
            type: 'downplay',
            seriesIndex: 0,
            dataIndex: 0
          });
        }
      });
      //检测鼠标移出后显示之前默认高亮的那块
      myChart.on('mouseout', function() {
        myChart.dispatchAction({
          type: 'highlight',
          seriesIndex: 0,
          dataIndex: 0
        });
      });
    }
  }
}
</script>

<style lang="scss" scoped>
.Echarts {
  display:flex;
  justify-content: flex-start;
}

.chart {
  height: $base-app-row1-height;
  width: $base-app-column-width;
  min-height: 400px;
  min-width: 575px;
}
</style>

