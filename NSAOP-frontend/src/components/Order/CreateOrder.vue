<template>
  <div class="container">
    <el-dialog
      title="创建新订单"
      :show-close="false"
      :visible.sync="dialogVisible"
      :before-close="handleClose"
      :close-on-click-modal="false"
    >
      <el-form
        ref="createForm"
        :model="info"
        :rules="rules"
      >
        <el-form-item prop="nickname" label="订单名称">
          <el-input
            v-model="info.nickname"
            suffix-icon="el-icon-office-building"
            @focus="getInputFocus($event)"
          />
        </el-form-item>
        <el-form-item prop="address" label="地址">
          <el-select v-model="info.address" placeholder="请选择">
            <el-option
              v-for="(item, index) in locations"
              :value="item.value"
              :label="item.name"
              :key="index">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item prop="paytype" label="支付方式">
          <el-radio v-model="info.paytype" label="month">
            按月支付
          </el-radio>
          <el-radio v-model="info.paytype" label="year">
            按年支付
          </el-radio>
        </el-form-item>
        <el-form-item prop="requirements" label="需求">
          <el-checkbox-group v-model="info.requirements">
            <el-checkbox label="private">
              私用
            </el-checkbox>
            <el-checkbox label="client">
              客户用
            </el-checkbox>
            <el-checkbox label="test">
              测试用
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item prop="detail" label="具体说明">
          <el-input
            v-model="info.detail"
            type="textarea"
            :rows="2"
          >
          </el-input>
        </el-form-item>
      </el-form>
      <div>
        <br>
        目前单价：1GB/元
      </div>
      <template #footer>
        <el-button @click="closeDialog">
          取消
        </el-button>
        <el-button
          :loading="loading"
          :disabled="debounce"
          @click="submitForm"
          type="primary"
        >
          提交
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {submitOrder} from "@/network/order";
import {getLocationByUser} from "@/network/location";
import {lengthValidator} from "@/common/validators";
import {reCAPTCHA} from "@/common/reCAPTCHA";

export default {
  name: "CreateOrder",
  props: {
    visible: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      info: {
        nickname: "新订单",
        detail: "",
        address: "",
        paytype: "",
        requirements: []
      },
      rules: {
        nickname: [
          {
            validator: this.nicknameValidator,
            required: true,
            trigger: "blur",
          }
        ],
        detail: [
          {
            validator: lengthValidator(1, 100),
            required: true,
            trigger: "blur",
          }
        ],
        address: [
          {
            required: true,
          }
        ],
        paytype: [
          {
            required: true,
            message: "请选择支付方式"
          }
        ],
        requirements: [
          {
            required: true,
            message: "请选择需求"
          }
        ]
      },
      locations: [],
      debounce: false,
      loading: false,
    }
  },
  computed: {
    dialogVisible: {
      get() {
        return this.visible
      },
      set(newValue) {
        this.$emit('update:visible', newValue)
      }
    }
  },
  async created() {
    const res = await this.getLocations(0, 10)
    if (res.locations !== null) {
      for (let location of res.locations) {
        const temp = {value: location.id, name: location.comment}
        this.locations.push(temp)
      }
    }
  },
  methods: {
    getInputFocus(event) {
      event.currentTarget.select();
    },
    getLocations(offset = 0, limit = 10) {
      const data = {
        offset,
        limit,
      }
      return getLocationByUser(data).then(res => res.data).catch(() => {})
    },
    nicknameValidator(rule, value, callback) {
      if (value === "") {
        callback(new Error("备注名不能为空"))
      } else if (value.length > 10) {
        callback(new Error("不能超过10个字符"))
      } else {
        callback()
      }
    },
    sendForm() {
      this.loading = true
      this.debounce = true
      reCAPTCHA("create_service").then(token => {
        let data = {
          comment: this.info.nickname,
          detail: this.info.detail,
          location: this.info.address,
          paytype: this.info.paytype,
          device: this.info.device,
          require: this.encodeRequirements(),
        }
        data['g_recaptcha_response'] = token
        submitOrder(data).then(() => {
          this.$message({
            message: "提交成功",
            type: "success",
            duration: 2000,
          })
          setTimeout(() => this.closeDialog(), 500)
          this.$emit('reload')
        }).catch(err => {
          this.$message({
            message: err,
            type: "error",
            duration: 2000,
          })
          // console.log(err.response);
          setTimeout(() => this.closeDialog(), 500)
        })
      })
    },
    submitForm() {
      this.$refs["createForm"].validate(valid => {
        if (valid) {
          this.sendForm()
        } else {
          this.$message({
            message: "请完善表单",
            type: "error",
            duration: 1000,
          })
          return false
        }
      })
    },
    closeDialog() {
      this.$emit('update:visible', false)
      this.debounce = false
      this.loading = false
    },
    handleClose(done) {
      this.$confirm("是否关闭? 未提交的内容不会保存")
        .then(() => done()).catch(() => {})
    },
    encodeRequirements() {
      let res = 0
      if (this.info.requirements.indexOf("private") !== -1) {
        res = res + 4
      }
      if (this.info.requirements.indexOf("client") !== -1) {
        res = res + 2
      }
      if (this.info.requirements.indexOf("test") !== -1) {
        res = res + 1
      }
      return res
    }
  }
}
</script>

<style scoped lang="scss">

.el-dialog__wrapper {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0 !important;
  overflow: auto;
  margin: 0;
  width: 100%;

  ::v-deep {
    .el-dialog {
      max-width: 500px !important;
      margin-top: 15vh !important;
      width: 80% !important;
    }
    .el-dialog__header {
      padding-top: 5%;
      text-align: center;
    }
  }
}

.el-form-item {
  width: 90%;
  padding-left: 8%;
}

.el-input {
  width: 70%;
}

.el-select{
  padding-left: 27px;
  width: 70%;
}


.el-form-item {
  padding-top: 15px;
}
</style>
