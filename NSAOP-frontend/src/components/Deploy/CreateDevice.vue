<template>
  <div class="container">
    <el-dialog
      title="创建新设备"
      :show-close="false"
      :visible.sync="dialogVisible"
      :before-close="handleClose"
      :close-on-click-modal="false"
      class="dialog"
    >
      <el-form
        ref="deviceForm"
        :model="info"
        :rules="rules"
      >
        <el-form-item prop="name" label="设备ID">
          <el-input v-model="info.name" placeholder="请输入英文字符" />
        </el-form-item>
        <el-form-item prop="type" label="类别">
          <el-select v-model="info.type" placeholder="请选择">
            <el-option
              v-for="(item, index) in models"
              :value="item.value"
              :label="item.name"
              :key="index">
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeDialog">
          取消
        </el-button>
        <el-button type="primary" :loading="loading" @click="submitForm">
          提交
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {consistsOfOnlyEnglish} from '@/common/regex'

export default {
  name: "CreateDevice",
  props: {
    visible: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      info: {
        name: "",
        type: "",
      },
      rules: {
        name: [
          {
            validator: this.nameValidator,
            required: true,
            trigger: "blur",
          }
        ],
        type: [
          {
            required: true,
            message: "请选择一个类型"
          }
        ],
      },
      models: [
        { name: 'AP4050DN', value: 'AP4050DN'}
      ],
      loading: false,
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
    nameValidator(rule, value, callback) {
      if (value === "") {
        callback(new Error("设备名不能为空"))
      } else if (value.length > 64) {
        callback(new Error("不能超过64个字符"))
      } else if (!consistsOfOnlyEnglish(value)) {
        callback(new Error("仅能填写英文字符"))
      } else {
        callback()
      }
    },
    submitForm() {
      this.loading = true
      this.$refs["deviceForm"].validate(async valid => {
        if (valid) {
          await this.$emit('send', this.info)
          await this.$message({
            message: "添加成功",
            type: "success",
            duration: 1000,
          })
          await setTimeout(() => this.closeDialog(), 250)
        } else {
          this.$message({
            message: "请完善表单",
            type: "error",
            duration: 1000,
          })
          this.loading = false
          return false
        }
      })
    },
    closeDialog() {
      this.$emit('update:visible', false)
      this.loading = false
    },
    handleClose(done) {
      this.$confirm("是否关闭? 未提交的内容不会保存")
        .then(() => done()).catch(() => {})
    },
  }
}
</script>

<style lang="scss" scoped>
.el-dialog__wrapper {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0 !important;
  overflow: auto;
  margin: 0;
  width: 100%;

  ::v-deep {
    .el-dialog {
      max-width: 500px !important;
      margin-top: 15vh !important;
      width: 80% !important;
    }
    .el-dialog__header {
      padding-top: 5%;
      text-align: center;
    }
  }
}
</style>
