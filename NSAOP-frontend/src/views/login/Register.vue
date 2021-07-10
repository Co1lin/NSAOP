<template>
  <div class="container">
    <el-row type="flex">
      <el-col>
        <el-card>
          <div class="header">
            <h3>注册</h3>
          </div>
          <div class="content">
            <el-form
              ref="registerForm"
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
              <el-form-item prop="repeatPassword">
                <el-input
                  v-model="credential.repeatPassword"
                  placeholder="repeat password"
                  suffix-icon="el-icon-lock"
                  type="password"
                  @copy.native.capture.prevent="falsy"
                  @cut.native.capture.prevent="falsy"
                />
              </el-form-item>
              <el-form-item prop="company">
                <el-input
                  v-model="credential.company"
                  placeholder="your company name"
                  suffix-icon="el-icon-office-building"
                />
              </el-form-item>
              <el-form-item prop="email">
                <el-input
                  v-model="credential.email"
                  placeholder="your email"
                  suffix-icon="el-icon-message"
                />
              </el-form-item>
              <el-form-item prop="tel">
                <el-input
                  v-model="credential.tel"
                  placeholder="your telephone number"
                  suffix-icon="el-icon-phone-outline"
                />
              </el-form-item>
              <el-form-item v-if="role !== 'customer'" prop="code">
                <el-input
                  v-model="credential.code"
                  placeholder="invitation code"
                  suffix-icon="el-icon-check"
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
import {usernameCheck, register} from "@/network/user";
import {encrypt} from "@/common/hash";
import {lengthValidator, validateTel, validateEmail} from "@/common/validators";
import {reCAPTCHA} from "@/common/reCAPTCHA";
import Attribution from "@/components/Attribution";

export default {
  name: "Register",
  components: {Attribution},
  data() {
    return {
      credential: {
        username: "",
        password: "",
        repeatPassword: "",
        company: "",
        email: "",
        tel: "",
        code: "",
      },
      rules: {
        username: [
          {
            validator: this.validateUsername,
            trigger: "blur"
          },
        ],
        password: [
          {
            validator: lengthValidator(),
            required: true,
          }
        ],
        repeatPassword: [
          {
            validator: this.validatePassword,
            trigger: "blur"
          }
        ],
        email: [
          {
            validator: validateEmail,
            type: "email",
            trigger: "blur"
          }
        ],
        tel: [
          {
            validator: validateTel,
            trigger: "blur",
          }
        ],
        company: [
          {
            validator: lengthValidator(1, 30),
            trigger: "blur",
          }
        ],
        code: [
          {
            validator: lengthValidator(0, 100),
            trigger: "blur",
          }
        ]
      },
      loading: false,
    }
  },
  computed: {
    role() {
      return this.$route.params.role
    },
    md5Cred() {
      return {
        username: this.credential.username,
        password: encrypt(this.credential.username, this.credential.password),
        company: this.credential.company,
        email: this.credential.email,
        phone: this.credential.tel,
        role: this.role,
        code: this.credential.code
      }
    }
  },
  beforeMount() {
    if (this.role !== 'customer' && this.role !== 'operator' && this.role !== 'engineer') {
      this.$router.push('/404')
    }
    this.loading = false
  },
  methods: {
    falsy() {
      return false
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
    validateUsername(rule, value, callback) {
      if (value === "") {
        callback(new Error("请输入用户名"))
      } else if (value.length < 5 || value.length > 20) {
        callback(new Error("用户名的长度应在5-20个字符之间"))
      } else {
        usernameCheck({username: value} ).then(() => callback()).catch(() => {
          callback(new Error("该用户名已被占用"))})
      }
    },
    submitForm() {
      this.loading = true
      reCAPTCHA("register").then((token) => {
        let data = this.md5Cred
        data['g_recaptcha_response'] = token
        this.$refs["registerForm"].validate(valid => {
          if (valid) {
            register(data).then(() => {
              this.$message({
                message: "注册成功，请登录",
                type: "success",
                duration: 1000,
              })
              this.$router.push("/login")
              this.loading = false
            }).catch(err => {
              this.$message({
                message: err,
                type: "error"
              })
              this.loading = false
            })
          } else {
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
      this.$refs["registerForm"].resetFields()
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
