<template>
  <div>
    <layout>
      <template #main-content>
        <div v-if="device === 'mobile'">
          <div class="mobile-container">
            <el-row>
              <el-card class="user-info">
                <img :src="avatar" class="icon-container">
                <p> {{ name }} </p>
                <el-table
                  :data="infoData"
                  :cell-style="cellStyle"
                  :show-header="false"
                >
                  <el-table-column
                    prop="header"
                  />
                  <el-table-column
                    prop="data"
                  />
                </el-table>
              </el-card>
            </el-row>
            <el-row>
              <el-card class="info-change">
                <div class="information-container">
                  <change-info
                    ref="change"
                    @changed="infoChanged"
                  />
                </div>
              </el-card>
            </el-row>
            <el-row>
              <el-card class="info-change">
                <div class="information-container">
                  <change-password
                    @changed="passwordChanged"
                  />
                </div>
              </el-card>
            </el-row>
          </div>
        </div>
        <div v-else>
          <div class="pc-container">
            <el-row>
              <el-col :span="8">
                <el-card class="user-info">
                  <img :src="avatar" class="icon-container">
                  <p> {{ name }} </p>
                  <el-table
                    :data="infoData"
                    :cell-style="cellStyle"
                    :show-header="false"
                  >
                    <el-table-column
                      prop="header"
                    />
                    <el-table-column
                      prop="data"
                    />
                  </el-table>
                </el-card>
              </el-col>
              <el-col :span="8">
                <el-card class="user-info">
                  <div class="information-container">
                    <change-info
                      ref="change"
                      @changed="infoChanged"
                    />
                  </div>
                </el-card>
              </el-col>
              <el-col :span="8">
                <el-card class="user-info">
                  <div class="information-container">
                    <change-password
                      @changed="passwordChanged"
                    />
                  </div>
                </el-card>
              </el-col>
            </el-row>
          </div>
        </div>
      </template>
    </layout>
  </div>
</template>

<script>
import Layout from "@/components/Layout";
import {authTest} from "@/network/user";
import {SET_COMPANY, SET_TEL, SET_EMAIL, SET_USERNAME} from "@/common/store";
import {mapGetters} from "vuex";
import ChangeInfo from "@/views/personalCenter/ChangeInfo";
import ChangePassword from "@/views/personalCenter/ChangePassword";

export default {
  name: "PersonalCenter",
  components: {
    ChangeInfo,
    ChangePassword,
    Layout,
  },
  data() {
    return {
      active: "0",
    }
  },
  computed: {
    ...mapGetters({
      layout: 'settings/layout',
      tabsBar: 'settings/tabsBar',
      collapse: 'settings/collapse',
      header: 'settings/header',
      device: 'settings/device',
    }),
    token() {
      return this.$store.state.user.token
    },
    name() {
      return this.$store.state.user.username
    },
    tel() {
      return this.$store.state.user.tel
    },
    email() {
      return this.$store.state.user.email
    },
    company() {
      return this.$store.state.user.company
    },
    avatar() {
      let username = this.name
      let hashCode = 0
      for(let i = 0; i < username.length; i++){
        hashCode = 37 * hashCode + username.charCodeAt(i)
      }
      hashCode = hashCode % 21 + 1
      return require("@/assets/img/user_image/bottts (" + hashCode + ").png")
    },
    infoData() {
      return [{
        header: "联系电话",
        data: this.tel,
      }, {
        header: "邮箱",
        data: this.email,
      }, {
        header: "公司名称",
        data: this.company,
      }]
    },
  },
  mounted() {
    if (this.token !== undefined && this.token !== '') {
      authTest().then(res => {
        this.$store.commit(SET_USERNAME, res.username)
        this.$store.commit(SET_EMAIL, res.email)
        this.$store.commit(SET_TEL,  res.phone)
        this.$store.commit(SET_COMPANY, res.company)
      })
    }
  },
  methods: {
    cellStyle({column}) {
      let cellStyle
      if(column.property === "header")
        cellStyle = "font-weight: 700; textAlign: center"
      else
        cellStyle = "textAlign: center"
      return cellStyle
    },
    loadInfo(tab) {
      if(tab.index === '1'){
        this.$refs.change.loadInfo()
        this.active = "1"
      }
    },
    infoChanged() {
      setTimeout(() => {
        this.active = "0"
        authTest().then(res => {
          this.$store.commit(SET_USERNAME, res.username)
          this.$store.commit(SET_EMAIL, res.email)
          this.$store.commit(SET_TEL,  res.phone)
          this.$store.commit(SET_COMPANY, res.company)
        })
      }, 400)
    },
    passwordChanged() {
      this.active = "2"
      setTimeout(() => {
        this.active = "0"
      }, 800)
    }
  }
}
</script>

<style lang="scss" scoped>
.pc-container {
  padding: 20px;
  font-size: large;
  color: #2f3447;
  line-height: 50px;

  /deep/
  .user-info {
    text-align: center;
    margin: 10px;
    height: $base-app-column2-height;
    min-height: 500px;
    border-radius: 10px;

    .el-card__body {
      position: relative;
      top: 50%;
      transform: translateY(-50%);
    }
  }

  .information-container {
    //margin-left: 50px;
    //margin-right: 50px;
  }

  .icon-container {
    height: 100px;
    border-radius: 100px;
  }
}

.mobile-container {
  padding: 20px;
  color: #2f3447;
  line-height: 50px;

  /deep/
  .user-info {
    text-align: center;
    margin: 10px;
    height: $base-app-column2-height;
    min-height: 500px;
    border-radius: 10px;

    .el-card__body {
      position: relative;
      top: 50%;
      transform: translateY(-50%);
    }
  }

  /deep/
  .info-change {
    text-align: center;
    margin: 10px;
    height: $base-app-column1-height;
    min-height: 400px;
    border-radius: 10px;

    .el-card__body {
      position: relative;
      top: 50%;
      transform: translateY(-40%);
    }
  }

  .information-container {
    //margin-left: 50px;
    //margin-right: 50px;
  }

  .icon-container {
    height: 100px;
    border-radius: 100px;
  }

  .information-container {
    margin-left: 10px;
    margin-right: 10px;
  }


}


</style>
