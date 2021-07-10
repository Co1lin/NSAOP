<template>
  <div class="container">
    <el-dialog
      title="创建新地址"
      :show-close="false"
      :visible.sync="dialogVisible"
      :before-close="handleClose"
      :close-on-click-modal="false"
    >
      <el-form
        ref="createLocation"
        :model="info"
        :rules="rules"
      >
        <el-form-item prop="comment" label="名称">
          <el-input
            v-model="info.comment"
            placeholder="address"
            suffix-icon="el-icon-info"
          />
        </el-form-item>
        <el-form-item prop="address" label="详细地址">
          <el-input
            v-model="info.address"
            placeholder="detailed address"
            suffix-icon="el-icon-location"
          />
        </el-form-item>
        <el-form-item prop="contact" label="联系人 ">
          <el-input
            v-model="info.contact"
            placeholder="contact"
            suffix-icon="el-icon-user"
          />
        </el-form-item>
        <el-form-item prop="phone" label="联系电话">
          <el-input
            v-model="info.phone"
            placeholder="phone"
            suffix-icon="el-icon-phone"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeDialog">
          取消
        </el-button>
        <el-button type="primary" :loading="loading" @click="submitForm">
          提交
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {createLocation} from "@/network/location";
import {lengthValidator, shortMessageValidator, validateTel} from "@/common/validators";
import {reCAPTCHA} from "@/common/reCAPTCHA";

export default {
  name: "CreateOrder",
  props: {
    visible: {
      type: Boolean,
      default: true
    }
  },
  data() {
    return {
      info: {
        address: "",
        comment: "",
        contact: "",
        phone: "",
      },
      rules: {
        address: [
          {
            validator: lengthValidator(),
            required: true,
            trigger: "blur",
          }
        ],
        comment: [
          {
            validator: shortMessageValidator,
            required: true,
            trigger: "blur",
          }
        ],
        contact: [
          {
            validator: shortMessageValidator,
            required: true,
            trigger: "blur",
          }
        ],
        phone: [
          {
            validator: validateTel,
            required: true,
            trigger: "blur",
          }
        ]
      },
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
  methods: {
    sendForm() {
      this.loading = true
      const data = {
        comment: this.info.comment,
        address: this.info.address,
        contact: this.info.contact,
        phone: this.info.phone
      }
      reCAPTCHA("create_location").then(token => {
        data['g_recaptcha_response'] = token
        createLocation(data).then(() => {
          this.$message({
            message: "提交成功",
            type: "success",
            duration: 2000,
          })
          setTimeout(() => this.closeDialog(), 500)
          this.$emit('reload')
        }).catch(err => {
          if(err === 'reCAPTCHA failed')
            this.$message({
              message: "人机验证失败，请重试",
              type: "error"
            })
          else
            this.$message({
              message: err,
              type: "error"
            })
          // console.log(err.response);
          setTimeout(() => this.closeDialog(), 500)
        })
      })
    },
    submitForm() {
      this.$refs["createLocation"].validate(valid => {
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
      this.loading = false
    },
    handleClose(done) {
      this.$confirm("是否关闭? 未提交的内容不会保存")
        .then(() => done()).catch(() => {})
    },
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

/deep/

.el-form-item {
  padding-left: 8%;
}

.el-input {
  float: right;
  width: 70%;
  //left: 50px;
}

.el-form {
  margin-right: 20px;
}

.el-form-item {
  padding-top: 15px;
}
</style>
