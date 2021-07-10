<template>
  <div>
    <layout>
      <template #main-content>
        <div class="location-list-container">
          <control-bar @reload="reload" />
          <div class="table-container">
            <el-table
              ref="tableSort"
              v-loading="listLoading"
              :data="locationList"
              :element-loading-text="elementLoadingText"
              :height="height"
              :row-style="{fontFamily: 'STHeiti'}"
              :header-cell-style="{backgroundColor: '#F0F0F0', color: '#191919'}"
            >
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
                label="名称"
                prop="comment"
                min-width="120"
                align="center"
              />
              <el-table-column
                show-overflow-tooltip
                label="详细地址"
                min-width="150"
                prop="address"
                align="center"
              />
              <el-table-column
                show-overflow-tooltip
                label="联系人"
                prop="contact"
                align="center"
              />
              <el-table-column
                show-overflow-tooltip
                label="联系电话"
                min-width="150"
                prop="phone"
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
                    @click.native.stop="handleDelete(row,)"
                  >
                    删除该地址
                  </el-button>
                  <!--                  <el-button-->
                  <!--                    v-if="orderStatus === 'waiting'"-->
                  <!--                    type="text"-->
                  <!--                    @click.native.stop="handleOperate(row,)"-->
                  <!--                  >-->
                  <!--                    设为默认-->
                  <!--                  </el-button>-->
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
        </div>
      </template>
    </layout>
  </div>
</template>

<script>
import Layout from "@/components/Layout"
import {loadingText} from "@/config/setting.config"
import ControlBar from "@/views/locationList/ControlBar";
import {SET_LOCATION_TARGET} from "@/common/store";
import {getLocationByUser, deleteLocationById} from "@/network/location";

export default {
  name: "LocationList",
  components: {
    Layout,
    ControlBar
  },
  data() {
    return {
      listLoading: false,
      locationList: [],
      elementLoadingText: loadingText,
      displayPageNO: 1,
      queryForm: {
        pageNo: 1,
        pageSize: 10,
        title: '',
      },
      layout: 'total, prev, pager, next, jumper',
      total: 0,
    }
  },
  computed: {
    height() {
      return this.$baseTableHeight() + 55
    },
  },
  mounted() {
    this.reload()
  },
  methods: {
    handleCurrentChange(val) {
      this.reload(val)
    },
    handleDelete(row) {
      this.$confirm("是否删除地址？").then(() => {
        let id = row.id
        deleteLocationById({
          id: id
        }).then(() => {
          this.reload()
          this.$message({
            message: "成功删除地址",
            type: "success",
            duration: 1000
          })
        }).catch(err => {
          this.$message({
            message: err,
            type: "warning",
            duration: 1000
          })
        })
      })
    },
    fetchData() {
      // this.locationList = []
      let promise = new Promise((resolve, reject) => {
        let target = this.$store.state.locationQuery.target
        let data = {
          offset: (this.queryForm.pageNo - 1) * this.queryForm.pageSize,
          limit: this.queryForm.pageSize,
        }
        if(target !== '') {
          this.$store.dispatch(SET_LOCATION_TARGET, {data: ''})
          data.search = target
        }
        getLocationByUser(data).then(res => {
          this.total = res.data.count
          if(res.data.locations !== null) {
            res.data.locations.forEach(location => {
              this.locationList.push(location)
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
      this.locationList = []
      let promise = new Promise((resolve, reject) => {
        this.listLoading = true
        this.queryForm.pageNo = page
        setTimeout(() => {
          setTimeout(() => {
            this.fetchData().then(() => {
              this.listLoading = false
            }).then(() => resolve()).catch(() => {
              this.$message({
                message: "加载地址列表失败，请刷新页面重试",
                type: "warning",
                duration: 1000,
              })
              this.listLoading = false
              this.locationList = []
              reject()
            })
          },200)
        }, 200)
      })
      promise.then(() => this.displayPageNO = page)
      return promise
    }
  },
}
</script>

<style scoped>
.location-list-container{
  margin-left: 20px;
  margin-right: 20px;
  padding: 20px;
  text-align: center;
}
</style>
