<template>
  <section>
    <el-container>
      <el-header>
        <el-menu
          router
          class="nav-item"
          mode="horizontal"
        >
          <el-menu-item index="/imp_xlsx">
            Import From Excel
          </el-menu-item>
          <el-menu-item index="/imp_docx">
            Import From Word
          </el-menu-item>
          <el-menu-item index="/exp_xlsx">
            Export To Excel
          </el-menu-item>
        </el-menu>

        <el-button
          class="setting-btn"
          type="info"
          :icon="Setting"
          circle
          @click="handleSetting"
        />
      </el-header>

      <el-main>
        <router-view
          v-if="isRouterAlive"
          v-slot="{ Component }"
        >
          <transition
            name="fade"
            mode="out-in"
          >
            <component
              :is="Component"
              :tabs="tabs"
            />
          </transition>
        </router-view>
      </el-main>
    </el-container>
    <el-dialog
      v-model="showDialog"
      title="基础配置"
      width="475px"
    >
      <el-form
        ref="settingForm"
        size="mini"
        :model="conf"
        label-width="100px"
      >
        <el-form-item
          prop="Type"
          label="数据库类型"
        >
          <el-radio-group
            v-model="conf.Type"
            @change="()=>hasChange=true"
          >
            <el-radio-button label="mysql" />
            <el-radio-button label="sqlite3" />
          </el-radio-group>
        </el-form-item>
        <el-form-item
          prop="ConnectionStr"
          label="连接串"
        >
          <el-input
            v-model="conf.ConnectionStr"
            @change="()=>hasChange=true"
          />
        </el-form-item>
        <el-form-item
          v-show="conf.Type==='mysql'"
          prop="DB"
          label="数据库"
        >
          <el-select
            v-model="conf.DB"
            filterable
            @change="()=>hasChange=true"
          >
            <el-option
              v-for="(item,i) in database"
              :key="i"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          prop="PoolSize"
          label="并发通道"
        >
          <el-input-number
            v-model="conf.PoolSize"
            @change="()=>hasChange=true"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button
          type="primary"
          size="mini"
          @click="handleSubmit"
        >
          保存
        </el-button>
      </template>
    </el-dialog>
  </section>
</template>

<script>
import { defineComponent, ref } from 'vue'
// import { useRouter } from 'vue-router'
// import { ElMessageBox } from 'element-plus'
import {
  Setting,
  Expand,
  Medal,
  Avatar,
  ArrowDown,
  Right
} from '@element-plus/icons'
import { get, post } from '@/plugin/http'

export default defineComponent({
  components: {
    Setting,
    Expand,
    Medal,
    Avatar,
    ArrowDown,
    Right
  },
  setup() {
    const isRouterAlive = ref(true)
    const hasChange = ref(false)
    const showDialog = ref(false)
    const database = ref('')
    const conf = ref({})
    const tabs = ref({})

    // const get = inject('$get')
    // const router = useRouter()
    const settingForm = ref(null)

    const handleSubmit = () => {
      if (hasChange.value) {
        post('setconf', conf.value).then(_ => {
          showDialog.value = false
          getConf()
          // isRouterAlive.value = false
          // nextTick(() => {
          //   isRouterAlive.value = true
          // })
        })
      } else {
        showDialog.value = false
      }
    }
    const handleSetting = () => {
      showDialog.value = true
    }

    const getConf = () => get('getconf').then(res => {
      if (res) {
        conf.value = res['Conf']
        database.value = res['Database']
        tabs.value = res['TabInfo']
      }
    })

    getConf()

    return {
      isRouterAlive,
      showDialog,
      database,
      conf,
      tabs,
      hasChange,
      Setting,
      settingForm,
      handleSetting,
      handleSubmit
    }
  }
})
</script>

<style scoped lang="scss">
section {
  .el-main {
    height: 100%;
  }
  .nav-item {
    width: 100%;
  }
  .setting-btn {
    position: relative;
    top: -50px;
    float: right;
  }
}
</style>
