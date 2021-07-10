<template>
  <div class="table-container">
    <el-dialog
      title="请扫码支付订单"
      :visible.sync="dialogVisible"
      class="dialog-container"
      width="30%"
    >
      <span>
        <div class="message">
          <p>订单id：{{ selectedOrder.id }}</p>
          <img
            class="logo-display"
            src="@/assets/img/default_qrcode.png"
            height="200"
          >
        </div>
      </span>
      <template #footer class="dialog-footer">
        <el-button @click="cancelPay">
          取消
        </el-button>
        <el-button type="primary" @click="submitPay">
          已支付
        </el-button>
      </template>
    </el-dialog>
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
      <!--      <el-table-column-->
      <!--        v-if="orderStatus !== 'waiting'"-->
      <!--        show-overflow-tooltip-->
      <!--        prop="device"-->
      <!--        label="设备"-->
      <!--        align="center"-->
      <!--      />-->
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
          <el-tooltip
            v-if="row.message !== undefined"
            class="item"
            effect="dark"
            :content="'驳回理由: ' + row.message"
            trigger
            placement="top">
            <div>
              <el-tag :type="type(row.status, )" :effect="effect(row.status, )">
                被驳回
              </el-tag>
            </div>
          </el-tooltip>
          <el-tag v-else :type="type(row.status, )" :effect="effect(row.status, )">
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
          <el-button
            type="text"
            @click.native.stop="handleDetail(row,)"
          >
            订单详情
          </el-button>
          <el-button
            v-if="orderStatus === 'first'"
            type="text"
            @click.native.stop="handleOperate(row, 'canceled',)"
          >
            取消订单
          </el-button>
          <el-button
            v-if="row.status === '运行中'"
            type="text"
            @click.native.stop="handleOperate(row, 'pause',)"
          >
            暂停订单
          </el-button>
          <el-button
            v-if="row.status === '暂停使用'"
            type="text"
            @click.native.stop="handleOperate(row, 'on', )"
          >
            恢复订单
          </el-button>
          <el-button
            v-if="row.status === '已欠费'"
            type="text"
            @click.native.stop="handlePay(row, )"
          >
            续费重启
          </el-button>
          <el-button
            v-if="row.status === '暂停使用'"
            type="text"
            @click.native.stop="handleOperate(row, 'retrieve', )"
          >
            取消订单
          </el-button>
          <el-button
            v-if="row.status === '待回收设备'"
            type="text"
            @click.native.stop="handleUndo(row, )"
          >
            撤销
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
import {getOrderByUser, changeOrder} from "@/network/order";
import {getEffect, getType, parseOrder} from "@/common/order";
import {SET_TARGET} from "@/common/store";
import {loadingText} from "@/config/setting.config"

export default {
  name: "OrderBlock",
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
      dialogVisible: false,
      displayPageNO: 1,
      queryForm: {
        pageNo: 1,
        pageSize: 10,
        title: '',
      },
      selectedOrder: {},
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
    handleCurrentChange(val) {
      this.reload(val)
    },
    handlePay(row) {
      this.selectedOrder = row
      this.dialogVisible = true
    },
    handleOperate(row, target) {
      let confirmMessage, successMessage, failedMessage = ''
      switch (target) {
      case "canceled":
        confirmMessage = "是否取消订单? "
        if(row.status === "待工程师部署")
          confirmMessage += "工程师正在准备您的订单，现在取消需要扣取30%初装费。"
        successMessage = "取消订单成功"
        failedMessage = "取消订单失败，请稍后重试"
        break
      case "on":
        confirmMessage = "是否恢复订单?"
        successMessage = "恢复订单成功"
        failedMessage = "恢复订单失败，请稍后重试"
        break
      case "pause":
        confirmMessage = "是否暂停订单?"
        successMessage = "暂停订单成功"
        failedMessage = "暂停订单失败，请稍后重试"
        break
      case "retrieve":
        confirmMessage = "是否取消已生效订单？将申请设备回收。您每月有两次机会使待回收设备的订单重新生效。"
        successMessage = "取消订单成功"
        failedMessage = "取消订单失败，请稍后重试"
      }
      this.$confirm(confirmMessage, {
        title: "温馨提示",
      })
        .then(() => {
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
    fetchData() {
      this.orderList = []
      let promise = new Promise((resolve, reject) => {
        let fetchStatus
        if(this.status === "first")
          fetchStatus = ["waiting", "pass"]
        else if(this.status === "second")
          fetchStatus = ["on", "pause", "suspend", "retrieve"]
        else if(this.status === "third")
          fetchStatus = ["canceled"]
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
      return promise
    },
    cancelPay() {
      this.dialogVisible = false
    },
    submitPay() {
      this.dialogVisible = false
      changeOrder({
        id: this.selectedOrder.id,
        target: "on"
      }).then(() => {
        this.reload()
        this.$message({
          message: "订单已重启",
          type: "success",
          duration: 1000,
        })
      }).catch(() => {
        this.$message({
          message: "重启失败，请重试",
          type: "warning",
          duration: 1000,
        })
      })
    },
    handleUndo(row) {
      this.$confirm("工程师正在回收设备的路上，是否将订单重新生效？每月您有两次机会使订单重新生效。", {
        title: "温馨提示"
      }).then(() => {
        changeOrder({
          id: row.id,
          target: "pause"
        }).then(() => {
          this.reload()
          this.$message({
            message: "订单已重新生效",
            type: "success",
            duration: 1000,
          })
        }).catch(err => {
          if(err === "undo limit exceed") {
            this.$message({
              message: "本月撤销次数已满！",
              type: "warning",
              duration: 1000,
            })
          } else {
            this.$message({
              message: "撤销失败，请重试",
              type: "warning",
              duration: 1000,
            })
          }
        })
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

<style lang="scss" scoped>
  .dialog-container {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    overflow: auto;
    margin: 0px;

    ::v-deep {
      .el-dialog {
        max-width: 500px !important;
        margin-top: 15vh !important;
        width: 80% !important;

        .el-dialog__body {
          padding: 0;
        }
      }
    }
  }
</style>
