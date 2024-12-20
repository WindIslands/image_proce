<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { NUpload, NSelect, NInputNumber, NButton, NIcon, NSwitch, NInput, useMessage, NPagination } from 'naive-ui'
import { CloudUpload } from '@vicons/ionicons5'

// 文件列表，存储待处理的文件
const fileList = ref([])

// 输出格式，默认为png
const outputFormat = ref('png')

// 图片质量，范围1-100，默认100
const quality = ref(100)

// 是否启用格式转换
const enableFormatConversion = ref(false)

// 输出目录路径
const outputPath = ref('')

// 消息提示实例
const message = useMessage()

// 当前页码，用于分页显示
const currentPage = ref(1)

// 每页显示的文件数量
const pageSize = ref(8)

// 是否保持原图尺寸
const keepOriginalSize = ref(true)

// 自定义宽度，当keepOriginalSize为false时生效
const width = ref(0)

// 自定义高度，当keepOriginalSize为false时生效
const height = ref(0)

// 支持的输出格式选项
const formatOptions = [
  { label: 'PNG', value: 'png' },
  { label: 'JPEG', value: 'jpeg' },
  { label: 'WEBP', value: 'webp' },
  { label: 'JPG', value: 'jpg' }
]

// 添加进度相关的状态
const processing = ref(false)
const progress = ref(0)
const currentFile = ref('')

// 添加workers设置
const workers = ref(4) // 默认4个协程

// 添加事件监听器
const setupEventListeners = () => {
  window.runtime.EventsOn("process-progress", (data) => {
    progress.value = (data.current / data.total) * 100
    currentFile.value = data.file
  })

  window.runtime.EventsOn("process-error", (data) => {
    message.error(`处理文件 ${data.file} 失败: ${data.error}`)
  })

  window.runtime.EventsOn("process-complete", () => {
    processing.value = false
    progress.value = 0
    currentFile.value = ''
    message.success('所有文件处理完成')
  })
}

// 在组件挂载时设置事件监听
onMounted(() => {
  setupEventListeners()
})

// 在组件卸载时清理事件监听
onUnmounted(() => {
  window.runtime.EventsOff("process-progress")
  window.runtime.EventsOff("process-error")
  window.runtime.EventsOff("process-complete")
})

// 选择输出目录
const handleSelectOutputDir = async () => {
  try {
    const res = await window.go.app.ImageApp.OpenDirectoryDialog()
    if (res) {  
      outputPath.value = res
    }
  } catch (error) {
    message.error('选择目录失败: ' + error.message)
  }
}

// 修改重置所有设置函数
const resetSettings = () => {
  outputFormat.value = 'png'
  quality.value = 100
  enableFormatConversion.value = false
  outputPath.value = ''
  fileList.value = [] // 清空文件列表
  currentPage.value = 1 // 重置页码
  keepOriginalSize.value = true // 重置保持原尺寸
  width.value = 0 // 重置宽度
  height.value = 0 // 重置高度
  message.success('已重置所有设置')
}

// 修改处理页码改变的函数
const handlePageChange = (page) => {
  currentPage.value = page
}

// 修改文件上传处理函数
const handleFileUpload = async ()=> {
  try {
    const res = await window.go.app.ImageApp.OpenMultipleFilesDialog()
    fileList.value = res
  } catch (error) {
    message.error('选择目录失败: ' + error.message)
  }
}

// 修改清空文件列表函数
const clearFileList = () => {
  fileList.value = []
  currentPage.value = 1 // 清空时重置到第一页
  message.success('已清空文件列表')
}

// 计算当前页显示多少文件
const displayedFiles = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return fileList.value.slice(start, end)
})

// 添加删除单个文件的函数
const handleDeleteFile = (fileId) => {
  fileList.value = fileList.value.filter(file => file.id !== fileId)
  // 如果当前页没有数据了，且不是第一页，则跳转到上一页
  if (displayedFiles.value.length === 0 && currentPage.value > 1) {
    currentPage.value--
  }
  message.success('已删除文件')
}

// 修改开始转换/压缩的处理函数
const handleConvert = async () => {
  if (!outputPath.value) {
    message.error('请先选择输出目录')
    return
  }
  
  if (fileList.value.length === 0) {
    message.error('请先添加待处理文件')
    return
  }

  try {
    processing.value = true
    progress.value = 0
    currentFile.value = ''

    const tasks = fileList.value.map(file => ({
      name: file.name,
      path: file.path,
      outputFormat: enableFormatConversion.value ? outputFormat.value : '',
      quality: quality.value,
      convertFormat: enableFormatConversion.value,
      width: keepOriginalSize.value ? 0 : width.value,
      height: keepOriginalSize.value ? 0 : height.value,
      keepOriginalSize: keepOriginalSize.value
    }))

    const options = {
      tasks,
      outputPath: outputPath.value,
      workers: workers.value
    }

    await window.go.app.ImageApp.ProcessImages(options)
  } catch (error) {
    message.error('处理过程出错: ' + error.message)
  } finally {
    processing.value = false
  }
}

// 计算按钮是否禁用
const isConvertDisabled = computed(() => {
  return !outputPath.value || fileList.value.length === 0 || processing.value
})

// 添加文件大小格式化函数
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}




</script>

<template>
  <div class="container">
    <div class="header">
      <h1>图片格式转换/压缩</h1>
      <p class="subtitle">轻松转换或压缩图片格式，批量处理更高效</p>
    </div>

    <div class="content">
      <!-- 左侧上传区域 -->
      <div class="upload-container">
        <div class="upload-area" @click="handleFileUpload">
          <div class="upload-box files">
            <div class="upload-icon-wrapper">
              <NIcon size="32" class="upload-icon">
                <CloudUpload />
              </NIcon>
            </div>
            <h3>文件上传</h3>
            <p>点击选择文件到此处</p>
            <div class="supported-formats">
              支持格式：PNG、JPG、JPEG、WEBP,支持同时选择多个文件
            </div>
          </div>
        </div>

        <!-- 文件列表 -->
        <div class="file-list">
          <div class="file-list-header">
            <h3>选中的文件 {{ fileList.length ? `(${fileList.length})` : '' }}</h3>
            <NButton text type="primary" size="small" @click="clearFileList">
              清空列表
            </NButton>
          </div>
          <div class="file-list-content">
            <div class="files-wrapper">
              <template v-if="fileList.length > 0">
                <div class="file-item" v-for="file in displayedFiles" :key="file.id">
                  <div class="file-info">
                    <span class="file-name" :title="file.name">{{ file.name }}</span>
                    <div class="file-actions">
                      <span class="file-size">{{ formatFileSize(file.size) }}</span>
                      <NButton text type="error" size="tiny" @click="handleDeleteFile(file.id)">
                        删除
                      </NButton>
                    </div>
                  </div>
                </div>
              </template>
              <div v-else class="empty-tip">
                暂无选中的文件，请添加文件
              </div>
            </div>
          </div>
          <div class="pagination-wrapper">
            <NPagination
              v-model:page="currentPage"
              :page-size="pageSize"
              :item-count="fileList.length"
              :page-slot="5"
              :show-size-picker="false"
              size="small"
              @update:page="handlePageChange"
            />
          </div>
        </div>
      </div>

      <!-- 右侧设置区域 -->
      <div class="settings-container">
        <div class="settings-header">
           <h3>设置</h3>
        </div>
      
        <div class="settings-content">
          <div class="setting-item">
            <div class="setting-header">
              <div class="dir-select">
                <NInput 
                  v-model:value="outputPath" 
                  placeholder="请选择输出目录" 
                  readonly 
                  class="dir-input"
                  :title="outputPath"
                />
                <NButton secondary type="primary" @click="handleSelectOutputDir">
                  选择
                </NButton>
              </div>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-header">
              <label>启用格式转换</label>
              <NSwitch v-model:value="enableFormatConversion" />
            </div>
            <NSelect v-model:value="outputFormat" :options="formatOptions" class="setting-input"
              :disabled="!enableFormatConversion" />
          </div>

          <div class="setting-item">
            <div class="setting-header">
              <label>图片压缩质量</label>
              <NInputNumber v-model:value="quality" :min="1" :max="100" class="quality-input" suffix="%" />
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-header">
              <label>启用原图尺寸</label>
              <NSwitch v-model:value="keepOriginalSize" />
            </div>
            <div class="size-inputs">
              <div class="size-input-group">
                <NInputNumber 
                  v-model:value="width" 
                  :min="1"
                  :max="9999"
                  placeholder="宽度"
                  class="size-input"
                  :disabled="keepOriginalSize"
                >
                  <template #suffix>px</template>
                </NInputNumber>
                <span class="size-separator">×</span>
                <NInputNumber 
                  v-model:value="height" 
                  :min="1"
                  :max="9999"
                  placeholder="高度"
                  class="size-input"
                  :disabled="keepOriginalSize"
                >
                  <template #suffix>px</template>
                </NInputNumber>
              </div>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-header">
              <label>并行处理数量</label>
              <NInputNumber 
                v-model:value="workers" 
                :min="1"
                :max="16"
                class="workers-input"
              />
            </div>
            <div class="setting-tip">
              设置同时处理的图片数量，建议不超过CPU核心数
            </div>
          </div>

          <!-- 在按钮上方添加进度显示 -->
          <div v-if="processing" class="progress-info">
            <div class="progress-bar">
              <div 
                class="progress-bar-inner" 
                :style="{ width: `${progress}%` }"
              ></div>
            </div>
            <div class="progress-text">
              正在处理: {{ currentFile }}
            </div>
          </div>

          <div class="buttons-container">
            <NButton 
              type="primary" 
              size="large" 
              class="convert-button"
              :disabled="isConvertDisabled"
              :loading="processing"
              @click="handleConvert"
            >
              {{ processing ? '处理中...' : '开始转换/压缩' }}
            </NButton>
            <NButton 
              type="error" 
              size="large" 
              class="reset-button" 
              @click="resetSettings"
            >
              重置所有设置
            </NButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

.container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.header {
  text-align: center;
  padding: 20px 0;
  margin-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.header h1 {
  margin: 0;
  color: #333;
  font-size: 24px;
  font-weight: 600;
}

.subtitle {
  margin: 6px 0 0;
  color: #666;
  font-size: 14px;
}

.content {
  display: flex;
  gap: 20px;
  padding: 0 20px 20px;
  height: calc(100% - 90px);
}

.upload-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: 100%;
}

.upload-area {
  width: 100%;
  height: 160px;
  display: block;
}

:deep(.n-upload) {
  width: 100%;
  height: 100%;
  display: block;
}

:deep(.n-upload-trigger) {
  height: 100%;
  display: block;
}

:deep(.n-upload-trigger > div) {
  height: 100%;
}

.upload-box {
  height: 100%;
  width: 100%;
  box-sizing: border-box;
  border: 2px dashed #e5e5e5;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
  padding: 16px;
  margin: 0;
}

.upload-box:hover {
  border-color: #2080f0;
  background: #f0f7ff;
}

.upload-box h3 {
  margin: 12px 0 4px;
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.upload-box p {
  margin: 0;
  color: #666;
  font-size: 13px;
}

.supported-formats {
  margin-top: 8px;
  font-size: 12px;
  color: #999;
  text-align: center;
}

.upload-icon-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: #edf4ff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-icon {
  color: #2080f0;
}

.file-list {
  flex: 1;
  background: #fafafa;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
  min-height: 320px;
}

.file-list-header {
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.file-list-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  display: flex;
  align-items: center;
  gap: 4px;
}

.file-list-content {
  flex: 1;
  overflow: hidden;
}

.files-wrapper {
  height: 100%;
  overflow-y: auto;
}

.file-item {
  padding: 8px 12px;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.3s;
  width: 100%;
  box-sizing: border-box;
}

.file-item:hover {
  background-color: #f5f5f5;
}

.file-item:last-child {
  border-bottom: none;
}

.file-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.file-name {
  font-size: 13px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: 180px;
  flex-shrink: 1;
}

.file-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
  min-width: 100px;
  justify-content: flex-end;
}

.file-size {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
}

.pagination-wrapper {
  padding: 8px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
  margin-top: auto;
}

.empty-tip {
  padding: 16px;
  text-align: center;
  color: #999;
  font-size: 13px;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0;
}

.settings-container {
  width: 280px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
}

.settings-header {
  padding: 0px;
  border-bottom: 1px solid #f0f0f0;
}

.settings-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: #333;
  text-align: center;
}

.settings-content {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.setting-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.setting-item label {
  font-size: 13px;
  color: #666;
  font-weight: 500;
}

.setting-input {
  width: 100%;
}

.convert-button {
  margin-top: 12px;
  width: 100%;
  height: 40px;
  font-size: 14px;
}

.setting-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.setting-header label {
  font-size: 13px;
  color: #666;
  font-weight: 500;
  white-space: nowrap;
}

.quality-input {
  width: 100px;
  /* 调整输入框宽度 */
}

.dir-select {
  display: flex;
  gap: 8px;
  flex: 1;
}

.dir-input {
  flex: 1;
}

/* 调整输出目录的setting-header样式 */
.setting-item:has(.dir-select) .setting-header {
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.setting-item:has(.dir-select) .setting-header label {
  margin-bottom: 4px;
}

.buttons-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
}

.convert-button,
.reset-button {
  width: 100%;
  height: 40px;
  font-size: 14px;
}

/* 修改删除按钮样式 */
:deep(.n-button.n-button--text-error) {
  padding: 0 4px;
  font-size: 12px;
  height: 22px;
  line-height: 22px;
}

/* 修改分页组件样式 */
:deep(.n-pagination) {
  display: flex;
  align-items: center;
  height: 24px;
  justify-content: center;
}

/* 确保分页按钮大小一致 */
:deep(.n-pagination .n-pagination-item) {
  min-width: 28px;
  height: 24px;
  line-height: 24px;
}

/* 隐藏快速跳转页码输入框 */
:deep(.n-pagination .n-pagination-quick-jumper) {
  display: none;
}

/* 添加禁用按样式 */
:deep(.n-button--disabled) {
  cursor: not-allowed;
  opacity: 0.6;
}

/* 添加新的样式 */
.size-inputs {
  margin-top: 8px;
}

.size-input-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.size-input {
  flex: 1;
}

:deep(.size-input.n-input-number--disabled) {
  opacity: 0.8;
}

.size-separator {
  color: #666;
  font-size: 13px;
  padding: 0 2px;
}

/* 调整禁用状态的输入框样式 */
:deep(.n-input-number.n-input-number--disabled .n-input-number-input) {
  color: #999 !important;
  -webkit-text-fill-color: #999 !important;
  background-color: #f5f5f5 !important;
}

/* 如果需要调整所有开关的大小，可以添加这个样式 */
:deep(.n-switch) {
  min-width: 40px;  /* 调整所有开关的宽度 */
}

.progress-info {
  margin-bottom: 12px;
}

.progress-bar {
  height: 4px;
  background: #f0f0f0;
  border-radius: 2px;
  overflow: hidden;
}

.progress-bar-inner {
  height: 100%;
  background: #2080f0;
  transition: width 0.3s ease;
}

.progress-text {
  margin-top: 4px;
  font-size: 12px;
  color: #666;
  text-align: center;
}

.workers-input {
  width: 100px;
}

.setting-tip {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
</style>