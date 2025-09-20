<template>
  <Layout>
    <div class="ratings">
      <div class="page-header">
        <h2>评价管理</h2>
        <el-button
          type="primary"
          @click="showCreateDialog = true"
          v-if="isStudent || isExpert"
        >
          <i class="el-icon-edit" /> 写评价
        </el-button>
      </div>

      <el-tabs v-model="activeTab" class="ratings-tabs">
        <el-tab-pane label="公开评价" name="public">
          <div class="tab-content">
            <div v-if="publicRatings.length === 0" class="empty-state">
              <i class="el-icon-star-off" />
              <p>暂无公开评价</p>
            </div>
            
            <div v-else class="rating-cards">
              <el-card
                v-for="rating in publicRatings"
                :key="rating.id"
                class="rating-card"
                shadow="hover"
              >
                <div class="rating-header">
                  <div class="rater-info">
                    <el-avatar :size="40" :src="rating.rater.avatar" icon="el-icon-user-solid" />
                    <div class="rater-details">
                      <h4>{{ rating.rater.name }}</h4>
                      <el-tag size="mini" :type="getRoleType(rating.rater.role)">
                        {{ getRoleText(rating.rater.role) }}
                      </el-tag>
                    </div>
                  </div>
                  <div class="rating-score">
                    <el-rate
                      v-model="rating.score"
                      disabled
                      show-score
                      text-color="#ff9900"
                      score-template="{value}"
                    />
                  </div>
                </div>
                
                <div class="rating-content">
                  <p>{{ rating.content }}</p>
                </div>
                
                <div class="rating-footer">
                  <div class="rating-meta">
                    <span class="rating-time">{{ formatTime(rating.created_at) }}</span>
                    <span v-if="rating.file" class="file-info">
                      <i class="el-icon-paperclip" />
                      {{ rating.file.title || rating.file.file_name }}
                    </span>
                  </div>
                  
                  <div class="rating-actions">
                    <el-button
                      size="mini"
                      type="text"
                      @click="giveFeedback(rating.id, true)"
                      :disabled="hasFeedback(rating)"
                    >
                      <i class="el-icon-thumb" />
                      有帮助 ({{ getHelpfulCount(rating.feedbacks) }})
                    </el-button>
                    <el-button
                      size="mini"
                      type="text"
                      @click="giveFeedback(rating.id, false)"
                      :disabled="hasFeedback(rating)"
                    >
                      <i class="el-icon-thumb" style="transform: rotate(180deg)" />
                      没帮助 ({{ getUnhelpfulCount(rating.feedbacks) }})
                    </el-button>
                  </div>
                </div>
              </el-card>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="我收到的评价" name="received" v-if="isStudent || isExpert">
          <div class="tab-content">
            <div v-if="receivedRatings.length === 0" class="empty-state">
              <i class="el-icon-star-off" />
              <p>还没有收到任何评价</p>
            </div>
            
            <div v-else class="rating-cards">
              <el-card
                v-for="rating in receivedRatings"
                :key="rating.id"
                class="rating-card"
                shadow="hover"
              >
                <div class="rating-header">
                  <div class="rater-info">
                    <el-avatar :size="40" :src="rating.rater.avatar" icon="el-icon-user-solid" />
                    <div class="rater-details">
                      <h4>{{ rating.rater.name }}</h4>
                      <el-tag size="mini" :type="getRoleType(rating.rater.role)">
                        {{ getRoleText(rating.rater.role) }}
                      </el-tag>
                    </div>
                  </div>
                  <div class="rating-score">
                    <el-rate
                      v-model="rating.score"
                      disabled
                      show-score
                      text-color="#ff9900"
                      score-template="{value}"
                    />
                  </div>
                </div>
                
                <div class="rating-content">
                  <p>{{ rating.content }}</p>
                </div>
                
                <div class="rating-footer">
                  <div class="rating-meta">
                    <span class="rating-time">{{ formatTime(rating.created_at) }}</span>
                    <span v-if="rating.file" class="file-info">
                      <i class="el-icon-paperclip" />
                      {{ rating.file.title || rating.file.file_name }}
                    </span>
                  </div>
                  
                  <div class="feedback-summary" v-if="rating.feedbacks && rating.feedbacks.length > 0">
                    <span class="helpful-count">
                      {{ getHelpfulCount(rating.feedbacks) }} 人觉得有帮助
                    </span>
                  </div>
                </div>
              </el-card>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="我给出的评价" name="given" v-if="isStudent || isExpert">
          <div class="tab-content">
            <div v-if="givenRatings.length === 0" class="empty-state">
              <i class="el-icon-edit-outline" />
              <p>还没有给出任何评价</p>
            </div>
            
            <div v-else class="rating-cards">
              <el-card
                v-for="rating in givenRatings"
                :key="rating.id"
                class="rating-card"
                shadow="hover"
              >
                <div class="rating-header">
                  <div class="rated-info">
                    <span class="rated-label">评价给：</span>
                    <el-avatar :size="32" :src="rating.rated.avatar" icon="el-icon-user-solid" />
                    <span class="rated-name">{{ rating.rated.name }}</span>
                  </div>
                  <div class="rating-score">
                    <el-rate
                      v-model="rating.score"
                      disabled
                      show-score
                      text-color="#ff9900"
                      score-template="{value}"
                    />
                  </div>
                </div>
                
                <div class="rating-content">
                  <p>{{ rating.content }}</p>
                </div>
                
                <div class="rating-footer">
                  <div class="rating-meta">
                    <span class="rating-time">{{ formatTime(rating.created_at) }}</span>
                    <span v-if="rating.file" class="file-info">
                      <i class="el-icon-paperclip" />
                      {{ rating.file.title || rating.file.file_name }}
                    </span>
                  </div>
                </div>
              </el-card>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 创建评价对话框 -->
    <el-dialog
      title="写评价"
      :visible.sync="showCreateDialog"
      width="600px"
    >
      <el-form :model="createForm" :rules="createRules" ref="createForm">
        <el-form-item label="被评价用户ID" prop="rated_id">
          <el-input
            v-model="createForm.rated_id"
            placeholder="请输入用户ID"
            type="number"
          />
          <div class="form-tip">提示：可以从用户列表或个人主页获取用户ID</div>
        </el-form-item>
        
        <el-form-item label="文件ID（可选）" prop="file_id">
          <el-input
            v-model="createForm.file_id"
            placeholder="如果是针对特定文件的评价，请输入文件ID"
            type="number"
          />
        </el-form-item>
        
        <el-form-item label="评分" prop="score">
          <el-rate v-model="createForm.score" show-text />
        </el-form-item>
        
        <el-form-item label="评价内容" prop="content">
          <el-input
            v-model="createForm.content"
            type="textarea"
            :rows="4"
            placeholder="请输入详细的评价内容..."
          />
        </el-form-item>
        
        <el-form-item>
          <el-checkbox v-model="createForm.is_anonymous">匿名评价</el-checkbox>
          <div class="form-tip">匿名评价不会显示您的身份信息</div>
        </el-form-item>
      </el-form>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="createRating" :loading="creating">提交评价</el-button>
      </div>
    </el-dialog>
  </Layout>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import Layout from '@/components/Layout.vue'

export default {
  name: 'Ratings',
  components: {
    Layout
  },
  data() {
    return {
      activeTab: 'public',
      showCreateDialog: false,
      creating: false,
      createForm: {
        rated_id: '',
        file_id: '',
        score: 5,
        content: '',
        is_anonymous: false
      },
      createRules: {
        rated_id: [
          { required: true, message: '请输入被评价用户ID', trigger: 'blur' }
        ],
        score: [
          { required: true, message: '请选择评分', trigger: 'change' }
        ],
        content: [
          { required: true, message: '请输入评价内容', trigger: 'blur' },
          { min: 10, message: '评价内容至少10个字符', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    ...mapGetters('auth', ['currentUser', 'isStudent', 'isExpert']),
    ...mapGetters('rating', ['publicRatings', 'receivedRatings', 'givenRatings'])
  },
  async created() {
    await this.loadData()
  },
  methods: {
    ...mapActions('rating', ['fetchPublicRatings', 'fetchUserRatings', 'createRating', 'createFeedback']),
    
    async loadData() {
      try {
        await Promise.all([
          this.fetchPublicRatings(),
          this.fetchUserRatings()
        ])
      } catch (error) {
        console.error('Failed to load ratings:', error)
      }
    },
    
    async createRating() {
      this.$refs.createForm.validate(async (valid) => {
        if (valid) {
          this.creating = true
          try {
            const ratingData = {
              rated_id: parseInt(this.createForm.rated_id),
              score: this.createForm.score,
              content: this.createForm.content,
              is_anonymous: this.createForm.is_anonymous
            }
            
            if (this.createForm.file_id) {
              ratingData.file_id = parseInt(this.createForm.file_id)
            }
            
            const result = await this.createRating(ratingData)
            if (result.success) {
              this.$message.success('评价提交成功')
              this.showCreateDialog = false
              this.resetCreateForm()
              await this.loadData()
            } else {
              this.$message.error(result.message)
            }
          } catch (error) {
            this.$message.error('提交失败，请重试')
          } finally {
            this.creating = false
          }
        }
      })
    },
    
    async giveFeedback(ratingId, isHelpful) {
      try {
        const result = await this.createFeedback({ ratingId, isHelpful })
        if (result.success) {
          this.$message.success('反馈成功')
          await this.fetchPublicRatings() // 刷新数据
        } else {
          this.$message.error(result.message)
        }
      } catch (error) {
        this.$message.error('反馈失败')
      }
    },
    
    resetCreateForm() {
      this.createForm = {
        rated_id: '',
        file_id: '',
        score: 5,
        content: '',
        is_anonymous: false
      }
      this.$refs.createForm?.resetFields()
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
    },
    
    getHelpfulCount(feedbacks) {
      return feedbacks ? feedbacks.filter(f => f.is_helpful).length : 0
    },
    
    getUnhelpfulCount(feedbacks) {
      return feedbacks ? feedbacks.filter(f => !f.is_helpful).length : 0
    },
    
    hasFeedback(rating) {
      // 检查当前用户是否已经对这个评价给过反馈
      return rating.feedbacks && rating.feedbacks.some(f => f.user_id === this.currentUser.id)
    }
  }
}
</script>

<style scoped>
.ratings {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  color: #333;
}

.ratings-tabs {
  background: white;
  border-radius: 8px;
  padding: 20px;
}

.tab-content {
  margin-top: 20px;
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

.rating-cards {
  display: grid;
  gap: 20px;
}

.rating-card {
  transition: transform 0.2s;
}

.rating-card:hover {
  transform: translateY(-2px);
}

.rating-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.rater-info, .rated-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.rater-details h4 {
  margin: 0 0 4px 0;
  font-size: 16px;
  color: #333;
}

.rated-label {
  color: #666;
  font-size: 14px;
}

.rated-name {
  font-weight: 500;
  color: #333;
}

.rating-content {
  margin-bottom: 15px;
}

.rating-content p {
  margin: 0;
  color: #666;
  line-height: 1.6;
  font-size: 14px;
}

.rating-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
}

.rating-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.rating-time {
  font-size: 12px;
  color: #999;
}

.file-info {
  font-size: 12px;
  color: #409eff;
}

.rating-actions {
  display: flex;
  gap: 10px;
}

.feedback-summary {
  font-size: 12px;
}

.helpful-count {
  color: #67c23a;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

/* 移动端响应式样式 */
@media (max-width: 768px) {
  .ratings {
    max-width: 100%;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .page-header h2 {
    font-size: 20px;
  }
  
  .ratings-tabs {
    padding: 15px;
  }
  
  .rating-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .rater-info, .rated-info {
    gap: 8px;
  }
  
  .rater-details h4 {
    font-size: 14px;
  }
  
  .rating-content p {
    font-size: 13px;
  }
  
  .rating-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .rating-actions {
    align-self: flex-end;
    gap: 5px;
  }
  
  .rating-actions .el-button {
    font-size: 11px;
    padding: 4px 8px;
  }
  
  .rating-meta {
    gap: 2px;
  }
}

@media (max-width: 480px) {
  .ratings-tabs {
    padding: 10px;
  }
  
  .page-header h2 {
    font-size: 18px;
  }
  
  .rating-card :deep(.el-card__body) {
    padding: 15px;
  }
  
  .rater-info, .rated-info {
    gap: 6px;
  }
  
  .rater-details h4 {
    font-size: 13px;
  }
  
  .rating-content p {
    font-size: 12px;
    line-height: 1.4;
  }
  
  .rating-time, .file-info {
    font-size: 11px;
  }
  
  .rating-actions .el-button {
    font-size: 10px;
    padding: 3px 6px;
  }
  
  .empty-state {
    padding: 40px 15px;
  }
  
  .empty-state i {
    font-size: 36px;
  }
}
</style>
