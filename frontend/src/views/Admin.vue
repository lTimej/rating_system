<template>
  <Layout>
    <div class="admin">
      <div class="page-header">
        <h2>系统管理</h2>
      </div>

      <!-- 统计卡片 -->
      <el-row :gutter="20" class="stats-section">
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-number">{{ stats.total_users || 0 }}</div>
              <div class="stat-label">总用户数</div>
            </div>
            <i class="el-icon-user stat-icon" />
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-number">{{ stats.total_students || 0 }}</div>
              <div class="stat-label">学生用户</div>
            </div>
            <i class="el-icon-user-solid stat-icon" />
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-number">{{ stats.total_experts || 0 }}</div>
              <div class="stat-label">专家用户</div>
            </div>
            <i class="el-icon-s-custom stat-icon" />
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-number">{{ stats.total_articles || 0 }}</div>
              <div class="stat-label">总文章数</div>
            </div>
            <i class="el-icon-document stat-icon" />
          </el-card>
        </el-col>
      </el-row>

      <el-tabs v-model="activeTab" class="admin-tabs">
        <el-tab-pane label="用户管理" name="users">
          <div class="tab-header">
            <div class="filters">
              <el-select v-model="userFilters.role" placeholder="筛选角色" clearable @change="fetchUsers">
                <el-option label="学生" value="student" />
                <el-option label="专家" value="expert" />
                <el-option label="管理员" value="admin" />
              </el-select>
            </div>
            <el-button type="primary" @click="showCreateUserDialog = true">
              <i class="el-icon-plus" /> 添加用户
            </el-button>
          </div>
          
          <el-table :data="users" v-loading="loadingUsers" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="username" label="用户名" />
            <el-table-column prop="name" label="姓名" />
            <el-table-column prop="email" label="邮箱" />
            <el-table-column prop="role" label="角色" width="100">
              <template slot-scope="scope">
                <el-tag :type="getRoleType(scope.row.role)">
                  {{ getRoleText(scope.row.role) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="is_active" label="状态" width="100">
              <template slot-scope="scope">
                <el-tag :type="scope.row.is_active ? 'success' : 'danger'">
                  {{ scope.row.is_active ? '活跃' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template slot-scope="scope">
                {{ formatTime(scope.row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="220">
              <template slot-scope="scope">
                <el-button size="mini" @click="editUser(scope.row)">编辑</el-button>
                <el-button
                  size="mini"
                  :type="scope.row.is_active ? 'warning' : 'success'"
                  @click="toggleUserStatus(scope.row)"
                >
                  {{ scope.row.is_active ? '禁用' : '启用' }}
                </el-button>
                <el-button
                  size="mini"
                  type="danger"
                  @click="deleteUser(scope.row.id)"
                  v-if="scope.row.id !== currentUser.id"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <div class="pagination">
            <el-pagination
              @current-change="handleUserPageChange"
              :current-page="userPagination.page"
              :page-size="userPagination.limit"
              :total="userPagination.total"
              layout="total, prev, pager, next"
            />
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="文章管理" name="articles">
          <div class="tab-header">
            <div class="filters">
              <el-select v-model="articleFilters.status" placeholder="筛选状态" clearable @change="fetchArticles">
                <el-option label="草稿" value="draft" />
                <el-option label="已发布" value="published" />
                <el-option label="已归档" value="archived" />
              </el-select>
            </div>
          </div>
          
          <el-table :data="articles" v-loading="loadingArticles" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="title" label="标题" show-overflow-tooltip />
            <el-table-column label="作者" width="120">
              <template slot-scope="scope">
                <div class="user-info">
                  <el-avatar :size="24" :src="scope.row.author.avatar" icon="el-icon-user-solid" />
                  <span>{{ scope.row.author.name || scope.row.author.username }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template slot-scope="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="view_count" label="浏览量" width="100" />
            <el-table-column prop="like_count" label="点赞数" width="100" />
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template slot-scope="scope">
                {{ formatTime(scope.row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template slot-scope="scope">
                <el-button
                  size="mini"
                  type="primary"
                  @click="pushArticle(scope.row)"
                >
                  推送
                </el-button>
                <el-button
                  size="mini"
                  @click="viewArticle(scope.row.id)"
                >
                  查看
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <div class="pagination">
            <el-pagination
              @current-change="handleArticlePageChange"
              :current-page="articlePagination.page"
              :page-size="articlePagination.limit"
              :total="articlePagination.total"
              layout="total, prev, pager, next"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 创建用户对话框 -->
    <el-dialog
      title="添加用户"
      :visible.sync="showCreateUserDialog"
      width="500px"
    >
      <el-form :model="createUserForm" :rules="createUserRules" ref="createUserForm">
        <el-form-item label="用户名" prop="username" for="create-username">
          <el-input v-model="createUserForm.username" id="create-username" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email" for="create-email">
          <el-input v-model="createUserForm.email" id="create-email" />
        </el-form-item>
        
        <el-form-item label="姓名" prop="name" for="create-name">
          <el-input v-model="createUserForm.name" id="create-name" />
        </el-form-item>
        
        <el-form-item label="角色" prop="role" for="create-role">
          <el-select v-model="createUserForm.role" id="create-role" style="width: 100%">
            <el-option label="学生" value="student" />
            <el-option label="专家" value="expert" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="密码" prop="password" for="create-password">
          <el-input v-model="createUserForm.password" id="create-password" type="password" />
        </el-form-item>
        
        <el-form-item label="用户状态" for="create-active">
          <el-checkbox v-model="createUserForm.is_active" id="create-active">激活用户</el-checkbox>
        </el-form-item>
      </el-form>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="showCreateUserDialog = false">取消</el-button>
        <el-button type="primary" @click="createUser" :loading="creatingUser">创建</el-button>
      </div>
    </el-dialog>

    <!-- 编辑用户对话框 -->
    <el-dialog
      title="编辑用户"
      :visible.sync="showEditUserDialog"
      width="500px"
    >
      <el-form :model="editUserForm" :rules="editUserRules" ref="editUserForm">
        <el-form-item label="用户名" prop="username" for="edit-username">
          <el-input v-model="editUserForm.username" id="edit-username" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email" for="edit-email">
          <el-input v-model="editUserForm.email" id="edit-email" />
        </el-form-item>
        
        <el-form-item label="姓名" prop="name" for="edit-name">
          <el-input v-model="editUserForm.name" id="edit-name" />
        </el-form-item>
        
        <el-form-item label="角色" prop="role" for="edit-role">
          <el-select v-model="editUserForm.role" id="edit-role" style="width: 100%">
            <el-option label="学生" value="student" />
            <el-option label="专家" value="expert" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="用户状态" for="edit-active">
          <el-checkbox v-model="editUserForm.is_active" id="edit-active">激活用户</el-checkbox>
        </el-form-item>
      </el-form>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="showEditUserDialog = false">取消</el-button>
        <el-button type="primary" @click="updateUser" :loading="updatingUser">保存</el-button>
      </div>
    </el-dialog>

    <!-- 推送文章对话框 -->
    <el-dialog
      title="推送文章"
      :visible.sync="showPushDialog"
      width="600px"
    >
      <div class="push-dialog-content">
        <div class="article-info" v-if="selectedArticle">
          <h4>{{ selectedArticle.title }}</h4>
          <p class="article-summary">{{ selectedArticle.summary || '暂无摘要' }}</p>
          <div class="article-meta">
            <span>作者：{{ selectedArticle.author.name || selectedArticle.author.username }}</span>
            <span>状态：{{ getStatusText(selectedArticle.status) }}</span>
          </div>
        </div>

        <el-divider>选择推送用户</el-divider>

        <div class="user-selection">
          <div class="selection-header">
            <el-input
              v-model="userSearchKeyword"
              placeholder="搜索用户..."
              prefix-icon="el-icon-search"
              @input="searchUsers"
              style="width: 300px; margin-right: 10px;"
            />
            <el-button @click="selectAllUsers" size="small">全选</el-button>
            <el-button @click="clearSelection" size="small">清空</el-button>
          </div>

          <div class="user-list">
            <el-checkbox-group v-model="selectedUserIds">
              <div
                v-for="user in filteredUsers"
                :key="user.id"
                class="user-item"
                :class="{ 'already-pushed': pushedUserIds.includes(user.id) }"
              >
                <el-checkbox :label="user.id" :disabled="pushedUserIds.includes(user.id)">
                  <div class="user-info">
                    <el-avatar :size="24" :src="user.avatar" icon="el-icon-user-solid" />
                    <span class="user-name">{{ user.name || user.username }}</span>
                    <el-tag size="mini" :type="getRoleType(user.role)">
                      {{ getRoleText(user.role) }}
                    </el-tag>
                    <el-tag v-if="pushedUserIds.includes(user.id)" size="mini" type="info">
                      已推送
                    </el-tag>
                  </div>
                </el-checkbox>
              </div>
            </el-checkbox-group>
          </div>
        </div>
      </div>

      <div slot="footer" class="dialog-footer">
        <el-button @click="showPushDialog = false">取消</el-button>
        <el-button 
          type="primary" 
          @click="submitPush"
          :loading="pushSubmitting"
          :disabled="selectedUserIds.length === 0"
        >
          推送给 {{ selectedUserIds.length }} 个用户
        </el-button>
      </div>
    </el-dialog>
  </Layout>
</template>

<script>
import { mapGetters } from 'vuex'
import Layout from '@/components/Layout.vue'

export default {
  name: 'Admin',
  components: {
    Layout
  },
  data() {
    return {
      activeTab: 'users',
      loadingUsers: false,
      loadingArticles: false,
      creatingUser: false,
      updatingUser: false,
      pushSubmitting: false,
      
      stats: {},
      users: [],
      articles: [],
      selectedArticle: null,
      selectedUserIds: [],
      userSearchKeyword: '',
      showPushDialog: false,
      pushedUserIds: [], // 已推送的用户ID列表
      
      userFilters: {
        role: ''
      },
      
      articleFilters: {
        status: ''
      },
      
      userPagination: {
        page: 1,
        limit: 20,
        total: 0
      },
      
      articlePagination: {
        page: 1,
        limit: 20,
        total: 0
      },
      
      showCreateUserDialog: false,
      showEditUserDialog: false,
      currentEditUser: null,
      
      createUserForm: {
        username: '',
        email: '',
        name: '',
        role: 'student',
        password: '',
        is_active: true
      },
      
      editUserForm: {
        username: '',
        email: '',
        name: '',
        role: '',
        is_active: true
      },
      
      createUserRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
        ],
        name: [
          { required: true, message: '请输入姓名', trigger: 'blur' }
        ],
        role: [
          { required: true, message: '请选择角色', trigger: 'change' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
        ]
      },
      
      editUserRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
        ],
        name: [
          { required: true, message: '请输入姓名', trigger: 'blur' }
        ],
        role: [
          { required: true, message: '请选择角色', trigger: 'change' }
        ]
      }
    }
  },
  computed: {
    ...mapGetters('auth', ['currentUser']),
    
    filteredUsers() {
      if (!this.userSearchKeyword) {
        return this.users
      }
      const keyword = this.userSearchKeyword.toLowerCase()
      return this.users.filter(user => 
        (user.name && user.name.toLowerCase().includes(keyword)) ||
        (user.username && user.username.toLowerCase().includes(keyword)) ||
        (user.email && user.email.toLowerCase().includes(keyword))
      )
    }
  },
  async created() {
    await this.loadData()
  },
  methods: {
    async loadData() {
      await Promise.all([
        this.fetchStats(),
        this.fetchUsers(),
        this.fetchArticles()
      ])
    },
    
    async fetchStats() {
      try {
        const response = await this.$http.get('/admin/stats')
        this.stats = response.data.stats
      } catch (error) {
        console.error('Failed to fetch stats:', error)
      }
    },
    
    async fetchUsers() {
      this.loadingUsers = true
      try {
        const params = {
          page: this.userPagination.page,
          limit: this.userPagination.limit
        }
        
        if (this.userFilters.role) {
          params.role = this.userFilters.role
        }
        
        const response = await this.$http.get('/admin/users', { params })
        this.users = response.data.users
        this.userPagination.total = response.data.total
      } catch (error) {
        this.$message.error('获取用户列表失败')
      } finally {
        this.loadingUsers = false
      }
    },
    
    async fetchArticles() {
      this.loadingArticles = true
      try {
        const params = {
          page: this.articlePagination.page,
          limit: this.articlePagination.limit
        }
        
        if (this.articleFilters.status) {
          params.status = this.articleFilters.status
        }
        
        const response = await this.$http.get('/admin/articles', { params })
        this.articles = response.data.articles
        this.articlePagination.total = response.data.total
      } catch (error) {
        this.$message.error('获取文章列表失败')
      } finally {
        this.loadingArticles = false
      }
    },
    
    async createUser() {
      this.$refs.createUserForm.validate(async (valid) => {
        if (valid) {
          this.creatingUser = true
          try {
            await this.$http.post('/admin/users', this.createUserForm)
            this.$message.success('用户创建成功')
            this.showCreateUserDialog = false
            this.resetCreateUserForm()
            await this.fetchUsers()
            await this.fetchStats()
          } catch (error) {
            this.$message.error(error.response?.data?.error || '创建用户失败')
          } finally {
            this.creatingUser = false
          }
        }
      })
    },
    
    editUser(user) {
      this.currentEditUser = user
      this.editUserForm = {
        username: user.username,
        email: user.email,
        name: user.name,
        role: user.role,
        is_active: user.is_active
      }
      this.showEditUserDialog = true
    },
    
    async updateUser() {
      this.$refs.editUserForm.validate(async (valid) => {
        if (valid) {
          this.updatingUser = true
          try {
            await this.$http.put(`/admin/users/${this.currentEditUser.id}`, this.editUserForm)
            this.$message.success('用户更新成功')
            this.showEditUserDialog = false
            await this.fetchUsers()
          } catch (error) {
            this.$message.error(error.response?.data?.error || '更新用户失败')
          } finally {
            this.updatingUser = false
          }
        }
      })
    },
    
    async toggleUserStatus(user) {
      try {
        await this.$http.put(`/admin/users/${user.id}`, {
          is_active: !user.is_active
        })
        this.$message.success(`用户已${user.is_active ? '禁用' : '启用'}`)
        await this.fetchUsers()
      } catch (error) {
        this.$message.error('操作失败')
      }
    },
    
    deleteUser(userId) {
      this.$confirm('确定要删除这个用户吗？此操作不可恢复。', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await this.$http.delete(`/admin/users/${userId}`)
          this.$message.success('用户删除成功')
          await this.fetchUsers()
          await this.fetchStats()
        } catch (error) {
          this.$message.error('删除用户失败')
        }
      }).catch(() => {})
    },
    
    deleteRating(ratingId) {
      this.$confirm('确定要删除这个评价吗？此操作不可恢复。', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await this.$http.delete(`/admin/ratings/${ratingId}`)
          this.$message.success('评价删除成功')
          await this.fetchRatings()
          await this.fetchStats()
        } catch (error) {
          this.$message.error('删除评价失败')
        }
      }).catch(() => {})
    },
    
    handleUserPageChange(page) {
      this.userPagination.page = page
      this.fetchUsers()
    },
    
    handleRatingPageChange(page) {
      this.ratingPagination.page = page
      this.fetchRatings()
    },
    
    resetCreateUserForm() {
      this.createUserForm = {
        username: '',
        email: '',
        name: '',
        role: 'student',
        password: '',
        is_active: true
      }
      this.$refs.createUserForm?.resetFields()
    },
    
    getRoleType(role) {
      const roleMap = {
        student: 'success',
        expert: 'warning',
        admin: 'danger'
      }
      return roleMap[role] || 'info'
    },
    
    getRoleText(role) {
      const roleMap = {
        student: '学生',
        expert: '专家',
        admin: '管理员'
      }
      return roleMap[role] || '未知'
    },
    
    // 文章相关方法
    handleArticlePageChange(page) {
      this.articlePagination.page = page
      this.fetchArticles()
    },

    getStatusType(status) {
      const statusMap = {
        draft: 'info',
        published: 'success',
        archived: 'warning'
      }
      return statusMap[status] || 'info'
    },

    getStatusText(status) {
      const statusMap = {
        draft: '草稿',
        published: '已发布',
        archived: '已归档'
      }
      return statusMap[status] || '未知'
    },

    async pushArticle(article) {
      this.selectedArticle = article
      this.selectedUserIds = []
      this.userSearchKeyword = ''
      this.pushedUserIds = []
      
      // 获取已推送的用户列表
      try {
        const response = await this.$http.get(`/admin/articles/${article.id}/pushed-users`)
        this.pushedUserIds = response.data.pushed_user_ids || []
      } catch (error) {
        console.error('Failed to get pushed users:', error)
      }
      
      this.showPushDialog = true
    },

    viewArticle(articleId) {
      this.$router.push(`/articles/${articleId}`)
    },

    searchUsers() {
      // 搜索功能由computed属性filteredUsers处理
    },

    selectAllUsers() {
      // 只选择未推送的用户
      this.selectedUserIds = this.filteredUsers
        .filter(user => !this.pushedUserIds.includes(user.id))
        .map(user => user.id)
    },

    clearSelection() {
      this.selectedUserIds = []
    },

    async submitPush() {
      if (this.selectedUserIds.length === 0) {
        this.$message.warning('请选择要推送的用户')
        return
      }

      this.pushSubmitting = true
      try {
        const response = await this.$http.post('/admin/articles/push', {
          article_id: this.selectedArticle.id,
          user_ids: this.selectedUserIds
        })

        // 根据推送结果显示不同的消息
        const { pushed_count, skipped_count, message } = response.data
        
        if (pushed_count > 0 && skipped_count > 0) {
          this.$message.success(message)
        } else if (pushed_count > 0) {
          this.$message.success(message)
        } else {
          this.$message.warning(message)
        }
        
        this.showPushDialog = false
      } catch (error) {
        this.$message.error(error.response?.data?.error || '推送失败')
      } finally {
        this.pushSubmitting = false
      }
    },

    formatTime(time) {
      return new Date(time).toLocaleString('zh-CN')
    }
  }
}
</script>

<style scoped>
.admin {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  color: #333;
}

.stats-section {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
  position: relative;
  overflow: hidden;
}

.stat-card :deep(.el-card__body) {
  padding: 20px;
}

.stat-content {
  position: relative;
  z-index: 2;
}

.stat-number {
  font-size: 32px;
  font-weight: bold;
  color: #333;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.stat-icon {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 40px;
  color: #e6e6e6;
  z-index: 1;
}

.admin-tabs {
  background: white;
  border-radius: 8px;
  padding: 20px;
}

.tab-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filters {
  display: flex;
  gap: 10px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pagination {
  margin-top: 20px;
  text-align: center;
}

/* 推送对话框样式 */
.push-dialog-content {
  max-height: 500px;
  overflow-y: auto;
}

.article-info {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 6px;
  margin-bottom: 20px;
}

.article-info h4 {
  margin: 0 0 8px 0;
  color: #333;
}

.article-summary {
  margin: 0 0 8px 0;
  color: #666;
  font-size: 14px;
}

.article-meta {
  display: flex;
  gap: 15px;
  font-size: 12px;
  color: #999;
}

.user-selection {
  max-height: 300px;
  overflow-y: auto;
}

.selection-header {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.user-list {
  border: 1px solid #e6e6e6;
  border-radius: 4px;
  max-height: 250px;
  overflow-y: auto;
}

.user-item {
  padding: 8px 12px;
  border-bottom: 1px solid #f0f0f0;
}

.user-item:last-child {
  border-bottom: none;
}

.user-item .user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-name {
  flex: 1;
  margin-right: 8px;
}

.dialog-footer {
  text-align: right;
}

.user-item.already-pushed {
  background-color: #f5f7fa;
  opacity: 0.7;
}

.user-item.already-pushed .user-info {
  color: #909399;
}

.user-item.already-pushed .el-checkbox {
  cursor: not-allowed;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .admin {
    padding: 0;
  }
  
  .page-header {
    padding: 15px;
    margin-bottom: 15px;
  }
  
  .page-header h2 {
    font-size: 20px;
  }
  
  .stats-section {
    padding: 0 15px;
    margin-bottom: 15px;
  }
  
  .stat-card {
    margin-bottom: 15px;
  }
  
  .admin-tabs {
    margin: 0;
    border-radius: 0;
    padding: 15px;
  }
  
  .tab-header {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }
  
  .filters {
    flex-direction: column;
    gap: 10px;
  }
  
  .filters .el-select,
  .filters .el-input {
    width: 100%;
  }
  
  /* 表格移动端优化 */
  .el-table {
    font-size: 12px;
  }
  
  .el-table th,
  .el-table td {
    padding: 8px 4px;
  }
  
  .user-info {
    flex-direction: column;
    gap: 4px;
    align-items: flex-start;
  }
  
  .user-info .el-avatar {
    align-self: center;
  }
  
  /* 推送对话框移动端优化 */
  .push-dialog-content {
    max-height: 60vh;
  }
  
  .article-info {
    padding: 12px;
  }
  
  .article-info h4 {
    font-size: 16px;
  }
  
  .selection-header {
    flex-direction: column;
    gap: 10px;
    align-items: stretch;
  }
  
  .selection-header .el-input {
    width: 100%;
  }
  
  .user-list {
    max-height: 200px;
  }
  
  .user-item {
    padding: 12px;
  }
  
  .user-item .user-info {
    flex-direction: row;
    align-items: center;
    gap: 8px;
  }
  
  .pagination {
    padding: 0 15px;
  }
}

@media (max-width: 480px) {
  .page-header h2 {
    font-size: 18px;
  }
  
  .stat-number {
    font-size: 24px;
  }
  
  .stat-label {
    font-size: 12px;
  }
  
  .admin-tabs {
    padding: 10px;
  }
  
  .el-table {
    font-size: 11px;
  }
  
  .el-table th,
  .el-table td {
    padding: 6px 2px;
  }
  
  .el-button--mini {
    padding: 4px 8px;
    font-size: 11px;
  }
  
  .push-dialog-content {
    max-height: 50vh;
  }
  
  .user-list {
    max-height: 150px;
  }
}
</style>
