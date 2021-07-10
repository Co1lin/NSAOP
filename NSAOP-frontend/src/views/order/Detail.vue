<template>
  <layout>
    <template #main-content>
      <display-detail :visible.sync="displayDevice">
        <template #attr>
          <div v-for="(item, name, index) of deviceDetail" :key="index" class="bicolumn">
            <div>
              <b>{{ item.name }}</b>
            </div>
            <div class="right">
              {{ item.value }}
            </div>
          </div>
        </template>
      </display-detail>
      <display-detail :visible.sync="displaySSID">
        <template #attr>
          <div v-for="(item, name, index) of ssidDetail" :key="index" class="bicolumn">
            <div>
              <b>{{ item.name }}</b>
            </div>
            <div class="right">
              {{ item.value }}
            </div>
          </div>
        </template>
      </display-detail>
      <div class="root">
        <div v-if="showParallel">
          <el-row>
            <div class="section1">
              <el-col :span="18">
                <el-card
                  v-loading="!trafficsAvailable"
                  element-loading-text="网络未运行，暂无流量统计数据"
                  element-loading-spinner="none"
                  element-loading-background="rgba(255, 255, 255, 0.8)"
                >
                  <echarts
                    v-loading="graphLoading"
                    :traffics="traffics"
                  />
                </el-card>
              </el-col>
              <el-col :span="6">
                <el-card id="steps" class="progress">
                  <h2>当前进度</h2>
                  <el-steps :space="chartStyle.height / 4" :active="activeStep" :direction="'vertical'">
                    <el-step title="提交订单" :description="info.createTime" :status="statusDisplay[0]" />
                    <el-step title="运营人员审核" :description="info.approveTime" :status="statusDisplay[1]" />
                    <el-step title="工程师部署" :description="info.deployTime" :status="statusDisplay[2]" />
                    <el-step title="网络正常运行" :description="currentDate" :status="statusDisplay[3]" />
                  </el-steps>
                </el-card>
              </el-col>
            </div>
          </el-row>
          <el-row>
            <div class="section2">
              <el-card class="info-container">
                <el-table
                  :data="infoData"
                  :cell-style="cellStyle"
                  :show-header="false"
                >
                  <el-table-column
                    prop="header"
                  />
                  <el-table-column
                    prop="data"
                  />
                </el-table>
              </el-card>
              <el-card class="info-container">
                <el-table
                  :data="orderData"
                  :cell-style="cellStyle"
                  :show-header="false"
                >
                  <el-table-column
                    prop="header"
                  />
                  <el-table-column
                    prop="data"
                  />
                </el-table>
              </el-card>
              <el-card v-loading="ssidLoading">
                <h3>SSID信息</h3>
                <div v-if="info.status!=='on'">
                  网络未搭设完成，暂无信息
                </div>
                <div
                  v-if="trafficsAvailable"
                  class="bicolumn2"
                >
                  <div v-for="(item, index) in ssids" :key="index">
                    <el-button type="warning" @click="showDetail(item, 'ssid')">
                      {{ item.name }}
                    </el-button>
                  </div>
                </div>
              </el-card>
              <el-card v-loading="deviceLoading">
                <h3>设备信息</h3>
                <div v-if="info.status!=='on'">
                  网络未搭设完成，暂无信息
                </div>
                <div v-if="trafficsAvailable" class="bicolumn2">
                  <div v-for="(item, index) in devices" :key="index">
                    <el-button type="primary" @click="showDetail(item, 'device')">
                      {{ item.name }}
                    </el-button>
                  </div>
                </div>
              </el-card>
            </div>
          </el-row>
        </div>
        <div v-else>
          <el-row>
            <div class="section1">
              <el-col :span="24">
                <el-card id="steps">
                  <h2>当前进度</h2>
                  <el-steps :space="chartStyle.height / 4" :active="activeStep" :direction="'horizontal'">
                    <el-step title="提交订单" :description="info.createTime" :status="statusDisplay[0]" />
                    <el-step title="运营人员审核" :description="info.approveTime" :status="statusDisplay[1]" />
                    <el-step title="工程师部署" :description="info.deployTime" :status="statusDisplay[2]" />
                    <el-step title="网络正常运行" :description="currentDate" :status="statusDisplay[3]" />
                  </el-steps>
                </el-card>
              </el-col>
            </div>
          </el-row>
          <el-row>
            <div class="section2">
              <el-card class="info-container-vertical">
                <el-table
                  :data="infoData"
                  :cell-style="cellStyle"
                  :show-header="false"
                >
                  <el-table-column
                    prop="header"
                  />
                  <el-table-column
                    prop="data"
                  />
                </el-table>
              </el-card>
            </div>
          </el-row>
          <el-row>
            <div class="section2">
              <el-card class="info-container-vertical">
                <el-table
                  :data="orderData"
                  :cell-style="cellStyle"
                  :show-header="false"
                >
                  <el-table-column
                    prop="header"
                  />
                  <el-table-column
                    prop="data"
                  />
                </el-table>
              </el-card>
            </div>
          </el-row>
          <el-row>
            <div class="section2">
              <el-card v-loading="ssidLoading">
                <h3>SSID信息</h3>
                <div v-if="info.status!=='on'">
                  等待工程师部署中，暂无信息
                </div>
                <div v-if="trafficsAvailable" class="bicolumn2" >
                  <div v-for="(item, index) in ssids" :key="index">
                    <el-button type="warning" @click="showDetail(item, 'ssid')">
                      {{ item.name }}
                    </el-button>
                  </div>
                </div>
              </el-card>
            </div>
          </el-row>
          <el-row>
            <div class="section2">
              <el-card v-loading="deviceLoading">
                <h3>设备信息</h3>
                <div v-if="info.status!=='on'">
                  网络未搭设完成，暂无信息
                </div>
                <div v-if="trafficsAvailable" class="bicolumn2">
                  <div v-for="(item, index) in devices" :key="index">
                    <el-button type="primary" @click="showDetail(item, 'device')">
                      {{ item.name }}
                    </el-button>
                  </div>
                </div>
              </el-card>
            </div>
          </el-row>
          <el-row>
            <div class="section3">
              <el-col :span="24">
                <el-card
                  v-loading="!trafficsAvailable"
                  element-loading-text="网络未运行，暂无流量统计数据"
                  element-loading-spinner="none"
                  element-loading-background="rgba(255, 255, 255, 0.8)"
                  class="vertical-chart"
                >
                  <echarts
                    v-loading="graphLoading"
                    :traffics="traffics"
                  />
                </el-card>
              </el-col>
            </div>
          </el-row>
        </div>
      </div>
    </template>
  </layout>
</template>

<script>
import {getOrderInfoById} from "@/network/order";
import layout from "@/components/Layout";
import DisplayDetail from "@/components/Deploy/DisplayDetail";
import {getDevices, getSSIDs} from "@/network/deploy";
import {getTraffic} from "@/network/traffic";
import Echarts from "@/components/Charts/Echarts";
import {NCE_ERROR} from "@/common/utils";

export default {
  name: "Detail",
  components: {
    Echarts,
    layout,
    DisplayDetail,
  },
  data() {
    return {
      ssidLoading: true,
      deviceLoading: true,
      graphLoading: true,
      screenWidth: window.innerWidth,
      info: {
        nickname: "",
        address: "",
        paytype: "",
        requirement: 0,
        status: "",
        createTime: "",
        approveTime: "",
        deployTime: "",
      },
      devices: [],
      ssids: [],
      displayDevice: false,
      displaySSID: false,
      deviceDetail: {
        name: "",
        type: "",
      },
      ssidDetail: {
        name: "",
        enable: true,
        mode: "",
        hide: false,
        radios: 7,
        userNum: 100,
        separation: true,
      },
      traffics: [
        104800.04306060415,154361.8992706353,5608.8439445213335,64814.66525230676,
        243859.55773116267,93093.10156667,148.04305285736936,40518.43229885954,20806.935173529742,
        29557.41470242722,123654.31946820758,7235.641275488192,6756.899792342021,111111.24028684436,
        52194.49819023299,34999.09011302613,43787.64898538983,40320.26966019834,2857.1192233989104,
        17537.193280653573,3287.746686057889,78339.15804771104,58654.02275387468,28951.615788742765,
        30748.731216004522,117720.29386141762,63489.09375131385,3228.7268291547175,21429.544458482924,
        138016.20738018677
      ],
      date: null,
    }
  },
  computed: {
    showParallel() {
      if(this.screenWidth >= 992)
        return true
      else
        return false
    },
    infoData() {
      let demand = ''
      for(let i = 0; i < this.requirements.length; i++){
        demand += this.requirements[i]
        demand += ' '
      }
      return[{
        header: "订单名称",
        data: this.info.nickname,
      }, {
        header: "地址",
        data: this.info.address.address,
      }, {
        header: "计费方式",
        data: this.chargeMethod,
      }, {
        header: "需求",
        data: demand,
      }, {
        header: '价格',
        data: '1GB/元'
      }]
    },
    orderData() {
      let orderData = [{
        header: "订单号",
        data: this.id,
      },{
        header: '联系人',
        data: this.info.address.contact
      },{
        header: '联系人电话',
        data: this.info.address.phone
      },{
        header: "运营人员电话",
        data: this.info.operatorTel,
      }]
      if(this.info.status !== 'waiting'){
        orderData.push({
          header: "工程师电话",
          data: this.info.engineerTel,
        })
      }
      return orderData
    },
    trafficsAvailable() {
      if (this.info.status === null || this.info.status === undefined) {
        return false
      }
      return ['on'].indexOf(this.info.status) !== -1
    },
    currentDate() {
      if (this.trafficsAvailable && this.date !== null && this.date !== undefined) {
        let ret = this.date.getFullYear() + '-'
        ret += this.appendHelper(this.date.getMonth() + 1) + '-'
        ret += this.appendHelper(this.date.getDate()) + ' '
        ret += this.appendHelper(this.date.getHours()) + ':'
        ret += this.appendHelper(this.date.getMinutes()) + ':'
        ret += this.appendHelper(this.date.getSeconds())
        return ret
      } else if (this.info.status === 'canceled') {
        return '订单被驳回/取消'
      } else if(this.info.status === 'pause') {
        return '订单暂停中'
      }
      return null
    },
    username() {
      return this.$store.getters.username
    },
    id() {
      return this.$route.params.id
    },
    cred() {
      return {
        id: this.id
      }
    },
    chartStyle() {
      return {
        width: document.documentElement.clientWidth * 0.5,
        height: document.documentElement.clientHeight * 0.5,
      }
    },
    chargeMethod() {
      if (this.info.paytype === "year") {
        return "按年计费"
      } else if (this.info.paytype === "month") {
        return "按月计费"
      } else {
        return "计费方式获取异常"
      }
    },
    activeStep() {
      if (this.info.status === "waiting") {
        return 2
      } else if (this.info.status === "pass") {
        return 3
      } else if (this.info.status === "on") {
        return 4
      } else if (this.info.status === "pause") {
        return 0
      } else if (this.info.status === "canceled") {
        return 1
      } else {
        return 1
      }
    },
    currentStep() {
      if (this.info.status === "waiting") {
        return "等待运营人员审核"
      } else if (this.info.status === "pass") {
        return "等待网络工程师架设"
      } else if (this.info.status === "on") {
        return "正常工作中"
      } else if (this.info.status === "pause") {
        return "服务暂停"
      } else if (this.info.status === "canceled") {
        return "已取消"
      } else {
        return "状态获取错误，请联系客服"
      }
    },
    statusDisplay() {
      if (this.currentStep === "提交成功") {
        return ["finish", "wait", "wait", "wait"]
      } else if (this.currentStep === "等待运营人员审核") {
        return ["success", "finish", "wait", "wait"]
      } else if (this.currentStep === "等待网络工程师架设") {
        return ["success", "success", "finish", "wait"]
      } else if (this.currentStep === "正常工作中") {
        return ["success", "success", "success", "success"]
      } else if (this.currentStep === "服务暂停") {
        return ["wait","wait","wait","wait"]
      } else if (this.currentStep === "已取消") {
        return ["error","error","error","error"]
      } else {
        return ["wait", "wait", "wait", "wait"]
      }
    },
    requirements() {
      let t = this.info.requirement
      let ret = []
      if (t >= 4) {
        ret.push('私用')
        t -= 4
      }
      if (t >= 2) {
        ret.push('客户使用')
        t -= 2
      }
      if (t >= 1) {
        ret.push('测试用')
      }
      return ret
    }
  },
  async created() {
    await this.getOrderInfo().then(async () => {
      this.date = new Date()
      if (this.info.status === 'on') {
        await this.getDevices()
        await this.getSSIDs()
        await this.getTraffic()
      }
      else {
        this.ssidLoading = false
        this.deviceLoading = false
      }
    })
  },
  mounted() {
    window.addEventListener('resize', () => this.screenWidth = window.innerWidth, false)
    setTimeout(() => {
      this.graphLoading = false
    }, 400)
  },
  methods: {
    cellStyle({column}) {
      let cellStyle
      if(column.property === "header")
        cellStyle = "font-weight: 700; textAlign: center"
      else
        cellStyle = "textAlign: center"
      return cellStyle
    },
    getOrderInfo() {
      return getOrderInfoById(this.cred).then(res => {
        this.info.device = res.device
        this.info.status = res.status
        this.info.createTime = res.create_at
        if (res.status !== "waiting" && res.status !== 'canceled') {
          this.info.approveTime = res.pass_at
        } else {
          this.info.approveTime = ""
        }
        if (this.trafficsAvailable) {
          this.info.deployTime = res.on_at
        } else {
          this.info.deployTime = ""
        }
        this.info.paytype = res.paytype
        this.info.requirement = res.require
        this.info.operatorTel = res.contact_operator
        this.info.engineerTel = res.contact_engineer
        this.info.nickname = res.comment
        this.info.address = res.location_info
      }).catch(() => {
        // console.log(err)
        this.$router.push('/404')
      })
    },
    getDevices() {
      getDevices(this.id).then(res => {
        this.devices = res.data
      }).then(() => {
        this.deviceLoading = false
      }).catch(err => {
        if (err === NCE_ERROR) {
          this.$message({
            message: "出现NCE错误，请联系管理员",
            type: "error",
            duration: 1000,
          })
        }
      })
    },
    getSSIDs() {
      getSSIDs(this.id).then(res => {
        this.ssids = res.data
      }).then(() => {
        this.ssidLoading = false
      }).catch(err => {
        if (err === NCE_ERROR) {
          this.$message({
            message: "出现NCE错误，请联系管理员",
            type: "error",
            duration: 1000,
          })
        }
      })
    },
    getTraffic() {
      getTraffic(this.id).then(res => {
        this.traffics = res.Traffic
      }).catch(() => {})
    },
    showDetail(item, type) {
      if (type === 'ssid') {
        this.ssidDetail = this.detailSSIDMap(item)
        this.displaySSID = true
      } else if (type === 'device') {
        this.deviceDetail = this.detailDeviceMap(item)
        this.displayDevice = true
      }
    },
    detailDeviceMap(item) {
      return {
        name: {name: '备注名', value: item.name},
        type: {name: '设备模型', value: item.deviceModel},
      }
    },
    detailSSIDMap(item) {
      return {
        name: {name: '备注名', value: item.name},
        mode: {name: '网络连接方式', value: item.connectionMode},
        hide: {name: '是否隐藏SSID', value: item.hidedEnable},
        userNum: {name: '最大用户数量', value: item.maxUserNumber},
        radios: {name: '射频类型', value: item.relativeRadios},
        separation: {name: '是否用户隔离', value:item.userSeparation},
      }
    },
    appendHelper(tail) {
      if (tail < 10) {
        return '0' + tail
      } else {
        return tail
      }
    },
    // loading(time = 3000) {
    //   const loading = this.$loading({
    //     lock: true,
    //     text: 'Loading',
    //     spinner: 'el-icon-loading',
    //     background: 'rgba(0, 0, 0, 0.7)'
    //   });
    //   setTimeout(() => {
    //     loading.close();
    //   }, time);
    // },
  }
}
</script>

<style scoped lang="scss">

.root {
  display: flex;
  flex-direction: column;
  padding: 10px;
}

.bicolumn {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  position: relative;
  font-size: 16px;
  div {
    margin-top: 20px;
    width: 49%;
  }
  .right {
    text-align: right;
  }
}

.section1 {
  display: flex;
  flex-direction: row;
  /deep/
  .el-step:last-child {
    flex-basis: 100px!important;
  }

  .progress {
    min-height: 500px;

    /deep/
    .el-card__body {
      height: 500px;
    }
  }

  #steps {
    flex: 1;
    padding-left: 30px;
  }
}

.section2 {
  display: flex;
  text-align: center;

  .info-container-vertical {
    height: 275px;

    /deep/
    .el-card__body {
      position: relative;
      top: 50%;
      transform: translateY(-50%);
    }
  }

  .info-container {
    height: 320px;

    /deep/
    .el-card__body {
      position: relative;
      top: 50%;
      transform: translateY(-50%);
    }
  }

  > * {
    flex: 1;
    margin: 10px;
  }
  .bicolumn2 {
    justify-content: center;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    position: relative;
    div {
      margin-top: 10px;
      width: 50%;
    }
    .el-button {
      width: 95%;
      border-radius: 0;
      font-size: 16px;
      line-height: 0.5;
    }
  }
}

.section3 {
  /deep/
  .vertical-chart {

  }

}

.traffic {
  /deep/ .el-dialog {
    width: 80%;
  }
}

/deep/ .el-card {
  margin: 7px;
  box-shadow: 0 4px 4px 0 rgba(0, 0, 0, 0.05), 0 6px 20px 0 rgba(0, 0, 0, 0.03);
  border-radius: 15px;
}
</style>
