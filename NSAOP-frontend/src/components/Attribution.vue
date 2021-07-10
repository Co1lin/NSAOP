<template>
  <a :href="link">
    background image attribute to {{ author }}
  </a>
</template>

<script>
export default {
  name: "Attribution",
  props: {
    link: {
      type: String,
      default: "https://walli.shanga.co/image/view/?id=1097",
    },
    author: {
      type: String,
      default: "Walli",
    }
  }
}
</script>

<style scoped>
a {
  position: fixed;
  font-size: 1px;
  bottom: 0;
  right: 0;
  text-decoration: none;
  color: #bbbbbb;
}
</style>
