<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <h2>注册评价系统</h2>
        <p>创建您的账户，开始使用评价系统</p>
      </div>
      
      <el-form :model="registerForm" :rules="registerRules" ref="registerForm" class="register-form">
        <el-form-item label="用户名" prop="username" for="register-username">
          <el-input
            v-model="registerForm.username"
            id="register-username"
            placeholder="用户名"
          
            size="large"
          />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email" for="register-email">
          <el-input
            v-model="registerForm.email"
            id="register-email"
            placeholder="邮箱"
      
            size="large"
          />
        </el-form-item>
        
        <el-form-item label="真实姓名" prop="name" for="register-name">
          <el-input
            v-model="registerForm.name"
            id="register-name"
            placeholder="真实姓名"
     
            size="large"
          />
        </el-form-item>
        
        <el-form-item label="角色" prop="role" for="register-role">
          <el-select
            v-model="registerForm.role"
            id="register-role"
            placeholder="选择角色"
            size="large"
            style="width: 100%"
          >
            <el-option label="学生" value="student" />
            <el-option label="专家" value="expert" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="密码" prop="password" for="register-password">
          <el-input
            v-model="registerForm.password"
            id="register-password"
            type="password"
            placeholder="密码"
        
            size="large"
          />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword" for="register-confirm-password">
          <el-input
            v-model="registerForm.confirmPassword"
            id="register-confirm-password"
            type="password"
            placeholder="确认密码"
      
            size="large"
            @keyup.enter.native="handleRegister"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleRegister"
            class="register-button"
          >
            {{ loading ? '注册中...' : '注册' }}
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="register-footer">
        <p>已有账户？ <router-link to="/login">立即登录</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  name: 'Register',
  data() {
    const validateConfirmPassword = (rule, value, callback) => {
      if (value !== this.registerForm.password) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    }
    
    return {
      loading: false,
      registerForm: {
        username: '',
        email: '',
        name: '',
        role: '',
        password: '',
        confirmPassword: ''
      },
      registerRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 3, max: 20, message: '用户名长度在3到20个字符', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
        ],
        name: [
          { required: true, message: '请输入真实姓名', trigger: 'blur' }
        ],
        role: [
          { required: true, message: '请选择角色', trigger: 'change' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
        ],
        confirmPassword: [
          { required: true, message: '请确认密码', trigger: 'blur' },
          { validator: validateConfirmPassword, trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    ...mapActions('auth', ['register']),
    
    handleRegister() {
      this.$refs.registerForm.validate(async (valid) => {
        if (valid) {
          this.loading = true
          try {
            const { confirmPassword, ...registerData } = this.registerForm
            const result = await this.register(registerData)
            if (result.success) {
              this.$message.success('注册成功')
              this.$router.push('/dashboard')
            } else {
              this.$message.error(result.message)
            }
          } catch (error) {
            this.$message.error('注册失败，请重试')
          } finally {
            this.loading = false
          }
        }
      })
    }
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  padding: 40px;
  width: 100%;
  max-width: 450px;
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-header h2 {
  color: #333;
  margin-bottom: 10px;
  font-size: 24px;
  font-weight: 600;
}

.register-header p {
  color: #666;
  font-size: 14px;
  margin: 0;
}

.register-form {
  margin-bottom: 20px;
}

.register-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 500;
}

.register-footer {
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.register-footer p {
  color: #666;
  margin: 0;
}

.register-footer a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.register-footer a:hover {
  text-decoration: underline;
}

/* 移动端响应式样式 */
@media (max-width: 768px) {
  .register-container {
    padding: 15px;
  }
  
  .register-card {
    padding: 30px;
    max-width: 100%;
  }
  
  .register-header h2 {
    font-size: 22px;
  }
  
  .register-header p {
    font-size: 13px;
  }
  
  .register-button {
    height: 44px;
    font-size: 15px;
  }
}

@media (max-width: 480px) {
  .register-container {
    padding: 10px;
  }
  
  .register-card {
    padding: 25px;
    border-radius: 8px;
  }
  
  .register-header {
    margin-bottom: 25px;
  }
  
  .register-header h2 {
    font-size: 20px;
  }
  
  .register-header p {
    font-size: 12px;
  }
  
  .register-form {
    margin-bottom: 15px;
  }
  
  .register-button {
    height: 42px;
    font-size: 14px;
  }
  
  .register-footer {
    padding-top: 15px;
  }
  
  .register-footer p {
    font-size: 13px;
  }
}
</style>
