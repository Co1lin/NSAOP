<template>
  <div>
    <layout>
      <template #main-content>
        <div v-if="showParallel" class="parallel-container">
          <el-row class="row1">
            <el-col :span="10">
              <el-card class="user-info">
                <img :src="avatar" class="icon-container">
                <p> {{ name }} </p>
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
            </el-col>
            <el-col :span="14">
              <el-card
                class="user-info"
                v-loading="!infoAvailable"
              >
                <orders id="orders" :pie-data="orders"/>
              </el-card>
            </el-col>
          </el-row>
          <el-row class="row2">
            <el-card
              class="traffic-chart"
              v-loading="!trafficsAvailable"
              element-loading-text="无运行订单，暂无流量统计数据"
              element-loading-spinner="none"
              element-loading-background="rgba(255, 255, 255, 0.8)"
            >
              <traffic :traffics="traffics" :in-home="true"/>
            </el-card>
          </el-row>
        </div>
        <div v-else class="parallel-container">
          <el-row class="row2">
            <el-card class="traffic-chart">
              <img :src="avatar" class="icon-container">
              <p> {{ name }} </p>
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
          </el-row>
          <el-row>
            <el-card
              class="user-info"
              v-loading="!infoAvailable"
            >
              <orders :pie-data="orders"/>
            </el-card>
          </el-row>
          <el-row class="row2">
            <el-card
              class="traffic-chart"
              v-loading="!trafficsAvailable"
              element-loading-text="无运行订单，暂无流量统计数据"
              element-loading-spinner="none"
              element-loading-background="rgba(255, 255, 255, 0.8)"
            >
              <traffic :traffics="traffics" :in-home="true"/>
            </el-card>
          </el-row>
        </div>
      </template>
    </layout>
  </div>
</template>

<script>
import Layout from "@/components/Layout";
import {authTest} from "@/network/user";
import {revokeToken} from "@/common/auth";
import {SET_COMPANY, SET_EMAIL, SET_TEL, SET_USERNAME} from "@/common/store";
import Orders from "@/components/Charts/Orders"
import Traffic from "@/components/Charts/Echarts"
import {mapGetters} from "vuex";
import {getOrderByUser} from "@/network/order";

export default {
  name: "Home",
  components: {
    Layout,
    Orders,
    Traffic
  },
  data() {
    return {
      trafficsAvailable: false,
      infoAvailable: false,
      screenWidth: window.innerWidth,
      traffics: [
        104800.04306060415,154361.8992706353,5608.8439445213335,64814.66525230676,
        243859.55773116267,93093.10156667,148.04305285736936,40518.43229885954,20806.935173529742,
        29557.41470242722,123654.31946820758,7235.641275488192,6756.899792342021,111111.24028684436,
        52194.49819023299,34999.09011302613,43787.64898538983,40320.26966019834,2857.1192233989104,
        17537.193280653573,3287.746686057889,78339.15804771104,58654.02275387468,28951.615788742765,
        30748.731216004522,117720.29386141762,63489.09375131385,3228.7268291547175,21429.544458482924,
        138016.20738018677
      ],
      orders: [
        // {value: 948, name: '运行中'},
        // {value: 484, name: '暂停中'},
        // {value: 300, name: '已取消'},
        // {value: 735, name: '待审批'},
        // {value: 580, name: '待部署'},
      ]
    }
  },
  computed: {
    ...mapGetters({
      collapse: "settings/collapse",
    }),
    isCollapse() {
      this.$refs["orders"].myEcharts()
      return this.collapse
    },
    avatar() {
      let username = this.name
      let hashCode = 0
      for(let i = 0; i < username.length; i++){
        hashCode = 37 * hashCode + username.charCodeAt(i)
      }
      hashCode = hashCode % 21 + 1
      return require("@/assets/img/user_image/bottts (" + hashCode + ").png")
    },
    infoData() {
      return [{
        header: "联系电话",
        data: this.tel,
      }, {
        header: "邮箱",
        data: this.email,
      }, {
        header: "公司名称",
        data: this.company,
      }]
    },
    token() {
      return this.$store.state.user.token
    },
    name() {
      return this.$store.state.user.username
    },
    tel() {
      return this.$store.state.user.tel
    },
    email() {
      return this.$store.state.user.email
    },
    company() {
      return this.$store.state.user.company
    },
    showParallel() {
      if(this.screenWidth >= 992)
        return true
      else
        return false
    }
  },
  mounted() {
    window.addEventListener('resize', () => this.screenWidth = window.innerWidth, false)
    if (this.token !== undefined && this.token !== '') {
      authTest().then(res => {
        this.$store.commit(SET_USERNAME, res.username)
        this.$store.commit(SET_EMAIL, res.email)
        this.$store.commit(SET_TEL,  res.phone)
        this.$store.commit(SET_COMPANY, res.company)

        let order = []
        let request = {
          offset: 0,
          limit: 1,
          status: ['on']
        }
        const promise1 = getOrderByUser(request).then(res => {
          order[0] = {
            value: res.data.count,
            name: '运行中'
          }
        }).catch(() => {})
        request.status = ['pause', 'suspend']
        const promise2 = getOrderByUser(request).then(res => {
          order[1] = {
            value: res.data.count,
            name: '暂停中'
          }
        }).catch(() => {})
        request.status = ['canceled', 'retrieve']
        const promise3 = getOrderByUser(request).then(res => {
          order[2] = {
            value: res.data.count,
            name: '已取消'
          }
        }).catch(() => {})
        request.status = ['waiting']
        const promise4 = getOrderByUser(request).then(res => {
          order[3] = {
            value: res.data.count,
            name: '待审批'
          }
        }).catch(() => {})
        request.status = ['pass']
        const promise5 = getOrderByUser(request).then(res => {
          order[4] = {
            value: res.data.count,
            name: '待部署'
          }
        }).catch(() => {})

        Promise.all([promise1, promise2, promise3, promise4, promise5]).then(() => {
          let count = 0
          for(let i = 0; i < 5; i++) {
            this.orders.push(order[i])
            if(i < 3)
              count += order[i].value
          }
          if(count > 0)
            this.trafficsAvailable = true
          this.infoAvailable = true
        })
      }).catch(() => {
        // console.log(err);
        revokeToken()
        this.$store.commit("SET_TOKEN", "")
        this.$router.push('/login')
      })
    }
  },
  updated() {
    window.onresize = () => {
      return (() => {
        this.handleWidthChange();
      })()
    }
  },
  methods: {
    handleWidthChange() {
      this.screenWidth = window.innerWidth
    },
    cellStyle({column}) {
      let cellStyle
      if(column.property === "header")
        cellStyle = "font-weight: 700; textAlign: center"
      else
        cellStyle = "textAlign: center"
      return cellStyle
    },
  }
}
</script>

<style lang="scss" scoped>
.parallel-container {
  padding: 10px;
  /deep/
  .user-info {
    text-align: center;
    margin: 10px;
    height: $base-app-column1-height;
    min-height: 380px;
    border-radius: 10px;

    .el-card__body {
      position: relative;
      top: 50%;
      transform: translateY(-50%);
    }
  }
  /deep/
  .traffic-chart {
    text-align: center;
    margin: 10px;
    height: $base-app-column2-height;
    min-height: 400px;
    border-radius: 10px;

    .el-card__body {
      position: relative;
      top: 50%;
      transform: translateY(-45%);
    }
  }

  .icon-container {
    height: 100px;
    border-radius: 100px;
  }

  .row1 {
    height: $base-app-row1-height;
    align-content: center;
  }

  .row2 {
    height: $base-app-row2-height;
    min-height: 420px;
    align-content: center;
  }
}

</style>
