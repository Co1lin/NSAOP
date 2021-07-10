<template>
  <el-form
    ref="infoForm"
    :model="user"
    :rules="rules"
    label-width="80px"
  >
    <el-form-item prop="phone" label="联系电话">
      <el-input
        v-model="user.phone"
        placeholder="telephone number"
        suffix-icon="el-icon-phone-outline"
      />
    </el-form-item>
    <el-form-item prop="email" label="邮箱">
      <el-input
        v-model="user.email"
        placeholder="email"
        suffix-icon="el-icon-message"
      />
    </el-form-item>
    <el-form-item prop="company" label="公司名称">
      <el-input
        v-model="user.company"
        placeholder="company name"
        suffix-icon="el-icon-office-building"
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
import {lengthValidator, validateEmail, validateTel} from "@/common/validators"

export default {
  name: "ChangeInfo",
  data() {
    return {
      isMounted: false,
      user: {
        company: this.$store.state.user.company,
        email: this.$store.state.user.email,
        phone: this.$store.state.user.tel,
      },
      rules: {
        company: [
          {
            validator: lengthValidator(1, 30),
            trigger: "blur"
          }
        ],
        email: [
          {
            validator: validateEmail,
            trigger: "blur"
          }
        ],
        phone: [
          {
            validator: validateTel,
            trigger: "blur",
          }
        ]
      },
    }
  },
  mounted(){
    this.isMounted = true;
  },
  watch: {
    company(newValue) {
      this.user.company = newValue
    },
    email(newValue) {
      this.user.email = newValue
    },
    phone(newValue) {
      this.user.phone = newValue
    },

  },
  computed: {
    company() {
      return this.$store.state.user.company
    },
    email() {
      return this.$store.state.user.email
    },
    phone() {
      return this.$store.state.user.tel
    },
    invalid() {
      if(this.user.company === this.$store.state.user.company
        && this.user.email === this.$store.state.user.email
        && this.user.phone === this.$store.state.user.tel
      ){
        return true
      }
      if(this.user.company === ''){
        return true
      }
      if(!this.isMounted)
        return
      let invalid = false
      this.$refs["infoForm"].validate(valid => {
        if (!valid) {
          invalid = true
        }
      })
      return invalid
    }
  },
  methods: {
    submitForm() {
      this.$refs["infoForm"].validate(valid => {
        if (valid) {
          changeInfo(this.user).then(() => {
            this.$message({
              message: "修改信息成功",
              type: "success",
              duration: 1000,
            })
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
  }
}
</script>

<style scoped>

</style>
