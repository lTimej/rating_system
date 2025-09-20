<template>
  <Layout>
    <div class="profile">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="24" :md="8" :lg="8">
          <el-card class="profile-card">
            <div class="profile-header">
              <el-avatar :size="80" :src="currentUser.avatar" icon="el-icon-user-solid" />
              <h3>{{ currentUser.name || currentUser.username }}</h3>
              <el-tag :type="roleType">{{ roleText }}</el-tag>
            </div>
            
            <div class="profile-info">
              <div class="info-item">
                <label>用户名：</label>
                <span>{{ currentUser.username }}</span>
              </div>
              <div class="info-item">
                <label>邮箱：</label>
                <span>{{ currentUser.email }}</span>
              </div>
              <div class="info-item" v-if="currentUser.bio">
                <label>个人简介：</label>
                <span>{{ currentUser.bio }}</span>
              </div>
              <div class="info-item" v-if="currentUser.links">
                <label>个人链接：</label>
                <span>{{ currentUser.links }}</span>
              </div>
            </div>
            
            <el-button type="primary" @click="showEditDialog = true" class="edit-btn">
              编辑资料
            </el-button>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :sm="24" :md="16" :lg="16">
          <el-tabs v-model="activeTab" class="profile-tabs">
            <el-tab-pane label="我的作品" name="files">
              <div class="tab-header">
                <h4>我的作品文件</h4>
                <el-button type="primary" @click="showUploadDialog = true">
                  <i class="el-icon-upload" /> 上传文件
                </el-button>
              </div>
              
              <div v-if="userFiles.length === 0" class="empty-state">
                <i class="el-icon-document" />
                <p>还没有上传任何作品</p>
              </div>
              
              <div v-else class="file-grid">
                <div
                  v-for="file in userFiles"
                  :key="file.id"
                  class="file-item"
                >
                  <div class="file-icon">
                    <i :class="getFileIcon(file.file_type)" />
                  </div>
                  <div class="file-info">
                    <h5>{{ file.title || file.file_name }}</h5>
                    <p v-if="file.description">{{ file.description }}</p>
                    <div class="file-meta">
                      <span>{{ formatFileSize(file.file_size) }}</span>
                      <span>{{ formatTime(file.created_at) }}</span>
                    </div>
                  </div>
                  <div class="file-actions">
                    <el-button size="mini" type="text" @click="downloadFile(file.id)">
                      下载
                    </el-button>
                    <el-button size="mini" type="text" @click="handleDeleteFile(file.id)">
                      删除
                    </el-button>
                  </div>
                </div>
              </div>
            </el-tab-pane>
            
            <el-tab-pane label="收到的评价" name="received">
              <div v-if="allReceivedRatings.length === 0" class="empty-state">
                <i class="el-icon-star-off" />
                <p>还没有收到任何评价</p>
              </div>
              
              <div v-else class="rating-list">
                <div
                  v-for="rating in allReceivedRatings"
                  :key="`${rating.type}-${rating.id}`"
                  class="rating-item"
                >
                  <div class="rating-header">
                    <div class="rater-info">
                      <el-avatar :size="32" :src="rating.rater.avatar" icon="el-icon-user-solid" />
                      <div class="rater-details">
                        <span class="rater-name">{{ rating.rater.name || rating.rater.username }}</span>
                        <el-tag size="mini" :type="rating.type === 'file' ? 'primary' : 'success'">
                          {{ rating.typeLabel }}
                        </el-tag>
                      </div>
                    </div>
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
                    <div class="rating-meta">
                      <span v-if="rating.type === 'file' && rating.file" class="file-info">
                        <i class="el-icon-paperclip" />
                        {{ rating.file.title || rating.file.file_name }}
                      </span>
                      <span v-if="rating.type === 'comment' && rating.comment" class="comment-info">
                        <i class="el-icon-chat-line-square" />
                        评论：{{ rating.comment.content.substring(0, 30) }}...
                      </span>
                    </div>
                    <div class="feedback-stats" v-if="rating.feedbacks && rating.feedbacks.length > 0">
                      <span>{{ getHelpfulCount(rating.feedbacks) }} 人觉得有帮助</span>
                    </div>
                  </div>
                </div>
              </div>
            </el-tab-pane>
            
            <el-tab-pane label="我的评价" name="given">
              <div v-if="allGivenRatings.length === 0" class="empty-state">
                <i class="el-icon-edit-outline" />
                <p>还没有给出任何评价</p>
              </div>
              
              <div v-else class="rating-list">
                <div
                  v-for="rating in allGivenRatings"
                  :key="`${rating.type}-${rating.id}`"
                  class="rating-item"
                >
                  <div class="rating-header">
                    <div class="rated-info">
                      <span class="rated-label">评价给：</span>
                      <span class="rated-name">{{ rating.rated.name || rating.rated.username }}</span>
                      <el-tag size="mini" :type="rating.type === 'file' ? 'primary' : 'success'">
                        {{ rating.typeLabel }}
                      </el-tag>
                    </div>
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
                    <div class="rating-meta">
                      <span v-if="rating.type === 'file' && rating.file" class="file-info">
                        <i class="el-icon-paperclip" />
                        {{ rating.file.title || rating.file.file_name }}
                      </span>
                      <span v-if="rating.type === 'comment' && rating.comment" class="comment-info">
                        <i class="el-icon-chat-line-square" />
                        评论：{{ rating.comment.content.substring(0, 30) }}...
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </el-row>
    </div>

    <!-- 编辑资料对话框 -->
    <el-dialog
      title="编辑个人资料"
      :visible.sync="showEditDialog"
      width="500px"
    >
      <el-form :model="editForm" :rules="editRules" ref="editForm">
        <el-form-item label="姓名" prop="name" for="edit-name">
          <el-input v-model="editForm.name" id="edit-name" />
        </el-form-item>
        
        <el-form-item label="个人简介" prop="bio" for="edit-bio">
          <el-input
            v-model="editForm.bio"
            id="edit-bio"
            type="textarea"
            :rows="3"
            placeholder="介绍一下自己..."
          />
        </el-form-item>
        
        <el-form-item label="头像链接" prop="avatar" for="edit-avatar">
          <el-input v-model="editForm.avatar" id="edit-avatar" placeholder="头像图片链接" />
        </el-form-item>
        
        <el-form-item label="个人链接" prop="links" for="edit-links">
          <el-input
            v-model="editForm.links"
            id="edit-links"
            placeholder="个人网站、社交媒体等链接"
          />
        </el-form-item>
      </el-form>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdateProfile" :loading="updating">保存</el-button>
      </div>
    </el-dialog>

    <!-- 文件上传对话框 -->
    <el-dialog
      title="上传文件"
      :visible.sync="showUploadDialog"
      width="500px"
    >
      <el-form :model="uploadForm" :rules="uploadRules" ref="uploadForm">
        <el-form-item label="文件" prop="file" for="upload-file">
          <el-upload
            ref="upload"
            id="upload-file"
            :auto-upload="false"
            :on-change="handleFileChange"
            :file-list="fileList"
            :limit="1"
            action=""
          >
            <el-button slot="trigger" size="small" type="primary">选择文件</el-button>
          </el-upload>
        </el-form-item>
        
        <el-form-item label="标题" prop="title" for="upload-title">
          <el-input v-model="uploadForm.title" id="upload-title" placeholder="文件标题" />
        </el-form-item>
        
        <el-form-item label="描述" prop="description" for="upload-description">
          <el-input
            v-model="uploadForm.description"
            id="upload-description"
            type="textarea"
            :rows="3"
            placeholder="文件描述..."
          />
        </el-form-item>
        
        <el-form-item label="公开设置" for="upload-public">
          <el-checkbox v-model="uploadForm.is_public" id="upload-public">公开文件</el-checkbox>
        </el-form-item>
      </el-form>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="showUploadDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUploadFile" :loading="uploading">上传</el-button>
      </div>
    </el-dialog>
  </Layout>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import Layout from '@/components/Layout.vue'

export default {
  name: 'Profile',
  components: {
    Layout
  },
  data() {
    return {
      activeTab: 'files',
      showEditDialog: false,
      showUploadDialog: false,
      updating: false,
      uploading: false,
      fileList: [],
      editForm: {
        name: '',
        bio: '',
        avatar: '',
        links: ''
      },
      uploadForm: {
        file: null,
        title: '',
        description: '',
        is_public: true
      },
      editRules: {
        name: [
          { required: true, message: '请输入姓名', trigger: 'blur' }
        ]
      },
      uploadRules: {
        title: [
          { required: true, message: '请输入文件标题', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    ...mapGetters('auth', ['currentUser']),
    ...mapGetters('user', ['userFiles']),
    ...mapGetters('rating', ['receivedRatings', 'givenRatings', 'receivedCommentRatings', 'givenCommentRatings']),
    
    // 合并所有评价（文章评价 + 评论评价）
    allReceivedRatings() {
      const fileRatings = this.receivedRatings.map(rating => ({
        ...rating,
        type: 'file',
        typeLabel: '文件评价'
      }))
      const commentRatings = this.receivedCommentRatings.map(rating => ({
        ...rating,
        type: 'comment',
        typeLabel: '评论评价'
      }))
      return [...fileRatings, ...commentRatings].sort((a, b) => 
        new Date(b.created_at) - new Date(a.created_at)
      )
    },
    
    allGivenRatings() {
      const fileRatings = this.givenRatings.map(rating => ({
        ...rating,
        type: 'file',
        typeLabel: '文件评价'
      }))
      const commentRatings = this.givenCommentRatings.map(rating => ({
        ...rating,
        type: 'comment',
        typeLabel: '评论评价'
      }))
      return [...fileRatings, ...commentRatings].sort((a, b) => 
        new Date(b.created_at) - new Date(a.created_at)
      )
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
    this.initEditForm()
  },
  methods: {
    ...mapActions('auth', ['updateProfile']),
    ...mapActions('user', ['fetchUserFiles', 'uploadFile', 'deleteFile']),
    ...mapActions('rating', ['fetchUserRatings', 'fetchUserCommentRatings']),
    
    async loadData() {
      try {
        await Promise.all([
          this.fetchUserFiles(),
          this.fetchUserRatings(),
          this.fetchUserCommentRatings()
        ])
      } catch (error) {
        console.error('Failed to load data:', error)
      }
    },
    
    initEditForm() {
      this.editForm = {
        name: this.currentUser.name || '',
        bio: this.currentUser.bio || '',
        avatar: this.currentUser.avatar || '',
        links: this.currentUser.links || ''
      }
    },
    
    async handleUpdateProfile() {
      this.$refs.editForm.validate(async (valid) => {
        if (valid) {
          this.updating = true
          try {
            const result = await this.updateProfile(this.editForm)
            if (result.success) {
              this.$message.success('资料更新成功')
              this.showEditDialog = false
            } else {
              this.$message.error(result.message)
            }
          } catch (error) {
            console.error('Profile update error:', error)
            this.$message.error('更新失败，请重试')
          } finally {
            this.updating = false
          }
        }
      })
    },
    
    handleFileChange(file) {
      this.uploadForm.file = file.raw
      if (!this.uploadForm.title) {
        this.uploadForm.title = file.name
      }
    },
    
    async handleUploadFile() {
      this.$refs.uploadForm.validate(async (valid) => {
        if (valid && this.uploadForm.file) {
          this.uploading = true
          try {
            const formData = new FormData()
            formData.append('file', this.uploadForm.file)
            formData.append('title', this.uploadForm.title)
            formData.append('description', this.uploadForm.description)
            formData.append('is_public', this.uploadForm.is_public)
            
            const result = await this.uploadFile(formData)
            if (result.success) {
              this.$message.success('文件上传成功')
              this.showUploadDialog = false
              this.resetUploadForm()
            } else {
              this.$message.error(result.message)
            }
          } catch (error) {
            console.error('File upload error:', error)
            this.$message.error('上传失败，请重试')
          } finally {
            this.uploading = false
          }
        } else if (!this.uploadForm.file) {
          this.$message.error('请选择要上传的文件')
        }
      })
    },
    
    async handleDeleteFile(fileId) {
      this.$confirm('确定要删除这个文件吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          const result = await this.deleteFile(fileId)
          if (result.success) {
            this.$message.success('文件删除成功')
          } else {
            this.$message.error(result.message)
          }
        } catch (error) {
          console.error('File delete error:', error)
          this.$message.error('删除失败')
        }
      }).catch(() => {})
    },
    
    async downloadFile(fileId) {
      try {
        const token = localStorage.getItem('token')
        const response = await fetch(`/api/files/${fileId}/download`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        
        // 获取文件名
        const contentDisposition = response.headers.get('Content-Disposition')
        let filename = `file_${fileId}`
        if (contentDisposition) {
          const matches = contentDisposition.match(/filename=(.+)/)
          if (matches) {
            filename = matches[1]
          }
        }
        
        // 创建blob并下载
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = filename
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
        window.URL.revokeObjectURL(url)
      } catch (error) {
        console.error('Download error:', error)
        this.$message.error('下载失败，请重试')
      }
    },
    
    resetUploadForm() {
      this.uploadForm = {
        file: null,
        title: '',
        description: '',
        is_public: true
      }
      this.fileList = []
      this.$refs.uploadForm?.resetFields()
    },
    
    getFileIcon(fileType) {
      const iconMap = {
        image: 'el-icon-picture',
        video: 'el-icon-video-camera',
        document: 'el-icon-document',
        other: 'el-icon-files'
      }
      return iconMap[fileType] || 'el-icon-files'
    },
    
    formatFileSize(bytes) {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    },
    
    formatTime(time) {
      return new Date(time).toLocaleString('zh-CN')
    },
    
    getHelpfulCount(feedbacks) {
      return feedbacks.filter(f => f.is_helpful).length
    }
  }
}
</script>

<style scoped>
.profile {
  max-width: 1200px;
  margin: 0 auto;
}

.profile-card {
  text-align: center;
}

.profile-header {
  margin-bottom: 20px;
}

.profile-header h3 {
  margin: 15px 0 10px 0;
  font-size: 20px;
  color: #333;
}

.profile-info {
  text-align: left;
  margin-bottom: 20px;
}

.info-item {
  margin-bottom: 10px;
  display: flex;
  align-items: flex-start;
}

.info-item label {
  font-weight: 500;
  color: #666;
  min-width: 80px;
  margin-right: 10px;
}

.info-item span {
  color: #333;
  word-break: break-all;
}

.edit-btn {
  width: 100%;
}

.profile-tabs {
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

.tab-header h4 {
  margin: 0;
  color: #333;
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

.file-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.file-item {
  border: 1px solid #e6e6e6;
  border-radius: 8px;
  padding: 15px;
  transition: box-shadow 0.3s;
}

.file-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.file-icon {
  text-align: center;
  margin-bottom: 10px;
}

.file-icon i {
  font-size: 32px;
  color: #409eff;
}

.file-info h5 {
  margin: 0 0 5px 0;
  color: #333;
  font-size: 16px;
}

.file-info p {
  margin: 0 0 10px 0;
  color: #666;
  font-size: 14px;
  line-height: 1.4;
}

.file-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
  margin-bottom: 10px;
}

.file-actions {
  text-align: right;
}

.rating-list {
  max-height: 500px;
  overflow-y: auto;
}

.rating-item {
  padding: 20px 0;
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

.rater-info, .rated-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.rater-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.rating-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-top: 8px;
}

.file-info, .comment-info {
  font-size: 12px;
  color: #666;
  display: flex;
  align-items: center;
  gap: 4px;
}

.rater-name, .rated-name {
  font-weight: 500;
  color: #333;
}

.rated-label {
  color: #666;
  font-size: 14px;
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
  color: #999;
}

.feedback-stats {
  color: #67c23a;
}

/* 移动端响应式样式 */
@media (max-width: 768px) {
  .profile {
    max-width: 100%;
  }
  
  .profile-card {
    margin-bottom: 20px;
  }
  
  .profile-header {
    padding: 20px 0;
  }
  
  .profile-header h3 {
    font-size: 18px;
  }
  
  .profile-info {
    margin: 15px 0;
  }
  
  .info-item {
    margin-bottom: 10px;
  }
  
  .info-item label {
    font-size: 13px;
  }
  
  .info-item span {
    font-size: 13px;
  }
  
  .edit-btn {
    width: 100%;
    font-size: 14px;
  }
  
  .tab-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .tab-header h4 {
    font-size: 16px;
  }
  
  .file-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }
  
  .file-card :deep(.el-card__body) {
    padding: 15px;
  }
  
  .file-content {
    gap: 10px;
  }
  
  .file-info h5 {
    font-size: 14px;
  }
  
  .file-info p {
    font-size: 13px;
  }
  
  .file-meta {
    font-size: 11px;
  }
  
  .rating-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .rating-content {
    font-size: 13px;
  }
  
  .rating-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}

@media (max-width: 480px) {
  .profile-header {
    padding: 15px 0;
  }
  
  .profile-header .el-avatar {
    width: 60px !important;
    height: 60px !important;
  }
  
  .profile-header h3 {
    font-size: 16px;
  }
  
  .profile-info {
    margin: 12px 0;
  }
  
  .info-item {
    margin-bottom: 8px;
  }
  
  .info-item label {
    font-size: 12px;
  }
  
  .info-item span {
    font-size: 12px;
  }
  
  .edit-btn {
    font-size: 13px;
    padding: 8px 15px;
  }
  
  .tab-header h4 {
    font-size: 15px;
  }
  
  .file-card :deep(.el-card__body) {
    padding: 12px;
  }
  
  .file-info h5 {
    font-size: 13px;
  }
  
  .file-info p {
    font-size: 12px;
  }
  
  .file-meta {
    font-size: 10px;
  }
  
  .file-actions .el-button {
    font-size: 11px;
    padding: 4px 8px;
  }
  
  .rating-item {
    padding: 15px 0;
  }
  
  .rating-content {
    font-size: 12px;
  }
  
  .rating-footer {
    font-size: 11px;
  }
  
  .empty-state {
    padding: 40px 15px;
  }
  
  .empty-state i {
    font-size: 36px;
  }
}
</style>
