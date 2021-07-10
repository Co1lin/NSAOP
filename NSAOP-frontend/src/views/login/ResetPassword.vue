<template>
  <div class="container">
    <el-row type="flex">
      <el-col>
        <el-card>
          <div class="header">
            <h3>重置密码</h3>
          </div>
          <div class="content">
            <el-form
              ref="resetForm"
              :model="credential"
              :rules="rules"
              @keyup.enter.native="submitForm"
            >
              <el-form-item prop="password">
                <el-input
                  v-model="credential.password"
                  placeholder="请输入新密码"
                  type="password"
                  @copy.native.capture.prevent="falsy"
                  @cut.native.capture.prevent="falsy"
                />
              </el-form-item>
              <el-form-item prop="confirm">
                <el-input
                  v-model="credential.confirm"
                  placeholder="请确认密码"
                  type="password"
                  @copy.native.capture.prevent="falsy"
                  @cut.native.capture.prevent="falsy"
                />
              </el-form-item>
              <el-form-item>
                <el-button id="forget-button" @click="toForgetPassword">
                  重新发送
                </el-button>
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
import {resetPassword} from "@/network/user";
import {encrypt} from "@/common/hash";
import {lengthValidator} from "@/common/validators";
import Attribution from "@/components/Attribution";

export default {
  name: "ResetPassword",
  components: {Attribution},
  data() {
    return {
      credential: {
        username: "",
        password: "",
        confirm: ""
      },
      rules: {
        password: [
          {
            validator: lengthValidator(),
            trigger: "blur",
            required: true,
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
      token: '',
      loading: false,
    }
  },
  computed: {
    cred() {
      return {
        username: this.credential.username,
        password: encrypt(this.credential.username, this.credential.password),
        token: this.token,
      }
    }
  },
  created() {
    this.credential.username = this.$route.query.username
    this.token = this.$route.query.token
  },
  methods: {
    submitForm() {
      this.loading = true
      this.$refs["resetForm"].validate(valid => {
        if (valid) {
          resetPassword(this.cred).then(() => {
            this.$message({
              message: "密码重置成功",
              type: "success",
              duration: 1000,
            })
            this.$router.push("/login")
            this.loading = false
          }).catch(() => {
            this.$message({
              message: "出于安全性考虑，链接不合法或已失效，请向邮箱重新发送密码重置请求",
              type: "error"
            })
            this.loading = false
          })
        } else {
          this.loading = false
          return false
        }
      })
    },
    resetForm() {
      this.$refs["resetForm"].resetFields()
    },
    validatePassword(rule, value, callback) {
      if (value === "") {
        callback(new Error("请再次输入密码"))
      } else if (value !== this.credential.password) {
        callback(new Error("两次密码不一致"))
      } else {
        callback()
      }
    },
    toForgetPassword() {
      this.$router.push('/login/forget')
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

  a {
    position: fixed;
    font-size: 1px;
    bottom: 0;
    right: 0;
    text-decoration: none;
    color: #bbbbbb;
  }

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
            min-width: 100px;
            width: 15vw;
            max-width: 150px;
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
