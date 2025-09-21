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
      :width="isMobile ? '95%' : '500px'"
      :fullscreen="isMobile"
      class="upload-dialog"
      :close-on-click-modal="false"
    >
      <el-form :model="uploadForm" :rules="uploadRules" ref="uploadForm" label-position="top">
        <el-form-item label="选择文件" prop="file" for="upload-file">
          <div class="upload-area">
            <el-upload
              ref="upload"
              id="upload-file"
              :auto-upload="false"
              :on-change="handleFileChange"
              :on-remove="handleFileRemove"
              :file-list="fileList"
              :limit="1"
              :accept="acceptedFileTypes"
              action=""
              drag
              class="upload-dragger"
            >
              <div class="upload-content">
                <i class="el-icon-upload upload-icon"></i>
                <div class="upload-text">
                  <p>点击或拖拽文件到此处上传</p>
                  <p class="upload-hint">支持文档、图片、视频等格式，最大50MB</p>
                </div>
              </div>
              <el-button slot="trigger" size="small" type="primary" class="upload-trigger-btn">
                <i class="el-icon-folder-opened"></i> 选择文件
              </el-button>
            </el-upload>
            
            <!-- 文件信息显示 -->
            <div v-if="uploadForm.file" class="file-preview">
              <div class="file-info">
                <i :class="getFileIconByName(uploadForm.file.name)" class="file-icon"></i>
                <div class="file-details">
                  <p class="file-name">{{ uploadForm.file.name }}</p>
                  <p class="file-size">{{ formatFileSize(uploadForm.file.size) }}</p>
                </div>
                <el-button 
                  type="text" 
                  icon="el-icon-close" 
                  @click="removeFile"
                  class="remove-btn"
                ></el-button>
              </div>
            </div>
          </div>
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
        
        <!-- 上传进度 -->
        <div v-if="uploading" class="upload-progress">
          <el-progress 
            :percentage="uploadProgress" 
            :status="uploadProgress === 100 ? 'success' : ''"
            :stroke-width="8"
          ></el-progress>
          <p class="progress-text">{{ uploadProgressText }}</p>
        </div>
      </el-form>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="cancelUpload" :disabled="uploading">{{ uploading ? '取消上传' : '取消' }}</el-button>
        <el-button 
          type="primary" 
          @click="handleUploadFile" 
          :loading="uploading"
          :disabled="!uploadForm.file || uploading"
        >
          {{ uploading ? '上传中...' : '开始上传' }}
        </el-button>
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
      },
      // 新增的数据属性
      uploadProgress: 0,
      uploadProgressText: '',
      uploadCancelToken: null,
      isMobile: false,
      acceptedFileTypes: '.pdf,.doc,.docx,.txt,.jpg,.jpeg,.png,.gif,.mp4,.avi,.mov,.zip,.rar'
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
  mounted() {
    this.checkMobile()
    window.addEventListener('resize', this.checkMobile)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.checkMobile)
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
    
    // 移动端检测
    checkMobile() {
      this.isMobile = window.innerWidth <= 768
    },
    
    // 文件选择处理
    handleFileChange(file) {
      // 文件大小检查 (50MB)
      const maxSize = 50 * 1024 * 1024 // 50MB
      if (file.size > maxSize) {
        this.$message.error('文件大小不能超过50MB')
        return false
      }
      
      this.uploadForm.file = file.raw
      if (!this.uploadForm.title) {
        this.uploadForm.title = file.name.replace(/\.[^/.]+$/, '') // 移除扩展名
      }
      
      // 更新文件列表显示
      this.fileList = [file]
    },
    
    // 文件移除处理
    handleFileRemove() {
      this.uploadForm.file = null
      this.uploadForm.title = ''
      this.fileList = []
    },
    
    // 手动移除文件
    removeFile() {
      this.handleFileRemove()
    },
    
    // 根据文件名获取图标
    getFileIconByName(fileName) {
      const ext = fileName.split('.').pop().toLowerCase()
      const iconMap = {
        // 文档类
        'pdf': 'el-icon-document',
        'doc': 'el-icon-document',
        'docx': 'el-icon-document',
        'txt': 'el-icon-document',
        'rtf': 'el-icon-document',
        // 图片类
        'jpg': 'el-icon-picture',
        'jpeg': 'el-icon-picture',
        'png': 'el-icon-picture',
        'gif': 'el-icon-picture',
        'bmp': 'el-icon-picture',
        'svg': 'el-icon-picture',
        // 视频类
        'mp4': 'el-icon-video-camera',
        'avi': 'el-icon-video-camera',
        'mov': 'el-icon-video-camera',
        'wmv': 'el-icon-video-camera',
        'flv': 'el-icon-video-camera',
        // 压缩包类
        'zip': 'el-icon-folder',
        'rar': 'el-icon-folder',
        '7z': 'el-icon-folder',
        'tar': 'el-icon-folder',
        // 其他
        'default': 'el-icon-document'
      }
      return iconMap[ext] || iconMap.default
    },
    
    async handleUploadFile() {
      this.$refs.uploadForm.validate(async (valid) => {
        if (valid && this.uploadForm.file) {
          this.uploading = true
          this.uploadProgress = 0
          this.uploadProgressText = '准备上传...'
          
          try {
            const formData = new FormData()
            formData.append('file', this.uploadForm.file)
            formData.append('title', this.uploadForm.title)
            formData.append('description', this.uploadForm.description)
            formData.append('is_public', this.uploadForm.is_public)
            
            // 模拟上传进度
            this.simulateUploadProgress()
            
            const result = await this.uploadFileWithProgress(formData)
            if (result.success) {
              this.uploadProgress = 100
              this.uploadProgressText = '上传完成！'
              
              setTimeout(() => {
                this.$message.success('文件上传成功')
                this.showUploadDialog = false
                this.resetUploadForm()
              }, 500)
            } else {
              this.$message.error(result.message || '上传失败')
            }
          } catch (error) {
            console.error('File upload error:', error)
            if (error.message === 'Upload cancelled') {
              this.$message.info('上传已取消')
            } else {
              this.$message.error('上传失败，请重试')
            }
          } finally {
            this.uploading = false
            this.uploadProgress = 0
            this.uploadProgressText = ''
          }
        } else if (!this.uploadForm.file) {
          this.$message.error('请选择要上传的文件')
        }
      })
    },
    
    // 带进度的文件上传
    async uploadFileWithProgress(formData) {
      return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest()
        
        // 保存取消令牌
        this.uploadCancelToken = () => {
          xhr.abort()
          reject(new Error('Upload cancelled'))
        }
        
        xhr.upload.addEventListener('progress', (event) => {
          if (event.lengthComputable) {
            const percentComplete = Math.round((event.loaded / event.total) * 100)
            this.uploadProgress = percentComplete
            this.uploadProgressText = `上传中... ${percentComplete}%`
          }
        })
        
        xhr.addEventListener('load', () => {
          if (xhr.status === 200) {
            try {
              const response = JSON.parse(xhr.responseText)
              // 更新store中的文件列表
              this.$store.commit('user/ADD_FILE', response.file)
              resolve({ success: true })
            } catch (error) {
              reject(error)
            }
          } else {
            try {
              const errorResponse = JSON.parse(xhr.responseText)
              resolve({ success: false, message: errorResponse.error })
            } catch (error) {
              resolve({ success: false, message: '上传失败' })
            }
          }
        })
        
        xhr.addEventListener('error', () => {
          reject(new Error('Network error'))
        })
        
        xhr.addEventListener('abort', () => {
          reject(new Error('Upload cancelled'))
        })
        
        // 获取token
        const token = localStorage.getItem('token')
        xhr.open('POST', `${process.env.VUE_APP_API_BASE_URL || 'http://localhost:8080'}/api/files`)
        xhr.setRequestHeader('Authorization', `Bearer ${token}`)
        xhr.send(formData)
      })
    },
    
    // 模拟上传进度（用于小文件快速上传的视觉反馈）
    simulateUploadProgress() {
      let progress = 0
      const interval = setInterval(() => {
        if (progress < 90 && this.uploading) {
          progress += Math.random() * 20
          if (progress > 90) progress = 90
          this.uploadProgress = Math.round(progress)
          this.uploadProgressText = `上传中... ${Math.round(progress)}%`
        } else {
          clearInterval(interval)
        }
      }, 200)
    },
    
    // 取消上传
    cancelUpload() {
      if (this.uploading && this.uploadCancelToken) {
        this.uploadCancelToken()
        this.uploadCancelToken = null
      } else {
        this.showUploadDialog = false
        this.resetUploadForm()
      }
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
      this.uploadProgress = 0
      this.uploadProgressText = ''
      this.uploadCancelToken = null
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
  
  /* 移动端文件上传优化 */
  .upload-dialog :deep(.el-dialog) {
    margin: 0 !important;
    width: 100% !important;
    height: 100% !important;
    max-width: none !important;
    border-radius: 0 !important;
  }
  
  .upload-dialog :deep(.el-dialog__body) {
    padding: 15px !important;
    max-height: calc(100vh - 120px);
    overflow-y: auto;
  }
  
  .upload-area {
    margin-bottom: 15px;
  }
  
  .upload-dragger :deep(.el-upload-dragger) {
    width: 100% !important;
    height: auto !important;
    min-height: 120px !important;
    border: 2px dashed #dcdfe6 !important;
    border-radius: 8px !important;
    padding: 20px 15px !important;
  }
  
  .upload-content {
    text-align: center;
  }
  
  .upload-icon {
    font-size: 40px !important;
    color: #c0c4cc !important;
    margin-bottom: 10px !important;
  }
  
  .upload-text p {
    margin: 5px 0 !important;
    font-size: 13px !important;
  }
  
  .upload-hint {
    color: #999 !important;
    font-size: 11px !important;
  }
  
  .upload-trigger-btn {
    margin-top: 10px !important;
    width: 100% !important;
    max-width: 200px !important;
  }
  
  .file-preview {
    margin-top: 15px;
    padding: 12px;
    background: #f5f7fa;
    border-radius: 6px;
  }
  
  .file-info {
    display: flex;
    align-items: center;
    gap: 10px;
  }
  
  .file-icon {
    font-size: 24px;
    color: #409EFF;
    flex-shrink: 0;
  }
  
  .file-details {
    flex: 1;
    min-width: 0;
  }
  
  .file-name {
    margin: 0 0 4px 0;
    font-size: 13px;
    font-weight: 500;
    color: #333;
    word-break: break-all;
  }
  
  .file-size {
    margin: 0;
    font-size: 11px;
    color: #666;
  }
  
  .remove-btn {
    color: #f56c6c !important;
    flex-shrink: 0;
  }
  
  .upload-progress {
    margin-top: 15px;
    padding: 15px;
    background: #f0f9ff;
    border-radius: 6px;
    border: 1px solid #e1f5fe;
  }
  
  .progress-text {
    margin: 8px 0 0 0;
    font-size: 12px;
    color: #666;
    text-align: center;
  }
  
  .dialog-footer {
    padding: 15px 0 0 0 !important;
    text-align: center;
  }
  
  .dialog-footer .el-button {
    width: 45% !important;
    margin: 0 2.5% !important;
  }
}

/* 文件上传样式优化 */
.upload-area {
  margin-bottom: 20px;
}

.upload-dragger :deep(.el-upload-dragger) {
  width: 100%;
  height: auto;
  min-height: 160px;
  border: 2px dashed #dcdfe6;
  border-radius: 12px;
  padding: 30px 20px;
  transition: all 0.3s ease;
}

.upload-dragger :deep(.el-upload-dragger:hover) {
  border-color: #409EFF;
  background-color: rgba(64, 158, 255, 0.05);
}

.upload-content {
  text-align: center;
}

.upload-icon {
  font-size: 48px;
  color: #c0c4cc;
  margin-bottom: 15px;
  display: block;
}

.upload-text p {
  margin: 8px 0;
  font-size: 14px;
  color: #606266;
}

.upload-hint {
  color: #909399;
  font-size: 12px;
}

.upload-trigger-btn {
  margin-top: 15px;
  min-width: 120px;
}

.file-preview {
  margin-top: 20px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-icon {
  font-size: 28px;
  color: #409EFF;
  flex-shrink: 0;
}

.file-details {
  flex: 1;
  min-width: 0;
}

.file-name {
  margin: 0 0 6px 0;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  word-break: break-all;
}

.file-size {
  margin: 0;
  font-size: 12px;
  color: #666;
}

.remove-btn {
  color: #f56c6c;
  flex-shrink: 0;
}

.remove-btn:hover {
  color: #f78989;
}

.upload-progress {
  margin-top: 20px;
  padding: 20px;
  background: #f0f9ff;
  border-radius: 8px;
  border: 1px solid #e1f5fe;
}

.progress-text {
  margin: 12px 0 0 0;
  font-size: 13px;
  color: #666;
  text-align: center;
}

/* 上传对话框样式 */
.upload-dialog :deep(.el-form-item__label) {
  font-weight: 500;
  color: #333;
}

.upload-dialog :deep(.el-form-item) {
  margin-bottom: 20px;
}

.upload-dialog :deep(.el-dialog__header) {
  padding: 20px 20px 10px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.upload-dialog :deep(.el-dialog__title) {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.upload-dialog :deep(.el-dialog__body) {
  padding: 20px;
}

.upload-dialog :deep(.el-dialog__footer) {
  padding: 15px 20px 20px 20px;
  border-top: 1px solid #f0f0f0;
}
</style>
