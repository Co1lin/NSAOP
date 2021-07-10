<template>
  <layout>
    <template #main-content>
      <create-device
        :visible.sync="deviceDialogVisible"
        @send="receiveDevice"
      />
      <create-ssid
        :visible.sync="ssidDialogVisible"
        @send="receiveSSID"
      />
      <display-detail :visible.sync="displayDevice">
        <template #attr>
          <div v-for="(item, name) of deviceDetail" :key="name">
            <h4> {{ item.name }} </h4>
            {{ item.value }}
          </div>
        </template>
      </display-detail>
      <display-detail :visible.sync="displaySSID">
        <template #attr>
          <div v-for="(item, name) of ssidDetail" :key="name">
            <h4> {{ item.name }} </h4>
            {{ item.value }}
          </div>
        </template>
      </display-detail>

      <div v-if="showParallel" class="root">
        <div class="section">
          <el-card>
            <template #header>
              <h2>详情栏</h2>
            </template>
            <el-button
              type="success"
              @click="finishDeploy"
              :disabled="info.status !== 'pass'"
            >
              完成部署
            </el-button>
            <div class="bicolumn">
              <div>
                <b>订单号</b>
              </div>
              <div class="right">
                {{ id }}
              </div>
              <div>
                <b>订单名称</b>
              </div>
              <div class="right">
                {{ info.nickname }}
              </div>
              <div>
                <b>地址</b>
              </div>
              <div class="right">
                {{ info.address.address }}
              </div>
              <div>
                <b>计费方式</b>
              </div>
              <div class="right">
                {{ chargeMethod }}
              </div>
              <div>
                <b>需求</b>
              </div>
              <div class="right">
                <span v-for="(item, index) in requirements" :key="index">
                  {{ item }}
                </span>
              </div>
              <div>
                <b>联系人</b>
              </div>
              <div class="right">
                <span>{{ info.address.contact }}</span>
              </div>
              <div>
                <b>客户电话</b>
                <i class="el-icon-phone" />
              </div>
              <div class="right">
                <span>{{ info.address.phone }}</span>
              </div>
            </div>
          </el-card>
        </div >
        <div class="section" >
          <el-card v-loading="loadDevice">
            <template #header>
              <h2>管理设备</h2>
            </template>
            <div class="buttons">
              <el-button type="danger" @click="addDevice">
                添加新设备
              </el-button>
              <el-button type="danger" :loading="loadDevice" @click="submitChange('device',)">
                保存更改
              </el-button>
            </div>
            <h3>现有设备</h3>
            <div class="tags" ref="parallelDevice">
              <el-tag
                v-for="(device, index) in existedDevices"
                :key="index + 'exist'"
                effect="dark"
                closable
                @close="handleExistedClose(device, 'device')"
                @click.native="showDetail(device, 'device',)"
              >
                {{ device.name.length>deviceLength?device.name.substring(0, deviceLength - 1)+'...':device.name}}
              </el-tag>
              <el-tag
                v-for="(device, index) in devices"
                :key="index"
                effect="dark"
                closable
                @close="handleClose(device, 'device',)"
                @click.native="showDetail(device, 'device',)"
              >
                {{ device.name.length>deviceLength?device.name.substring(0, deviceLength - 1)+'...':device.name}}
              </el-tag>
            </div>
          </el-card>
        </div>
        <div class="section" >
          <el-card v-loading="loadSSID">
            <template #header>
              <h2>管理SSID</h2>
            </template>
            <div class="buttons">
              <el-button type="danger" @click="addSSID">
                添加新SSID
              </el-button>
              <el-button type="danger" :loading="loadSSID" @click="submitChange('ssid',)" >
                保存更改
              </el-button>
            </div>
            <h3>现有SSID</h3>
            <div class="tags" ref="parallelSSID">
              <el-tag
                v-for="(ssid, index) in existedSSIDs"
                :key="index + 'exist'"
                type="info"
                effect="dark"
                closable
                @close="handleExistedClose(ssid, 'ssid')"
                @click.native="showDetail(ssid, 'ssid',)"
              >
                {{ ssid.name.length > ssidLength ? ssid.name.substring(0, ssidLength - 1) + '...' : ssid.name}}
              </el-tag>
              <el-tag
                v-for="(ssid, index) in ssids"
                :key="index"
                type="info"
                effect="dark"
                closable
                @close="handleClose(ssid, 'ssid',)"
                @click.native="showDetail(ssid, 'ssid',)"
              >
                {{ ssid.name.length > ssidLength ? ssid.name.substring(0, ssidLength - 1) + '...' : ssid.name}}
              </el-tag>
            </div>
          </el-card>
        </div>
      </div>
      <div v-else>
        <el-row>
          <div class="section">
            <el-card>
              <template #header>
                <h2>详情栏</h2>
              </template>
              <el-button
                type="success"
                @click="finishDeploy"
                :disabled="info.status !== 'pass'"
              >
                完成部署
              </el-button>
              <div class="bicolumn">
                <div>
                  <b>订单号</b>
                </div>
                <div class="right">
                  {{ id }}
                </div>
                <div>
                  <b>订单名称</b>
                </div>
                <div class="right">
                  {{ info.nickname }}
                </div>
                <div>
                  <b>地址</b>
                </div>
                <div class="right">
                  {{ info.address.address }}
                </div>
                <div>
                  <b>计费方式</b>
                </div>
                <div class="right">
                  {{ chargeMethod }}
                </div>
                <div>
                  <b>需求</b>
                </div>
                <div class="right">
                  <span v-for="(item, index) in requirements" :key="index">
                    {{ item }}
                  </span>
                </div>
                <div>
                  <b>联系人</b>
                </div>
                <div class="right">
                  <span>{{ info.address.contact }}</span>
                </div>
                <div>
                  <b>客户电话</b>
                  <i class="el-icon-phone" />
                </div>
                <div class="right">
                  <span>{{ info.address.phone }}</span>
                </div>
              </div>
            </el-card>
          </div >
        </el-row>
        <el-row>
          <div class="section" >
            <el-card v-loading="loadDevice">
              <template #header>
                <h2>管理设备</h2>
              </template>
              <div class="buttons">
                <el-button type="danger" @click="addDevice">
                  添加新设备
                </el-button>
                <el-button type="danger" :loading="loadDevice" @click="submitChange('device',)">
                  保存更改
                </el-button>
              </div>
              <h3>现有设备</h3>
              <div class="tags" ref="verticalDevice">
                <el-tag
                  v-for="(device, index) in existedDevices"
                  :key="index + 'exist'"
                  effect="dark"
                  closable
                  @close="handleExistedClose(device, 'device')"
                  @click.native="showDetail(device, 'device',)"
                >
                  {{ device.name.length>deviceLength?device.name.substring(0, deviceLength - 1)+'...':device.name}}
                </el-tag>
                <el-tag
                  v-for="(device, index) in devices"
                  :key="index"
                  effect="dark"
                  closable
                  @close="handleClose(device, 'device',)"
                  @click.native="showDetail(device, 'device',)"
                >
                  {{ device.name.length>deviceLength?device.name.substring(0, deviceLength - 1)+'...':device.name}}
                </el-tag>
              </div>
            </el-card>
          </div>
        </el-row>
        <el-row>
          <div class="section" >
            <el-card v-loading="loadSSID">
              <template #header>
                <h2>管理SSID</h2>
              </template>
              <div class="buttons">
                <el-button type="danger" @click="addSSID">
                  添加新SSID
                </el-button>
                <el-button type="danger" :loading="loadSSID" @click="submitChange('ssid',)" >
                  保存更改
                </el-button>
              </div>
              <h3>现有SSID</h3>
              <div class="tags" ref="verticalSSID">
                <el-tag
                  v-for="(ssid, index) in existedSSIDs"
                  :key="index + 'exist'"
                  type="info"
                  effect="dark"
                  closable
                  @close="handleExistedClose(ssid, 'ssid')"
                  @click.native="showDetail(ssid, 'ssid',)"
                >
                  {{ ssid.name.length > ssidLength ? ssid.name.substring(0, ssidLength - 1) + '...' : ssid.name}}
                </el-tag>
                <el-tag
                  v-for="(ssid, index) in ssids"
                  :key="index"
                  type="info"
                  effect="dark"
                  closable
                  @close="handleClose(ssid, 'ssid',)"
                  @click.native="showDetail(ssid, 'ssid',)"
                >
                  {{ ssid.name.length > ssidLength ? ssid.name.substring(0, ssidLength - 1) + '...' : ssid.name}}
                </el-tag>
              </div>
            </el-card>
          </div>
        </el-row>
      </div>
    </template>
  </layout>
</template>

<script>
import {changeOrder, getOrderInfoById} from "@/network/order";
import layout from "@/components/Layout";
import CreateDevice from "@/components/Deploy/CreateDevice"
import CreateSsid from "@/components/Deploy/CreateSSID";
import DisplayDetail from "@/components/Deploy/DisplayDetail";
import {deepcopy, NCE_ERROR, filterNCEError} from "@/common/utils";
import {createDevices, createSSIDs, deleteDevices, deleteSSIDs, getDevices, getSSIDs} from "@/network/deploy";

export default {
  name: "Detail",
  components: {
    DisplayDetail,
    layout,
    CreateDevice,
    CreateSsid,
  },
  data() {
    return {
      screenWidth: window.innerWidth,
      info: {
        nickname: "",
        address: "",
        paytype: "",
        requirement: 0,
        status: "",
        createTime: "",
        userTel: "",
      },
      devices: [],
      ssids: [],
      existedDevices: [],
      existedSSIDs: [],
      devicesToDelete: [],
      ssidsToDelete: [],
      deviceDialogVisible: false,
      ssidDialogVisible: false,
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
      displayDevice: false,
      displaySSID: false,
      loadSSID: false,
      loadDevice: false,
    }
  },
  computed: {
    ssidLength() {
      let clientWidth
      if(this.showParallel) {
        if(this.$refs.parallelSSID !== undefined)
          clientWidth = this.$refs.parallelSSID.clientWidth
      }
      else {
        if (this.$refs.verticalSSID !== undefined)
          clientWidth = this.$refs.verticalSSID.clientWidth
      }
      return clientWidth / 25
    },
    deviceLength() {
      let clientWidth = 100
      if(this.showParallel) {
        if(this.$refs.parallelDevice !== undefined)
          clientWidth = this.$refs.parallelDevice.clientWidth
      }
      else {
        if(this.$refs.verticalDevice !== undefined)
          clientWidth = this.$refs.verticalDevice.clientWidth
      }
      return clientWidth / 25
    },
    role() {
      return this.$store.getters.role
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
    chargeMethod() {
      if (this.info.paytype === "year") {
        return "按年计费"
      } else if (this.info.paytype === "month") {
        return "按月计费"
      } else {
        return "计费方式获取异常"
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
    },
    devicesMap() {
      return this.devices.map(device => {
        return {
          device_model: device.type,
          name: device.name,
        }
      })
    },
    ssidsMap() {
      return this.ssids.map(ssid => {
        return {
          connection_mode: ssid.mode,
          enable: true,
          hided_enable: ssid.hide,
          max_user_number: ssid.userNum,
          name: ssid.name,
          relative_radios: ssid.radios,
          user_separation: ssid.separation,
        }
      })
    },
    showParallel() {
      if(this.screenWidth >= 992)
        return true
      else
        return false
    },

  },
  watch: {
    role() {
      if (this.role !== undefined && this.role !== '') {
        if (this.role !== 'engineer') {
          this.$router.push('/404')
        }
      }
    }
  },
  async created() {
    await this.getOrderInfo()
    this.getDevices()
    this.getSSIDs()
  },
  mounted() {
    window.addEventListener('resize', () => this.screenWidth = window.innerWidth, false)
  },
  methods: {
    getOrderInfo() {
      getOrderInfoById(this.cred).then(res => {
        this.info.device = res.device
        this.info.createTime = res.create_at
        this.info.paytype = res.paytype
        this.info.status = res.status
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
      this.loadDevice = true
      getDevices(this.id).then(res => {
        this.existedDevices = res.data
        this.loadDevice = false
      }).catch(err => {
        this.loadDevice = false
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
      this.loadSSID = true
      getSSIDs(this.id).then(res => {
        this.existedSSIDs = res.data
        this.loadSSID = false
      }).catch(err => {
        this.loadSSID = false
        if (err === NCE_ERROR) {
          this.$message({
            message: "出现NCE错误，请联系管理员",
            type: "error",
            duration: 1000,
          })
        }
      })
    },
    addDevice() {
      this.deviceDialogVisible = true
    },
    addSSID() {
      this.ssidDialogVisible = true
    },
    handleClose(item, type) {
      if (type === "device") {
        this.devices.splice(this.devices.indexOf(item), 1)
      } else if (type === "ssid") {
        this.ssids.splice(this.devices.indexOf(item), 1)
      }
    },
    handleExistedClose(item, type) {
      if (type === "device") {
        this.devicesToDelete.push(item)
        this.existedDevices.splice(this.existedDevices.indexOf(item), 1)
      } else if (type === "ssid") {
        this.ssidsToDelete.push(item)
        this.existedSSIDs.splice(this.existedSSIDs.indexOf(item), 1)
      }
    },
    receiveDevice(device) {
      const temp = deepcopy(device)
      this.devices.push(temp)
    },
    receiveSSID(ssid) {
      const temp = deepcopy(ssid)
      this.ssids.push(temp)
    },
    showDetail(item, type) {
      if (type === 'device' && item.type !== undefined) {
        item.deviceModel = item.type
      }
      if (type === 'ssid' && item.userNum !== undefined) {
        item.hidedEnable = item.enable
        item.hidedEnable = item.hide
        item.connectionMode = item.mode
        item.relativeRadios = item.radios
        item.userSeparation = item.separation
        item.maxUserNumber = item.userNum
      }
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
    async submitChange(type) {
      let flag, errors = []
      if (type === 'device') {
        this.loadDevice = true
        await Promise.all([
          createDevices({service_id: this.id, devices: this.devicesMap})
            .then(res => res)
            .catch(err => errors.push(err)),
          deleteDevices({service_id: this.id, device_ids: this.devicesToDelete.map(item => item.id)})
            .then(res => res)
            .catch(err => errors.push(err)),
        ]).then(() => flag = true).catch(() => flag = false)
        setTimeout(() => this.loadDevice = false, 700)
      } else if (type === 'ssid') {
        this.loadSSID = true
        await Promise.all([
          createSSIDs(this.id, this.ssidsMap)
            .then(res => res)
            .catch(err => errors.push(err)),
          deleteSSIDs({service_id: this.id, ssid_ids: this.ssidsToDelete.map(item => item.id)})
            .then(res => res)
            .catch(err => errors.push(err))
        ]).then(() => flag = true).catch(() => flag = false)
        setTimeout(() => this.loadSSID = false, 700)
      }
      if (flag) {
        this.$message({
          message: "提交成功",
          type: "success",
          duration: 1000,
        })
      } else if (filterNCEError(errors).catch(err => err === NCE_ERROR)) {
        this.$message({
          message: "出现NCE错误，请联系管理员",
          type: "error",
          duration: 1000,
        })
      } else {
        this.$message({
          message: "提交失败，请重试",
          type: "error",
          duration: 1000,
        })
      }
    },
    finishDeploy() {
      this.$confirm("是否完成部署? 未提交的内容不会保存", {
        title: "确认信息"
      }).then(() => {
        changeOrder({
          id: this.id,
          target: 'on',
        }).then(() => {
          this.$message({
            message: "部署成功",
            type: "success",
            duration: 1000,
          })
          sessionStorage.setItem('currentTab', "third")
          this.$router.push('/orders')
        }).catch(() => {
          this.$message({
            message: "部署失败",
            type: "error",
            duration: 1000,
          })
        })
      })
    },
  }
}
</script>

<style scoped lang="scss">

.root {
  display: flex;
  flex-direction: row;
  margin: 10px;
  .section {
    flex: 1;
  }
}

.el-card {
  height: 100%;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  .el-tag:hover {
    cursor: pointer;
  }
  /deep/ .el-tag {
    padding-top: 18px;
    padding-bottom: 18px;
    padding-left: 20px;
    margin: 2% 5px 2% 5px;
    border-radius: 4px;
    font-size: 20px;
    width: 45%;
    display: flex;
    align-items: center;

    .el-tag__close {
      margin-top: 4px;
      margin-left: auto;
    }
  }

  /deep/ .el-tag:last-child {
    margin-right: 0!important;
  }
}



.el-button {
  width: 48%;
}

.buttons {
  display: flex;
  flex-direction: row;
  position: relative;
  bottom: 0;
  .el-button {
    flex: 1;
    width: 48%;
    margin: 0 5px 10px 0;
  }
}


.bicolumn {
  font-size: 18px;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  position: relative;
  div {
    margin-top: 20px;
    width: 49%;
  }
  .right {
    text-align: right;
  }
}

/deep/ .el-card {
  margin: 7px;
  box-shadow: 0 4px 4px 0 rgba(0, 0, 0, 0.05), 0 6px 20px 0 rgba(0, 0, 0, 0.03);
  border-radius: 15px;
}
</style>
