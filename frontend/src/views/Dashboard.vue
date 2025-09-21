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
        <el-col :xs="12" :sm="12" :md="6" :lg="6" v-if="isStudent || isExpert">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-number">{{ totalReceivedRatings }}</div>
              <div class="stat-label">收到的评价</div>
            </div>
            <i class="el-icon-star-on stat-icon" />
          </el-card>
        </el-col>
        
        <el-col :xs="12" :sm="12" :md="6" :lg="6" v-if="isStudent || isExpert">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-number">{{ totalGivenRatings }}</div>
              <div class="stat-label">给出的评价</div>
            </div>
            <i class="el-icon-edit stat-icon" />
          </el-card>
        </el-col>
        
        <el-col :xs="12" :sm="12" :md="6" :lg="6" v-if="isExpert">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-number">{{ following.length }}</div>
              <div class="stat-label">关注的学生</div>
            </div>
            <i class="el-icon-user stat-icon" />
          </el-card>
        </el-col>
        
        <el-col :xs="12" :sm="12" :md="6" :lg="6">
          <el-card class="stat-card clickable-card" @click.native="showFollowersDialog">
            <div class="stat-content">
              <div class="stat-number">{{ followers.length }}</div>
              <div class="stat-label">关注我的用户</div>
            </div>
            <i class="el-icon-view stat-icon" />
          </el-card>
        </el-col>
      </el-row>

      <div class="content-section">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="24" :md="12" :lg="12">
            <el-card class="content-card">
              <div slot="header" class="card-header">
                <span>{{ isStudent ? '最新推送文章评论' : '最新公开文章评论' }}</span>
                <el-button type="text" @click="$router.push('/articles')">查看更多</el-button>
              </div>
              
              <div v-if="latestComments.length === 0" class="empty-state">
                <i class="el-icon-chat-line-square" />
                <p>{{ isStudent ? '暂无推送文章评论' : '暂无文章评论' }}</p>
                <p v-if="isStudent" class="empty-hint">管理员推送文章后会显示相关评论</p>
              </div>
              
              <div v-else class="comment-list">
                <div
                  v-for="comment in latestComments.slice(0, 5)"
                  :key="comment.id"
                  class="comment-item"
                >
                  <div class="comment-header">
                    <span class="commenter-name">{{ comment.user.name || comment.user.username }}</span>
                    <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
                  </div>
                  <p class="comment-content">{{ comment.content }}</p>
                  <div class="comment-footer">
                    <span class="article-title" @click="viewArticle(comment.article_id)">
                      文章：{{ comment.article ? comment.article.title : '未知文章' }}
                    </span>
                  </div>
                </div>
              </div>
            </el-card>
          </el-col>
          
          <el-col :xs="24" :sm="24" :md="12" :lg="12">
            <el-card class="content-card">
              <div slot="header" class="card-header">
                <span>快速操作</span>
              </div>
              
              <div class="quick-actions">
                <el-button
                  type="primary"
                  icon="el-icon-edit"
                  @click="$router.push('/articles/create')"
                >
                  创建文章
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

      <!-- 最新文章展示区域 -->
      <div class="articles-section">
        <el-card class="articles-card">
          <div slot="header" class="card-header">
            <span>最新文章</span>
            <el-button type="text" @click="$router.push('/articles')">查看更多</el-button>
          </div>
          
          <div v-if="latestArticles.length === 0" class="empty-state">
            <i class="el-icon-document" />
            <p>暂无文章</p>
          </div>
          
          <div v-else class="articles-list">
            <div
              v-for="article in latestArticles.slice(0, 5)"
              :key="article.id"
              class="article-item"
              @click="viewArticle(article.id)"
            >
              <div class="article-info">
                <h4 class="article-title">{{ article.title }}</h4>
                <p class="article-summary">{{ article.summary || '暂无摘要' }}</p>
                <div class="article-meta">
                  <span class="author">{{ article.author.name || article.author.username }}</span>
                  <span class="publish-time">{{ formatTime(article.published_at || article.created_at) }}</span>
                  <div class="article-stats">
                    <span><i class="el-icon-view" /> {{ article.view_count }}</span>
                    <span><i class="el-icon-star-off" /> {{ article.like_count }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 用户作品展示区域 -->
      <div class="works-section">
        <el-card class="works-card">
          <div slot="header" class="card-header">
            <span>最新用户作品</span>
            <el-button type="text" @click="refreshPublicFiles">刷新</el-button>
          </div>
          
          <div v-if="publicFiles.length === 0" class="empty-state">
            <i class="el-icon-folder-opened" />
            <p>暂无公开作品</p>
          </div>
          
          <div v-else class="works-grid">
            <div
              v-for="file in publicFiles.slice(0, 8)"
              :key="file.id"
              class="work-item"
              @click="viewFile(file)"
            >
              <div class="work-preview">
                <div class="file-icon">
                  <i :class="getFileIcon(file.file_type)" />
                </div>
                <div v-if="isImageFile(file)" class="image-preview">
                  <img :src="getFilePreview(file)" :alt="file.title || file.file_name" />
                </div>
              </div>
              
              <div class="work-info">
                <h4 class="work-title">{{ file.title || file.file_name }}</h4>
                <p class="work-description" v-if="file.description">{{ file.description }}</p>
                <div class="work-meta">
                  <span class="author">{{ file.user.name || file.user.username }}</span>
                  <span class="upload-time">{{ formatTime(file.created_at) }}</span>
                </div>
                <div class="work-stats">
                  <span class="file-size">{{ formatFileSize(file.file_size) }}</span>
                  <span class="file-type">{{ getFileTypeText(file.file_type) }}</span>
                </div>
              </div>
              
              <div class="work-actions">
                <el-button size="mini" type="text" @click.stop="showComments(file)">
                  <i class="el-icon-chat-line-square" /> 评论
                </el-button>
                <el-button size="mini" type="text" @click.stop="viewUserProfile(file.user.id)">
                  <i class="el-icon-user" /> 作者
                </el-button>
              </div>
            </div>
          </div>
        </el-card>
      </div>
    </div>


    <!-- 评论对话框 -->
    <el-dialog
      :title="`${selectedFile?.title || selectedFile?.file_name || '文件'} - 评论`"
      :visible.sync="showCommentDialog"
      width="700px"
      :before-close="closeCommentDialog"
    >
      <div class="comment-section">
        <!-- 评论输入框 -->
        <div class="comment-input-section">
          <el-input
            v-model="commentForm.content"
            type="textarea"
            :rows="3"
            placeholder="写下你的评论..."
            maxlength="500"
            show-word-limit
          />
          <div class="comment-actions">
            <el-button 
              type="primary" 
              size="small" 
              @click="submitComment"
              :loading="commentSubmitting"
              :disabled="!commentForm.content.trim()"
            >
              发表评论
            </el-button>
          </div>
        </div>

        <!-- 评论列表 -->
        <div class="comments-list">
          <div class="comments-header">
            <h4>评论 ({{ fileComments.length }})</h4>
            <el-button size="mini" type="text" @click="refreshComments">
              <i class="el-icon-refresh" /> 刷新
            </el-button>
          </div>

          <div v-if="fileComments.length === 0" class="empty-comments">
            <i class="el-icon-chat-line-square" />
            <p>暂无评论，快来抢沙发吧！</p>
          </div>

          <div v-else class="comment-list">
            <div
              v-for="comment in fileComments"
              :key="comment.id"
              class="comment-item"
            >
              <div class="comment-header">
                <div class="comment-user">
                  <span class="username">{{ comment.user.name || comment.user.username }}</span>
                  <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
                </div>
                <div class="comment-actions">
                  <el-button size="mini" type="text" @click="replyToComment(comment)">
                    回复
                  </el-button>
                  <el-button 
                    v-if="comment.user_id === currentUser.id"
                    size="mini" 
                    type="text" 
                    @click="deleteComment(comment.id)"
                  >
                    删除
                  </el-button>
                </div>
              </div>
              
              <div class="comment-content">{{ comment.content }}</div>
              
              <!-- 回复列表 -->
              <div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
                <div
                  v-for="reply in comment.replies"
                  :key="reply.id"
                  class="reply-item"
                >
                  <div class="reply-header">
                    <div class="reply-user">
                      <span class="username">{{ reply.user.name || reply.user.username }}</span>
                      <span class="reply-time">{{ formatTime(reply.created_at) }}</span>
                    </div>
                    <el-button 
                      v-if="reply.user_id === currentUser.id"
                      size="mini" 
                      type="text" 
                      @click="deleteComment(reply.id)"
                    >
                      删除
                    </el-button>
                  </div>
                  <div class="reply-content">{{ reply.content }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 粉丝对话框 -->
    <el-dialog
      title="关注我的用户"
      :visible.sync="showFollowersDialogVisible"
      width="800px"
      :before-close="closeFollowersDialog"
    >
      <div class="followers-dialog">
        <div class="followers-stats">
          <span>共 {{ followers.length }} 人关注了我</span>
        </div>

        <div v-if="followersLoading" class="followers-loading">
          <el-skeleton :rows="3" animated />
        </div>
        
        <div v-else-if="followers.length === 0" class="followers-empty">
          <i class="el-icon-user" />
          <p>还没有人关注你</p>
          <p class="tip">发布优质内容来吸引更多关注者吧！</p>
        </div>
        
        <div v-else class="followers-dialog-list">
          <div
            v-for="follow in followers"
            :key="follow.id"
            class="follower-dialog-item"
          >
            <div class="follower-dialog-avatar">
              <img 
                v-if="follow.follower.avatar" 
                :src="follow.follower.avatar" 
                :alt="follow.follower.name || follow.follower.username"
              />
              <div v-else class="default-dialog-avatar">
                <i class="el-icon-user" />
              </div>
            </div>
            
            <div class="follower-dialog-info">
              <div class="follower-dialog-basic">
                <h4 class="follower-dialog-name">{{ follow.follower.name || follow.follower.username }}</h4>
                <el-tag :type="getRoleType(follow.follower.role)" size="mini">
                  {{ getRoleText(follow.follower.role) }}
                </el-tag>
              </div>
              
              <p class="follower-dialog-bio" v-if="follow.follower.bio">{{ follow.follower.bio }}</p>
              
              <div class="follow-dialog-info">
                <span class="follow-dialog-time">
                  <i class="el-icon-time" /> 
                  关注时间：{{ formatTime(follow.created_at) }}
                </span>
              </div>
            </div>
            
            <div class="follower-dialog-actions">
              <el-button
                size="small"
                :type="isFollowing(follow.follower.id) ? 'success' : 'primary'"
                @click="toggleFollow(follow.follower.id)"
                :loading="followLoading[follow.follower.id]"
              >
                {{ isFollowing(follow.follower.id) ? '已关注' : '回关' }}
              </el-button>
            </div>
          </div>
        </div>
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
      showCommentDialog: false,
      commentSubmitting: false,
      selectedFile: null,
      fileComments: [],
      latestArticles: [],
      latestComments: [],
      commentForm: {
        content: '',
        parent_id: null
      },
      showFollowersDialogVisible: false,
      followersLoading: false,
      followLoading: {}
    }
  },
  computed: {
    ...mapGetters('auth', ['currentUser', 'isStudent', 'isExpert', 'isAdmin']),
    ...mapGetters('rating', ['receivedRatings', 'givenRatings', 'receivedCommentRatings', 'givenCommentRatings']),
    ...mapGetters('user', ['following', 'followers', 'publicFiles']),
    
    // 计算总评价数（包括文章评价和评论评价）
    totalReceivedRatings() {
      return this.receivedRatings.length + this.receivedCommentRatings.length
    },
    
    totalGivenRatings() {
      return this.givenRatings.length + this.givenCommentRatings.length
    },
    
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
    ...mapActions('rating', ['fetchUserRatings', 'fetchUserCommentRatings']),
    ...mapActions('user', ['fetchFollowing', 'fetchFollowers', 'fetchPublicFiles', 'followUser', 'unfollowUser']),
    
    async loadData() {
      try {
        await Promise.all([
          this.fetchUserRatings(),
          this.fetchUserCommentRatings(),
          this.fetchFollowing(),
          this.fetchFollowers(),
          this.fetchPublicFiles(),
          this.loadLatestArticles(),
          this.loadLatestComments()
        ])
      } catch (error) {
        console.error('Failed to load data:', error)
      }
    },

    async loadLatestArticles() {
      try {
        const response = await this.$http.get('/articles', {
          params: { page: 1, page_size: 5 }
        })
        this.latestArticles = response.data.articles || []
      } catch (error) {
        console.error('Failed to load articles:', error)
        this.latestArticles = []
      }
    },

    viewArticle(articleId) {
      this.$router.push(`/articles/${articleId}`)
    },

    async loadLatestComments() {
      try {
        // 使用新的API获取最新评论，已经应用了匿名逻辑
        const response = await this.$http.get('/comments/latest')
        this.latestComments = (response.data.comments || []).slice(0, 5)
      } catch (error) {
        console.error('Failed to load comments:', error)
        this.latestComments = []
      }
    },
    
    
    formatTime(time) {
      return new Date(time).toLocaleString('zh-CN')
    },

    // 作品相关方法
    async refreshPublicFiles() {
      try {
        await this.fetchPublicFiles()
        this.$message.success('刷新成功')
      } catch (error) {
        this.$message.error('刷新失败')
      }
    },

    getFileIcon(fileType) {
      const iconMap = {
        'image': 'el-icon-picture-outline',
        'video': 'el-icon-video-camera',
        'document': 'el-icon-document',
        'other': 'el-icon-files'
      }
      return iconMap[fileType] || 'el-icon-files'
    },

    getFileTypeText(fileType) {
      const typeMap = {
        'image': '图片',
        'video': '视频',
        'document': '文档',
        'other': '其他'
      }
      return typeMap[fileType] || '未知'
    },

    isImageFile(file) {
      return file.file_type === 'image'
    },

    getFilePreview(file) {
      // 这里应该返回文件的预览URL，暂时返回占位符
      return `/api/files/${file.id}/preview`
    },

    formatFileSize(bytes) {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    },

    viewFile(file) {
      // 查看文件详情
      this.$message.info(`查看文件: ${file.title || file.file_name}`)
    },

    // 评论相关方法
    async showComments(file) {
      this.selectedFile = file
      this.showCommentDialog = true
      await this.loadFileComments(file.id)
    },

    async loadFileComments(fileId) {
      try {
        const response = await this.$http.get(`/files/${fileId}/comments`)
        this.fileComments = response.data.comments || []
      } catch (error) {
        this.$message.error('加载评论失败')
        this.fileComments = []
      }
    },

    async submitComment() {
      if (!this.commentForm.content.trim()) {
        this.$message.warning('请输入评论内容')
        return
      }

      this.commentSubmitting = true
      try {
        const response = await this.$http.post(`/files/${this.selectedFile.id}/comments`, {
          content: this.commentForm.content.trim(),
          parent_id: this.commentForm.parent_id
        })

        if (response.data.comment) {
          this.$message.success('评论发表成功')
          this.commentForm.content = ''
          this.commentForm.parent_id = null
          await this.loadFileComments(this.selectedFile.id)
        }
      } catch (error) {
        this.$message.error(error.response?.data?.error || '评论发表失败')
      } finally {
        this.commentSubmitting = false
      }
    },

    async refreshComments() {
      if (this.selectedFile) {
        await this.loadFileComments(this.selectedFile.id)
        this.$message.success('评论刷新成功')
      }
    },

    replyToComment(comment) {
      this.commentForm.parent_id = comment.id
      this.commentForm.content = `@${comment.user.name || comment.user.username} `
      // 聚焦到输入框
      this.$nextTick(() => {
        const textarea = this.$el.querySelector('.comment-input-section textarea')
        if (textarea) {
          textarea.focus()
          textarea.setSelectionRange(textarea.value.length, textarea.value.length)
        }
      })
    },

    async deleteComment(commentId) {
      try {
        await this.$confirm('确定要删除这条评论吗？', '确认删除', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        await this.$http.delete(`/comments/${commentId}`)
        this.$message.success('评论删除成功')
        await this.loadFileComments(this.selectedFile.id)
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除评论失败')
        }
      }
    },

    closeCommentDialog() {
      this.showCommentDialog = false
      this.selectedFile = null
      this.fileComments = []
      this.commentForm.content = ''
      this.commentForm.parent_id = null
    },

    // 粉丝对话框相关方法
    async showFollowersDialog() {
      this.showFollowersDialogVisible = true
      await this.loadFollowersData()
    },

    async loadFollowersData() {
      this.followersLoading = true
      try {
        await Promise.all([
          this.fetchFollowers(),
          this.fetchFollowing()
        ])
      } catch (error) {
        this.$message.error('加载粉丝数据失败')
      } finally {
        this.followersLoading = false
      }
    },

    closeFollowersDialog() {
      this.showFollowersDialogVisible = false
    },

    isFollowing(userId) {
      return this.following.some(follow => follow.followed_id === userId)
    },

    async toggleFollow(userId) {
      this.$set(this.followLoading, userId, true)
      try {
        if (this.isFollowing(userId)) {
          const result = await this.unfollowUser(userId)
          if (result.success) {
            this.$message.success('取消关注成功')
          } else {
            this.$message.error(result.message)
          }
        } else {
          const result = await this.followUser(userId)
          if (result.success) {
            this.$message.success('关注成功')
          } else {
            this.$message.error(result.message)
          }
        }
      } catch (error) {
        this.$message.error('操作失败')
      } finally {
        this.$set(this.followLoading, userId, false)
      }
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
      return roleMap[role] || '用户'
    },

    viewUserProfile(userId) {
      // 查看用户资料
      this.$message.info(`查看用户 ${userId} 的资料`)
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

.clickable-card {
  cursor: pointer;
  transition: all 0.3s ease;
}

.clickable-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px 0 rgba(0, 0, 0, 0.15);
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

.articles-section {
  margin-bottom: 20px;
}

.articles-card {
  min-height: 300px;
}

.articles-list {
  max-height: 400px;
  overflow-y: auto;
}

.article-item {
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: all 0.3s ease;
}

.article-item:hover {
  background-color: #f8f9fa;
  padding-left: 10px;
}

.article-item:last-child {
  border-bottom: none;
}

.article-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.article-summary {
  font-size: 13px;
  color: #666;
  margin: 0 0 8px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #999;
}

.article-meta .author {
  color: #667eea;
  font-weight: 500;
}

.article-stats {
  display: flex;
  gap: 10px;
}

.works-section {
  margin-bottom: 20px;
}

.works-card {
  min-height: 400px;
}

.works-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  max-height: 600px;
  overflow-y: auto;
}

.work-item {
  border: 1px solid #e6e6e6;
  border-radius: 8px;
  padding: 15px;
  cursor: pointer;
  transition: all 0.3s ease;
  background: white;
}

.work-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.work-preview {
  position: relative;
  height: 120px;
  background: #f8f9fa;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
  overflow: hidden;
}

.file-icon {
  font-size: 48px;
  color: #999;
}

.image-preview {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 6px;
}

.work-info {
  margin-bottom: 12px;
}

.work-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.work-description {
  font-size: 13px;
  color: #666;
  margin: 0 0 8px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.work-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 12px;
  color: #999;
}

.author {
  font-weight: 500;
  color: #667eea;
}

.work-stats {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #999;
}

.work-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.work-actions .el-button {
  font-size: 12px;
  padding: 4px 8px;
}

/* 评论相关样式 */
.comment-section {
  max-height: 500px;
  overflow-y: auto;
}

.comment-input-section {
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.comment-actions {
  margin-top: 10px;
  text-align: right;
}

.comments-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.comments-header h4 {
  margin: 0;
  color: #333;
  font-size: 16px;
}

.empty-comments {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.empty-comments i {
  font-size: 48px;
  margin-bottom: 15px;
  display: block;
}

.comment-list {
  max-height: 300px;
  overflow-y: auto;
}

.comment-item {
  padding: 15px 0;
  border-bottom: 1px solid #f5f5f5;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.comment-user .username {
  font-weight: 500;
  color: #333;
  margin-right: 10px;
}

.comment-time {
  font-size: 12px;
  color: #999;
}

.comment-actions {
  display: flex;
  gap: 5px;
}

.comment-content {
  color: #666;
  line-height: 1.5;
  margin-bottom: 10px;
}

.replies-list {
  margin-left: 20px;
  padding-left: 15px;
  border-left: 2px solid #f0f0f0;
  margin-top: 10px;
}

.reply-item {
  padding: 10px 0;
  border-bottom: 1px solid #f8f8f8;
}

.reply-item:last-child {
  border-bottom: none;
}

.reply-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}

.reply-user .username {
  font-weight: 500;
  color: #333;
  margin-right: 8px;
  font-size: 14px;
}

.reply-time {
  font-size: 11px;
  color: #999;
}

.reply-content {
  color: #666;
  line-height: 1.4;
  font-size: 14px;
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

.comment-list {
  max-height: 300px;
  overflow-y: auto;
}

.comment-item {
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.commenter-name {
  font-weight: 500;
  color: #667eea;
}

.comment-time {
  font-size: 12px;
  color: #999;
}

.comment-content {
  margin: 8px 0;
  color: #666;
  line-height: 1.5;
  font-size: 14px;
}

.comment-footer {
  margin-top: 8px;
}

.article-title {
  font-size: 12px;
  color: #409eff;
  cursor: pointer;
  text-decoration: none;
}

.article-title:hover {
  text-decoration: underline;
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
  margin-left: 0;
}

/* 移动端响应式样式 */
@media (max-width: 768px) {
  .dashboard {
    max-width: 100%;
  }
  
  .welcome-content h2 {
    font-size: 20px;
  }
  
  .welcome-content p {
    font-size: 14px;
  }
  
  .stats-section {
    margin-bottom: 15px;
  }
  
  .stat-card :deep(.el-card__body) {
    padding: 15px;
  }
  
  .stat-number {
    font-size: 24px;
  }
  
  .stat-label {
    font-size: 12px;
  }
  
  .stat-icon {
    font-size: 30px;
    right: 15px;
  }
  
  .content-card {
    height: auto;
    margin-bottom: 15px;
  }
  
  .rating-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .rating-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .rating-feedback {
    align-self: flex-end;
  }
  
  .quick-actions {
    gap: 10px;
  }
  
  .quick-actions .el-button {
    height: 40px;
    font-size: 13px;
  }
  
  .works-grid {
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 15px;
  }
  
  .work-item {
    padding: 12px;
  }
  
  .work-preview {
    height: 100px;
  }
  
  .file-icon {
    font-size: 36px;
  }
  
  .work-title {
    font-size: 14px;
  }
  
  .work-description {
    font-size: 12px;
  }
}

@media (max-width: 480px) {
  .welcome-card :deep(.el-card__body) {
    padding: 20px;
  }
  
  .welcome-content h2 {
    font-size: 18px;
  }
  
  .stat-card :deep(.el-card__body) {
    padding: 12px;
  }
  
  .stat-number {
    font-size: 20px;
  }
  
  .stat-icon {
    font-size: 24px;
    right: 12px;
  }
  
  .rating-item {
    padding: 12px 0;
  }
  
  .rating-content {
    font-size: 13px;
  }
  
  .rating-feedback .el-button {
    font-size: 11px;
    padding: 4px 8px;
  }
  
  .works-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .work-item {
    padding: 10px;
  }
  
  .work-preview {
    height: 80px;
  }
  
  .file-icon {
    font-size: 28px;
  }
  
  .work-title {
    font-size: 13px;
  }
  
  .work-description {
    font-size: 11px;
  }
  
  .work-meta,
  .work-stats {
    font-size: 11px;
  }
  
  .work-actions .el-button {
    font-size: 11px;
    padding: 2px 6px;
  }
}

/* 粉丝对话框样式 */
.followers-dialog {
  max-height: 600px;
  overflow-y: auto;
}

.followers-stats {
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e6e6e6;
  color: #666;
  font-size: 14px;
}

.followers-loading {
  padding: 20px;
}

.followers-empty {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.followers-empty i {
  font-size: 48px;
  margin-bottom: 15px;
  color: #ddd;
}

.followers-empty p {
  margin: 8px 0;
  font-size: 14px;
}

.followers-empty .tip {
  font-size: 12px;
  color: #999;
}

.followers-dialog-list {
  padding: 0;
}

.follower-dialog-item {
  display: flex;
  align-items: flex-start;
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
}

.follower-dialog-item:last-child {
  border-bottom: none;
}

.follower-dialog-avatar {
  width: 50px;
  height: 50px;
  margin-right: 12px;
  flex-shrink: 0;
}

.follower-dialog-avatar img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.default-dialog-avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 20px;
}

.follower-dialog-info {
  flex: 1;
  margin-right: 12px;
}

.follower-dialog-basic {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.follower-dialog-name {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.follower-dialog-bio {
  margin: 0 0 6px 0;
  color: #666;
  font-size: 12px;
  line-height: 1.4;
}

.follow-dialog-info {
  font-size: 11px;
  color: #999;
}

.follow-dialog-time {
  display: flex;
  align-items: center;
  gap: 3px;
}

.follower-dialog-actions {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex-shrink: 0;
}

.follower-dialog-actions .el-button {
  width: 70px;
  font-size: 12px;
  padding: 5px 10px;
}

/* 粉丝对话框移动端样式 */
@media (max-width: 768px) {
  .followers-dialog {
    max-height: 500px;
  }
  
  .follower-dialog-item {
    padding: 12px 0;
  }
  
  .follower-dialog-avatar {
    width: 40px;
    height: 40px;
    margin-right: 10px;
  }
  
  .follower-dialog-name {
    font-size: 13px;
  }
  
  .follower-dialog-bio {
    font-size: 11px;
  }
  
  .follower-dialog-actions .el-button {
    width: 60px;
    font-size: 11px;
    padding: 4px 8px;
  }
}

.empty-hint {
  font-size: 12px !important;
  color: #999 !important;
  margin-top: 5px !important;
}
</style>

