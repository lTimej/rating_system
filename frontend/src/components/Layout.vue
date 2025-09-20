<template>
  <div class="layout">
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <h1>评价系统</h1>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" :src="currentUser.avatar" icon="el-icon-user-solid" />
              <span class="username">{{ currentUser.name || currentUser.username }}</span>
              <i class="el-icon-arrow-down" />
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item command="profile">个人资料</el-dropdown-item>
              <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-container>
        <el-aside width="200px" class="sidebar">
          <el-menu
            :default-active="$route.path"
            router
            class="sidebar-menu"
          >
            <el-menu-item index="/dashboard">
              <i class="el-icon-s-home" />
              <span>首页</span>
            </el-menu-item>
            <el-menu-item index="/profile">
              <i class="el-icon-user" />
              <span>个人主页</span>
            </el-menu-item>
            <el-menu-item index="/ratings">
              <i class="el-icon-star-on" />
              <span>评价管理</span>
            </el-menu-item>
            <el-menu-item v-if="isAdmin" index="/admin">
              <i class="el-icon-setting" />
              <span>系统管理</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        
        <el-main class="main-content">
          <slot />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  name: 'Layout',
  computed: {
    ...mapGetters('auth', ['currentUser', 'isAdmin'])
  },
  methods: {
    ...mapActions('auth', ['logout']),
    
    handleCommand(command) {
      switch (command) {
        case 'profile':
          this.$router.push('/profile')
          break
        case 'logout':
          this.handleLogout()
          break
      }
    },
    
    handleLogout() {
      this.$confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.logout()
        this.$message.success('已退出登录')
        this.$router.push('/login')
      }).catch(() => {})
    }
  }
}
</script>

<style scoped>
.layout {
  height: 100vh;
}

.header {
  background: white;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left h1 {
  margin: 0;
  color: #333;
  font-size: 20px;
  font-weight: 600;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.username {
  margin: 0 8px;
  color: #333;
  font-size: 14px;
}

.sidebar {
  background: #f8f9fa;
  border-right: 1px solid #e6e6e6;
}

.sidebar-menu {
  border: none;
  background: transparent;
}

.main-content {
  background: #f5f5f5;
  padding: 20px;
}
</style>
