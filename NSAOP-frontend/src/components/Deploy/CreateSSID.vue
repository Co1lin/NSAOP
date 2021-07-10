<template>
  <div class="container">
    <el-dialog
      title="创建新SSID"
      :show-close="false"
      :visible.sync="dialogVisible"
      :before-close="handleClose"
      :close-on-click-modal="false"
    >
      <el-form
        ref="ssidForm"
        :model="info"
        :rules="rules"
      >
        <el-form-item prop="name" label="设备ID">
          <el-input v-model="info.name" placeholder="请输入英文字符" />
        </el-form-item>
        <el-form-item prop="mode" label="网络连接方式">
          <el-select v-model="info.mode" placeholder="请选择">
            <el-option
              v-for="(item, index) in modes"
              :value="item.value"
              :label="item.name"
              :key="index">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item prop="hide" label="是否隐藏SSID">
          <el-checkbox v-model="info.hide">
            hide_enable
          </el-checkbox>
        </el-form-item>
        <el-form-item prop="separation" label="是否用户隔离">
          <el-checkbox v-model="info.separation">
            user_separation
          </el-checkbox>
        </el-form-item>
        <el-form-item prop="radios" label="射频类型">
          <el-slider
            v-model="info.radios"
            :step="1"
            show-stops
            show-input
            :min="1"
            :max="7"
          />
        </el-form-item>
        <el-form-item prop="userNum" label="最大用户数量">
          <el-slider
            v-model="info.userNum"
            :step="32"
            show-input
            :min="1"
            :max="512"
          />
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
import {ssidNameValid} from '@/common/regex'

export default {
  name: "CreateSsid",
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
        enable: true,
        mode: "",
        hide: false,
        radios: 7,
        userNum: 100,
        separation: true,
      },
      rules: {
        name: [
          {
            validator: this.nameValidator,
            required: true,
            trigger: "blur",
          }
        ],
        mode: [
          {
            required: true,
            message: "请选择一个类型"
          }
        ],
        hide: [
          {
            required: true,
            message: "请选择是否允许隐藏"
          }
        ],
        radios: [
          {
            required: true,
            message: "请选择一个基站数量"
          }
        ],
        userNum: [
          {
            required: true,
            message: "请选择最大用户数量"
          }
        ],
        separation: [
          {
            required: true,
            message: "请选择是否允许分离"
          }
        ],
      },
      modes: [
        { name: 'bridge', value: 'bridge'},
        { name: 'nat', value: 'nat'},
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
      } else if (value.length > 32) {
        callback(new Error("不能超过32个字符"))
      } else if (!ssidNameValid(value)) {
        callback(new Error("不能包含特殊符号&、=、?、 #、%、+。不以空格开头或结尾；不以引号开头"))
      } else {
        callback()
      }
    },
    submitForm() {
      this.loading = true
      this.$refs["ssidForm"].validate(valid => {
        if (valid) {
          this.$emit('send', this.info)
          this.$message({
            message: "成功添加",
            type: "success",
            duration: 1000,
          })
          setTimeout(() => this.closeDialog(), 250)
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
.el-slider {
  margin-top: 30px;
}
</style>
