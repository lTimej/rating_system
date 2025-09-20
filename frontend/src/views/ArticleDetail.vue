<template>
  <Layout>
    <div class="article-detail">
      <div v-if="loading" class="loading">
        <el-skeleton :rows="8" animated />
      </div>
      
      <div v-else-if="article" class="article-content">
        <!-- 文章头部 -->
        <div class="article-header">
          <h1 class="article-title">{{ article.title }}</h1>
          
          <div class="article-meta">
            <div class="meta-left">
              <div class="author-info">
                <span class="author-name">{{ article.author.name || article.author.username }}</span>
                <span class="publish-time">{{ formatTime(article.published_at || article.created_at) }}</span>
              </div>
              <div class="article-stats">
                <span class="view-count">
                  <i class="el-icon-view" /> {{ article.view_count }} 阅读
                </span>
                <span class="like-count">
                  <i class="el-icon-star-off" /> {{ article.like_count }} 点赞
                </span>
              </div>
            </div>
            
            <div class="meta-right">
              <el-button
                v-if="article.author_id !== currentUser.id"
                size="small"
                :type="isFollowing(article.author_id) ? 'success' : 'info'"
                @click="toggleFollow(article.author_id)"
                :loading="followLoading"
              >
                {{ isFollowing(article.author_id) ? '已关注' : '关注作者' }}
              </el-button>
              <el-button
                v-if="article.author_id === currentUser.id"
                size="small"
                @click="editArticle"
              >
                编辑
              </el-button>
              <el-button
                size="small"
                :type="isLiked ? 'primary' : 'default'"
                @click="toggleLike"
                :loading="likeLoading"
              >
                <i class="el-icon-star-off" /> {{ isLiked ? '已点赞' : '点赞' }}
              </el-button>
            </div>
          </div>
          
          <!-- 分类和标签 -->
          <div class="article-tags-section" v-if="article.category || article.tags">
            <el-tag v-if="article.category" type="primary" size="small">
              {{ getCategoryText(article.category) }}
            </el-tag>
            <el-tag
              v-for="tag in parseTagsArray(article.tags)"
              :key="tag"
              size="small"
              type="info"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>

        <!-- 封面图片 -->
        <div class="article-cover" v-if="article.cover_image">
          <img :src="article.cover_image" :alt="article.title" />
        </div>

        <!-- 文章内容 -->
        <div class="article-body">
          <div class="markdown-content" v-html="renderedContent"></div>
        </div>

        <!-- 评论区域 -->
        <div class="comments-section">
          <div class="comments-header">
            <h3>评论 ({{ comments.length }})</h3>
            <el-button size="small" @click="loadComments">
              <i class="el-icon-refresh" /> 刷新
            </el-button>
          </div>

          <!-- 评论输入框 -->
          <div class="comment-input-section">
            <el-input
              v-model="commentForm.content"
              type="textarea"
              :rows="4"
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
            <div v-if="comments.length === 0" class="empty-comments">
              <i class="el-icon-chat-line-square" />
              <p>暂无评论，快来抢沙发吧！</p>
            </div>

            <div v-else class="comment-list">
              <div
                v-for="comment in comments"
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
      </div>
      
      <div v-else class="error-state">
        <i class="el-icon-warning" />
        <p>文章不存在或已被删除</p>
        <el-button @click="$router.go(-1)">返回</el-button>
      </div>
    </div>
  </Layout>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import Layout from '@/components/Layout.vue'

export default {
  name: 'ArticleDetail',
  components: {
    Layout
  },
  data() {
    return {
      article: null,
      comments: [],
      loading: false,
      likeLoading: false,
      followLoading: false,
      commentSubmitting: false,
      isLiked: false,
      commentForm: {
        content: '',
        parent_id: null
      }
    }
  },
  computed: {
    ...mapGetters('auth', ['currentUser']),
    ...mapGetters('user', ['following']),
    
    renderedContent() {
      if (!this.article?.content) return ''
      // 简单的Markdown渲染，实际项目中建议使用专门的Markdown库
      return this.article.content
        .replace(/\n/g, '<br>')
        .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
        .replace(/\*(.*?)\*/g, '<em>$1</em>')
    }
  },
  async created() {
    await this.loadArticle()
    await this.loadComments()
    await this.loadFollowing()
  },
  methods: {
    ...mapActions('user', ['fetchFollowing', 'followUser', 'unfollowUser']),
    async loadArticle() {
      this.loading = true
      try {
        const articleId = this.$route.params.id
        const response = await this.$http.get(`/articles/${articleId}`)
        this.article = response.data.article
      } catch (error) {
        this.$message.error('加载文章失败')
        this.article = null
      } finally {
        this.loading = false
      }
    },

    async loadComments() {
      if (!this.article) return
      
      try {
        const response = await this.$http.get(`/articles/${this.article.id}/comments`)
        this.comments = response.data.comments || []
      } catch (error) {
        this.$message.error('加载评论失败')
        this.comments = []
      }
    },

    async toggleLike() {
      if (!this.article) return
      
      this.likeLoading = true
      try {
        const response = await this.$http.post(`/articles/${this.article.id}/like`)
        this.isLiked = response.data.liked
        
        // 更新点赞数
        if (this.isLiked) {
          this.article.like_count++
        } else {
          this.article.like_count--
        }
        
        this.$message.success(response.data.message)
      } catch (error) {
        this.$message.error('操作失败')
      } finally {
        this.likeLoading = false
      }
    },

    async submitComment() {
      if (!this.commentForm.content.trim()) {
        this.$message.warning('请输入评论内容')
        return
      }

      this.commentSubmitting = true
      try {
        await this.$http.post(`/articles/${this.article.id}/comments`, {
          content: this.commentForm.content.trim(),
          parent_id: this.commentForm.parent_id
        })

        this.$message.success('评论发表成功')
        this.commentForm.content = ''
        this.commentForm.parent_id = null
        await this.loadComments()
      } catch (error) {
        this.$message.error(error.response?.data?.error || '评论发表失败')
      } finally {
        this.commentSubmitting = false
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

        await this.$http.delete(`/article-comments/${commentId}`)
        this.$message.success('评论删除成功')
        await this.loadComments()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除评论失败')
        }
      }
    },

    editArticle() {
      this.$router.push(`/articles/${this.article.id}/edit`)
    },

    formatTime(time) {
      if (!time) return ''
      return new Date(time).toLocaleString('zh-CN')
    },

    getCategoryText(category) {
      const categoryMap = {
        tech: '技术',
        life: '生活',
        study: '学习',
        other: '其他'
      }
      return categoryMap[category] || category
    },

    parseTagsArray(tags) {
      if (!tags) return []
      try {
        return JSON.parse(tags)
      } catch {
        return tags.split(',').map(tag => tag.trim()).filter(tag => tag)
      }
    },

    async loadFollowing() {
      try {
        await this.fetchFollowing()
      } catch (error) {
        console.error('Failed to load following:', error)
      }
    },

    isFollowing(userId) {
      return this.following.some(follow => follow.followed_id === userId)
    },

    async toggleFollow(userId) {
      this.followLoading = true
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
        this.followLoading = false
      }
    }
  }
}
</script>

<style scoped>
.article-detail {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.loading {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 30px;
}

.article-content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.article-header {
  padding: 30px;
  border-bottom: 1px solid #f0f0f0;
}

.article-title {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  margin: 0 0 20px 0;
  line-height: 1.3;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
}

.author-info {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.author-name {
  font-size: 16px;
  font-weight: 500;
  color: #667eea;
}

.publish-time {
  font-size: 14px;
  color: #999;
}

.article-stats {
  display: flex;
  gap: 15px;
  margin-top: 8px;
  font-size: 14px;
  color: #666;
}

.meta-right {
  display: flex;
  gap: 10px;
}

.article-tags-section {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.article-cover {
  width: 100%;
  max-height: 400px;
  overflow: hidden;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.article-body {
  padding: 30px;
}

.markdown-content {
  font-size: 16px;
  line-height: 1.8;
  color: #333;
}

.comments-section {
  border-top: 1px solid #f0f0f0;
  padding: 30px;
}

.comments-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.comments-header h3 {
  margin: 0;
  color: #333;
}

.comment-input-section {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.comment-actions {
  margin-top: 10px;
  text-align: right;
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

.comment-item {
  padding: 20px 0;
  border-bottom: 1px solid #f5f5f5;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
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
  line-height: 1.6;
  margin-bottom: 10px;
}

.replies-list {
  margin-left: 20px;
  padding-left: 15px;
  border-left: 2px solid #f0f0f0;
  margin-top: 15px;
}

.reply-item {
  padding: 15px 0;
  border-bottom: 1px solid #f8f8f8;
}

.reply-item:last-child {
  border-bottom: none;
}

.reply-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
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
  line-height: 1.5;
  font-size: 14px;
}

.error-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  color: #999;
}

.error-state i {
  font-size: 48px;
  margin-bottom: 20px;
}

/* 移动端响应式 */
@media (max-width: 768px) {
  .article-detail {
    padding: 15px;
  }
  
  .article-header {
    padding: 20px;
  }
  
  .article-title {
    font-size: 24px;
  }
  
  .article-meta {
    flex-direction: column;
    gap: 15px;
  }
  
  .article-body {
    padding: 20px;
  }
  
  .comments-section {
    padding: 20px;
  }
  
  .comment-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .replies-list {
    margin-left: 10px;
    padding-left: 10px;
  }
}
</style>
