<template>
  <el-dropdown @command="handleCommand">
    <span class="avatar-dropdown">
      <img class="user-avatar" :src="avatar">
      <div class="user-name">
        {{ role }} {{ username }}
        <i class="el-icon-arrow-down el-icon--right" />
      </div>
    </span>
    <el-dropdown-menu v-slot: dropdown>
      <el-dropdown-item command="logout">
        登出
      </el-dropdown-item>
      <el-dropdown-item command="personalCenter" divided>
        用户信息
      </el-dropdown-item>
    </el-dropdown-menu>
  </el-dropdown>
</template>

<script>
import { MessageBox } from "element-ui";
import { LOGOUT } from "@/common/store";

export default {
  name: 'Avatar',
  computed: {
    avatar() {
      let username = this.username
      let hashCode = 0
      for(let i = 0; i < username.length; i++){
        hashCode = 37 * hashCode + username.charCodeAt(i)
      }
      hashCode = hashCode % 21 + 1
      return require("@/assets/img/user_image/bottts (" + hashCode + ").png")
    },
    username() {
      return this.$store.state.user.username
    },
    role() {
      if(this.$store.state.user.role === 'customer'){
        return '客户'
      }
      else if(this.$store.state.user.role === 'operator'){
        return '运营工程师'
      }
      else if(this.$store.state.user.role === 'engineer'){
        return '网络工程师'
      }
      return ''
    }
  },
  methods: {
    handleCommand(command) {
      switch (command) {
      case 'logout':
        this.logout()
        break
      case 'personalCenter':
        this.personalCenter()
        break
      }
    },
    personalCenter() {
      this.$router.push('/profile')
    },
    logout() {
      MessageBox.confirm('您确定要退出账号 ' + this.username + ' 吗?', null || '温馨提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        closeOnClickModal: false,
        distinguishCancelAndClose: true,
        type: 'message',
      }).then(async () => {
        await this.$store.dispatch(LOGOUT)
        await this.$router.push('/login')
        sessionStorage.setItem('currentTab', '')
      }).catch(() => {
      })
    },
  },
}
</script>

<style lang="scss" scoped>
  .avatar-dropdown {
    display: flex;
    align-content: center;
    align-items: center;
    justify-content: center;
    justify-items: center;
    height: 50px;
    padding: 0;

    .user-avatar {
      width: 40px;
      height: 40px;
      border-radius: 50%;
    }

    .user-name {
      position: relative;
      margin-left: 5px;
    }
  }
</style>
