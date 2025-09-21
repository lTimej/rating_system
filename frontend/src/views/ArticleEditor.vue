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

          <el-row :gutter="20" class="form-row">
            <el-col :xs="24" :sm="12" :span="12">
              <el-form-item label="分类">
                <el-select v-model="articleForm.category" placeholder="选择分类" style="width: 100%">
                  <el-option label="技术" value="tech" />
                  <el-option label="生活" value="life" />
                  <el-option label="学习" value="study" />
                  <el-option label="其他" value="other" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :span="12">
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
                <div class="toolbar-section">
                  <el-button-group class="format-buttons">
                    <el-button size="mini" @click="insertMarkdown('**', '**')" title="粗体">
                      <i class="el-icon-bold"></i>
                      <span class="button-text">粗体</span>
                    </el-button>
                    <el-button size="mini" @click="insertMarkdown('*', '*')" title="斜体">
                      <i class="el-icon-italic"></i>
                      <span class="button-text">斜体</span>
                    </el-button>
                    <el-button size="mini" @click="insertMarkdown('`', '`')" title="代码">
                      <i class="el-icon-tickets"></i>
                      <span class="button-text">代码</span>
                    </el-button>
                  </el-button-group>
                  <el-button-group class="structure-buttons">
                    <el-button size="mini" @click="insertMarkdown('### ', '')" title="标题">
                      <i class="el-icon-s-grid"></i>
                      <span class="button-text">标题</span>
                    </el-button>
                    <el-button size="mini" @click="insertMarkdown('- ', '')" title="列表">
                      <i class="el-icon-menu"></i>
                      <span class="button-text">列表</span>
                    </el-button>
                    <el-button size="mini" @click="insertMarkdown('[链接文字](', ')')" title="链接">
                      <i class="el-icon-link"></i>
                      <span class="button-text">链接</span>
                    </el-button>
                  </el-button-group>
                </div>
                <div class="toolbar-actions">
                  <!-- <el-button 
                    size="mini" 
                    :type="showSideBySide ? 'primary' : ''" 
                    @click="toggleSideBySide" 
                    title="分屏预览"
                  >
                    <i class="el-icon-view"></i>
                    <span class="button-text">{{ showSideBySide ? '关闭预览' : '分屏预览' }}</span>
                  </el-button> -->
                  <el-button size="mini" type="primary" @click="showPreview = true" title="全屏预览">
                    <i class="el-icon-view"></i>
                    <span class="button-text">预览</span>
                  </el-button>
                </div>
              </div>
              
              <!-- 分屏编辑器 -->
              <div class="editor-content-wrapper" :class="{ 'side-by-side': showSideBySide }">
                <div class="editor-panel">
                  <div class="editor-header" v-if="showSideBySide">
                    <span>编辑</span>
                  </div>
                  <el-input
                    ref="contentEditor"
                    v-model="articleForm.content"
                    type="textarea"
                    :rows="showSideBySide ? 25 : 20"
                    placeholder="请输入文章内容，支持Markdown语法&#10;&#10;支持的语法：&#10;**粗体** 或 __粗体__&#10;*斜体* 或 _斜体_&#10;`代码`&#10;### 标题&#10;- 列表项&#10;[链接文字](URL)"
                    class="markdown-editor"
                    @input="onContentChange"
                  />
                </div>
                
                <div class="preview-panel" v-if="showSideBySide">
                  <div class="editor-header">
                    <span>预览</span>
                  </div>
                  <div class="preview-content" v-html="renderedContent"></div>
                </div>
              </div>
            </div>
          </el-form-item>

          <el-form-item label="设置">
            <el-checkbox v-model="articleForm.is_public">公开文章</el-checkbox>
          </el-form-item>
        </el-form>
      </div>

      <!-- 移动端固定底部操作栏 -->
      <div class="mobile-bottom-actions" v-if="isMobile">
        <el-button @click="$router.go(-1)" size="small">取消</el-button>
        <el-button @click="saveDraft" :loading="saving" size="small">草稿</el-button>
        <el-button @click="showPreview = true" type="info" size="small">预览</el-button>
        <el-button type="primary" @click="publishArticle" :loading="publishing" size="small">
          {{ isEdit ? '更新' : '发布' }}
        </el-button>
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
      },
      isMobile: false,
      showSideBySide: false
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
      if (!this.articleForm.content) return '<p class="empty-hint">在左侧输入内容，这里会显示预览效果</p>'
      
      let content = this.articleForm.content
      
      // 处理标题
      content = content.replace(/^### (.*?)$/gm, '<h3>$1</h3>')
      content = content.replace(/^## (.*?)$/gm, '<h2>$1</h2>')
      content = content.replace(/^# (.*?)$/gm, '<h1>$1</h1>')
      
      // 处理粗体 (支持 ** 和 __ 两种语法)
      content = content.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      content = content.replace(/__(.*?)__/g, '<strong>$1</strong>')
      
      // 处理斜体 (支持 * 和 _ 两种语法，但要避免与粗体冲突)
      content = content.replace(/(?<!\*)\*([^*]+?)\*(?!\*)/g, '<em>$1</em>')
      content = content.replace(/(?<!_)_([^_]+?)_(?!_)/g, '<em>$1</em>')
      
      // 处理行内代码
      content = content.replace(/`([^`]+?)`/g, '<code>$1</code>')
      
      // 处理代码块
      content = content.replace(/```([\s\S]*?)```/g, '<pre><code>$1</code></pre>')
      
      // 处理链接
      content = content.replace(/\[([^\]]+?)\]\(([^)]+?)\)/g, '<a href="$2" target="_blank">$1</a>')
      
      // 处理列表
      content = content.replace(/^- (.*?)$/gm, '<li>$1</li>')
      content = content.replace(/(<li>.*<\/li>)/s, '<ul>$1</ul>')
      
      // 处理换行
      content = content.replace(/\n\n/g, '</p><p>')
      content = content.replace(/\n/g, '<br>')
      
      // 包装段落
      if (content && !content.startsWith('<')) {
        content = '<p>' + content + '</p>'
      }
      
      return content
    }
  },
  async created() {
    if (this.isEdit) {
      await this.loadArticle()
    }
    this.checkMobile()
  },
  mounted() {
    window.addEventListener('resize', this.checkMobile)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.checkMobile)
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
    },

    checkMobile() {
      this.isMobile = window.innerWidth <= 768
      // 移动端自动关闭分屏预览
      if (this.isMobile && this.showSideBySide) {
        this.showSideBySide = false
      }
    },

    toggleSideBySide() {
      if (this.isMobile) {
        this.$message.info('移动端请使用全屏预览')
        return
      }
      this.showSideBySide = !this.showSideBySide
    },

    onContentChange() {
      // 内容变化时的处理，可以用于自动保存等功能
      // 这里暂时留空，预留给未来功能
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

.toolbar-section {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  align-items: center;
}

.toolbar-actions {
  display: flex;
  gap: 5px;
  align-items: center;
}

.format-buttons,
.structure-buttons {
  display: flex;
  gap: 2px;
  align-items: center;
}

.editor-toolbar .el-button {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 10px;
  font-size: 12px;
  height: 28px;
}

.editor-toolbar .el-button i {
  font-size: 14px;
}

.button-text {
  font-size: 12px;
}

/* 复选框桌面端样式 */
.el-checkbox {
  font-size: 14px;
  line-height: 1.5;
}

.el-checkbox__label {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.el-checkbox__input {
  margin-right: 8px;
}

.el-checkbox__inner {
  width: 16px;
  height: 16px;
  border: 2px solid #dcdfe6;
  border-radius: 3px;
}

.el-checkbox__inner:hover {
  border-color: #409EFF;
}

.mobile-bottom-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  border-top: 1px solid #e6e6e6;
  padding: 12px 15px;
  display: flex;
  gap: 8px;
  z-index: 1000;
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(10px);
}

.mobile-bottom-actions .el-button {
  flex: 1;
  height: 40px;
  font-size: 14px;
}

/* 分屏编辑器样式 */
.editor-content-wrapper {
  display: flex;
  flex-direction: column;
}

.editor-content-wrapper.side-by-side {
  flex-direction: row;
  height: 600px;
}

.editor-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.preview-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  border-left: 1px solid #dcdfe6;
}

.editor-header {
  background: #f8f9fa;
  padding: 8px 12px;
  border-bottom: 1px solid #dcdfe6;
  font-size: 12px;
  font-weight: 600;
  color: #666;
}

.side-by-side .markdown-editor {
  height: 100%;
  flex: 1;
}

.side-by-side .markdown-editor .el-textarea__inner {
  height: 100% !important;
  resize: none;
  border: none;
  border-radius: 0;
}

.preview-content {
  flex: 1;
  padding: 15px;
  overflow-y: auto;
  background: white;
  font-size: 14px;
  line-height: 1.6;
}

/* 预览内容样式 */
.preview-content h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 20px 0 15px 0;
  color: #333;
  border-bottom: 2px solid #eee;
  padding-bottom: 8px;
}

.preview-content h2 {
  font-size: 20px;
  font-weight: 600;
  margin: 18px 0 12px 0;
  color: #333;
}

.preview-content h3 {
  font-size: 16px;
  font-weight: 600;
  margin: 15px 0 10px 0;
  color: #333;
}

.preview-content p {
  margin: 10px 0;
  color: #555;
}

.preview-content strong {
  font-weight: 600;
  color: #333;
}

.preview-content em {
  font-style: italic;
  color: #666;
}

.preview-content code {
  background: #f1f1f1;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Courier New', monospace;
  font-size: 13px;
  color: #e83e8c;
}

.preview-content pre {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  padding: 12px;
  overflow-x: auto;
  margin: 15px 0;
}

.preview-content pre code {
  background: none;
  padding: 0;
  color: #333;
  font-size: 13px;
}

.preview-content ul {
  margin: 10px 0;
  padding-left: 20px;
}

.preview-content li {
  margin: 5px 0;
  color: #555;
}

.preview-content a {
  color: #409EFF;
  text-decoration: none;
}

.preview-content a:hover {
  text-decoration: underline;
}

.preview-content .empty-hint {
  color: #999;
  font-style: italic;
  text-align: center;
  padding: 50px 20px;
}

/* Markdown编辑器增强样式 */
.markdown-editor .el-textarea__inner {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  color: #333;
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
    padding: 10px;
  }
  
  .editor-header {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
    padding: 15px;
    margin-bottom: 15px;
  }
  
  .editor-header h2 {
    font-size: 20px;
    text-align: center;
    margin: 0;
  }
  
  .header-actions {
    display: none;
  }
  
  .header-actions .el-button {
    flex: 1;
    min-height: 40px;
    font-size: 14px;
  }
  
  .editor-content {
    padding: 15px;
    padding-bottom: 120px; /* 为底部操作栏留出更多空间 */
    background: white;
    border-radius: 8px;
  }
  
  /* 表单优化 */
  .el-form {
    margin: 0;
  }
  
  .el-form-item {
    margin-bottom: 20px;
  }
  
  .el-form-item__label {
    font-size: 14px;
    font-weight: 600;
    line-height: 1.5;
    padding-bottom: 8px;
  }
  
  .el-input__inner {
    height: 44px;
    font-size: 16px;
    line-height: 44px;
  }
  
  .el-textarea__inner {
    font-size: 16px;
    line-height: 1.5;
    padding: 12px;
  }
  
  .el-select {
    width: 100%;
  }
  
  /* 分类和标签行优化 */
  .el-row {
    margin: 0 !important;
  }
  
  .el-col {
    padding: 0 !important;
    margin-bottom: 15px;
  }
  
  /* 编辑器工具栏优化 */
  .editor-container {
    border: 1px solid #dcdfe6;
    border-radius: 6px;
    overflow: hidden;
  }
  
  .editor-toolbar {
    background: #f8f9fa;
    padding: 10px;
    border-bottom: 1px solid #dcdfe6;
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }
  
  .toolbar-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 8px;
  }
  
  .toolbar-actions {
    display: flex;
    /* justify-content: center; */
    gap: 8px;
  }
  
  .format-buttons,
  .structure-buttons {
    display: flex;
    gap: 4px;
  }
  
  .editor-toolbar .el-button {
    min-width: 55px;
    height: 38px;
    font-size: 11px;
    padding: 4px 6px;
    border-radius: 4px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 2px;
    flex: 1;
    max-width: 70px;
  }
  
  .editor-toolbar .el-button .button-text {
    font-size: 10px;
    line-height: 1;
    white-space: nowrap;
  }
  
  .editor-toolbar .el-button i {
    font-size: 14px;
  }
  
  /* 内容编辑器优化 */
  .editor-container .el-textarea {
    border: none;
  }
  
  .editor-container .el-textarea__inner {
    border: none;
    border-radius: 0;
    resize: vertical;
    min-height: 300px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Courier New', monospace;
    font-size: 14px;
    line-height: 1.6;
  }
  
  /* 移动端禁用分屏预览 */
  .editor-content-wrapper.side-by-side {
    flex-direction: column;
    height: auto;
  }
  
  .preview-panel {
    border-left: none;
    border-top: 1px solid #dcdfe6;
    max-height: 300px;
  }
  
  /* 复选框优化 */
  .el-checkbox {
    font-size: 14px;
  }
  
  .el-checkbox__label {
    font-size: 14px;
    color: #333 !important;
  }
  
  .el-checkbox__input {
    margin-right: 8px;
  }
  
  .el-checkbox__inner {
    width: 16px;
    height: 16px;
    border: 2px solid #dcdfe6;
  }
}

@media (max-width: 480px) {
  .article-editor {
    padding: 5px;
  }
  
  .editor-header {
    padding: 10px;
    margin-bottom: 10px;
  }
  
  .editor-header h2 {
    font-size: 18px;
  }
  
  .header-actions .el-button {
    min-height: 36px;
    font-size: 13px;
    padding: 0 8px;
  }
  
  .editor-content {
    padding: 10px;
    padding-bottom: 100px; /* 小屏幕也需要足够的底部空间 */
  }
  
  .el-form-item {
    margin-bottom: 15px;
  }
  
  .el-form-item__label {
    font-size: 13px;
    padding-bottom: 6px;
  }
  
  .el-input__inner {
    height: 40px;
    font-size: 16px;
    line-height: 40px;
  }
  
  .el-col {
    margin-bottom: 10px;
  }
  
  .editor-toolbar {
    padding: 8px;
  }
  
  .editor-toolbar .el-button {
    min-width: 48px;
    height: 36px;
    font-size: 10px;
    padding: 2px 4px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1px;
    flex: 1;
    max-width: 60px;
  }
  
  .editor-toolbar .el-button .button-text {
    font-size: 9px;
    line-height: 1;
    white-space: nowrap;
  }
  
  .editor-toolbar .el-button i {
    font-size: 12px;
  }
  
  .editor-container .el-textarea__inner {
    min-height: 250px;
    font-size: 14px;
    padding: 10px;
  }
}

/* 预览对话框移动端优化 */
@media (max-width: 768px) {
  :deep(.el-dialog) {
    width: 95% !important;
    margin: 0 auto !important;
  }
  
  .article-preview {
    padding: 15px;
    font-size: 14px;
    line-height: 1.6;
  }
  
  .article-preview h1 {
    font-size: 22px;
    margin-bottom: 15px;
  }
  
  .preview-meta {
    font-size: 12px;
    margin-bottom: 15px;
  }
  
  .preview-content {
    font-size: 14px;
    line-height: 1.6;
  }
  
  .preview-content h1 { font-size: 20px; }
  .preview-content h2 { font-size: 18px; }
  .preview-content h3 { font-size: 16px; }
  .preview-content h4 { font-size: 15px; }
  .preview-content h5 { font-size: 14px; }
  .preview-content h6 { font-size: 13px; }
  
  .preview-content code {
    font-size: 12px;
    padding: 2px 4px;
  }
  
  .preview-content pre {
    font-size: 12px;
    padding: 10px;
    overflow-x: auto;
  }
}
</style>
