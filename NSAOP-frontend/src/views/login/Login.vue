<template>
  <div :class="showParallel ? 'parallel-container' : 'vertical-container'">
    <div class="image" v-if="showParallel">
      <img src="~@/assets/img/white_logo.png" alt="" class="img">
      <div>
        <h3>NSAOP</h3>
        <div>网 络 服 务 接 入 平 台</div>
      </div>
    </div>
    <div>
      <el-card>
        <div class="header">
          <h3>登录</h3>
        </div>
        <div class="content">
          <el-form
            ref="loginForm"
            label-position="top"
            :model="credential"
            :rules="rules"
            @keyup.enter.native="submitForm"
          >
            <el-form-item prop="username">
              <el-input
                v-model="credential.username"
                placeholder="username"
                suffix-icon="el-icon-user"
              />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="credential.password"
                placeholder="password"
                suffix-icon="el-icon-lock"
                type="password"
                @copy.native.capture.prevent="falsy"
                @cut.native.capture.prevent="falsy"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                id="login-button"
                type="primary"
                :loading="loginLoading"
                @click="submitForm"
              >
                登录
              </el-button>
              <el-button id="register-button" type="primary" @click="toRegister">
                注册
              </el-button>
              <el-button id="forget-button" type="warning" @click="toForgetPassword">
                忘记密码
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-card>
    </div>
    <attribution/>
  </div>
</template>

<script>
import {LOGIN} from "@/common/store"
import {encrypt} from "@/common/hash";
import {lengthValidator} from "@/common/validators";
import {reCAPTCHA} from "@/common/reCAPTCHA";
import Attribution from "@/components/Attribution";

export default {
  name: "Login",
  components: {Attribution},
  data() {
    return {
      screenWidth: window.innerWidth,
      credential: {
        username: "",
        password: ""
      },
      rules: {
        username: [
          {
            validator: lengthValidator(5, 20),
            required: true,
            trigger: "blur"
          },
        ],
        password: [
          {
            required: true,
            validator: lengthValidator()
          }
        ]
      },
      loginLoading: false,
    }
  },
  computed: {
    md5Cred() {
      return {
        username: this.credential.username,
        password: encrypt(this.credential.username, this.credential.password),
      }
    },
    showParallel() {
      if(this.screenWidth >= 992)
        return true
      else
        return false
    }
  },
  mounted() {
    const that = this
    window.onresize = () => {
      return (() => {
        window.screenWidth = document.body.clientWidth
        that.screenWidth = window.screenWidth
      })()
    }
  },
  beforeCreate() {
    if (this.$store.getters.token !== undefined && this.$store.getters.token !== '') {
      this.$router.push('/home')
    }
  },
  created() {
    this.loginLoading = false
  },
  methods: {
    handleWidthChange() {
      this.screenWidth = window.innerWidth
    },
    falsy() {
      return false
    },
    submitForm() {
      this.loginLoading = true
      reCAPTCHA("login").then(token => {
        this.$refs["loginForm"].validate(valid => {
          if (valid) {
            let data = this.md5Cred
            data['g_recaptcha_response'] = token
            this.$store.dispatch(LOGIN, data).then(() => {
              this.$message({
                message: "登录成功",
                type: "success",
                duration: 1000,
              })
              setTimeout(() => this.$router.push("/home"), 300)
            }).catch(err => {
              // console.log(err);
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
              this.resetForm()
            })
            this.loginLoading = false
          } else {
            this.loginLoading = false
            return false
          }
        })
      }).catch(() => {
        this.$message({
          message: "reCAPTCHA连接错误，请检查网络后重试",
          type: "error"
        })
      })
    },
    resetForm() {
      this.$refs["loginForm"].resetFields()
    },
    toRegister() {
      this.$router.push('/register/customer')
    },
    toForgetPassword() {
      this.$router.push('/login/forget')
    }
  }
}
</script>

<style scoped lang="scss">

.parallel-container{
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  text-align: center;
  height: 100vh;
  background: url('~@/assets/img/background2.png') no-repeat fixed;
  background-size: 100% 100%;

  /deep/
  .el-card {
    border-radius: 0 5px 5px 0;
    box-shadow: none;
    width: 100%;
    background-color: rgb(255, 255, 255);
    color: #000000;
  }
  .el-input {
    width: 310px;
  }

  #login-button {
    width: 98px;
  }

  #register-button {
    width: 98px;
  }

  #forget-button {
    width: 98px;
  }
}

.vertical-container{
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  text-align: center;
  height: 100vh;
  background: url('~@/assets/img/background2.png') no-repeat fixed;
  background-size: 100% 100%;

  /deep/
  .el-card {
    border-radius: 5px;
    box-shadow: none;
    width: 100%;
    background-color: rgb(255, 255, 255);
    color: #000000;
  }

  .el-input {
    width: 40vw;
    min-width: 200px;
  }

  #login-button {
    width: 12vw;
    min-width: 75px;
  }

  #register-button {
    width: 12vw;
    min-width: 75px;
  }

  #forget-button {
    width: 12vw;
    min-width: 100px;
  }
}

.el-card__body {
  height: 400px;
  .div {
    margin-top: 10px;
  }
}

.header {
  line-height: 30px;
  border-bottom: 1px solid #ffffff;
}

.el-form {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.image {
  border-radius: 5px 0 0 5px;
  font-size: 20px;
  display: flex;
  flex-direction: column;
  height: 336.4px;
  width: 400px;
  background-color: rgba(37,77,110,0.9);
  align-items: center;
  color: #ffffff;
  img {
    margin-top: 15%;
    width: 25%;
  }
}




</style>
