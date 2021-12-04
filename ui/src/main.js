import { createApp } from 'vue'
import router from './router'
import 'element-plus/theme-chalk/src/index.scss'
import App from './App.vue'

const app = createApp(App)
app.use(router).mount('#app')

app.config.productionTip = false
