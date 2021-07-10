<template>
  <el-scrollbar class="side-bar-container" :class="{ 'is-collapse': collapse, }">
    <el-menu
      :background-color="variables['menu-background']"
      :text-color="variables['menu-color']"
      :active-text-color="variables['menu-color-active']"
      :default-active="activeMenu"
      :collapse="collapse"
      :collapse-transition="false"
      :default-openeds="defaultOpens"
      :unique-opened="uniqueOpened"
      :router="true"
      mode="vertical"
    >
      <div :class="'logo-container-' + layout">
        <router-link to="/">
          <img
            class="logo-display"
            src="@/assets/img/logo.png"
          >
          <span
            class="title"
            :class="{ 'hidden-xs-only': layout === 'horizontal', }"
            :title="title"
          >
            {{ title }}
          </span>
        </router-link>
      </div>
      <el-menu-item index="/home">
        <i class="el-icon-s-home" />
        <template #title>
          主页
        </template>
      </el-menu-item>
      <el-menu-item
        v-if="login === 'login'"
        index="/orders"
      >
        <i class="el-icon-menu" />
        <template #title>
          订单管理
        </template>
      </el-menu-item>
      <el-menu-item
        v-if="this.$store.state.user.role === 'customer'"
        index="/locations"
      >
        <i class="el-icon-location" />
        <template #title>
          地址管理
        </template>
      </el-menu-item>
    </el-menu>
  </el-scrollbar>
</template>
<script>
import variables from "@/styles/variables.scss";
import { mapGetters } from "vuex";
import { defaultOopeneds, uniqueOpened } from "@/config";

export default {
  name: "VabSideBar",
  data() {
    return {
      uniqueOpened,
      title: this.$baseTitle,
    };
  },
  computed: {
    ...mapGetters({
      collapse: "settings/collapse",
      layout: 'settings/layout',
    }),
    login() {
      if(this.$store.state.user.token === "" || this.$store.state.user.token === undefined) {
        return "logout"
      } else {
        return "login"
      }
    },
    defaultOpens() {
      // eslint-disable-next-line no-empty
      if (this.collapse) {

      }
      return defaultOopeneds;
    },
    activeMenu() {
      const route = this.$route;
      const { meta, path } = route;
      if (meta.activeMenu) {
        return meta.activeMenu;
      }
      return path;
    },
    variables() {
      return variables;
    },
  },
};
</script>
<style lang="scss" scoped>
@mixin active {
  &:hover {
    color: $base-color-white;
    background-color: $base-menu-background-active !important;
  }

  &.is-active {
    color: $base-color-white;
    background-color: $base-menu-background-active !important;
  }
}
@mixin container {
  position: relative;
  height: $base-top-bar-height;
  overflow: hidden;
  line-height: $base-top-bar-height;
  background: $base-menu-background;
}

@mixin title {
  margin-left: 5px;
  display: inline-block;
  overflow: hidden;
  font-size: 18px;
  line-height: 55px;
  color: $base-title-color;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.logo-display {
  //height: $base-menu-item-height;
  display: inline-block;
  width: 28px;
  height: 28px;
  margin-right: 0;
  color: $base-title-color;
  vertical-align: middle;
}

.side-bar-container {
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  z-index: $base-z-index;
  width: $base-left-menu-width;
  height: 100vh;
  overflow: hidden;
  background: $base-menu-background;
  box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
  transition: width $base-transition-time;

  &.is-collapse {
    width: $base-left-menu-width-min;
    border-right: 0;

    ::v-deep {
      .el-menu {
        transition: width $base-transition-time;
      }

      .el-menu--collapse {
        border-right: 0;

        .el-submenu__icon-arrow {
          right: 10px;
          margin-top: -3px;
        }
      }
    }
  }

  ::v-deep {
    .el-scrollbar__wrap {
      overflow-x: hidden;
    }

    .el-menu {
      border: 0;

      .vab-fas-icon {
        padding-right: 3px;
        font-size: $base-font-size-default;
      }

      .vab-remix-icon {
        padding-right: 3px;
        font-size: $base-font-size-default + 2;
      }
    }

    .el-menu-item,
    .el-submenu__title {
      height: $base-menu-item-height;
      overflow: hidden;
      line-height: $base-menu-item-height;
      text-overflow: ellipsis;
      white-space: nowrap;
      vertical-align: middle;
    }

    .el-menu-item {
      @include active;
    }
  }
}

.logo-container-horizontal {
  @include container;

  .title {
    @include title;
  }
}

.logo-container-vertical {
  @include container;

  height: $base-logo-height;
  line-height: $base-logo-height;
  text-align: center;

  .title {
    @include title;

    max-width: calc(#{$base-left-menu-width} - 60px);
  }
}
</style>
