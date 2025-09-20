<template>
  <div class="layout">
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-button
            class="mobile-menu-btn"
            type="text"
            icon="el-icon-menu"
            @click="toggleMobileMenu"
          />
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
        <el-aside 
          :width="sidebarWidth" 
          class="sidebar"
          :class="{ 'mobile-sidebar': isMobile, 'mobile-sidebar-open': mobileMenuOpen }"
        >
          <div class="sidebar-overlay" @click="closeMobileMenu" v-if="isMobile && mobileMenuOpen"></div>
          <div class="sidebar-content">
            <el-menu
              :default-active="$route.path"
              router
              class="sidebar-menu"
              :collapse="false"
            >
              <el-menu-item index="/dashboard" @click="handleMenuClick">
                <i class="el-icon-s-home" />
                <span>首页</span>
              </el-menu-item>
              <el-menu-item index="/profile" @click="handleMenuClick">
                <i class="el-icon-user" />
                <span>个人主页</span>
              </el-menu-item>
              <el-menu-item index="/ratings" @click="handleMenuClick">
                <i class="el-icon-star-on" />
                <span>评价管理</span>
              </el-menu-item>
              <el-menu-item v-if="isAdmin" index="/admin" @click="handleMenuClick">
                <i class="el-icon-setting" />
                <span>系统管理</span>
              </el-menu-item>
            </el-menu>
          </div>
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
  data() {
    return {
      mobileMenuOpen: false,
      isMobile: false
    }
  },
  computed: {
    ...mapGetters('auth', ['currentUser', 'isAdmin']),
    
    sidebarWidth() {
      if (this.isMobile) {
        return this.mobileMenuOpen ? '200px' : '0px'
      }
      return '200px'
    }
  },
  mounted() {
    this.checkMobile()
    window.addEventListener('resize', this.checkMobile)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.checkMobile)
  },
  methods: {
    ...mapActions('auth', ['logout']),
    
    checkMobile() {
      this.isMobile = window.innerWidth <= 768
      if (!this.isMobile) {
        this.mobileMenuOpen = false
      }
    },
    
    toggleMobileMenu() {
      this.mobileMenuOpen = !this.mobileMenuOpen
    },
    
    closeMobileMenu() {
      this.mobileMenuOpen = false
    },
    
    handleMenuClick() {
      if (this.isMobile) {
        this.mobileMenuOpen = false
      }
    },
    
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
  position: relative;
  z-index: 1001;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left h1 {
  margin: 0;
  color: #333;
  font-size: 20px;
  font-weight: 600;
}

.mobile-menu-btn {
  display: none;
  font-size: 18px;
  padding: 8px;
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
  transition: all 0.3s ease;
  position: relative;
  z-index: 1000;
}

.sidebar-content {
  height: 100%;
  overflow-y: auto;
}

.sidebar-menu {
  border: none;
  background: transparent;
}

.sidebar-overlay {
  display: none;
}

.main-content {
  background: #f5f5f5;
  padding: 20px;
  transition: all 0.3s ease;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .header {
    padding: 0 15px;
  }
  
  .header-left h1 {
    font-size: 18px;
  }
  
  .mobile-menu-btn {
    display: block;
  }
  
  .username {
    display: none;
  }
  
  .user-info {
    padding: 8px;
  }
  
  .mobile-sidebar {
    position: fixed;
    top: 60px;
    left: 0;
    height: calc(100vh - 60px);
    width: 200px !important;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
    z-index: 1000;
    box-shadow: 2px 0 8px rgba(0, 0, 0, 0.15);
  }
  
  .mobile-sidebar-open {
    transform: translateX(0);
  }
  
  .mobile-sidebar-open .sidebar-overlay {
    display: block;
    position: fixed;
    top: 0;
    left: 200px;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 999;
  }
  
  .main-content {
    padding: 15px;
    width: 100%;
  }
}

@media (max-width: 480px) {
  .header {
    padding: 0 10px;
  }
  
  .header-left h1 {
    font-size: 16px;
  }
  
  .main-content {
    padding: 10px;
  }
}
</style>
