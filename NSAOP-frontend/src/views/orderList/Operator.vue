<template>
  <div class="table-container">
    <el-dialog
      title="确认订单信息"
      :visible.sync="dialogVisible"
      class="dialog-container"
      width="30%"
    >
      <span>
        <el-table
          :data="tableData"
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
      </span>
      <template #footer class="dialog-footer">
        <el-button @click="cancelOrder">
          驳回订单
        </el-button>
        <el-button
          type="primary"
          :disabled="debounce"
          @click="passOrder"
          :loading="passLoading"
        >
          审核通过
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
                已驳回
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
        width="180px"
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
            @click.native.stop="handleOperate(row,)"
          >
            审核订单
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
import {changeOrder, getOrderByUser, getOrderInfoById} from "@/network/order";
import {parseOrder, getType, getEffect} from "@/common/order";
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
      displayPageNO: 1,
      dialogVisible: false,
      queryForm: {
        pageNo: 1,
        pageSize: 10,
        title: '',
      },
      selectedOrder: {
      },
      debounce: false,
      passLoading: false,
    }
  },
  computed: {
    height() {
      return this.$baseTableHeight()
    },
    tableData() {
      let selectedOrder = this.selectedOrder
      return [{
        header:  "订单id",
        data: selectedOrder.id
      }, {
        header:  "结算方式",
        data: selectedOrder.paytype
      }, {
        header:  "创建时间",
        data: selectedOrder.createTime
      }, {
        header:  "订单详情",
        data: selectedOrder.detail
      }, {
        header:  "用户名称",
        data: selectedOrder.contact
      }, {
        header:  "用户联系方式",
        data: selectedOrder.phone
      }, {
        header:  "订单地址",
        data: selectedOrder.location
      }]
    }
  },
  mounted() {
    this.reload()
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
    handleDetail(row) {
      let id = row.id
      this.$router.push({
        name: 'detail',
        params: {id: id}
      })
    },
    handleOperate(row) {
      this.selectedOrder = {}
      getOrderInfoById({
        id: row.id
      }).then(res => {
        row.detail = res.detail
        row.contact = res.location_info.contact
        row.phone = res.location_info.phone
        row.location = res.location_info.address
        this.dialogVisible = true
        this.debounce = false
      }).then(() => {
        this.selectedOrder = row
      }).catch(() => {
        // console.log(err)
        this.$message({
          message: "获取订单详情失败，请重试",
          type: "error",
          duration: 1000
        })
      })
    },
    handleCurrentChange(val) {
      this.reload(val)
    },
    type(status) {
      return getType(status)
    },
    effect(status) {
      return getEffect(status)
    },
    cancelOrder() {
      this.debounce = true
      this.$prompt('请输入驳回理由', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /[\s\S]{1,200}/,
        inputErrorMessage: '字段长度需在1-200之间'
      }).then(value => {
        changeOrder({
          id: this.selectedOrder.id,
          target: "canceled",
          message: value.value
        }).then(() => {
          this.reload()
          this.$message({
            message: "已成功驳回订单",
            type: "success",
            duration: 1000,
          })
        }).catch(() => {
          this.$message({
            message: "驳回失败，请重试",
            type: "warning",
            duration: 1000,
          })
        })
      })
      this.dialogVisible = false
    },
    passOrder() {
      this.debounce = true
      this.passLoading = true
      changeOrder({
        id: this.selectedOrder.id,
        target: "pass"
      }).then(() => {
        this.reload()
        this.passLoading = false
        this.dialogVisible = false
        this.$message({
          message: "已成功审核订单",
          type: "success",
          duration: 1000,
        })
      }).catch(err => {
        if(err.response.status === 500 || err.response.status === 400) {
          this.$message({
            message: "出现NCE错误，请联系管理员",
            type: "error",
            duration: 1000,
          })
        } else if(err === "permission denied"){
          this.$message({
            message: "该订单信息有变更，请重试",
            type: "warning",
            duration: 1000,
          })
          this.reload()
        } else {
          this.$message({
            message: "审核失败，请重试",
            type: "warning",
            duration: 1000,
          })
        }
        this.passLoading = false
        this.dialogVisible = false
      })
    },
    fetchData() {
      this.orderList = []
      let promise = new Promise((resolve, reject) => {
        let fetchStatus
        if(this.status === "first")
          fetchStatus = ["waiting"]
        else if(this.status === "second")
          fetchStatus = ["pass", "on", "suspend", "pause", "retrieve", "canceled"]
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
      }
    }
  }
</style>
