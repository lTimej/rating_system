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
          <el-button
            class="back-btn"
            type="text"
            icon="el-icon-arrow-left"
            @click="goBack"
            v-if="canGoBack"
            title="返回上一页"
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
        <!-- 移动端遮罩层 -->
        <div 
          class="mobile-overlay" 
          v-if="isMobile && mobileMenuOpen" 
          @click="closeMobileMenu"
        ></div>
        
        <el-aside 
          :width="sidebarWidth" 
          class="sidebar"
          :class="{ 'mobile-sidebar': isMobile, 'mobile-sidebar-open': mobileMenuOpen }"
          @click.stop
        >
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
              <el-menu-item index="/articles" @click="handleMenuClick">
                <i class="el-icon-document" />
                <span>文章</span>
              </el-menu-item>
              <!-- <el-menu-item index="/ratings" @click="handleMenuClick">
                <i class="el-icon-star-on" />
                <span>评价管理</span>
              </el-menu-item> -->
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
      isMobile: false,
      canGoBack: false
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
    this.updateCanGoBack()
    window.addEventListener('resize', this.checkMobile)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.checkMobile)
  },
  watch: {
    '$route': {
      handler() {
        this.updateCanGoBack()
      },
      immediate: true
    }
  },
  methods: {
    ...mapActions('auth', ['logout']),
    
    checkMobile() {
      this.isMobile = window.innerWidth <= 768
      if (!this.isMobile) {
        this.mobileMenuOpen = false
      }
    },

    updateCanGoBack() {
      // 检查浏览器历史记录是否可以返回
      const currentPath = this.$route.path
      const currentName = this.$route.name
      
      // 排除的页面：登录页面、注册页面和首页
      const excludedPages = ['/login', '/register', '/dashboard', '/']
      const excludedNames = ['Login', 'Register', 'Dashboard']
      
      const isExcludedPage = excludedPages.includes(currentPath) || excludedNames.includes(currentName)
      
      // 如果不是排除的页面，则显示返回按钮
      this.canGoBack = !isExcludedPage
      
      // 临时调试信息
      if (process.env.NODE_ENV === 'development') {
        console.log('updateCanGoBack:', {
          currentPath,
          currentName,
          isExcludedPage,
          canGoBack: this.canGoBack
        })
      }
    },

    goBack() {
      const currentPath = this.$route.path
      const currentName = this.$route.name
      
      // 根据当前页面智能返回
      if (currentName === 'ArticleDetail') {
        // 文章详情页返回文章列表
        this.$router.push('/articles')
      } else if (currentName === 'ArticleCreate' || currentName === 'ArticleEdit') {
        // 文章编辑页返回文章列表
        this.$router.push('/articles')
      } else if (currentPath.startsWith('/articles/')) {
        // 其他文章相关页面返回文章列表
        this.$router.push('/articles')
      } else {
        // 其他页面使用浏览器返回
        if (window.history.length > 1) {
          this.$router.go(-1)
        } else {
          // 如果没有历史记录，返回首页
          this.$router.push('/dashboard')
        }
      }
    },
    toggleMobileMenu() {
      this.mobileMenuOpen = !this.mobileMenuOpen
    },
    
    closeMobileMenu() {
      if (this.isMobile && this.mobileMenuOpen) {
        this.mobileMenuOpen = false
      }
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
      // 移动端使用更友好的确认方式
      if (this.isMobile) {
        this.$msgbox({
          title: '退出登录',
          message: '确定要退出登录吗？',
          showCancelButton: true,
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
          center: true,
          customClass: 'mobile-confirm-dialog'
        }).then(() => {
          this.logout()
          this.$message.success('已退出登录')
          this.$router.push('/login')
        }).catch(() => {})
      } else {
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
  flex: 1;
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

.back-btn {
  font-size: 18px;
  padding: 8px;
  color: #409EFF;
  transition: all 0.3s ease;
  border-radius: 6px;
  min-width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-btn:hover {
  color: #66b1ff;
  background-color: rgba(64, 158, 255, 0.1);
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

/* 移动端遮罩层 */
.mobile-overlay {
  position: fixed;
  top: 60px;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
  backdrop-filter: blur(2px);
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
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
  
  .header-left {
    gap: 15px;
  }
  
  .header-left h1 {
    font-size: 18px;
    flex: 1;
  }
  
  .mobile-menu-btn {
    display: block;
  }
  
  .back-btn {
    font-size: 20px;
    padding: 12px;
    min-width: 48px;
    height: 48px;
    background-color: rgba(64, 158, 255, 0.1);
    border: 1px solid rgba(64, 158, 255, 0.3);
    border-radius: 12px;
    color: #409EFF;
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
    -webkit-tap-highlight-color: transparent;
  }
  
  .back-btn:active {
    transform: scale(0.95);
    background-color: rgba(64, 158, 255, 0.2);
    box-shadow: 0 1px 4px rgba(64, 158, 255, 0.4);
  }
  
  .back-btn:hover {
    background-color: rgba(64, 158, 255, 0.15);
    box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
    transform: translateY(-1px);
  }
  
  /* 返回按钮进入动画 */
  .back-btn {
    animation: slideInLeft 0.3s ease-out;
  }
  
  @keyframes slideInLeft {
    from {
      opacity: 0;
      transform: translateX(-20px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
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
  
  
  .main-content {
    padding: 15px;
    width: 100%;
  }
}

@media (max-width: 480px) {
  .header {
    padding: 0 10px;
  }
  
  .header-left {
    gap: 12px;
  }
  
  .header-left h1 {
    font-size: 16px;
    flex: 1;
  }
  
  .back-btn {
    font-size: 18px;
    padding: 10px;
    min-width: 44px;
    height: 44px;
    background-color: rgba(64, 158, 255, 0.12);
    border: 1px solid rgba(64, 158, 255, 0.3);
    border-radius: 10px;
    box-shadow: 0 2px 6px rgba(64, 158, 255, 0.25);
  }
  
  .main-content {
    padding: 10px;
  }
}

/* 移动端确认对话框样式 */
:deep(.mobile-confirm-dialog) {
  width: 90% !important;
  max-width: 400px !important;
  margin: 0 auto !important;
}

:deep(.mobile-confirm-dialog .el-message-box__header) {
  padding: 20px 20px 10px !important;
}

:deep(.mobile-confirm-dialog .el-message-box__content) {
  padding: 10px 20px !important;
  font-size: 16px !important;
}

:deep(.mobile-confirm-dialog .el-message-box__btns) {
  padding: 10px 20px 20px !important;
}

:deep(.mobile-confirm-dialog .el-button) {
  min-width: 80px !important;
  height: 40px !important;
  font-size: 16px !important;
}

/* 移动端下拉菜单优化 */
@media (max-width: 768px) {
  :deep(.el-dropdown-menu) {
    min-width: 120px !important;
  }
  
  :deep(.el-dropdown-menu .el-dropdown-menu__item) {
    padding: 12px 16px !important;
    font-size: 16px !important;
    line-height: 1.5 !important;
  }
}
</style>
