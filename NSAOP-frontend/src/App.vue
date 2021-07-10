<template>
  <div>
    <router-view />
  </div>
</template>

<script>
export default {
  name: "App",
}
</script>

<style>
body{
  font-family: Avenir,Helvetica,"PingFang SC","Hiragino Sans GB","黑体", Arial,sans-serif;;
  margin:0;
  padding:0;
  background-color: #e9e9e9;
}

.el-message-box {
  max-width: 80vw !important;
}

.el-dialog {
  max-width: 80vw !important;
}

</style>
