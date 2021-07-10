<template>
  <div class="table-container">
    <el-row :gutter="0" class="order-list-header">
      <el-col
        class="left-column"
        :xs="span"
        :sm="span"
        :md="span"
        :lg="span"
        :xl="span"
      >
        <div class="left-panel">
          <CreateOrder
            v-if="this.$store.state.user.role === 'customer'"
            :visible.sync="visible"
            @reload="reload"
          />
          <el-button
            v-if="this.$store.state.user.role === 'customer'"
            icon="el-icon-plus"
            type="primary"
            @click="showDialog"
          >
            添加新订单
          </el-button>
        </div>
      </el-col>
      <el-col
        class="right-column"
        :xs="span"
        :sm="span"
        :md="span"
        :lg="span"
        :xl="span"
      >
        <div class="right-panel">
          <el-form
            ref="forms"
            :model="queryForm"
            :inline="true"
            :rules="rules"
            class="right-panel"
            @submit.native.prevent
          >
            <el-form-item prop="title">
              <el-input v-model="queryForm.title" placeholder="订单名称/订单号" />
            </el-form-item>
            <el-form-item>
              <el-button
                icon="el-icon-search"
                type="primary"
                native-type="submit"
                @click="handleQuery"
                :disabled="invalid"
              >
                查询
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
    </el-row>
  </div>
</template>


<script>
import CreateOrder from "@/components/Order/CreateOrder";
import {authTest} from "@/network/user";
import {SET_TOKEN, SET_TARGET} from "@/common/store";
import {revokeToken} from "@/common/auth";
import {lengthValidator} from "@/common/validators"

export default {
  name: 'ControlBar',
  components: {
    CreateOrder
  },
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'gray',
        deleted: 'danger',
      }
      return statusMap[status]
    },
  },
  props: {
    span: {
      type: Number,
      default: 12,
    },
  },
  data() {
    return {
      isMounted: false,
      queryForm: {
        title: '',
      },
      rules: {
        title: [
          {
            validator: lengthValidator(0, 40),
            trigger: "change"
          }
        ]
      },
      visible: false,
    }
  },
  computed: {
    invalid() {
      if(!this.isMounted)
        return
      let invalid = false
      this.$refs["forms"].validate(valid => {
        if(!valid) {
          invalid = true
        }
      })
      return invalid
    },
    height() {
      return this.$baseTableHeight()
    },
  },
  beforeDestroy() {},
  mounted() {
    this.isMounted = true
  },
  methods: {
    showDialog() {
      authTest().catch(() => {
        // console.log(err.response);
        this.$message({
          message: "请重新登录",
          type: "error",
          duration: 1000,
        })
        revokeToken()
        this.$store.commit(SET_TOKEN, "")
        this.$router.push('/login')
      })
      this.visible = true
    },
    reload() {
      this.$emit('reload')
    },
    handleQuery() {
      this.$store.dispatch(SET_TARGET, {data: this.queryForm.title})
      this.$emit("reload")
    }
  },
}
</script>

<style lang="scss" scoped>
@mixin panel {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-start;
}

/deep/ .el-dialog__wrapper {
  left: 20%;
}

.order-list-header {
  margin-bottom: 10px;

  .left-column {
    float: left;
  }
  .right-column {
    float: right;
  }

  ::v-deep {
    .left-panel {
      @include panel;

      > .el-button,
      .el-form-item {
        margin: 5px;
      }
    }

    .right-panel {
      @include panel;

      justify-content: flex-end;

      .el-form-item {
        margin: 5px;
      }
    }
  }
}
</style>

