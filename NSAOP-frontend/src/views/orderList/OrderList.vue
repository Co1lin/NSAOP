<template>
  <div>
    <layout>
      <template #main-content>
        <div class="order-list-container">
          <control-bar @reload="reload" />
          <el-tabs
            v-if="$store.state.user.role === 'customer'"
            tab-position="top"
            @tab-click="handleClick"
            v-model="status"
          >
            <el-tab-pane label="处理中" name="first"/>
            <el-tab-pane label="使用中" name="second"/>
            <el-tab-pane label="已取消" name="third"/>
            <customer
              ref="list"
              :status="status"
            />
          </el-tabs>
          <el-tabs
            v-if="$store.state.user.role === 'operator'"
            tab-position="top"
            @tab-click="handleClick"
            v-model="status"
          >
            <el-tab-pane label="待审批" name="first"/>
            <el-tab-pane label="审批完成" name="second"/>
            <operator
              ref="list"
              :status="status"
            />
          </el-tabs>
          <el-tabs
            v-if="$store.state.user.role === 'engineer'"
            tab-position="top"
            @tab-click="handleClick"
            v-model="status"
          >
            <el-tab-pane label="待部署" name="first"/>
            <el-tab-pane label="设备回收" name="second"/>
            <el-tab-pane label="部署完成" name="third"/>
            <engineer
              ref="list"
              :status="status"
            />
          </el-tabs>
        </div>
      </template>
    </layout>
  </div>
</template>

<script>
import Layout from "@/components/Layout";
import Customer from "./Customer"
import ControlBar from "./ControlBar"
import Operator from "./Operator"
import Engineer from "./Engineer"

export default {
  name: "PersonalCenter",
  components: {
    Layout,
    Customer,
    Operator,
    ControlBar,
    Engineer
  },
  data() {
    return {
      status: ""
    }
  },
  computed: {
    token() {
      return this.$store.state.user.token
    },
  },
  mounted() {
    let tabs = sessionStorage.getItem('currentTab')
    if(tabs !== null && tabs !== ''){
      this.status = tabs
    } else {
      this.status = "first"
    }
  },
  methods: {
    handleClick(tab) {
      if (tab.index === '0') {
        this.status = "first"
      } else if (tab.index === '1') {
        this.status = "second"
      } else if (tab.index === '2') {
        this.status = "third"
      }
      sessionStorage.setItem('currentTab', this.status)
      this.reload()
    },
    reload() {
      this.$refs.list.reload()
    }
  },
}
</script>

<style scoped>
.order-list-container{
  margin-left: 20px;
  margin-right: 20px;
  padding: 20px;
  text-align: center;
}

.information-container{
  margin: 50px;
}
</style>
