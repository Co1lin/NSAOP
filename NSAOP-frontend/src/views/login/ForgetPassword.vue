<template>
  <div class="container">
    <el-row type="flex">
      <el-col>
        <el-card>
          <div class="header">
            <h3>忘记密码</h3>
          </div>
          <div class="content">
            <p>
              请输入您的用户名 <br>
              我们会向您的邮箱发送重置密码的邮件
            </p>
            <el-form
              ref="resetForm"
              :model="credential"
              :rules="rules"
              @keyup.enter.native="submitForm"
            >
              <el-form-item prop="username">
                <el-input
                  v-model="credential.username"
                  placeholder="username"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="loading" @click="submitForm">
                  提交
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <attribution/>
  </div>
</template>

<script>
import {usernameCheck, requestResetPassword} from "@/network/user";
import {reCAPTCHA} from "@/common/reCAPTCHA";
import Attribution from "@/components/Attribution";
import {deepcopy} from "@/common/utils";

export default {
  name: "ResetPassword",
  components: {Attribution},
  data() {
    return {
      credential: {
        username: "",
      },
      rules: {
        username: [
          {
            validator: this.validateUsername,
            trigger: "blur"
          },
        ],
      },
      loading: false,
    }
  },
  computed: {
    cred() {
      return {
        username: this.credential.username
      }
    }
  },
  beforeMount() {
  },
  methods: {
    validateUsername(rule, value, callback) {
      if (value === "") {
        callback(new Error("请输入用户名"))
      } else if (value.length < 5 || value.length > 20) {
        callback(new Error("用户名的长度应在5-20个字符之间"))
      } else {
        usernameCheck({username: value} ).then(() => {
          callback(new Error("该用户名不存在"))
        }).catch(err => {
          // console.log(err.response)
          if (err.response.status === 409) {
            callback()
          } else {
            callback(new Error("网络错误"))
          }
        })
      }
    },
    submitForm() {
      this.loading = true
      reCAPTCHA("request_passwd_reset").then(token => {
        let data = deepcopy(this.cred)
        data['g-recaptcha-response'] = token
        this.$refs["resetForm"].validate(valid => {
          if (valid) {
            requestResetPassword(data).then(() => {
              this.$message({
                message: "邮件发送成功",
                type: "success",
                duration: 1000,
              })
              this.loading = false
            }).catch(() => {
              this.$message({
                message: "邮件发送失败，请联系管理员",
                type: "error"
              })
              this.loading = false
            })
          } else {
            this.$message({
              message: "请完善表单",
              type: "error"
            })
            this.loading = false
            return false
          }
        })
      }).catch(() => {
        this.$message({
          message: "reCAPTCHA连接错误，请检查网络后重试",
          type: "error"
        })
        this.loading = false
      })
    },
    resetForm() {
      this.$refs["resetForm"].resetFields()
    }
  }
}
</script>

<style scoped lang="scss">
.container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  height: 100vh;
  background: url('~@/assets/img/background2.png') no-repeat fixed;
  background-size: 100% 100%;

  .el-card {
    margin-top: 30px;
    width: 50vw;
    min-width: 260px;
    max-width: 400px;
    background-color: rgb(255, 255, 255);
    color: #000000;
    border-radius: 5px;
    .header {
      line-height: 0px;
      border-bottom: 1px solid #ffffff;
      padding-bottom: 10px;
    }

    .content {
      p{
        margin: 0 20px 0 20px;
        text-align: center;
      }
      .el-form {
        display: flex;
        flex-direction: column;
        align-items: center;

        .el-form-item {
          width: 100%;

          .el-input {
            max-width: 310px;
            min-width: 200px;
            width: 40vw;
          }

          .el-button {
            max-width: 310px;
            min-width: 200px;
            width: 40vw;
          }

          /deep/ .el-form-item__error {
            padding-left: 25px;
          }
        }

        .el-form-item:first-child {
          margin-top: 20px;
        }

      }
    }
  }
}
</style>
