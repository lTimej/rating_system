<template>
  <Layout>
    <div class="article-editor">
      <div class="editor-header">
        <h2>{{ isEdit ? '编辑文章' : '写文章' }}</h2>
        <div class="header-actions">
          <el-button @click="$router.go(-1)">取消</el-button>
          <el-button @click="saveDraft" :loading="saving">保存草稿</el-button>
          <el-button type="primary" @click="publishArticle" :loading="publishing">
            {{ isEdit ? '更新文章' : '发布文章' }}
          </el-button>
        </div>
      </div>

      <div class="editor-content">
        <el-form :model="articleForm" :rules="articleRules" ref="articleForm" label-width="80px">
          <el-form-item label="标题" prop="title">
            <el-input
              v-model="articleForm.title"
              placeholder="请输入文章标题"
              maxlength="100"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="摘要" prop="summary">
            <el-input
              v-model="articleForm.summary"
              type="textarea"
              :rows="3"
              placeholder="请输入文章摘要（可选）"
              maxlength="200"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="封面图片">
            <el-input
              v-model="articleForm.cover_image"
              placeholder="请输入封面图片URL（可选）"
            />
          </el-form-item>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="分类">
                <el-select v-model="articleForm.category" placeholder="选择分类">
                  <el-option label="技术" value="tech" />
                  <el-option label="生活" value="life" />
                  <el-option label="学习" value="study" />
                  <el-option label="其他" value="other" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="标签">
                <el-input
                  v-model="tagsInput"
                  placeholder="输入标签，用逗号分隔"
                  @blur="updateTags"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="内容" prop="content">
            <div class="editor-container">
              <div class="editor-toolbar">
                <el-button-group>
                  <el-button size="mini" @click="insertMarkdown('**', '**')">粗体</el-button>
                  <el-button size="mini" @click="insertMarkdown('*', '*')">斜体</el-button>
                  <el-button size="mini" @click="insertMarkdown('`', '`')">代码</el-button>
                  <el-button size="mini" @click="insertMarkdown('### ', '')">标题</el-button>
                  <el-button size="mini" @click="insertMarkdown('- ', '')">列表</el-button>
                  <el-button size="mini" @click="insertMarkdown('[链接文字](', ')')">链接</el-button>
                </el-button-group>
              </div>
              <el-input
                ref="contentEditor"
                v-model="articleForm.content"
                type="textarea"
                :rows="20"
                placeholder="请输入文章内容，支持Markdown语法"
              />
            </div>
          </el-form-item>

          <el-form-item label="设置">
            <el-checkbox v-model="articleForm.is_public">公开文章</el-checkbox>
          </el-form-item>
        </el-form>
      </div>

      <!-- 预览对话框 -->
      <el-dialog title="文章预览" :visible.sync="showPreview" width="80%">
        <div class="article-preview">
          <h1>{{ articleForm.title }}</h1>
          <div class="preview-meta">
            <span>分类：{{ getCategoryText(articleForm.category) }}</span>
            <span v-if="parsedTags.length > 0">
              标签：
              <el-tag v-for="tag in parsedTags" :key="tag" size="mini" type="info">
                {{ tag }}
              </el-tag>
            </span>
          </div>
          <div class="preview-content" v-html="renderedContent"></div>
        </div>
      </el-dialog>
    </div>
  </Layout>
</template>

<script>
import { mapGetters } from 'vuex'
import Layout from '@/components/Layout.vue'

export default {
  name: 'ArticleEditor',
  components: {
    Layout
  },
  data() {
    return {
      articleForm: {
        title: '',
        content: '',
        summary: '',
        cover_image: '',
        category: '',
        is_public: true
      },
      tagsInput: '',
      saving: false,
      publishing: false,
      showPreview: false,
      articleRules: {
        title: [
          { required: true, message: '请输入文章标题', trigger: 'blur' }
        ],
        content: [
          { required: true, message: '请输入文章内容', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    ...mapGetters('auth', ['currentUser']),
    
    isEdit() {
      return !!this.$route.params.id
    },
    
    parsedTags() {
      if (!this.tagsInput) return []
      return this.tagsInput.split(',').map(tag => tag.trim()).filter(tag => tag)
    },
    
    renderedContent() {
      if (!this.articleForm.content) return ''
      return this.articleForm.content
        .replace(/\n/g, '<br>')
        .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
        .replace(/\*(.*?)\*/g, '<em>$1</em>')
        .replace(/`(.*?)`/g, '<code>$1</code>')
        .replace(/### (.*?)$/gm, '<h3>$1</h3>')
    }
  },
  async created() {
    if (this.isEdit) {
      await this.loadArticle()
    }
  },
  methods: {
    async loadArticle() {
      try {
        const articleId = this.$route.params.id
        const response = await this.$http.get(`/articles/${articleId}`)
        const article = response.data.article
        
        // 检查权限
        if (article.author_id !== this.currentUser.id) {
          this.$message.error('您没有权限编辑此文章')
          this.$router.go(-1)
          return
        }
        
        this.articleForm = {
          title: article.title,
          content: article.content,
          summary: article.summary || '',
          cover_image: article.cover_image || '',
          category: article.category || '',
          is_public: article.is_public
        }
        
        // 解析标签
        if (article.tags) {
          try {
            const tags = JSON.parse(article.tags)
            this.tagsInput = tags.join(', ')
          } catch {
            this.tagsInput = article.tags
          }
        }
      } catch (error) {
        this.$message.error('加载文章失败')
        this.$router.go(-1)
      }
    },

    updateTags() {
      // 标签会在提交时处理
    },

    insertMarkdown(before, after) {
      const textarea = this.$refs.contentEditor.$refs.textarea
      const start = textarea.selectionStart
      const end = textarea.selectionEnd
      const selectedText = this.articleForm.content.substring(start, end)
      
      const newText = before + selectedText + after
      this.articleForm.content = 
        this.articleForm.content.substring(0, start) + 
        newText + 
        this.articleForm.content.substring(end)
      
      this.$nextTick(() => {
        textarea.focus()
        textarea.setSelectionRange(start + before.length, start + before.length + selectedText.length)
      })
    },

    async saveDraft() {
      if (!this.articleForm.title.trim()) {
        this.$message.warning('请输入文章标题')
        return
      }

      this.saving = true
      try {
        const articleData = {
          ...this.articleForm,
          tags: JSON.stringify(this.parsedTags),
          status: 'draft'
        }

        if (this.isEdit) {
          await this.$http.put(`/articles/${this.$route.params.id}`, articleData)
          this.$message.success('草稿保存成功')
        } else {
          const response = await this.$http.post('/articles', articleData)
          this.$message.success('草稿保存成功')
          // 跳转到编辑页面
          this.$router.replace(`/articles/${response.data.article.id}/edit`)
        }
      } catch (error) {
        this.$message.error(error.response?.data?.error || '保存失败')
      } finally {
        this.saving = false
      }
    },

    async publishArticle() {
      this.$refs.articleForm.validate(async (valid) => {
        if (valid) {
          this.publishing = true
          try {
            const articleData = {
              ...this.articleForm,
              tags: JSON.stringify(this.parsedTags),
              status: 'published'
            }

            if (this.isEdit) {
              await this.$http.put(`/articles/${this.$route.params.id}`, articleData)
              this.$message.success('文章更新成功')
            } else {
              const response = await this.$http.post('/articles', articleData)
              this.$message.success('文章发布成功')
              this.$router.push(`/articles/${response.data.article.id}`)
              return
            }
            
            this.$router.push(`/articles/${this.$route.params.id}`)
          } catch (error) {
            this.$message.error(error.response?.data?.error || '发布失败')
          } finally {
            this.publishing = false
          }
        }
      })
    },

    getCategoryText(category) {
      const categoryMap = {
        tech: '技术',
        life: '生活',
        study: '学习',
        other: '其他'
      }
      return categoryMap[category] || category
    }
  }
}
</script>

<style scoped>
.article-editor {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #e6e6e6;
}

.editor-header h2 {
  margin: 0;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.editor-content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 30px;
}

.editor-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

.editor-toolbar {
  background: #f5f7fa;
  padding: 10px;
  border-bottom: 1px solid #dcdfe6;
}

.editor-container .el-textarea {
  border: none;
}

.editor-container .el-textarea__inner {
  border: none;
  border-radius: 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
}

.article-preview {
  max-height: 600px;
  overflow-y: auto;
  padding: 20px;
}

.article-preview h1 {
  color: #333;
  margin-bottom: 15px;
}

.preview-meta {
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e6e6e6;
  font-size: 14px;
  color: #666;
}

.preview-meta span {
  margin-right: 15px;
}

.preview-content {
  font-size: 16px;
  line-height: 1.8;
  color: #333;
}

.preview-content h3 {
  color: #333;
  margin: 20px 0 10px 0;
}

.preview-content code {
  background: #f1f1f1;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

/* 移动端响应式 */
@media (max-width: 768px) {
  .article-editor {
    padding: 15px;
  }
  
  .editor-header {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }
  
  .header-actions {
    justify-content: center;
  }
  
  .editor-content {
    padding: 20px;
  }
  
  .editor-toolbar .el-button-group {
    flex-wrap: wrap;
  }
}
</style>
