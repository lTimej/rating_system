<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
#app {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 20px;
  margin-bottom: 20px;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px 0;
  margin-bottom: 30px;
}

.header h1 {
  margin: 0;
  text-align: center;
  font-size: 28px;
  font-weight: 300;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 6px;
  padding: 12px 24px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
}

.form-control {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.3s ease;
}

.form-control:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.text-center {
  text-align: center;
}

.mb-20 {
  margin-bottom: 20px;
}

.mt-20 {
  margin-top: 20px;
}

/* 全局移动端响应式样式 */
@media (max-width: 768px) {
  .container {
    max-width: 100%;
    padding: 15px;
  }
  
  .card {
    padding: 15px;
    margin-bottom: 15px;
  }
  
  .header {
    padding: 15px 0;
    margin-bottom: 20px;
  }
  
  .header h1 {
    font-size: 24px;
  }
  
  .btn-primary {
    padding: 10px 20px;
    font-size: 13px;
  }
  
  .form-control {
    padding: 10px;
    font-size: 13px;
  }
  
  .form-group {
    margin-bottom: 15px;
  }
  
  .form-group label {
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .container {
    padding: 10px;
  }
  
  .card {
    padding: 12px;
    margin-bottom: 12px;
    border-radius: 6px;
  }
  
  .header {
    padding: 12px 0;
    margin-bottom: 15px;
  }
  
  .header h1 {
    font-size: 20px;
  }
  
  .btn-primary {
    padding: 8px 16px;
    font-size: 12px;
  }
  
  .form-control {
    padding: 8px;
    font-size: 12px;
  }
  
  .form-group {
    margin-bottom: 12px;
  }
  
  .form-group label {
    font-size: 12px;
    margin-bottom: 6px;
  }
}

/* Element UI 组件移动端优化 */
@media (max-width: 768px) {
  .el-dialog {
    width: 95% !important;
    margin-top: 5vh !important;
  }
  
  .el-dialog__body {
    padding: 15px 20px !important;
  }
  
  .el-form-item__label {
    font-size: 13px !important;
  }
  
  .el-input__inner {
    font-size: 13px !important;
  }
  
  .el-button {
    font-size: 13px !important;
  }
  
  .el-tabs__item {
    font-size: 13px !important;
  }
  
  .el-card__header {
    padding: 15px 20px !important;
  }
  
  .el-card__body {
    padding: 15px 20px !important;
  }
}

@media (max-width: 480px) {
  .el-dialog {
    width: 98% !important;
    margin: 2vh auto !important;
  }
  
  .el-dialog__header {
    padding: 15px 20px 10px !important;
  }
  
  .el-dialog__body {
    padding: 10px 20px !important;
  }
  
  .el-dialog__footer {
    padding: 10px 20px 20px !important;
  }
  
  .el-form-item__label {
    font-size: 12px !important;
  }
  
  .el-input__inner {
    font-size: 12px !important;
    padding: 8px 12px !important;
  }
  
  .el-button {
    font-size: 12px !important;
    padding: 8px 15px !important;
  }
  
  .el-tabs__item {
    font-size: 12px !important;
    padding: 0 15px !important;
  }
  
  .el-card__header {
    padding: 12px 15px !important;
  }
  
  .el-card__body {
    padding: 12px 15px !important;
  }
  
  .el-rate {
    font-size: 14px !important;
  }
  
  .el-rate__text {
    font-size: 12px !important;
  }
}
</style>
