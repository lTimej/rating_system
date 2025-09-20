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
              <div class="stat-number">{{ stats.total_ratings || 0 }}</div>
              <div class="stat-label">总评价数</div>
            </div>
            <i class="el-icon-star-on stat-icon" />
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
            <el-table-column label="操作" width="200">
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
        
        <el-tab-pane label="评价管理" name="ratings">
          <el-table :data="ratings" v-loading="loadingRatings" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column label="评价者" width="120">
              <template slot-scope="scope">
                <div class="user-info">
                  <el-avatar :size="24" :src="scope.row.rater.avatar" icon="el-icon-user-solid" />
                  <span>{{ scope.row.rater.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="被评价者" width="120">
              <template slot-scope="scope">
                <div class="user-info">
                  <el-avatar :size="24" :src="scope.row.rated.avatar" icon="el-icon-user-solid" />
                  <span>{{ scope.row.rated.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="score" label="评分" width="100">
              <template slot-scope="scope">
                <el-rate :value="scope.row.score" disabled show-score text-color="#ff9900" />
              </template>
            </el-table-column>
            <el-table-column prop="content" label="评价内容" show-overflow-tooltip />
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template slot-scope="scope">
                {{ formatTime(scope.row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template slot-scope="scope">
                <el-button
                  size="mini"
                  type="danger"
                  @click="deleteRating(scope.row.id)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <div class="pagination">
            <el-pagination
              @current-change="handleRatingPageChange"
              :current-page="ratingPagination.page"
              :page-size="ratingPagination.limit"
              :total="ratingPagination.total"
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
        <el-form-item label="用户名" prop="username">
          <el-input v-model="createUserForm.username" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="createUserForm.email" />
        </el-form-item>
        
        <el-form-item label="姓名" prop="name">
          <el-input v-model="createUserForm.name" />
        </el-form-item>
        
        <el-form-item label="角色" prop="role">
          <el-select v-model="createUserForm.role" style="width: 100%">
            <el-option label="学生" value="student" />
            <el-option label="专家" value="expert" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input v-model="createUserForm.password" type="password" />
        </el-form-item>
        
        <el-form-item>
          <el-checkbox v-model="createUserForm.is_active">激活用户</el-checkbox>
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
        <el-form-item label="用户名" prop="username">
          <el-input v-model="editUserForm.username" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="editUserForm.email" />
        </el-form-item>
        
        <el-form-item label="姓名" prop="name">
          <el-input v-model="editUserForm.name" />
        </el-form-item>
        
        <el-form-item label="角色" prop="role">
          <el-select v-model="editUserForm.role" style="width: 100%">
            <el-option label="学生" value="student" />
            <el-option label="专家" value="expert" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        
        <el-form-item>
          <el-checkbox v-model="editUserForm.is_active">激活用户</el-checkbox>
        </el-form-item>
      </el-form>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="showEditUserDialog = false">取消</el-button>
        <el-button type="primary" @click="updateUser" :loading="updatingUser">保存</el-button>
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
      loadingRatings: false,
      creatingUser: false,
      updatingUser: false,
      
      stats: {},
      users: [],
      ratings: [],
      
      userFilters: {
        role: ''
      },
      
      userPagination: {
        page: 1,
        limit: 20,
        total: 0
      },
      
      ratingPagination: {
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
    ...mapGetters('auth', ['currentUser'])
  },
  async created() {
    await this.loadData()
  },
  methods: {
    async loadData() {
      await Promise.all([
        this.fetchStats(),
        this.fetchUsers(),
        this.fetchRatings()
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
    
    async fetchRatings() {
      this.loadingRatings = true
      try {
        const params = {
          page: this.ratingPagination.page,
          limit: this.ratingPagination.limit
        }
        
        const response = await this.$http.get('/admin/ratings', { params })
        this.ratings = response.data.ratings
        this.ratingPagination.total = response.data.total
      } catch (error) {
        this.$message.error('获取评价列表失败')
      } finally {
        this.loadingRatings = false
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
</style>
