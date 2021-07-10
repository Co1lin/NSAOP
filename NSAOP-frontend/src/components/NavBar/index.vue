<template>
  <div class="nav-bar-container">
    <el-row :gutter="15">
      <el-col
        :xs="4"
        :sm="12"
        :md="12"
        :lg="12"
        :xl="12"
      >
        <div class="left-panel">
          <i
            :class="collapse ? 'el-icon-s-unfold' : 'el-icon-s-fold'"
            :title="collapse ? '展开' : '收起'"
            class="fold-unfold"
            @click="handleCollapse"
          />
          <!--          <vab-breadcrumb class="hidden-xs-only" />-->
        </div>
      </el-col>
      <el-col
        :xs="20"
        :sm="12"
        :md="12"
        :lg="12"
        :xl="12"
      >
        <div class="right-panel">
          <div class="mail-container" v-if="role === 'operator'" @click="sendMail">
            <img src="@/assets/img/email.png" class="icon">
          </div>
          <avatar v-if="userStatus" />
          <div v-if="!userStatus">
            <div>
              <el-button
                type="primary"
                class="button"
                @click="login"
              >
                Login
              </el-button>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex'
import Avatar from '@/components/Avatar'
import {sendEmail} from "@/network/user";

export default {
  name: 'NavBar',
  components: {
    Avatar
  },
  data() {
    return {
      pulse: false,
    }
  },
  computed: {
    ...mapGetters({
      collapse: 'settings/collapse',
    }),
    userStatus(){
      if(this.$store.state.user.token === "" || this.$store.state.user.token === undefined) {
        return false
      } else {
        return true
      }
    },
    role() {
      return this.$store.state.user.role
    }
  },
  methods: {
    ...mapActions({
      changeCollapse: 'settings/changeCollapse',
    }),
    handleCollapse() {
      this.changeCollapse()
    },
    login() {
      this.$router.push('/login')
    },
    sendMail() {
      this.$confirm("是否向您管理的所有用户发送订单统计邮件？", "温馨提示"
      ).then(() => {
        sendEmail().then(() => this.$message({
          message: "邮件发送成功",
          type: "success",
          duration: 1000,
        })).catch(() => this.$message({
          message: "邮件发送失败",
          type: "error",
          duration: 1000,
        }))
      }).catch(() => {})
    }
  },
}
</script>

<style lang="scss" scoped>
.nav-bar-container {
  position: relative;
  height: $base-nav-bar-height;
  padding-right: $base-padding;
  padding-left: $base-padding;
  overflow: hidden;
  user-select: none;
  background: $base-color-white;
  box-shadow: $base-box-shadow;

  .left-panel {
    display: flex;
    align-items: center;
    justify-items: center;
    height: $base-nav-bar-height;

    .fold-unfold {
      color: $base-color-gray;
      cursor: pointer;
    }

    ::v-deep {
      .breadcrumb-container {
        margin-left: 10px;
      }
    }
  }

  .right-panel {
    display: flex;
    align-content: center;
    align-items: center;
    justify-content: flex-end;
    height: $base-nav-bar-height;

    ::v-deep {

      .mail-container {
        width: 50px;
        text-align: left;
        padding-left: 50px;
        display: flex;
        align-items: center;
      }

      .icon {
        padding-top: 3px;
        cursor: pointer;
        height: 25px;
      }

      svg {
        width: 1em;
        height: 1em;
        margin-right: 15px;
        font-size: $base-font-size-small;
        color: $base-color-gray;
        cursor: pointer;
        fill: $base-color-gray;
      }

      button {
        svg {
          margin-right: 0;
          color: $base-color-white;
          cursor: pointer;
          fill: $base-color-white;
        }
      }

      .el-badge {
        margin-right: 15px;
      }
    }
  }
}
</style>
