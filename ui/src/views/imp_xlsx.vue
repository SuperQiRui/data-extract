<template>
  <section>
    <el-form
      style="width: 420px;"
      size="mini"
      label-width="100px"
    >
      <el-form-item>
        <el-select
          v-model="table"
          placeholder="选择目标物理表"
          filterable
        >
          <el-option
            v-for="(t,name) in tabs"
            :key="name"
            :label="`${name} ${t.Comment}`"
            :value="name"
          />
        </el-select>
        <el-tooltip
          v-if="table"
          style="margin-left:10px"
          effect="dark"
          content="下载模板"
          placement="right"
        >
          <el-button
            :icon="Download"
            circle
            @click="handleDownloadTemplate"
          />
        </el-tooltip>
      </el-form-item>
      <el-form-item>
        <el-switch
          v-model="autoID"
          active-value="Y"
          inactive-value="N"
          active-text="自动生成ID"
        />
        <!-- <el-tooltip
          :model-value="true"
          manual
          content="是否自动生成ID"
          placement="right"
        >
          <el-switch
            v-model="autoID"
            inline-prompt
            active-value="Y"
            inactive-value="N"
            active-text="是"
            inactive-text="否"
          />
        </el-tooltip> -->
      </el-form-item>

      <el-form-item>
        <el-upload
          ref="upload"
          multiple
          action=""
          accept=".xlsx"
          :auto-upload="false"
          :on-change="handleFileChange"
          :before-remove="handleFileChange"
          :file-list="fileList"
        >
          <el-button
            type="primary"
            style="margin-left:20px"
            @click="uploadAction"
          >
            开始导入
          </el-button>
          <template #trigger>
            <el-button
              type="primary"
            >
              上传Excel
            </el-button>
          </template>
        </el-upload>
      </el-form-item>
    </el-form>
  </section>
</template>

<script>

import { defineComponent, ref } from 'vue'
import { post, download } from '@/plugin/http'

import { ElMessage } from 'element-plus'
import { Download } from '@element-plus/icons'
// import { get } from '../plugin/http'
export default defineComponent({
  props: {
    tabs: {
      type: Object,
      default: () => {}
    }
  },
  setup(props) {
    const upload = ref(null)
    const table = ref(null)
    const autoID = ref('Y')
    const fileList = ref([])

    const uploadAction = () => {
      if (!table.value) {
        ElMessage.warning(`请先指定表名`)
        return
      }
      if (fileList.value.length === 0) {
        ElMessage.warning('请先上传Excel')
        return
      }

      var formData = new FormData()
      formData.append('Table', table.value)
      formData.append('AutoID', autoID.value)
      fileList.value.forEach(file => {
        formData.append('XLSX', file.raw)
      })
      post('xlsximp', formData).then((res) => {
        upload.value?.clearFiles()
      })
    }
    const handleFileChange = (_, files) => {
      fileList.value = files
    }
    const handleDownloadTemplate = () => {
      let filename = props.tabs[table.value].Comment
      if (filename === '') {
        filename = table.value
      }
      download(`exporttemplate/${table.value}`, `${filename}.xlsx`)
    }

    return {
      upload,
      table,
      autoID,
      fileList,
      Download,
      uploadAction,
      handleFileChange,
      handleDownloadTemplate
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
}
</style>>
