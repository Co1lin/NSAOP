<template>
  <div class="table-container">
    <el-table
      ref="tableSort"
      v-loading="listLoading"
      :data="orderList"
      :element-loading-text="elementLoadingText"
      :height="height"
      :row-style="{fontFamily: 'STHeiti'}"
      :header-cell-style="{backgroundColor: '#F0F0F0', color: '#191919'}"
    >
      <!--    @selection-change="setSelectRows"-->
      <!--    @sort-change="tableSortChange"-->
      <el-table-column
        show-overflow-tooltip
        label="序号"
        width="80"
        align="center"
      >
        <template #default="scope">
          {{ scope.$index + 1 + (displayPageNO - 1) * queryForm.pageSize }}
        </template>
      </el-table-column>
      <el-table-column
        show-overflow-tooltip
        label="备注名"
        min-width="180"
        prop="comment"
        align="center"
      />
      <el-table-column
        show-overflow-tooltip
        label="结算方式"
        prop="paytype"
        align="center"
      />
      <el-table-column
        show-overflow-tooltip
        label="订单状态"
        min-width="120"
        align="center"
      >
        <template #default="{ row, }">
          <!--          <el-tag :type="row.status | statusFilter">-->
          <el-tag :type="type(row.status, )" :effect="effect(row.status, )">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        show-overflow-tooltip
        label="创建时间"
        min-width="150"
        prop="createTime"
        align="center"
      />
      <el-table-column
        show-overflow-tooltip
        label="操作"
        width="250px"
        align="center"
      >
        <template #default="{ row, }">
<!--          <el-button
            type="text"
            @click.native.stop="handleDetail(row,)"
          >
            订单详情
          </el-button>-->
          <el-button
            v-if="orderStatus === 'first'"
            type="text"
            @click.native.stop="handleDeploy(row)"
          >
            部署订单
          </el-button>
          <el-button
            v-else
            type="text"
            @click.native.stop="handleDeploy(row)"
          >
            设备变更
          </el-button>
<!--          <el-button-->
<!--            v-if="orderStatus === 'first'"-->
<!--            type="text"-->
<!--            @click.native.stop="handleOperate(row, 'on',)"-->
<!--          >-->
<!--            确认部署-->
<!--          </el-button>-->
          <el-button
            v-if="row.status === '待回收设备'"
            type="text"
            @click.native.stop="handleOperate(row, 'canceled',)"
          >
            确认回收
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      :background="true"
      :current-page="queryForm.pageNo"
      :layout="layout"
      :total="total"
      :page-size="queryForm.pageSize"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script>
import {changeOrder, getOrderByUser} from "@/network/order";
import {getEffect, getType, parseOrder} from "@/common/order";
import {SET_TARGET} from "@/common/store";
import {loadingText} from "@/config/setting.config"

export default {
  name: "EngineerVue",
  props: {
    status: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      orderStatus: '',
      orderList: [],
      total: 0,
      listLoading: false,
      layout: 'total, prev, pager, next, jumper',
      selectRows: '',
      elementLoadingText: loadingText,
      displayPageNO: 1,
      queryForm: {
        pageNo: 1,
        pageSize: 10,
        title: '',
      },
    }
  },
  computed: {
    height() {
      return this.$baseTableHeight()
    },
  },
  mounted() {
    this.reload()
  },
  methods: {
    handleDetail(row) {
      let id = row.id
      this.$router.push({
        name: 'detail',
        params: {id: id}
      })
    },
    type(status) {
      return getType(status)
    },
    effect(status) {
      return getEffect(status)
    },
    handleDeploy(row) {
      this.$router.push('/orders/deploy/' + row.id)
    },
    handleOperate(row, target) {
      let confirmMessage, successMessage, failedMessage = ''
      if(target === 'on') {
        confirmMessage = "是否部署订单？"
        successMessage = "部署订单成功"
        failedMessage = "部署订单失败，请稍后重试"
      } else if(target === "canceled") {
        confirmMessage = "是否已成功回收设备？"
        successMessage = "订单已取消，感谢您为公司的付出！"
        failedMessage = "请稍后重试"
      }
      this.$confirm(confirmMessage, {
        title: "确认信息"
      }).then(() => {
        changeOrder({
          id: row.id,
          target: target
        }).then(() => {
          this.reload()
          this.$message({
            message: successMessage,
            type: "success",
            duration: 1000,
          })
        }).catch(() => {
          this.$message({
            message: failedMessage,
            type: "warning",
            duration: 1000,
          })
        })
      })
    },
    handleCurrentChange(val) {
      this.reload(val)
    },
    fetchData() {
      return new Promise((resolve, reject) => {
        let fetchStatus
        if(this.status === "first") {
          fetchStatus = ["pass"]
        } else if(this.status === "second") {
          fetchStatus = ["retrieve"]
        } else if(this.status === "third")
          fetchStatus = ["on", "pause"]
        let target = this.$store.state.orderQuery.target
        let data = {
          offset: (this.queryForm.pageNo - 1) * this.queryForm.pageSize,
          limit: this.queryForm.pageSize,
          status: fetchStatus,
        }
        if(target !== '') {
          this.$store.dispatch(SET_TARGET, {data: ''})
          data.search = target
        }
        getOrderByUser(data).then(res => {
          this.total = res.data.count
          if(res.data.services !== null) {
            res.data.services.forEach(order => {
              this.orderList.push(parseOrder(order))
            })
            resolve()
          }
          else {
            resolve()
          }
        }).catch(() => reject())
      })
    },
    reload(page = 1) {
      this.orderList = []
      let promise = new Promise((resolve, reject) => {
        this.listLoading = true
        this.queryForm.pageNo = page
        setTimeout(() => {
          this.orderStatus = this.status
          setTimeout(() => {
            this.fetchData().then(() => {
              this.listLoading = false
            }).then(resolve()).catch(() => {
              this.$message({
                message: "加载订单失败，请刷新页面重试",
                type: "failed",
                duration: 1000,
              })
              this.listLoading = false
              this.orderList = []
              reject()
            })
          },200)
        }, 200)
      })
      promise.then(() => this.displayPageNO = page)
      return promise
    }
  }
}
</script>

<style scoped>

</style>
