<template>
  <section>
    <el-form
      style="width: 780px"
      size="mini"
      label-width="100px"
    >
      <el-form-item>
        <el-input
          v-model="sql"
          :rows="16"
          type="textarea"
          placeholder="输入sql...（多条sql之间以';'分隔）"
        />
        <el-upload
          class="load-sql"
          action=""
          :show-file-list="false"
          accept=".sql"
          :auto-upload="false"
          :on-change="handleChange"
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
          class="download-sql"
          title="保存脚本"
          size="24"
          color="#8E8E8E"
          @click="handleDownload"
        >
          <Download />
        </el-icon>
      </el-form-item>

      <el-form-item>
        <el-button
          type="primary"
          @click="downloadAction"
        >
          开始导出
        </el-button>
      </el-form-item>
    </el-form>
  </section>
</template>

<script>
import { defineComponent, ref } from 'vue'
import { download, saveAs } from '@/plugin/http'

import { ElMessage } from 'element-plus'
import { Paperclip, Download } from '@element-plus/icons'

export default defineComponent({
  components: {
    Paperclip, Download
  },
  setup() {
    const sql = ref('')

    const downloadAction = () => {
      if (!sql.value) {
        ElMessage.warning('请提交sql')
        return
      }

      download('xlsxexp', 'query_result.xlsx', sql.value.split(/\s*;\n\s*/))
    }

    const handleDownload = () => {
      if (sql.value) {
        saveAs(new Blob([sql.value]), 'query.sql')
      } else {
        ElMessage.warning('未发现sql')
      }
    }

    const handleChange = (req) => {
      if (req) {
        const fr = new FileReader()
        fr.onload = () => {
          sql.value = fr.result
        }
        fr.readAsText(req.raw)
      }
    }

    return {
      sql,
      downloadAction,
      handleDownload,
      handleChange
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
  .load-sql, .download-sql {
    position: absolute;
    right: -26px;
  }
  .load-sql {
    top: 0px;
  }
  .download-sql {
    top: 40px;
    cursor: pointer;
  }
}
</style>>
