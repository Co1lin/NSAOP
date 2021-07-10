<template>
  <el-form
    ref="passwordForm"
    :model="password"
    :rules="rules"
    label-width="80px"
  >
    <el-form-item prop="origin" label="原密码">
      <el-input
        v-model="password.origin"
        placeholder="origin password"
        suffix-icon="el-icon-lock"
        type="password"
      />
    </el-form-item>
    <el-form-item prop="new" label="新密码">
      <el-input
        v-model="password.new"
        placeholder="new password"
        suffix-icon="el-icon-lock"
        type="password"
      />
    </el-form-item>
    <el-form-item prop="confirm" label="确认密码">
      <el-input
        v-model="password.confirm"
        placeholder="confirm password"
        suffix-icon="el-icon-lock"
        type="password"
      />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" :disabled="invalid" @click="submitForm">
        确认
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import {changeInfo} from "@/network/user";
import {checkPassword} from "@/network/user";
import {lengthValidator} from "@/common/validators";
import {encrypt} from "@/common/hash";


export default {
  name: "ChangePassword",
  data() {
    return {
      password: {
        origin: '',
        new: '',
        confirm: '',
      },
      rules: {
        origin: [
          {
            validator: this.correctPassword,
            trigger: "blur",
            required: true
          }
        ],
        new: [
          {
            validator: lengthValidator(),
            trigger: "blur",
            required: true
          }
        ],
        confirm: [
          {
            validator: this.validatePassword,
            trigger: "blur",
            required: true
          }
        ]
      },
    }
  },
  computed: {
    invalid() {
      let invalid = false
      if(this.password.origin === '' || this.password.new === '' || this.password.confirm === '')
        return true
      if(this.$refs["passwordForm"] !== undefined) {
        this.$refs["passwordForm"].validate(valid => {
          if (!valid) {
            invalid = true
          }
        })
      }
      return invalid
    },
    username() {
      return this.$store.state.user.username
    }
  },
  methods: {
    submitForm() {
      this.$refs["passwordForm"].validate(valid => {
        if (valid) {
          changeInfo({
            password: encrypt(this.username, this.password.new),
            old_password: encrypt(this.username, this.password.origin)
          }).then(() => {
            this.$message({
              message: "修改密码成功",
              type: "success",
              duration: 1000,
            })
            this.password.origin = ''
            this.password.new = ''
            this.password.confirm = ''
            this.$emit("changed")
          }).catch(err => {
            this.$message({
              message: err,
              type: "error"
            })
          })
        } else {
          return false
        }
      })
    },
    correctPassword(rule, value, callback) {
      let message = {
        password: encrypt(this.username, value),
      }
      checkPassword(message).then(() => {
        callback()
      }).catch(() => {
        callback(new Error("原密码错误"))
      })
    },
    validatePassword(rule, value, callback) {
      if (value === "") {
        callback(new Error("请再次输入密码"))
      } else if (value !== this.password.new) {
        callback(new Error("两次密码不一致"))
      } else {
        callback()
      }
    },
  }
}
</script>

<style scoped>

</style>
