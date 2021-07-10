<template>
  <div class="error-container">
    <el-card>
      <h2> 404 Page Not Found </h2>
      <h3> 当前页面不存在... </h3>
      <p>
        请检查您输入的网址是否正确，或点击下面的按钮返回首页。
      </p>
      <p id="backLink" class="el-icon-arrow-left" @click="returnHome">
        <span>Back to our site</span>
      </p>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'Page404',
  data() {
    return {
    }
  },
  methods: {
    returnHome() {
      this.$router.push('/')
    }
  },
}
</script>

<style lang="scss" scoped>
.error-container {
  background-color: rgb(14, 30, 37);
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  h3 {
    color: #606266;
  }
  #backLink {
    color: #007067;
    font-weight: 600;
    line-height: 10px;
    span:hover {
      border-bottom: 2px solid darkgoldenrod;
    }
  }
  #backLink:hover {
    cursor: pointer;
  }
  .el-card {
    height: 250px;
  }
}
</style>
