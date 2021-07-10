<template>
  <div class="container">
    <el-dialog
      title="详情"
      :visible.sync="dialogVisible"
      :before-close="handleClose"
    >
      <div class="">
        <slot name="attr"></slot>
      </div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: "DisplayDetail",
  props: {
    visible: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    dialogVisible: {
      get() {
        return this.visible
      },
      set(newValue) {
        this.$emit('update:visible', newValue)
      }
    },
  },
  methods: {
    closeDialog() {
      this.$emit('update:visible', false)
    },
    handleClose(done) {
      done()
    },
  }
}
</script>

<style scoped>

/deep/ .el-dialog {
  width: 500px;
}

</style>
