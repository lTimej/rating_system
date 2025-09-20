<template>
  <Layout>
    <div class="dashboard">
      <div class="welcome-section">
        <el-card class="welcome-card">
          <div class="welcome-content">
            <h2>欢迎回来，{{ currentUser.name || currentUser.username }}！</h2>
            <p>您的角色：<el-tag :type="roleType">{{ roleText }}</el-tag></p>
          </div>
        </el-card>
      </div>

      <el-row :gutter="20" class="stats-section">
        <el-col :span="6" v-if="isStudent || isExpert">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-number">{{ receivedRatings.length }}</div>
              <div class="stat-label">收到的评价</div>
            </div>
            <i class="el-icon-star-on stat-icon" />
          </el-card>
        </el-col>
        
        <el-col :span="6" v-if="isStudent || isExpert">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-number">{{ givenRatings.length }}</div>
              <div class="stat-label">给出的评价</div>
            </div>
            <i class="el-icon-edit stat-icon" />
          </el-card>
        </el-col>
        
        <el-col :span="6" v-if="isExpert">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-number">{{ following.length }}</div>
              <div class="stat-label">关注的学生</div>
            </div>
            <i class="el-icon-user stat-icon" />
          </el-card>
        </el-col>
        
        <el-col :span="6" v-if="isStudent">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-number">{{ followers.length }}</div>
              <div class="stat-label">关注我的专家</div>
            </div>
            <i class="el-icon-view stat-icon" />
          </el-card>
        </el-col>
      </el-row>

      <div class="content-section">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-card class="content-card">
              <div slot="header" class="card-header">
                <span>最新公开评价</span>
                <el-button type="text" @click="$router.push('/ratings')">查看更多</el-button>
              </div>
              
              <div v-if="publicRatings.length === 0" class="empty-state">
                <i class="el-icon-document" />
                <p>暂无公开评价</p>
              </div>
              
              <div v-else class="rating-list">
                <div
                  v-for="rating in publicRatings.slice(0, 5)"
                  :key="rating.id"
                  class="rating-item"
                >
                  <div class="rating-header">
                    <span class="rater-name">{{ rating.rater.name }}</span>
                    <el-rate
                      v-model="rating.score"
                      disabled
                      show-score
                      text-color="#ff9900"
                      score-template="{value}"
                    />
                  </div>
                  <p class="rating-content">{{ rating.content }}</p>
                  <div class="rating-footer">
                    <span class="rating-time">{{ formatTime(rating.created_at) }}</span>
                    <div class="rating-feedback">
                      <el-button
                        size="mini"
                        type="text"
                        @click="giveFeedback(rating.id, true)"
                      >
                        <i class="el-icon-thumb" /> 有帮助
                      </el-button>
                      <el-button
                        size="mini"
                        type="text"
                        @click="giveFeedback(rating.id, false)"
                      >
                        <i class="el-icon-thumb" style="transform: rotate(180deg)" /> 没帮助
                      </el-button>
                    </div>
                  </div>
                </div>
              </div>
            </el-card>
          </el-col>
          
          <el-col :span="12">
            <el-card class="content-card">
              <div slot="header" class="card-header">
                <span>快速操作</span>
              </div>
              
              <div class="quick-actions">
                <el-button
                  type="primary"
                  icon="el-icon-edit"
                  @click="showRatingDialog = true"
                  v-if="isStudent || isExpert"
                >
                  写评价
                </el-button>
                
                <el-button
                  type="success"
                  icon="el-icon-upload"
                  @click="$router.push('/profile')"
                >
                  上传作品
                </el-button>
                
                <el-button
                  type="info"
                  icon="el-icon-user"
                  @click="$router.push('/profile')"
                >
                  编辑资料
                </el-button>
                
                <el-button
                  v-if="isAdmin"
                  type="warning"
                  icon="el-icon-setting"
                  @click="$router.push('/admin')"
                >
                  系统管理
                </el-button>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 评价对话框 -->
    <el-dialog
      title="写评价"
      :visible.sync="showRatingDialog"
      width="600px"
    >
      <el-form :model="ratingForm" :rules="ratingRules" ref="ratingForm">
        <el-form-item label="被评价用户" prop="rated_id">
          <el-input
            v-model="ratingForm.rated_id"
            placeholder="请输入用户ID"
            type="number"
          />
        </el-form-item>
        
        <el-form-item label="评分" prop="score">
          <el-rate v-model="ratingForm.score" show-text />
        </el-form-item>
        
        <el-form-item label="评价内容" prop="content">
          <el-input
            v-model="ratingForm.content"
            type="textarea"
            :rows="4"
            placeholder="请输入评价内容"
          />
        </el-form-item>
        
        <el-form-item>
          <el-checkbox v-model="ratingForm.is_anonymous">匿名评价</el-checkbox>
        </el-form-item>
      </el-form>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="showRatingDialog = false">取消</el-button>
        <el-button type="primary" @click="submitRating" :loading="submitting">提交</el-button>
      </div>
    </el-dialog>
  </Layout>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import Layout from '@/components/Layout.vue'

export default {
  name: 'Dashboard',
  components: {
    Layout
  },
  data() {
    return {
      showRatingDialog: false,
      submitting: false,
      ratingForm: {
        rated_id: '',
        score: 5,
        content: '',
        is_anonymous: false
      },
      ratingRules: {
        rated_id: [
          { required: true, message: '请输入被评价用户ID', trigger: 'blur' }
        ],
        score: [
          { required: true, message: '请选择评分', trigger: 'change' }
        ],
        content: [
          { required: true, message: '请输入评价内容', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    ...mapGetters('auth', ['currentUser', 'isStudent', 'isExpert', 'isAdmin']),
    ...mapGetters('rating', ['publicRatings', 'receivedRatings', 'givenRatings']),
    ...mapGetters('user', ['following', 'followers']),
    
    roleType() {
      const roleMap = {
        student: 'success',
        expert: 'warning',
        admin: 'danger'
      }
      return roleMap[this.currentUser.role] || 'info'
    },
    
    roleText() {
      const roleMap = {
        student: '学生',
        expert: '专家',
        admin: '管理员'
      }
      return roleMap[this.currentUser.role] || '未知'
    }
  },
  async created() {
    await this.loadData()
  },
  methods: {
    ...mapActions('rating', ['fetchPublicRatings', 'fetchUserRatings', 'createRating', 'createFeedback']),
    ...mapActions('user', ['fetchFollowing', 'fetchFollowers']),
    
    async loadData() {
      try {
        await Promise.all([
          this.fetchPublicRatings(),
          this.fetchUserRatings(),
          this.fetchFollowing(),
          this.fetchFollowers()
        ])
      } catch (error) {
        console.error('Failed to load data:', error)
      }
    },
    
    async submitRating() {
      this.$refs.ratingForm.validate(async (valid) => {
        if (valid) {
          this.submitting = true
          try {
            const result = await this.createRating({
              rated_id: parseInt(this.ratingForm.rated_id),
              score: this.ratingForm.score,
              content: this.ratingForm.content,
              is_anonymous: this.ratingForm.is_anonymous
            })
            
            if (result.success) {
              this.$message.success('评价提交成功')
              this.showRatingDialog = false
              this.resetRatingForm()
              await this.fetchPublicRatings()
            } else {
              this.$message.error(result.message)
            }
          } catch (error) {
            this.$message.error('提交失败，请重试')
          } finally {
            this.submitting = false
          }
        }
      })
    },
    
    async giveFeedback(ratingId, isHelpful) {
      try {
        const result = await this.createFeedback({ ratingId, isHelpful })
        if (result.success) {
          this.$message.success('反馈成功')
        } else {
          this.$message.error(result.message)
        }
      } catch (error) {
        this.$message.error('反馈失败')
      }
    },
    
    resetRatingForm() {
      this.ratingForm = {
        rated_id: '',
        score: 5,
        content: '',
        is_anonymous: false
      }
      this.$refs.ratingForm?.resetFields()
    },
    
    formatTime(time) {
      return new Date(time).toLocaleString('zh-CN')
    }
  }
}
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.welcome-section {
  margin-bottom: 20px;
}

.welcome-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.welcome-card :deep(.el-card__body) {
  padding: 30px;
}

.welcome-content h2 {
  margin: 0 0 10px 0;
  font-size: 24px;
  font-weight: 600;
}

.welcome-content p {
  margin: 0;
  font-size: 16px;
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

.content-section {
  margin-bottom: 20px;
}

.content-card {
  height: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.empty-state i {
  font-size: 48px;
  margin-bottom: 20px;
}

.rating-list {
  max-height: 300px;
  overflow-y: auto;
}

.rating-item {
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
}

.rating-item:last-child {
  border-bottom: none;
}

.rating-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.rater-name {
  font-weight: 500;
  color: #333;
}

.rating-content {
  margin: 10px 0;
  color: #666;
  line-height: 1.5;
}

.rating-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.rating-time {
  color: #999;
}

.rating-feedback {
  display: flex;
  gap: 10px;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.quick-actions .el-button {
  justify-content: flex-start;
  height: 48px;
  font-size: 14px;
}
</style>
