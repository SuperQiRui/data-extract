<template>
  <section>
    <el-form
      style="width: 780px"
      size="mini"
      label-width="100px"
    >
      <el-form-item>
        <el-select
          v-model="table"
          placeholder="查看物理表信息"
          filterable
        >
          <el-option
            v-for="(t,name) in tabs"
            :key="name"
            :label="`${name} ${t.Comment}`"
            :value="name"
          />
        </el-select>
        <el-alert
          v-if="tabs[table]?.Col"
          style="margin-top:10px"
          :closable="false"
          title="字段信息"
          :description="tabs[table].Col.map(c=>c.Comment?`${c.Name}(${c.Comment})`:c.Name).join(', ')"
          type="info"
        />
      </el-form-item>
      <el-form-item>
        <el-input
          v-model="ankScript"
          :rows="16"
          type="textarea"
          placeholder="输入AnkoScript..."
        />
        <el-upload
          class="load-ank"
          action=""
          :show-file-list="false"
          accept=".ank"
          :auto-upload="false"
          :on-change="handleAnkoChange"
          :limit="1"
        >
          <el-icon
            title="加载脚本"
            color="#8E8E8E"
            size="24"
          >
            <Paperclip />
          </el-icon>
        </el-upload>
        <el-icon
          class="download-ank"
          title="保存脚本"
          size="24"
          color="#8E8E8E"
          @click="handleDownloadAnko"
        >
          <Download />
        </el-icon>
      </el-form-item>

      <el-form-item>
        <el-upload
          ref="upload"
          multiple
          action=""
          accept=".docx"
          :auto-upload="false"
          :on-change="handleFileChange"
          :before-remove="handleFileChange"
          :file-list="fileList"
        >
          <el-button
            type="primary"
            style="margin-left: 20px"
            @click="uploadAction"
          >
            开始导入
          </el-button>
          <template #trigger>
            <el-button type="primary">
              上传Word
            </el-button>
          </template>
        </el-upload>
      </el-form-item>
    </el-form>
  </section>
</template>

<script>
import { defineComponent, ref } from 'vue'
import { post, saveAs } from '@/plugin/http'

import { ElMessage } from 'element-plus'
import { Paperclip, Download } from '@element-plus/icons'

export default defineComponent({
  components: {
    Paperclip, Download
  },
  props: {
    tabs: {
      type: Object,
      default: () => {}
    }
  },
  setup() {
    const upload = ref(null)
    const table = ref({})
    const ankScript = ref('')
    const fileList = ref([])

    const uploadAction = () => {
      if (fileList.value.length === 0) {
        ElMessage.warning('请先上传Excel')
        return
      }
      if (!ankScript.value) {
        ElMessage.warning('请提交AnkScript')
        return
      }

      const formData = new FormData()
      formData.append('AnkoScript', ankScript.value)
      fileList.value.forEach((file) => {
        formData.append('DOCX', file.raw)
      })
      post('docximp', formData).then((res) => {
        upload.value?.clearFiles()
      })
    }
    const handleFileChange = (_, files) => {
      fileList.value = files
    }

    const handleDownloadAnko = () => {
      if (ankScript.value) {
        saveAs(new Blob([ankScript.value]), 'docx_import.ank')
      } else {
        ElMessage.warning('未发现AnkoScript')
      }
    }

    const handleAnkoChange = (req) => {
      if (req) {
        const fr = new FileReader()
        fr.onload = () => {
          ankScript.value = fr.result
        }
        fr.readAsText(req.raw)
      }
    }

    return {
      upload,
      table,
      ankScript,
      fileList,
      uploadAction,
      handleDownloadAnko,
      handleAnkoChange,
      handleFileChange
    }
  }
})
</script>

<style lang="scss" scoped>
section {
  display: flex;
  padding-top: 40px;
  align-content: center;
  justify-content: center;
  :deep textarea {
    font-size: 18px;
    font-family: serif;
  }
  .load-ank, .download-ank {
    position: absolute;
    right: -26px;
  }
  .load-ank {
    top: 0px;
  }
  .download-ank {
    top: 40px;
    cursor: pointer;
  }
}
</style>>
