<template>
  <Layout>
    <div class="articles">
      <div class="articles-header">
        <h2>{{ isStudent ? '推送文章' : '文章列表' }}</h2>
        <p v-if="isStudent" class="student-notice">
          <i class="el-icon-info"></i>
          以下是管理员推送给您的文章
        </p>
      </div>

      <!-- 筛选器 - 只对非学生用户显示 -->
      <div v-if="!isStudent" class="filters">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-select v-model="filters.category" placeholder="选择分类" clearable @change="loadArticles">
              <el-option label="技术" value="tech" />
              <el-option label="生活" value="life" />
              <el-option label="学习" value="study" />
              <el-option label="其他" value="other" />
            </el-select>
          </el-col>
          <el-col :span="6">
            <el-input
              v-model="filters.authorId"
              placeholder="作者ID"
              clearable
              @clear="loadArticles"
              @keyup.enter="loadArticles"
            />
          </el-col>
          <el-col :span="4">
            <el-button @click="loadArticles">搜索</el-button>
          </el-col>
        </el-row>
      </div>

      <!-- 文章列表 -->
      <div class="articles-list">
        <div v-if="loading" class="loading">
          <el-skeleton :rows="5" animated />
        </div>
        
        <div v-else-if="articles.length === 0" class="empty-state">
          <i class="el-icon-document" />
          <p>{{ isStudent ? '暂无推送文章' : '暂无文章' }}</p>
          <p v-if="isStudent" class="empty-hint">管理员还没有推送文章给您</p>
        </div>
        
        <div v-else>
          <div
            v-for="article in articles"
            :key="article.id"
            class="article-item"
            @click="viewArticle(article.id)"
          >
            <div class="article-cover" v-if="article.cover_image">
              <img :src="article.cover_image" :alt="article.title" />
            </div>
            
            <div class="article-content">
              <h3 class="article-title">{{ article.title }}</h3>
              <p class="article-summary">{{ article.summary || '暂无摘要' }}</p>
              
              <div class="article-meta">
                <div class="meta-left">
                  <span class="author">
                    <i class="el-icon-user" /> {{ article.author.name || article.author.username }}
                  </span>
                  <span class="publish-time">
                    <i class="el-icon-time" /> {{ formatTime(article.published_at || article.created_at) }}
                  </span>
                  <span class="category" v-if="article.category">
                    <i class="el-icon-collection-tag" /> {{ getCategoryText(article.category) }}
                  </span>
                </div>
                
                <div class="meta-right">
                  <span class="view-count">
                    <i class="el-icon-view" /> {{ article.view_count }}
                  </span>
                  <span class="like-count">
                    <i class="el-icon-star-off" /> {{ article.like_count }}
                  </span>
                  <el-button
                    v-if="article.author.id !== currentUser.id"
                    size="mini"
                    :type="isFollowing(article.author.id) ? 'success' : 'primary'"
                    @click.stop="toggleFollow(article.author.id)"
                    :loading="followLoading[article.author.id]"
                  >
                    {{ isFollowing(article.author.id) ? '已关注' : '关注' }}
                  </el-button>
                </div>
              </div>
              
              <div class="article-tags" v-if="article.tags">
                <el-tag
                  v-for="tag in parseTagsArray(article.tags)"
                  :key="tag"
                  size="mini"
                  type="info"
                >
                  {{ tag }}
                </el-tag>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination" v-if="total > 0">
        <el-pagination
          @current-change="handlePageChange"
          :current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next, total"
        />
      </div>
    </div>
  </Layout>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import Layout from '@/components/Layout.vue'

export default {
  name: 'Articles',
  components: {
    Layout
  },
  data() {
    return {
      articles: [],
      loading: false,
      total: 0,
      currentPage: 1,
      pageSize: 10,
      followLoading: {},
      filters: {
        category: '',
        authorId: ''
      }
    }
  },
  computed: {
    ...mapGetters('auth', ['currentUser']),
    ...mapGetters('user', ['following']),
    
    isStudent() {
      return this.currentUser && this.currentUser.role === 'student'
    }
  },
  async created() {
    await this.loadArticles()
    await this.loadFollowing()
  },
  methods: {
    ...mapActions('user', ['fetchFollowing', 'followUser', 'unfollowUser']),
    async loadArticles() {
      this.loading = true
      try {
        let response
        
        if (this.isStudent) {
          // 学生用户获取推送的文章
          const params = {
            page: this.currentPage,
            limit: this.pageSize
          }
          response = await this.$http.get('/my/pushed-articles', { params })
        } else {
          // 专家和管理员获取所有公开文章
          const params = {
            page: this.currentPage,
            page_size: this.pageSize
          }
          
          if (this.filters.category) {
            params.category = this.filters.category
          }
          if (this.filters.authorId) {
            params.author_id = this.filters.authorId
          }

          response = await this.$http.get('/articles', { params })
        }
        
        this.articles = response.data.articles || []
        this.total = response.data.total || 0
      } catch (error) {
        this.$message.error('加载文章失败')
        this.articles = []
        this.total = 0
      } finally {
        this.loading = false
      }
    },

    handlePageChange(page) {
      this.currentPage = page
      this.loadArticles()
    },

    viewArticle(articleId) {
      this.$router.push(`/articles/${articleId}`)
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
    }
  }
}
</script>

<style scoped>
.articles {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.articles-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.articles-header h2 {
  margin: 0;
  color: #333;
}

.filters {
  margin-bottom: 20px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.articles-list {
  min-height: 400px;
}

.loading {
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  color: #999;
}

.empty-state i {
  font-size: 48px;
  margin-bottom: 20px;
}

.article-item {
  display: flex;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.article-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px 0 rgba(0, 0, 0, 0.15);
}

.article-cover {
  width: 200px;
  height: 120px;
  margin-right: 20px;
  border-radius: 6px;
  overflow: hidden;
  flex-shrink: 0;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.article-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.article-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 10px 0;
  line-height: 1.4;
}

.article-summary {
  color: #666;
  line-height: 1.5;
  margin: 0 0 15px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  font-size: 13px;
  color: #999;
}

.meta-left {
  display: flex;
  gap: 15px;
}

.meta-right {
  display: flex;
  gap: 15px;
  align-items: center;
}

.author {
  color: #667eea;
  font-weight: 500;
}

.article-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.pagination {
  margin-top: 30px;
  text-align: center;
}

/* 移动端响应式 */
@media (max-width: 768px) {
  .articles {
    padding: 15px;
  }
  
  .articles-header {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }
  
  .filters .el-row {
    flex-direction: column;
    gap: 10px;
  }
  
  .article-item {
    flex-direction: column;
    padding: 15px;
  }
  
  .article-cover {
    width: 100%;
    height: 200px;
    margin-right: 0;
    margin-bottom: 15px;
  }
  
  .article-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .meta-left,
  .meta-right {
    flex-direction: column;
    gap: 5px;
    align-items: flex-start;
  }
  
  .meta-right .el-button {
    margin-top: 5px;
  }
}

.student-notice {
  margin: 10px 0 0 0;
  color: #409EFF;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 5px;
}

.empty-hint {
  font-size: 14px !important;
  color: #999 !important;
  margin-top: 5px !important;
}
</style>
